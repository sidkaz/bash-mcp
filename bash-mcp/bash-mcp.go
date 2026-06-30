package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// -------- MCP minimal types --------

type Request struct {
	Jsonrpc string          `json:"jsonrpc"`
	ID      any             `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type Response struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      any    `json:"id,omitempty"`
	Result  any    `json:"result,omitempty"`
	Error   any    `json:"error,omitempty"`
}

// -------- MCP tool schema --------

func handleToolsList(id any) Response {
	return Response{
		Jsonrpc: "2.0",
		ID:      id,
		Result: map[string]any{
			"tools": []any{
				map[string]any{
					"name":        "bash_run",
					"description": "Execute a bash command",
					"inputSchema": map[string]any{
						"type": "object",
						"properties": map[string]any{
							"command": map[string]any{
								"type": "string",
							},
						},
						"required": []string{"command"},
					},
				},
			},
		},
	}
}

// -------- bash execution --------

func runBash(command string) (string, string, int) {
	ctx := context.Background()

	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", command)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}

	return stdout.String(), stderr.String(), exitCode
}

// -------- tool handler --------

func handleToolsCall(id any, params json.RawMessage) Response {
	var p struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}

	_ = json.Unmarshal(params, &p)

	if p.Name != "bash_run" {
		return Response{Jsonrpc: "2.0", ID: id, Error: "unknown tool"}
	}

	var args struct {
		Command string `json:"command"`
	}

	_ = json.Unmarshal(p.Arguments, &args)

	stdout, stderr, code := runBash(args.Command)

	return Response{
		Jsonrpc: "2.0",
		ID:      id,
		Result: map[string]any{
			"stdout":   stdout,
			"stderr":   stderr,
			"exitCode": code,
		},
	}
}

// -------- main MCP loop --------

func main() {
	decoder := json.NewDecoder(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)

	for {
		var req Request
		if err := decoder.Decode(&req); err != nil {
			if err == io.EOF {
				return
			}
			continue
		}

		switch req.Method {

		case "initialize":
			_ = encoder.Encode(Response{
				Jsonrpc: "2.0",
				ID:      req.ID,
				Result: map[string]any{
					"protocolVersion": "2024-11-05",
					"capabilities": map[string]any{
						"tools": map[string]any{},
					},
					"serverInfo": map[string]any{
						"name":    "minimal-bash-mcp",
						"version": "0.1.0",
					},
				},
			})

		case "tools/list":
			_ = encoder.Encode(handleToolsList(req.ID))

		case "tools/call":
			_ = encoder.Encode(handleToolsCall(req.ID, req.Params))

		default:
			// Notifications have no id; don't respond to them
			if req.ID == nil {
				continue
			}
			_ = encoder.Encode(Response{
				Jsonrpc: "2.0",
				ID:      req.ID,
				Error:   fmt.Sprintf("unknown method: %s", req.Method),
			})
		}
	}
}
