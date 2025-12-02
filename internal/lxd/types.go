package lxd

type GetInstanceArgs struct {
	Name string `json:"name" jsonschema:"Name of the instance"`
}

type CreateInstanceArgs struct {
	Name  string `json:"name" jsonschema:"Name of the instance"`
	Image string `json:"image" jsonschema:"Image to use (e.g. 'ubuntu:24.04' or 'images:alpine/3.18')"`
	Type  string `json:"type" jsonschema:"Instance type: 'container' or 'virtual-machine'"`
	Start *bool  `json:"start" jsonschema:"Start the instance after creation"`
}

type DeleteInstanceArgs struct {
	Name  string `json:"name" jsonschema:"Name of the instance"`
	Force bool   `json:"force" jsonschema:"Force deletion even if running"`
}

type StopInstanceArgs struct {
	Name  string `json:"name" jsonschema:"Name of the instance"`
	Force bool   `json:"force" jsonschema:"Force stop (immediate shutdown)"`
}

type RestartInstanceArgs struct {
	Name  string `json:"name" jsonschema:"Name of the instance"`
	Force bool   `json:"force" jsonschema:"Force restart"`
}

type RenameInstanceArgs struct {
	Name    string `json:"name" jsonschema:"Current name of the instance"`
	NewName string `json:"new_name" jsonschema:"New name for the instance"`
}

type UpdateInstanceArgs struct {
	Name   string            `json:"name" jsonschema:"Name of the instance"`
	Config map[string]string `json:"config" jsonschema:"Configuration key-value pairs (e.g. {'limits.cpu': '2' 'limits.memory': '2GiB'})"`
}

type ExecInstanceArgs struct {
	Name    string   `json:"name" jsonschema:"Name of the instance"`
	Command []string `json:"command" jsonschema:"Command to execute as array of strings"`
}
