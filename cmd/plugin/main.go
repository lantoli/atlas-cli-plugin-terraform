package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/mongodb-labs/atlas-cli-plugin-terraform/internal/cli/clu2adv"
	"github.com/mongodb-labs/atlas-cli-plugin-terraform/internal/convert"
	"github.com/spf13/cobra"
)

func main() {
	terraformCmd := &cobra.Command{
		Use:     "terraform",
		Short:   "Utilities for Terraform's MongoDB Atlas Provider",
		Aliases: []string{"tf"},
	}
	terraformCmd.AddCommand(clu2adv.Builder())

	mcpCmd := &cobra.Command{
		Use:   "mcp",
		Short: "Start an MCP server exposing the clu2adv tool (for LLM/tool integration)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMCPServer()
		},
	}
	terraformCmd.AddCommand(mcpCmd)

	completionOption := &cobra.CompletionOptions{
		DisableDefaultCmd:   true,
		DisableNoDescFlag:   true,
		DisableDescriptions: true,
		HiddenDefaultCmd:    true,
	}
	rootCmd := &cobra.Command{
		DisableFlagParsing: true,
		DisableAutoGenTag:  true,
		DisableSuggestions: true,
		CompletionOptions:  *completionOption,
	}
	rootCmd.AddCommand(terraformCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runMCPServer() error {
	s := server.NewMCPServer(
		"Atlas MCP Server",
		"1.0.0",
	)

	clu2advTool := mcp.NewTool(
		"clu2adv",
		mcp.WithDescription("Convert mongodbatlas_cluster config to mongodbatlas_advanced_cluster config"),
		mcp.WithString("cluster_config",
			mcp.Required(),
			mcp.Description("Terraform configuration for mongodbatlas_cluster resource as a string"),
		),
	)

	s.AddTool(clu2advTool, clu2advHandler)

	if err := server.ServeStdio(s); err != nil {
		return fmt.Errorf("Server error: %v", err)
	}
	return nil
}

func clu2advHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	clusterConfig, ok := request.Params.Arguments["cluster_config"].(string)
	if !ok {
		return nil, fmt.Errorf("cluster_config must be a string")
	}

	converted, err := convert.ClusterToAdvancedCluster([]byte(clusterConfig), true)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("conversion error: %v", err)), nil
	}

	return mcp.NewToolResultText(string(converted)), nil
}
