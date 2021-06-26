package git

import "github.com/ZachiNachshon/anchor/pkg/utils/shell"

type Git interface {
	Clone(s shell.Shell, clonePath string) error
	Init(s shell.Shell, path string) error
	AddOrigin(s shell.Shell, path string, url string) error
	FetchShallow(s shell.Shell, path string, url string, branch string) error
	Reset(s shell.Shell, path string, revision string) error
	Clean(s shell.Shell, path string) error
}

type gitImpl struct {
	Git
}
