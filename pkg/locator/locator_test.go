package locator

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/config"
	"github.com/ZachiNachshon/anchor/internal/errors"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/pkg/extractor"
	"github.com/ZachiNachshon/anchor/pkg/models"
	"github.com/ZachiNachshon/anchor/pkg/parser"
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
			Name: "fail on already initialized",
			Func: FailOnAlreadyInitialized,
		},
		{
			Name: "not return anchor folder if missing from scan",
			Func: NotReturnCommandFolderIfMissingFromScan,
		},
		{
			Name: "scan anchorfiles test repo and find expected anchor folders",
			Func: ScanAndFindExpectedCommandFolders,
		},
		{
			Name: "not locate any anchor folders due bad YAML",
			Func: NotLocateAnyCommandFoldersDueToBadYaml,
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
				locatorErr := l.Scan("/invalid/anchorfiles/path",
					extractor.CreateFakeExtractor(),
					parser.CreateFakeParser())
				assert.NotNil(t, locatorErr, "expected to fail on invalid anchorfiles local path")
				assert.Contains(t, locatorErr.GoError().Error(), "invalid anchorfile local path")
			})
		})
	})
}

var FailOnAlreadyInitialized = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := config.GetDefaultTestConfigText()
			with.Config(ctx, yamlConfigText, func(config *config.AnchorConfig) {
				with.HarnessAnchorfilesTestRepo(ctx)
				l := New()
				fakeExtractor := extractor.CreateFakeExtractor()
				fakeExtractor.ExtractCommandFolderInfoMock = func(dirPath string, p parser.Parser) (*models.CommandFolderInfo, error) {
					return nil, nil
				}
				locatorErr := l.Scan(ctx.AnchorFilesPath(),
					fakeExtractor,
					parser.CreateFakeParser())
				locatorErr = l.Scan(ctx.AnchorFilesPath(),
					fakeExtractor,
					parser.CreateFakeParser())
				assert.NotNil(t, locatorErr.Code(), errors.AlreadyInitialized)
			})
		})
	})
}

var NotReturnCommandFolderIfMissingFromScan = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := config.GetDefaultTestConfigText()
			with.Config(ctx, yamlConfigText, func(config *config.AnchorConfig) {
				l := New()
				result := l.CommandFolderByName("not-exists")
				assert.Nil(t, result, "should not identify application after scan took place")
			})
		})
	})
}

var ScanAndFindExpectedCommandFolders = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := config.GetDefaultTestConfigText()
			with.Config(ctx, yamlConfigText, func(config *config.AnchorConfig) {
				with.HarnessAnchorfilesTestRepo(ctx)
				l := New()
				err := l.Scan(ctx.AnchorFilesPath(),
					extractor.New(),
					parser.New())
				assert.Nil(t, err, "expect locator to scan successfully")
				commandFolders := l.CommandFolders()
				assert.Equal(t, 3, len(commandFolders), "expected 3 anchor folders but found %v", len(commandFolders))

				// Anchor Folders
				commandFoldersAsMap := l.CommandFoldersAsMap()
				assert.Equal(t, 3, len(commandFoldersAsMap), "expected map of 3 anchor folders")
				assert.NotNil(t, l.CommandFolderByName("app"))
				assert.Nil(t, l.CommandFolderByName("app-ignored"))
				assert.NotNil(t, l.CommandFolderByName("controller"))
				assert.NotNil(t, l.CommandFolderByName("k8s"))

				// Anchor Folder: App
				appCommandFolder := l.CommandFolderByName("app")
				assert.Equal(t, "app", appCommandFolder.Name)
				assert.NotNil(t, appCommandFolder.Command)
				assert.Equal(t, "app", appCommandFolder.Command.Use)
				assert.Equal(t, "Application commands", appCommandFolder.Command.Short)

				// Anchor Folder Items: App
				appItems := l.CommandFolderItems("app")
				assert.NotNil(t, appItems, "expected to have valid items for anchor folder: app")
				assert.NotNil(t, 2, len(appItems), "expected 2 items for anchor folder: app")
				firstAppName := "first-app"
				assert.Equal(t, appItems[0].Name, firstAppName)
				secondAppName := "second-app"
				assert.Equal(t, appItems[1].Name, secondAppName)
			})
		})
	})
}

var NotLocateAnyCommandFoldersDueToBadYaml = func(t *testing.T) {
	with.Context(func(ctx common.Context) {
		with.Logging(ctx, t, func(logger logger.Logger) {
			yamlConfigText := config.GetDefaultTestConfigText()
			with.Config(ctx, yamlConfigText, func(config *config.AnchorConfig) {
				with.HarnessAnchorfilesTestRepo(ctx)
				fakeExtractor := extractor.CreateFakeExtractor()
				fakeExtractor.ExtractCommandFolderInfoMock = func(dirPath string, p parser.Parser) (*models.CommandFolderInfo, error) {
					return nil, fmt.Errorf("failed to extract anchor folder info")
				}
				l := New()
				err := l.Scan(ctx.AnchorFilesPath(), fakeExtractor, parser.CreateFakeParser())
				assert.Nil(t, err, "expect locator to scan successfully")
				commandFolders := l.CommandFolders()
				assert.Equal(t, 0, len(commandFolders), "expected no anchor folders to be found")
			})
		})
	})
}
