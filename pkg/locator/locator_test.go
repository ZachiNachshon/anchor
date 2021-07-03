package locator

import (
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/config"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/ZachiNachshon/anchor/test/with"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_LocatorShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "fail on invalid anchorfiles local path",
			Func: FailOnInvalidAnchorfilesLocalPath,
		},
		{
			Name: "scan anchorfiles test repo and find expected applications",
			Func: ScanAndFindExpectedApplications,
		},
	}
	harness.RunTests(t, tests)
}

var FailOnInvalidAnchorfilesLocalPath = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := config.GetDefaultTestConfigText()
			with.Config(ctx, yamlConfigText, func(config config.AnchorConfig) {
				// Given I create a new locator with invalid anchorfiles path
				l := New()
				// When I scan the anchorfiles repo
				err := l.Scan("/invalid/anchorfiles/path")
				// Then I expect to fail on invalid anchorfiles path
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
			with.Config(ctx, yamlConfigText, func(config config.AnchorConfig) {
				// Given I prepare an anchorfiles test repo
				harness.HarnessAnchorfilesTestRepo(ctx)
				// And I create a new locator
				l := New()
				// When I scan the anchorfiles repo
				assert.Nil(t, l.Scan(ctx.AnchorFilesPath()), "expect locator to scan successfully")
				// Then I expect to find 2 applications
				applications := l.Applications()
				assert.Equal(t, 2, len(applications), "expected 2 applications but found %v", len(applications))
				// And the 1st application must be named "first-app"
				firstAppName := "first-app"
				assert.NotNil(t, l.Application(firstAppName), "expected application to exist. Name: %s", firstAppName)
				assert.Equal(t, firstAppName, l.Application(firstAppName).Name, "expected application %s but found %s",
					firstAppName, l.Application(firstAppName).Name)
				// And the 2nd application must be named "second-app"
				secondAppName := "second-app"
				assert.NotNil(t, l.Application(secondAppName), "expected application to exist. Name: %s", secondAppName)
				assert.Equal(t, secondAppName, l.Application(secondAppName).Name, "expected application %s but found %s",
					secondAppName, l.Application(secondAppName).Name)
			})
		})
	})
}
