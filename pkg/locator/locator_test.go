package locator

import (
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_LocatorShould(t *testing.T) {
	tests := []harness.TestsHarness{
		//{
		//	Name: "fail on invalid anchorfiles local path",
		//	Func: FailOnInvalidAnchorfilesLocalPath,
		//},
		//{
		//	Name: "scan anchorfiles test repo and find expected applications",
		//	Func: ScanAndFindExpectedApplications,
		//},
		{
			Name: "test paths1",
			Func: TestPath1,
		},
	}
	harness.RunTests(t, tests)
}

var FailOnInvalidAnchorfilesLocalPath = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := config.GetDefaultTestConfigText()
			with.Config(ctx, yamlConfigText, func(config *config.AnchorConfig) {
				l := New()
				err := l.Scan("/invalid/anchorfiles/path")
				assert.NotNil(t, err, "expected to fail on invalid anchorfiles local path")
				assert.Contains(t, err.Error(), "invalid anchorfile local path")
			})
		})
	})
}

var ScanAndFindExpectedApplications = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := config.GetDefaultTestConfigText()
			with.Config(ctx, yamlConfigText, func(config *config.AnchorConfig) {
				with.HarnessAnchorfilesTestRepo(ctx)
				l := New()
				assert.Nil(t, l.Scan(ctx.AnchorFilesPath()), "expect locator to scan successfully")
				applications := l.Applications()
				assert.Equal(t, 2, len(applications), "expected 2 applications but found %v", len(applications))
				firstAppName := "first-app"
				assert.NotNil(t, l.Application(firstAppName), "expected application to exist. Name: %s", firstAppName)
				assert.Equal(t, firstAppName, l.Application(firstAppName).Name, "expected application %s but found %s",
					firstAppName, l.Application(firstAppName).Name)
				secondAppName := "second-app"
				assert.NotNil(t, l.Application(secondAppName), "expected application to exist. Name: %s", secondAppName)
				assert.Equal(t, secondAppName, l.Application(secondAppName).Name, "expected application %s but found %s",
					secondAppName, l.Application(secondAppName).Name)
			})
		})
	})
}

var TestPath1 = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.LoggingVerbose(ctx, t, func(logger logger.Logger) {
			with.Config(ctx, config.GetDefaultTestConfigText(), func(config *config.AnchorConfig) {
				with.HarnessAnchorfilesTestRepo(ctx)
				assert.NotEmpty(t, ctx.AnchorFilesPath())
				logger.Debugf("test context (ctx.AnchorFilesPath()): %s", ctx.AnchorFilesPath())
				logger.Debugf("GOPATH: %s", os.Getenv("GOPATH"))
				getwd, _ := os.Getwd()
				logger.Debugf("CWD: %s", getwd)
				assert.DirExists(t, ctx.AnchorFilesPath())
			})
		})
	})
}
