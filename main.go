package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	lxd "github.com/canonical/lxd/client"
	"github.com/canonical/lxd/shared/api"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const version = "0.1.0"

type LXDTools struct {
	client lxd.InstanceServer
}

func NewLXDTools() (*LXDTools, error) {
	client, err := lxd.ConnectLXDUnix("", nil)
	if err != nil {
		return nil, fmt.Errorf("Failed connecting to LXD: %w", err)
	}

	return &LXDTools{client: client}, nil
}

func (lt *LXDTools) parseImageSource(image string) (string, string, string) {
	if len(image) == 0 {
		return "", "", ""
	}

	if len(image) > 7 && image[0:7] == "ubuntu:" {
		return "https://cloud-images.ubuntu.com/releases", image[7:], "simplestreams"
	}

	if len(image) > 7 && image[0:7] == "images:" {
		return "https://images.lxd.canonical.com", image[7:], "simplestreams"
	}

	return "", image, ""
}

func (lt *LXDTools) ListInstances(ctx context.Context, req *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
	instances, err := lt.client.GetInstances(api.InstanceTypeAny)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed listing instances: %w", err)
	}

	data, err := json.MarshalIndent(instances, "", "  ")
	if err != nil {
		return nil, nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}

type GetInstanceArgs struct {
	Name string `json:"name" jsonschema:"Name of the instance"`
}

func (lt *LXDTools) GetInstance(ctx context.Context, req *mcp.CallToolRequest, args GetInstanceArgs) (*mcp.CallToolResult, any, error) {
	instance, _, err := lt.client.GetInstance(args.Name)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed getting instance: %w", err)
	}

	data, err := json.MarshalIndent(instance, "", "  ")
	if err != nil {
		return nil, nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}

func (lt *LXDTools) GetInstanceState(ctx context.Context, req *mcp.CallToolRequest, args GetInstanceArgs) (*mcp.CallToolResult, any, error) {
	state, _, err := lt.client.GetInstanceState(args.Name)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed getting instance state: %w", err)
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return nil, nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}

type CreateInstanceArgs struct {
	Name  string `json:"name" jsonschema:"Name of the instance"`
	Image string `json:"image" jsonschema:"Image to use (e.g. 'ubuntu:24.04' or 'images:alpine/3.18')"`
	Type  string `json:"type" jsonschema:"Instance type: 'container' or 'virtual-machine'"`
	Start *bool  `json:"start" jsonschema:"Start the instance after creation"`
}

func (lt *LXDTools) CreateInstance(ctx context.Context, req *mcp.CallToolRequest, args CreateInstanceArgs) (*mcp.CallToolResult, any, error) {
	if args.Image == "" {
		return nil, nil, fmt.Errorf("Image parameter is required")
	}

	if args.Type == "" {
		args.Type = "container"
	}

	shouldStart := true
	if args.Start != nil {
		shouldStart = *args.Start
	}

	instanceType := api.InstanceTypeContainer
	if args.Type == "virtual-machine" {
		instanceType = api.InstanceTypeVM
	}

	imgServer, imgAlias, protocol := lt.parseImageSource(args.Image)

	apiReq := api.InstancesPost{
		Name: args.Name,
		Type: instanceType,
		Source: api.InstanceSource{
			Type:     "image",
			Server:   imgServer,
			Protocol: protocol,
			Alias:    imgAlias,
		},
	}

	op, err := lt.client.CreateInstance(apiReq)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed creating instance: %w", err)
	}

	if err := op.Wait(); err != nil {
		return nil, nil, fmt.Errorf("Failed waiting for instance creation: %w", err)
	}

	result := fmt.Sprintf("Instance %s created successfully", args.Name)

	if shouldStart {
		startReq := api.InstanceStatePut{
			Action:  "start",
			Timeout: -1,
		}

		op, err := lt.client.UpdateInstanceState(args.Name, startReq, "")
		if err != nil {
			return nil, nil, fmt.Errorf("Instance created but failed starting: %w", err)
		}

		if err := op.Wait(); err != nil {
			return nil, nil, fmt.Errorf("Instance created but failed waiting for start: %w", err)
		}

		result += " and started"
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: result},
		},
	}, nil, nil
}

type DeleteInstanceArgs struct {
	Name  string `json:"name" jsonschema:"Name of the instance"`
	Force bool   `json:"force" jsonschema:"Force deletion even if running"`
}

