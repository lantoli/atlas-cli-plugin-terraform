package web

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"github.com/mongodb-labs/atlas-cli-plugin-terraform/internal/convert"
	"github.com/spf13/cobra"
)

//go:embed index.html
var indexHTML string

type TransformRequest struct {
	Type         string `json:"type"`
	Source       string `json:"source"`
	IncludeMoved bool   `json:"includeMoved"`
}

type TransformResponse struct {
	Result string `json:"result"`
	Error  string `json:"error,omitempty"`
}

func Builder() *cobra.Command {
	var port int
	var openBrowser bool

	cmd := &cobra.Command{
		Use:   "web",
		Short: "Launch web interface for Terraform transformations",
		Long: "Launch a web server with an interactive interface for converting Terraform configurations " +
			"between different formats (cluster to advanced_cluster and advanced_cluster v1 to v2)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWebServer(port, openBrowser)
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the web server on")
	cmd.Flags().BoolVarP(&openBrowser, "open", "o", true, "Automatically open browser")

	return cmd
}

func runWebServer(port int, openBrowser bool) error {
	mux := http.NewServeMux()

	// Serve the HTML page
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, indexHTML)
	})

	// API endpoint for transformations
	mux.HandleFunc("/api/transform", handleTransform)

	// Health check endpoint
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// Get available port if default is taken
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		// Try to find an available port
		listener, err = net.Listen("tcp", ":0")
		if err != nil {
			return fmt.Errorf("failed to start server: %w", err)
		}
	}
	defer listener.Close()

	actualPort := listener.Addr().(*net.TCPAddr).Port
	url := fmt.Sprintf("http://localhost:%d", actualPort)

	fmt.Printf("\n🚀 Terraform Transformation Web Interface\n")
	fmt.Printf("   Server running at: %s\n", url)
	fmt.Printf("   Press Ctrl+C to stop\n\n")

	// Open browser if requested
	if openBrowser {
		go func() {
			time.Sleep(100 * time.Millisecond) // Give server time to start
			openURL(url)
		}()
	}

	server := &http.Server{
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return server.Serve(listener)
}

func handleTransform(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TransformRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var result []byte
	var err error

	switch req.Type {
	case "clu2adv":
		result, err = convert.ClusterToAdvancedCluster([]byte(req.Source), req.IncludeMoved)
	case "adv2v2":
		result, err = convert.AdvancedClusterToV2([]byte(req.Source))
	default:
		http.Error(w, "Invalid transformation type", http.StatusBadRequest)
		return
	}

	resp := TransformResponse{}
	if err != nil {
		resp.Error = err.Error()
	} else {
		resp.Result = string(result)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func openURL(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default: // Linux and others
		cmd = "xdg-open"
		args = []string{url}
	}

	if err := exec.Command(cmd, args...).Start(); err != nil {
		log.Printf("Failed to open browser: %v", err)
	}
}