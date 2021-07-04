package git

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
)

func New(s shell.Shell) Git {
	return &gitImpl{
		shell: s,
	}
}

func (g *gitImpl) Clone(url string, branch string, clonePath string) error {
	err := g.Init(clonePath)
	if err != nil {
		return err
	}

	err = g.AddOrigin(clonePath, url)
	if err != nil {
		return err
	}

	err = g.FetchShallow(clonePath, branch)
	if err != nil {
		return err
	}

	err = g.Clean(clonePath)
	if err != nil {
		return err
	}
	return nil
}

func (g *gitImpl) Init(path string) error {
	logger.Infof("Git init a new index. path: %s", path)
	script := fmt.Sprintf("git init %s", path)
	return g.shell.Execute(script)
}

func (g *gitImpl) AddOrigin(path string, url string) error {
	logger.Infof("Git add remote origin. url: %s", url)
	script := fmt.Sprintf("git -C %s remote add origin %s", path, url)
	return g.shell.Execute(script)
}

func (g *gitImpl) FetchShallow(path string, branch string) error {
	logger.Infof("Git fetching branch with shallow refs. branch: %s, since: 4 weeks ago", branch)
	script := fmt.Sprintf(`git -C %s fetch --shallow-since="4 weeks ago" --force origin refs/heads/%s:refs/remotes/origin/%s`, path, branch, branch)
	return g.shell.Execute(script)
}

func (g *gitImpl) Reset(path string, revision string) error {
	logger.Infof("Git reset to a specific revision. commit: %s", revision)
	script := fmt.Sprintf(`git -C %s reset --hard "%s"`, path, revision)
	return g.shell.Execute(script)
}

func (g *gitImpl) Checkout(path string, branch string) error {
	logger.Infof("Git checkout. branch: %s", branch)
	script := fmt.Sprintf(`git -C %s checkout %s`, path, branch)
	return g.shell.Execute(script)
}

func (g *gitImpl) Clean(path string) error {
	logger.Infof("Git cleaning untracked files from index. path: %s", path)
	script := fmt.Sprintf(`git -C %s clean -xdf`, path)
	return g.shell.Execute(script)
}

func (g *gitImpl) GetHeadCommitHash(branch string) (string, error) {
	logger.Infof("Git reading HEAD latest commit hash. branch: %s", branch)
	script := fmt.Sprintf(`git ls-remote origin -h refs/heads/%s`, branch)
	return g.shell.ExecuteWithOutput(script)
}
