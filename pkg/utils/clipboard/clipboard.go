package clipboard

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/logger"

	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
	"runtime"
)

type Clipboard interface {
	Load(content string) error
}

type clipboardImpl struct {
	shell shell.Shell
}

func New(shell shell.Shell) Clipboard {
	return &clipboardImpl{
		shell: shell,
	}
}

func (c clipboardImpl) Load(content string) error {
	switch runtime.GOOS {
	case "darwin":
		{
			if err := c.shell.Execute(fmt.Sprintf("echo \"%v\" | pbcopy", content)); err != nil {
				logger.Info("Failed setting value to clipboard using 'pbcopy ...'")
				return err
			}
			break
		}
	case "linux":
		{
			if err := c.shell.Execute(fmt.Sprintf("xclip -selection \"%v\"", content)); err != nil {
				logger.Info("Failed setting value to clipboard using 'xclip -selection ...'")
				return err
			}
			break
		}
	}
	return nil
}
