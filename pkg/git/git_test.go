package git

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/logger"

	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_GitShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "init successfully",
			Func: InitSuccessfully,
		},
		{
			Name: "add origin successfully",
			Func: AddOriginSuccessfully,
		},
		{
			Name: "fetch shallow successfully",
			Func: FetchShallowSuccessfully,
		},
		{
			Name: "reset successfully",
			Func: ResetSuccessfully,
		},
		{
			Name: "checkout successfully",
			Func: CheckoutSuccessfully,
		},
		{
			Name: "clean successfully",
			Func: CleanSuccessfully,
		},
		{
			Name: "get HEAD commit hash successfully",
			Func: GetHeadCommitHashSuccessfully,
		},
		{
			Name: "get local origin commit hash successfully",
			Func: GetLocalOriginCommitHashSuccessfully,
		},
		{
			Name: "log revisions diff pretty successfully",
			Func: LogRevisionsDiffPrettySuccessfully,
		},
		{
			Name: "fail clone due to init failure",
			Func: GitCloneFailsOnInit,
		},
		{
			Name: "fail clone due to add-origin failure",
			Func: GitCloneFailsOnAddOrigin,
		},
		{
			Name: "fail clone due to fetch-shallow failure",
			Func: GitCloneFailsOnFetchShallow,
		},
		{
			Name: "fail clone due to clean failure",
			Func: GitCloneFailsOnClean,
		},
		{
			Name: "clone successfully",
			Func: CloneSuccessfully,
		},
	}
	harness.RunTests(t, tests)
}

var InitSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			repoName := "my-repo"
			fakeShell := shell.CreateFakeShell()
			fakeShell.ExecuteSilentlyMock = func(script string) error {
				assert.Equal(t, fmt.Sprintf("git init %s", repoName), script)
				return nil
			}
			git := New(fakeShell)
			_ = git.Init(repoName)
		})
	})
}

var AddOriginSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			clonePath := "/some/path"
			url := "git@some-repo"
			fakeShell := shell.CreateFakeShell()
			fakeShell.ExecuteSilentlyMock = func(script string) error {
				expected := fmt.Sprintf("git -C %s remote add origin %s", clonePath, url)
				assert.Equal(t, expected, script)
				return nil
			}
			git := New(fakeShell)
			_ = git.AddOrigin(clonePath, url)
		})
	})
}

var FetchShallowSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			clonePath := "/some/path"
			branch := "my-branch"
			fakeShell := shell.CreateFakeShell()
			fakeShell.ExecuteSilentlyMock = func(script string) error {
				expected := fmt.Sprintf(`git -C %s fetch --shallow-since="4 weeks ago" --force origin refs/heads/%s:refs/remotes/origin/%s`, clonePath, branch, branch)
				assert.Equal(t, expected, script)
				return nil
			}
			git := New(fakeShell)
			_ = git.FetchShallow(clonePath, branch)
		})
	})
}

var ResetSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			clonePath := "/some/path"
			revision := "l33tf4k3c0mm1757r1n6"
			fakeShell := shell.CreateFakeShell()
			fakeShell.ExecuteSilentlyMock = func(script string) error {
				expected := fmt.Sprintf(`git -C %s reset --hard "%s"`, clonePath, revision)
				assert.Equal(t, expected, script)
				return nil
			}
			git := New(fakeShell)
			_ = git.Reset(clonePath, revision)
		})
	})
}

var CheckoutSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			clonePath := "/some/path"
			branch := "my-branch"
			fakeShell := shell.CreateFakeShell()
			fakeShell.ExecuteSilentlyMock = func(script string) error {
				expected := fmt.Sprintf(`git -C %s checkout %s`, clonePath, branch)
				assert.Equal(t, expected, script)
				return nil
			}
			git := New(fakeShell)
			_ = git.Checkout(clonePath, branch)
		})
	})
}

var CleanSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			clonePath := "/some/path"
			fakeShell := shell.CreateFakeShell()
			fakeShell.ExecuteSilentlyMock = func(script string) error {
				expected := fmt.Sprintf(`git -C %s clean -xdf`, clonePath)
				assert.Equal(t, expected, script)
				return nil
			}
			git := New(fakeShell)
			_ = git.Clean(clonePath)
		})
	})
}

var GetHeadCommitHashSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			clonePath := "/some/path"
			branch := "my-branch"
			url := "git@some-repo"
			revision := "l33tf4k3c0mm1757r1n6"
			fakeShell := shell.CreateFakeShell()
			fakeShell.ExecuteReturnOutputMock = func(script string) (string, error) {
				expected := fmt.Sprintf(`git -C %s ls-remote %s %s`, clonePath, url, branch)
				assert.Equal(t, expected, script)
				return revision, nil
			}
			git := New(fakeShell)
			rev, _ := git.GetRemoteHeadCommitHash(clonePath, url, branch)
			assert.Equal(t, rev, revision)
		})
	})
}

var GetLocalOriginCommitHashSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			clonePath := "/some/path"
			branch := "my-branch"
			revision := "l33tf4k3c0mm1757r1n6"
			fakeShell := shell.CreateFakeShell()
			fakeShell.ExecuteReturnOutputMock = func(script string) (string, error) {
				expected := fmt.Sprintf(`git -C %s rev-parse origin/%s`, clonePath, branch)
				assert.Equal(t, expected, script)
				return revision, nil
			}
			git := New(fakeShell)
			rev, _ := git.GetLocalOriginCommitHash(clonePath, branch)
			assert.Equal(t, rev, revision)
		})
	})
}

var LogRevisionsDiffPrettySuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			clonePath := "/some/path"
			prevRevision := "l33tf4k3c0mm1757r1n6"
			headRevision := "head-revision"
			fakeShell := shell.CreateFakeShell()
			fakeShell.ExecuteMock = func(script string) error {
				expected := fmt.Sprintf(`git -C %s --no-pager log \
--oneline \
--graph \
--pretty='%%C(Yellow)%%h%%Creset %%<(5) %%C(auto)%%s%%Creset %%Cgreen(%%ad) %%C(bold blue)<%%an>%%Creset' \
--date=short %s..%s`, clonePath, prevRevision, headRevision)
				assert.Equal(t, expected, script)
				return nil
			}
			git := New(fakeShell)
			_ = git.LogRevisionsDiffPretty(clonePath, prevRevision, headRevision)
		})
	})
}

var GitCloneFailsOnInit = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			clonePath := "/some/path"
			url := "git@some-repo"
			branch := "my-branch"
			fakeShell := shell.CreateFakeShell()
			fakeShell.ExecuteSilentlyMock = func(script string) error {
				if strings.Contains(script, "git init") {
					return fmt.Errorf("failed to init")
				}
				return nil
			}
			git := New(fakeShell)
			err := git.Clone(url, branch, clonePath)
			assert.NotNil(t, err, "expect to fail")
			assert.Equal(t, "failed to init", err.Error())
		})
	})
}

var GitCloneFailsOnAddOrigin = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			clonePath := "/some/path"
			url := "git@some-repo"
			branch := "my-branch"
			fakeShell := shell.CreateFakeShell()
			fakeShell.ExecuteSilentlyMock = func(script string) error {
				if strings.Contains(script, "git init") {
					return nil
				} else if strings.Contains(script, "add origin") {
					return fmt.Errorf("failed to add origin")
				}
				return nil
			}
			git := New(fakeShell)
			err := git.Clone(url, branch, clonePath)
			assert.NotNil(t, err, "expect to fail")
			assert.Equal(t, "failed to add origin", err.Error())
		})
	})
}

var GitCloneFailsOnFetchShallow = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			clonePath := "/some/path"
			url := "git@some-repo"
			branch := "my-branch"
			fakeShell := shell.CreateFakeShell()
			fakeShell.ExecuteSilentlyMock = func(script string) error {
				if strings.Contains(script, "git init") {
					return nil
				} else if strings.Contains(script, "add origin") {
					return nil
				} else if strings.Contains(script, "fetch --shallow-since") {
					return fmt.Errorf("failed to fetch")
				}
				return nil
			}
			git := New(fakeShell)
			err := git.Clone(url, branch, clonePath)
			assert.NotNil(t, err, "expect to fail")
			assert.Equal(t, "failed to fetch", err.Error())
		})
	})
}

var GitCloneFailsOnClean = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			clonePath := "/some/path"
			url := "git@some-repo"
			branch := "my-branch"
			fakeShell := shell.CreateFakeShell()
			fakeShell.ExecuteSilentlyMock = func(script string) error {
				if strings.Contains(script, "git init") {
					return nil
				} else if strings.Contains(script, "add origin") {
					return nil
				} else if strings.Contains(script, "fetch --shallow-since") {
					return nil
				} else if strings.Contains(script, "clean -xdf") {
					return fmt.Errorf("failed to clean")
				}
				return nil
			}
			git := New(fakeShell)
			err := git.Clone(url, branch, clonePath)
			assert.NotNil(t, err, "expect to fail")
			assert.Equal(t, "failed to clean", err.Error())
		})
	})
}

var CloneSuccessfully = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			clonePath := "/some/path"
			url := "git@some-repo"
			branch := "my-branch"
			fakeShell := shell.CreateFakeShell()
			fakeShell.ExecuteSilentlyMock = func(script string) error {
				return nil
			}
			git := New(fakeShell)
			err := git.Clone(url, branch, clonePath)
			assert.Nil(t, err, "expect to succeed")
		})
	})
}
