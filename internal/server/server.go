package server

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/kadinsayani/lxd-mcp/internal/lxd"
)

const Version = "0.1.0"

func Setup(t *lxd.Tools) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "lxd-mcp-server",
		Version: Version,
	}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_instances",
		Description: "List all LXD instances (containers and VMs)",
	}, t.ListInstances)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_instance",
		Description: "Get detailed information about a specific instance",
	}, t.GetInstance)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_instance_state",
		Description: "Get current state and resource usage of an instance",
	}, t.GetInstanceState)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_instance",
		Description: "Create a new LXD instance (container or VM)",
	}, t.CreateInstance)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_instance",
		Description: "Delete an instance",
	}, t.DeleteInstance)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "start_instance",
		Description: "Start a stopped instance",
	}, t.StartInstance)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "stop_instance",
		Description: "Stop a running instance",
	}, t.StopInstance)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "restart_instance",
		Description: "Restart a running instance",
	}, t.RestartInstance)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "freeze_instance",
		Description: "Freeze (pause) a running instance",
	}, t.FreezeInstance)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "unfreeze_instance",
		Description: "Unfreeze (resume) a frozen instance",
	}, t.UnfreezeInstance)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "rename_instance",
		Description: "Rename an instance",
	}, t.RenameInstance)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "update_instance",
		Description: "Update instance configuration (e.g., CPU, memory limits)",
	}, t.UpdateInstance)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "exec_instance",
		Description: "Execute a command in an instance",
	}, t.ExecInstance)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_images",
		Description: "List available images",
	}, t.ListImages)

	return server
}
