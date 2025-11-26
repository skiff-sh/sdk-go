package skiff

import (
	"context"
	"io/fs"

	"github.com/skiff-sh/api/go/skiff/plugin/v1alpha1"
)

type Context struct {
	Ctx context.Context
	// The root of the project.
	Root fs.FS
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

func GetPlugin() Plugin {
	return plugin
}
