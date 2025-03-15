package clu2adv

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/mongodb-labs/atlas-cli-plugin-terraform/internal/convert"
	"github.com/mongodb-labs/atlas-cli-plugin-terraform/internal/file"
	"github.com/openai/openai-go"
	"github.com/spf13/afero"
)

type opts struct {
	fs            afero.Fs
	file          string
	output        string
	replaceOutput bool
	watch         bool
	includeMoved  bool
}

func (o *opts) PreRun() error {
	if err := file.MustExist(o.fs, o.file); err != nil {
		return err
	}
	if !o.replaceOutput {
		return file.MustNotExist(o.fs, o.output)
	}
	return nil
}

func (o *opts) Run() error {
	if err := o.generateFile(false); err != nil {
		return err
	}
	if o.watch {
		return o.watchFile()
	}
	return nil
}

func (o *opts) generateFile(allowParseErrors bool) error {
	inConfig, err := afero.ReadFile(o.fs, o.file)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", o.file, err)
	}
	outConfig, err := convert.ClusterToAdvancedCluster(inConfig, o.includeMoved)
	if err != nil {
		if allowParseErrors {
			outConfig = []byte("# CONVERT ERROR: " + err.Error() + "\n\n")
			outConfig = append(outConfig, inConfig...)
		} else {
			return err
		}
	}

	promptPrefix := "# prompt: "
	inStr := string(inConfig)
	outStr := string(outConfig)
	if strings.HasPrefix(inStr, promptPrefix) {
		prompt := strings.SplitN(inStr, "\n", 2)[0][len(promptPrefix):]
		promptOut := fmt.Sprintf(`
			Input File: """%s"""
			Output File: """%s"""
			We want to transform the Input File that is an HCL Terraform file with mongodbatlas_cluster resources into an Output File that is also an HCL Terraform file with mongodbatlas_advanced_cluster resources.
			We will ignore any resources that are not mongodbatlas_cluster or mongodbatlas_advanced_cluster and keep as it is.
			Your response must be a valid HCL Terraform file, please make sure to keep the syntax correct and write in the content of the Output File using comments with #.
			This is what I want to do in the Output File: """%s""".
		`, inStr, outStr, prompt)
		client := openai.NewClient()
		chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
			Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(promptOut),
			}),
			Model: openai.F(openai.ChatModelGPT4o),
		})
		if err == nil {
			outConfig = []byte(chatCompletion.Choices[0].Message.Content)
		} else {
			fmt.Println("Failed to get completion from OpenAI: ", err)
		}
	}

	if err := afero.WriteFile(o.fs, o.output, outConfig, 0o600); err != nil {
		return fmt.Errorf("failed to write file %s: %w", o.output, err)
	}
	return nil
}

func (o *opts) watchFile() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil
	}
	defer watcher.Close()
	if err := watcher.Add(o.file); err != nil {
		return err
	}
	for {
		if err := o.waitForFileEvent(watcher); err != nil {
			return err
		}
	}
}

func (o *opts) waitForFileEvent(watcher *fsnotify.Watcher) error {
	watcherError := errors.New("watcher has been closed")
	select {
	case event, ok := <-watcher.Events:
		if !ok {
			return watcherError
		}
		if event.Has(fsnotify.Write) {
			if err := o.generateFile(true); err != nil {
				return err
			}
		}
	case err, ok := <-watcher.Errors:
		if !ok {
			return watcherError
		}
		return err
	}
	return nil
}
