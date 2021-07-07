package git

import "github.com/ZachiNachshon/anchor/pkg/utils/shell"

type Git interface {
	Clone(url string, branch string, clonePath string) error
	Init(path string) error
	AddOrigin(path string, url string) error
	FetchShallow(path string, branch string) error
	Reset(path string, revision string) error
	Checkout(path string, branch string) error
	Clean(path string) error
	GetHeadCommitHash(path string, branch string) (string, error)
}

type gitImpl struct {
	Git
	shell shell.Shell
}