func (lt *LXDTools) DeleteInstance(ctx context.Context, req *mcp.CallToolRequest, args DeleteInstanceArgs) (*mcp.CallToolResult, any, error) {
	if args.Force {
		state, _, err := lt.client.GetInstanceState(args.Name)
		if err != nil {
			return nil, nil, fmt.Errorf("Failed getting instance state: %w", err)
		}

		if state.Status == "Running" {
			stopReq := api.InstanceStatePut{
				Action:  "stop",
				Timeout: -1,
				Force:   true,
			}

			op, err := lt.client.UpdateInstanceState(args.Name, stopReq, "")
			if err != nil {
				return nil, nil, fmt.Errorf("Failed stopping instance: %w", err)
			}

			if err := op.Wait(); err != nil {
				return nil, nil, fmt.Errorf("Failed waiting for instance stop: %w", err)
			}
		}
	}

	op, err := lt.client.DeleteInstance(args.Name)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed deleting instance: %w", err)
	}

	if err := op.Wait(); err != nil {
		return nil, nil, fmt.Errorf("Failed waiting for instance deletion: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Instance %s deleted successfully", args.Name)},
		},
	}, nil, nil
}

func (lt *LXDTools) StartInstance(ctx context.Context, req *mcp.CallToolRequest, args GetInstanceArgs) (*mcp.CallToolResult, any, error) {
	apiReq := api.InstanceStatePut{
		Action:  "start",
		Timeout: -1,
	}

	op, err := lt.client.UpdateInstanceState(args.Name, apiReq, "")
	if err != nil {
		return nil, nil, fmt.Errorf("Failed starting instance: %w", err)
	}

	if err := op.Wait(); err != nil {
		return nil, nil, fmt.Errorf("Failed waiting for instance start: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Instance %s started successfully", args.Name)},
		},
	}, nil, nil
}

type StopInstanceArgs struct {
	Name  string `json:"name" jsonschema:"Name of the instance"`
	Force bool   `json:"force" jsonschema:"Force stop (immediate shutdown)"`
}

func (lt *LXDTools) StopInstance(ctx context.Context, req *mcp.CallToolRequest, args StopInstanceArgs) (*mcp.CallToolResult, any, error) {
	apiReq := api.InstanceStatePut{
		Action:  "stop",
		Timeout: -1,
		Force:   args.Force,
	}

	op, err := lt.client.UpdateInstanceState(args.Name, apiReq, "")
	if err != nil {
		return nil, nil, fmt.Errorf("Failed stopping instance: %w", err)
	}

	if err := op.Wait(); err != nil {
		return nil, nil, fmt.Errorf("Failed waiting for instance stop: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Instance %s stopped successfully", args.Name)},
		},
	}, nil, nil
}

type RestartInstanceArgs struct {
	Name  string `json:"name" jsonschema:"Name of the instance"`
	Force bool   `json:"force" jsonschema:"Force restart"`
}

func (lt *LXDTools) RestartInstance(ctx context.Context, req *mcp.CallToolRequest, args RestartInstanceArgs) (*mcp.CallToolResult, any, error) {
	apiReq := api.InstanceStatePut{
		Action:  "restart",
		Timeout: -1,
		Force:   args.Force,
	}

	op, err := lt.client.UpdateInstanceState(args.Name, apiReq, "")
	if err != nil {
		return nil, nil, fmt.Errorf("Failed restarting instance: %w", err)
	}

	if err := op.Wait(); err != nil {
		return nil, nil, fmt.Errorf("Failed waiting for instance restart: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Instance %s restarted successfully", args.Name)},
		},
	}, nil, nil
}

func (lt *LXDTools) FreezeInstance(ctx context.Context, req *mcp.CallToolRequest, args GetInstanceArgs) (*mcp.CallToolResult, any, error) {
	apiReq := api.InstanceStatePut{
		Action:  "freeze",
		Timeout: -1,
	}

	op, err := lt.client.UpdateInstanceState(args.Name, apiReq, "")
	if err != nil {
		return nil, nil, fmt.Errorf("Failed freezing instance: %w", err)
	}

	if err := op.Wait(); err != nil {
		return nil, nil, fmt.Errorf("Failed waiting for instance freeze: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Instance %s frozen successfully", args.Name)},
		},
	}, nil, nil
}

func (lt *LXDTools) UnfreezeInstance(ctx context.Context, req *mcp.CallToolRequest, args GetInstanceArgs) (*mcp.CallToolResult, any, error) {
	apiReq := api.InstanceStatePut{
		Action:  "unfreeze",
		Timeout: -1,
	}

	op, err := lt.client.UpdateInstanceState(args.Name, apiReq, "")
	if err != nil {
		return nil, nil, fmt.Errorf("Failed unfreezing instance: %w", err)
	}

	if err := op.Wait(); err != nil {
		return nil, nil, fmt.Errorf("Failed waiting for instance unfreeze: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Instance %s unfrozen successfully", args.Name)},
		},
	}, nil, nil
}

