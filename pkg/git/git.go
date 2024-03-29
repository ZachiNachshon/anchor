package git

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/logger"

	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
	"strings"
)

type Git interface {
	Clone(url string, branch string, clonePath string) error
	Init(path string) error
	AddOrigin(path string, url string) error
	FetchShallow(path string, url string, branch string) error
	Reset(path string, revision string) error
	Checkout(path string, branch string) error
	Clean(path string) error
	GetRemoteHeadCommitHash(path string, repoUrl string, branch string) (string, error)
	GetLocalOriginCommitHash(path string, branch string) (string, error)
	LogRevisionsDiffPretty(path string, prevRevision string, newRevision string) error
}

type gitImpl struct {
	Git
	shell shell.Shell
}

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

	err = g.FetchShallow(clonePath, url, branch)
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

func (g *gitImpl) FetchShallow(path string, url string, branch string) error {
	limitationOption := `--shallow-since="4 weeks ago"`
	if strings.Contains(url, "https") {
		limitationOption = "--depth 1"
	}
	logger.Debugf("Git shallow fetching a branch. branch: %s, limitation: %s", branch, limitationOption)
	script := fmt.Sprintf(`git -C %s fetch %s --force origin refs/heads/%s:refs/remotes/origin/%s`,
		path, limitationOption, branch, branch)
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

func (g *gitImpl) GetRemoteHeadCommitHash(path string, repoUrl string, branch string) (string, error) {
	logger.Debugf("Git reading HEAD latest commit hash. branch: %s", branch)
	script := fmt.Sprintf(`git -C %s ls-remote %s %s`, path, repoUrl, branch)
	output, err := g.shell.ExecuteReturnOutput(script)
	if err == nil {
		fields := strings.Fields(output)
		if len(fields) > 0 {
			return strings.TrimSuffix(fields[0], "\n"), nil
		}
	}
	return output, err
}

func (g *gitImpl) GetLocalOriginCommitHash(path string, branch string) (string, error) {
	logger.Debugf("Git reading HEAD latest commit hash. branch: %s", branch)
	script := fmt.Sprintf(`git -C %s rev-parse origin/%s`, path, branch)
	output, err := g.shell.ExecuteReturnOutput(script)
	if err == nil {
		output = strings.TrimSuffix(output, "\n")
	}
	return output, err
}

func (g *gitImpl) LogRevisionsDiffPretty(path string, prevRevision string, newRevision string) error {
	logger.Debugf("Git logging revisions changed (pretty). prevRev: %s, newRev: %s", prevRevision, newRevision)
	script := fmt.Sprintf(`git -C %s --no-pager log \
--oneline \
--graph \
--pretty='%%C(Yellow)%%h%%Creset %%<(5) %%C(auto)%%s%%Creset %%Cgreen(%%ad) %%C(bold blue)<%%an>%%Creset' \
--date=short %s..%s`, path, prevRevision, newRevision)
	fmt.Print("\nCommits:\n\n")
	if err := g.shell.Execute(script); err != nil {
		return err
	}
	return nil
}
