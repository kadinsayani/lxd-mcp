package lxd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	lxdclient "github.com/canonical/lxd/client"
	"github.com/canonical/lxd/shared/api"
	"github.com/modelcontextprotocol/go-sdk/mcp"

)

func (t *Tools) ExecInstance(ctx context.Context, req *mcp.CallToolRequest, args ExecInstanceArgs) (*mcp.CallToolResult, any, error) {
	cmdStr, _ := json.Marshal(args.Command)
	log.Printf("Executing command on instance %s: %s", args.Name, string(cmdStr))

	apiReq := api.InstanceExecPost{
		Command:     args.Command,
		WaitForWS:   true,
		Interactive: false,
	}

	var stdout, stderr bytes.Buffer
	execArgs := lxdclient.InstanceExecArgs{
		Stdout: &stdout,
		Stderr: &stderr,
	}

	op, err := t.Client.ExecInstance(args.Name, apiReq, &execArgs)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed executing command: %w", err)
	}

	if err := op.Wait(); err != nil {
		return nil, nil, fmt.Errorf("Failed waiting for command execution: %w", err)
	}

	opAPI := op.Get()
	exitCode := int(opAPI.Metadata["return"].(float64))

	result := fmt.Sprintf("Executed command: %s\nExit code: %d", string(cmdStr), exitCode)
	if stdout.Len() > 0 {
		result += fmt.Sprintf("\nStdout:\n%s", stdout.String())
	}
	if stderr.Len() > 0 {
		result += fmt.Sprintf("\nStderr:\n%s", stderr.String())
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: result},
		},
	}, nil, nil
}
