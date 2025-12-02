# LXD MCP Server

A Model Context Protocol (MCP) server for managing LXD containers and virtual machines.

![demo.gif](doc/demo.gif)

## Features

This MCP server provides the following tools:

### Instance Management
- **list_instances** - List all LXD instances (containers and VMs)
- **get_instance** - Get detailed information about a specific instance
- **get_instance_state** - Get current state and resource usage of an instance
- **create_instance** - Create a new LXD instance (container or VM)
- **delete_instance** - Delete an instance
- **rename_instance** - Rename an instance
- **update_instance** - Update instance configuration (e.g., CPU, memory limits)

### Instance Lifecycle
- **start_instance** - Start a stopped instance
- **stop_instance** - Stop a running instance (supports force stop)
- **restart_instance** - Restart a running instance
- **freeze_instance** - Freeze (pause) a running instance
- **unfreeze_instance** - Unfreeze (resume) a frozen instance

### Instance Operations
- **exec_instance** - Execute commands in an instance

### Image Management
- **list_images** - List available images

## Prerequisites

- LXD installed and running on the system
- Go 1.25.4 or higher
- Access to LXD Unix socket (typically requires user to be in `lxd` group)

## Installation

1. Build the server:
```bash
make build
```

2. Make sure the binary is in your PATH or note its full path for configuration.

## Usage

### With GitHub Copilot CLI

The easiest way to use this MCP server is with GitHub Copilot CLI:

1. Run the following command to add the server:
```bash
/mcp add lxd-mcp-server
```

2. When prompted, provide the full path to the `lxd-mcp-server` binary (e.g., `/home/user/git/lxd-mcp/lxd-mcp-server`)

3. The server will be automatically configured and you can start using LXD tools in your Copilot CLI session

4. To verify the server is running, you can ask Copilot to list your LXD instances or perform other LXD operations

### Direct Usage

The MCP server can also communicate directly via JSON-RPC 2.0 over stdin/stdout:

```bash
./lxd-mcp-server
```

### Example Tool Calls

#### List all instances
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "list_instances",
    "arguments": {}
  }
}
```

#### Create an instance
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "create_instance",
    "arguments": {
      "name": "my-container",
      "image": "ubuntu:24.04",
      "type": "container",
      "start": true
    }
  }
}
```

#### Get instance details
```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "get_instance",
    "arguments": {
      "name": "my-container"
    }
  }
}
```

#### Update instance configuration
```json
{
  "jsonrpc": "2.0",
  "id": 4,
  "method": "tools/call",
  "params": {
    "name": "update_instance",
    "arguments": {
      "name": "my-container",
      "config": {
        "limits.cpu": "2",
        "limits.memory": "2GiB"
      }
    }
  }
}
```

#### Start an instance
```json
{
  "jsonrpc": "2.0",
  "id": 5,
  "method": "tools/call",
  "params": {
    "name": "start_instance",
    "arguments": {
      "name": "my-container"
    }
  }
}
```

#### Execute a command
```json
{
  "jsonrpc": "2.0",
  "id": 6,
  "method": "tools/call",
  "params": {
    "name": "exec_instance",
    "arguments": {
      "name": "my-container",
      "command": ["ls", "-la", "/"]
    }
  }
}
```

#### Delete an instance
```json
{
  "jsonrpc": "2.0",
  "id": 7,
  "method": "tools/call",
  "params": {
    "name": "delete_instance",
    "arguments": {
      "name": "my-container",
      "force": true
    }
  }
}
```

## Architecture

The server uses the LXD Go client library (`github.com/canonical/lxd/client`) to interact with LXD via the Unix socket at `/var/snap/lxd/common/lxd/unix.socket` (or `/var/lib/lxd/unix.socket` for non-snap installations).

### Code Structure

- `main.go` - Main MCP server implementation
  - `LXDTools` - Struct holding LXD client connection
  - `NewLXDTools()` - Initializes connection to LXD Unix socket
  - Tool methods - Individual functions for each LXD operation:
    - `ListInstances()`, `GetInstance()`, `GetInstanceState()`
    - `CreateInstance()`, `DeleteInstance()`, `RenameInstance()`, `UpdateInstance()`
    - `StartInstance()`, `StopInstance()`, `RestartInstance()`
    - `FreezeInstance()`, `UnfreezeInstance()`
    - `ExecInstance()`, `ListImages()`
  - `main()` - Initializes MCP server and registers all tools

## Future Enhancements

Potential additions:
- Snapshot management
- File push/pull operations
- Network configuration
- Storage pool management
- Profile management
- Log access
- Image import/export

## Security Considerations

- The server requires LXD access permissions (user must be in `lxd` group)
- Command execution in containers has security implications
- Consider implementing authentication for production use
- Validate all input parameters
