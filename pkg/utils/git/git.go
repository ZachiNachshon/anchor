package git

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
)

func New() Git {
	return gitImpl{}
}

func (g *gitImpl) GitClone(s shell.Shell, url string, branch string, clonePath string) error {
	err := g.GitInit(s, clonePath)
	if err != nil {
		return err
	}

	err = g.GitAddOrigin(s, clonePath, url)
	if err != nil {
		return err
	}

	err = g.GitFetchShallow(s, clonePath, url, branch)
	if err != nil {
		return err
	}

	err = g.GitClean(s, clonePath)
	if err != nil {
		return err
	}
	return nil
}

func (g *gitImpl) GitInit(s shell.Shell, path string) error {
	logger.Infof("Git init a new index. path: %s", path)
	script := fmt.Sprintf("git init %s", path)
	return s.Execute(script)
}

func (g *gitImpl) GitAddOrigin(s shell.Shell, path string, url string) error {
	logger.Infof("Git add remote origin. url: %s", url)
	script := fmt.Sprintf("git -C %s remote add origin %s", path, url)
	return s.Execute(script)
}

func (g *gitImpl) GitFetchShallow(s shell.Shell, path string, url string, branch string) error {
	logger.Infof("Git fetching branch with shallow refs. branch: %s, since: 4 weeks ago", branch)
	script := fmt.Sprintf(`git -C %s fetch --shallow-since="4 weeks ago" --force origin refs/heads/%s:refs/remotes/origin/%s`, path, url, branch)
	return s.Execute(script)
}

func (g *gitImpl) GitReset(s shell.Shell, path string, revision string) error {
	logger.Infof("Git reset to a specific revision. commit: %s", revision)
	script := fmt.Sprintf(`git -C %s reset --hard "%s"`, path, revision)
	return s.Execute(script)
}

func (g *gitImpl) GitClean(s shell.Shell, path string) error {
	logger.Infof("Git cleaning untracked files from index. path: %s", path)
	script := fmt.Sprintf(`git -C %s clean -xdf`, path)
	return s.Execute(script)
}
