package git

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
	"strings"
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
	logger.Debugf("Git init a new index. path: %s", path)
	script := fmt.Sprintf("git init %s", path)
	return g.shell.ExecuteSilently(script)
}

func (g *gitImpl) AddOrigin(path string, url string) error {
	logger.Debugf("Git add remote origin. url: %s", url)
	script := fmt.Sprintf("git -C %s remote add origin %s", path, url)
	return g.shell.ExecuteSilently(script)
}

func (g *gitImpl) FetchShallow(path string, branch string) error {
	logger.Debugf("Git fetching branch with shallow refs. branch: %s, since: 4 weeks ago", branch)
	script := fmt.Sprintf(`git -C %s fetch --shallow-since="4 weeks ago" --force origin refs/heads/%s:refs/remotes/origin/%s`, path, branch, branch)
	return g.shell.ExecuteSilently(script)
}

func (g *gitImpl) Reset(path string, revision string) error {
	logger.Debugf("Git reset to a specific revision. commit: %s", revision)
	script := fmt.Sprintf(`git -C %s reset --hard "%s"`, path, revision)
	return g.shell.ExecuteSilently(script)
}

func (g *gitImpl) Checkout(path string, branch string) error {
	logger.Debugf("Git checkout. branch: %s", branch)
	script := fmt.Sprintf(`git -C %s checkout %s`, path, branch)
	return g.shell.ExecuteSilently(script)
}

func (g *gitImpl) Clean(path string) error {
	logger.Debugf("Git cleaning untracked files from index. path: %s", path)
	script := fmt.Sprintf(`git -C %s clean -xdf`, path)
	return g.shell.ExecuteSilently(script)
}

func (g *gitImpl) GetHeadCommitHash(path string, branch string) (string, error) {
	logger.Debugf("Git reading HEAD latest commit hash. branch: %s", branch)
	script := fmt.Sprintf(`git -C %s rev-parse origin/%s`, path, branch)
	output, err := g.shell.ExecuteWithOutput(script)
	if err == nil {
		output = strings.TrimSuffix(output, "\n")
	}
	return output, err
}
