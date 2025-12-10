package skiff

import (
	"context"
	"io/fs"

	"github.com/skiff-sh/api/go/skiff/plugin/v1alpha1"
)

type Context struct {
	Ctx context.Context

	// The current working directory of the user. This is only set if the user explicitly provides access.
	CWD *VolumeMount

	// The data provided by the user by name. The available names can be seen in the schema field of the package.
	Data map[string]*v1alpha1.Value

	// The metadata of the request.
	Metadata *v1alpha1.RequestMetadata
}

type VolumeMount struct {
	// The fs.FS representing the mount for the plugin from the host (user's machine).
	FS fs.FS

	// The absolute path to the host path for the volume. Useful for logging.
	HostPath string
}

// Plugin is capable of performing operations on files
type Plugin interface {
	// WriteFile handle a request to edit a singular file within the user's project.
	WriteFile(ctx *Context, req *v1alpha1.WriteFileRequest) (*v1alpha1.WriteFileResponse, error)
}

var plugin Plugin

func Register(p Plugin) {
	plugin = p
}