type RenameInstanceArgs struct {
	Name    string `json:"name" jsonschema:"Current name of the instance"`
	NewName string `json:"new_name" jsonschema:"New name for the instance"`
}

func (lt *LXDTools) RenameInstance(ctx context.Context, req *mcp.CallToolRequest, args RenameInstanceArgs) (*mcp.CallToolResult, any, error) {
	apiReq := api.InstancePost{
		Name: args.NewName,
	}

	op, err := lt.client.RenameInstance(args.Name, apiReq)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed renaming instance: %w", err)
	}

	if err := op.Wait(); err != nil {
		return nil, nil, fmt.Errorf("Failed waiting for instance rename: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Instance %s renamed to %s successfully", args.Name, args.NewName)},
		},
	}, nil, nil
}

type UpdateInstanceArgs struct {
	Name   string            `json:"name" jsonschema:"Name of the instance"`
	Config map[string]string `json:"config" jsonschema:"Configuration key-value pairs (e.g. {'limits.cpu': '2' 'limits.memory': '2GiB'})"`
}

func (lt *LXDTools) UpdateInstance(ctx context.Context, req *mcp.CallToolRequest, args UpdateInstanceArgs) (*mcp.CallToolResult, any, error) {
	instance, etag, err := lt.client.GetInstance(args.Name)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed getting instance: %w", err)
	}

	if instance.Config == nil {
		instance.Config = make(map[string]string)
	}

	for key, value := range args.Config {
		instance.Config[key] = value
	}

	op, err := lt.client.UpdateInstance(args.Name, instance.Writable(), etag)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed updating instance: %w", err)
	}

	if err := op.Wait(); err != nil {
		return nil, nil, fmt.Errorf("Failed waiting for instance update: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Instance %s updated successfully", args.Name)},
		},
	}, nil, nil
}

type ExecInstanceArgs struct {
	Name    string   `json:"name" jsonschema:"Name of the instance"`
	Command []string `json:"command" jsonschema:"Command to execute as array of strings"`
}

func (lt *LXDTools) ExecInstance(ctx context.Context, req *mcp.CallToolRequest, args ExecInstanceArgs) (*mcp.CallToolResult, any, error) {
	apiReq := api.InstanceExecPost{
		Command:     args.Command,
		WaitForWS:   true,
		Interactive: false,
	}

	var stdout, stderr bytes.Buffer
	execArgs := lxd.InstanceExecArgs{
		Stdout: &stdout,
		Stderr: &stderr,
	}

	op, err := lt.client.ExecInstance(args.Name, apiReq, &execArgs)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed executing command: %w", err)
	}

	if err := op.Wait(); err != nil {
		return nil, nil, fmt.Errorf("Failed waiting for command execution: %w", err)
	}

	opAPI := op.Get()
	exitCode := int(opAPI.Metadata["return"].(float64))

	result := fmt.Sprintf("Command executed with exit code: %d", exitCode)
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

func (lt *LXDTools) ListImages(ctx context.Context, req *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
	images, err := lt.client.GetImages()
	if err != nil {
		return nil, nil, fmt.Errorf("Failed listing images: %w", err)
	}

	data, err := json.MarshalIndent(images, "", "  ")
	if err != nil {
		return nil, nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}

func main() {
	tools, err := NewLXDTools()
	if err != nil {
		log.Fatalf("Failed creating LXD tools: %v", err)
	}

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "lxd-mcp-server",
		Version: version,
	}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_instances",
		Description: "List all LXD instances (containers and VMs)",
	}, tools.ListInstances)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_instance",
		Description: "Get detailed information about a specific instance",
	}, tools.GetInstance)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_instance_state",
		Description: "Get current state and resource usage of an instance",
	}, tools.GetInstanceState)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_instance",
		Description: "Create a new LXD instance (container or VM)",
	}, tools.CreateInstance)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_instance",
		Description: "Delete an instance",
	}, tools.DeleteInstance)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "start_instance",
		Description: "Start a stopped instance",
	}, tools.StartInstance)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "stop_instance",
		Description: "Stop a running instance",
	}, tools.StopInstance)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "restart_instance",
		Description: "Restart a running instance",
	}, tools.RestartInstance)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "freeze_instance",
		Description: "Freeze (pause) a running instance",
	}, tools.FreezeInstance)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "unfreeze_instance",
		Description: "Unfreeze (resume) a frozen instance",
	}, tools.UnfreezeInstance)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "rename_instance",
		Description: "Rename an instance",
	}, tools.RenameInstance)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "update_instance",
		Description: "Update instance configuration (e.g., CPU, memory limits)",
	}, tools.UpdateInstance)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "exec_instance",
		Description: "Execute a command in an instance",
	}, tools.ExecInstance)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_images",
		Description: "List available images",
	}, tools.ListImages)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
