package git

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
)

func New(s shell.Shell) Git {
	return &gitImpl{
		Shell: s,
	}
}

func (g *gitImpl) GitClone(url string, branch string, clonePath string) error {
	err := g.GitInit(clonePath)
	if err != nil {
		return err
	}

	err = g.GitAddOrigin(clonePath, url)
	if err != nil {
		return err
	}

	err = g.GitFetchShallow(clonePath, url, branch)
	if err != nil {
		return err
	}

	err = g.GitClean(clonePath)
	if err != nil {
		return err
	}
	return nil
}

func (g *gitImpl) GitInit(path string) error {
	logger.Infof("Git init a new index. path: %s", path)
	script := fmt.Sprintf("git init %s", path)
	return g.Shell.Execute(script)
}

func (g *gitImpl) GitAddOrigin(path string, url string) error {
	logger.Infof("Git add remote origin. url: %s", url)
	script := fmt.Sprintf("git -C %s remote add origin %s", path, url)
	return g.Shell.Execute(script)
}

func (g *gitImpl) GitFetchShallow(path string, url string, branch string) error {
	logger.Infof("Git fetching branch with shallow refs. branch: %s, since: 4 weeks ago", branch)
	script := fmt.Sprintf(`git -C %s fetch --shallow-since="4 weeks ago" --force origin refs/heads/%s:refs/remotes/origin/%s`, path, url, branch)
	return g.Shell.Execute(script)
}

func (g *gitImpl) GitReset(path string, revision string) error {
	logger.Infof("Git reset to a specific revision. commit: %s", revision)
	script := fmt.Sprintf(`git -C %s reset --hard "%s"`, path, revision)
	return g.Shell.Execute(script)
}

func (g *gitImpl) GitClean(path string) error {
	logger.Infof("Git cleaning untracked files from index. path: %s", path)
	script := fmt.Sprintf(`git -C %s clean -xdf`, path)
	return g.Shell.Execute(script)
}

func (g *gitImpl) GetHeadCommitHash(branch string) error {
	logger.Infof("Git reading HEAD latest commit hash. branch: %s", branch)
	script := fmt.Sprintf(`git ls-remote origin -h refs/heads/%s`, branch)
	return g.Shell.Execute(script)
}
