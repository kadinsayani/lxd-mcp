package lxd

import (
	"fmt"

	lxd "github.com/canonical/lxd/client"
)

// Tools provides MCP tools for LXD operations.
type Tools struct {
	Client lxd.InstanceServer
}

// NewTools creates a new [Tools] struct connected to the local LXD daemon.
func NewTools() (*Tools, error) {
	client, err := lxd.ConnectLXDUnix("", nil)
	if err != nil {
		return nil, fmt.Errorf("Failed connecting to LXD: %w", err)
	}

	return &Tools{Client: client}, nil
}

func (t *Tools) ParseImageSource(image string) (string, string, string) {
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
