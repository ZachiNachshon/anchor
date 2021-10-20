package config

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/internal/registry"
	"github.com/ZachiNachshon/anchor/pkg/utils/ioutils"
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func Test_ConfigShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "setup file loader: fail to get config filepath",
			Func: SetupFileLoaderFailToGetConfigFilepath,
		},
		{
			Name: "setup file loader: fail to set config path",
			Func: SetupFileLoaderFailToSetConfigPath,
		},
		{
			Name: "setup file loader: fail to load config from file",
			Func: SetupFileLoaderFailToLoadConfigFromFile,
		},
		{
			Name: "setup file loader: fail to set env vars",
			Func: SetupFileLoaderFailToSetEnvVars,
		},
		{
			Name: "setup file loader: successful setup",
			Func: SetupFileLoaderSuccessfulSetup,
		},
		{
			Name: "override config: update all config entries",
			Func: OverrideConfigUpdateAllConfigEntries,
		},
		{
			Name: "override config: update single config entry",
			Func: OverrideConfigUpdateSingleConfigEntry,
		},
		{
			Name: "switch active config: switch successfully",
			Func: SwitchActiveConfigSwitchSuccessfully,
		},
		{
			Name: "switch active config: fail to get config context",
			Func: SwitchActiveConfigFailToGetConfigContext,
		},
		{
			Name: "create config object: create successfully",
			Func: CreateConfigObjectCreateSuccessfully,
		},
		{
			Name: "create config object: fail to merge",
			Func: CreateConfigObjectFailToMerge,
		},
		{
			Name: "create config object: fail to validate",
			Func: CreateConfigObjectFailToValidate,
		},
		//{
		//	Name: "config file path: fail to get user home dir",
		//	Func: ConfigFilePathFailToGetUserHomeDir,
		//},
		{
			Name: "config file path: get user home dir",
			Func: ConfigFilePathGetUserHomeDir,
		},
		//{
		//	Name: "default repo clone path: fail to get user home dir",
		//	Func: DefaultRepoClonePathFailToGetUserHomeDir,
		//},
		{
			Name: "default repo clone path: get clone path",
			Func: DefaultRepoClonePathGetClonePath,
		},
		{
			Name: "config text: unmarshal from YAML",
			Func: ConfigTextUnmarshalFromYaml,
		},
		{
			Name: "config text: fail to unmarshal",
			Func: ConfigTextFailToUnmarshal,
		},
		{
			Name: "config object to text: missing object",
			Func: ConfigObjectToTextMissingObject,
		},
		//{
		//	Name: "config object to text: fail to unmarshal to YAML",
		//	Func: ConfigObjectToTextFailToUnmarshalToYaml,
		//},
		{
			Name: "config object to text: unmarshal to YAML",
			Func: ConfigObjectToTextUnmarshalToYaml,
		},
		{
			Name: "context: empty config",
			Func: ContextEmtpyConfig,
		},
		{
			Name: "context: retrieve successfully",
			Func: ContextRetrieveSuccessfully,
		},
		{
			Name: "context: set config in context",
			Func: ContextSetConfigInContext,
		},
		{
			Name: "config context: empty contexts",
			Func: ConfigContextEmptyContexts,
		},
		{
			Name: "config context: context not found",
			Func: ConfigContextContextNotFound,
		},
		{
			Name: "config context: found context",
			Func: ConfigContextFoundContext,
		},
		{
			Name: "validations: verify mandatory config entries exist",
			Func: ValidationsVerifyMandatoryConfigEntriesExist,
		},
		{
			Name: "path: get config file path",
			Func: PathGetConfigFilePath,
		},
		{
			Name: "defaults: skip when config do not override on explicit values",
			Func: DefaultsDoNotOverrideOnExplicitValues,
		},
		{
			Name: "defaults: skip on missing remote config",
			Func: DefaultsSkipOnMissingRemoteConfig,
		},
		{
			Name: "defaults: do not override on explicit values",
			Func: DefaultsDoNotOverrideOnExplicitValues,
		},
		{
			Name: "defaults: fail on creation of default clone path for config context",
			Func: DefaultsFailOnCreationOfDefaultClonePathForConfigContext,
		},
		{
			Name: "defaults: set empty entries with default values on all contexts",
			Func: DefaultsSetEmptyEntriesWithDefaultValuesOnAllContexts,
		},
		{
			Name: "append new empty config context successfully",
			Func: AppendNewEmptyConfigContextSuccessfully,
		},
	}
	harness.RunTests(t, tests)
}

var SetupFileLoaderFailToGetConfigFilepath = func(t *testing.T) {
	withContext(func(ctx common.Context) {
		withLogging(ctx, t, false, func(logger logger.Logger) {
			cfgManager := NewManager()
			getCfgFilePathCallCount := 0
			cfgManager.getConfigFilePathFunc = func() (string, error) {
				getCfgFilePathCallCount++
				return "", fmt.Errorf("failed to get config path")
			}
			err := cfgManager.SetupConfigFileLoader()
			assert.NotNil(t, err, "expected a failure")
			assert.Equal(t, "failed to get config path", err.Error())
			assert.Equal(t, 1, getCfgFilePathCallCount)
		})
	})
}

var SetupFileLoaderFailToSetConfigPath = func(t *testing.T) {
	withContext(func(ctx common.Context) {
		withLogging(ctx, t, false, func(logger logger.Logger) {
			fakeAdapter := CreateFakeViperConfigAdapter()
			cfgManager := NewManager()
			cfgManager.adapter = fakeAdapter
			getCfgFilePathCallCount := 0
			cfgManager.getConfigFilePathFunc = func() (string, error) {
				getCfgFilePathCallCount++
				return "", nil
			}
			setCfgPathCallCount := 0
			fakeAdapter.SetConfigPathMock = func(path string) error {
				setCfgPathCallCount++
				return fmt.Errorf("failed to set config path")
			}
			err := cfgManager.SetupConfigFileLoader()
			assert.NotNil(t, err, "expected a failure")
			assert.Equal(t, "failed to set config path", err.Error())
			assert.Equal(t, 1, getCfgFilePathCallCount)
			assert.Equal(t, 1, setCfgPathCallCount)
		})
	})
}

var SetupFileLoaderFailToLoadConfigFromFile = func(t *testing.T) {
	withContext(func(ctx common.Context) {
		withLogging(ctx, t, false, func(logger logger.Logger) {
			fakeAdapter := CreateFakeViperConfigAdapter()
			cfgManager := NewManager()
			cfgManager.adapter = fakeAdapter
			fakeAdapter.SetConfigPathMock = func(path string) error {
				return nil
			}
			setCfgPathCallCount := 0
			fakeAdapter.SetConfigPathMock = func(path string) error {
				setCfgPathCallCount++
				return nil
			}
			loadCfgFromFileCallCount := 0
			fakeAdapter.LoadConfigFromFileMock = func() error {
				loadCfgFromFileCallCount++
				return fmt.Errorf("failed to load config from file")
			}
			err := cfgManager.SetupConfigFileLoader()
			assert.NotNil(t, err, "expected a failure")
			assert.Equal(t, "failed to load config from file", err.Error())
			assert.Equal(t, 1, setCfgPathCallCount)
			assert.Equal(t, 1, loadCfgFromFileCallCount)
		})
	})
}

var SetupFileLoaderFailToSetEnvVars = func(t *testing.T) {
	withContext(func(ctx common.Context) {
		withLogging(ctx, t, false, func(logger logger.Logger) {
			fakeAdapter := CreateFakeViperConfigAdapter()
			cfgManager := NewManager()
			cfgManager.adapter = fakeAdapter
			getCfgFilePathCallCount := 0
			cfgManager.getConfigFilePathFunc = func() (string, error) {
				getCfgFilePathCallCount++
				return "", nil
			}
			setCfgPathCallCount := 0
			fakeAdapter.SetConfigPathMock = func(path string) error {
				setCfgPathCallCount++
				return nil
			}
			loadCfgFromFileCallCount := 0
			fakeAdapter.LoadConfigFromFileMock = func() error {
				loadCfgFromFileCallCount++
				return nil
			}
			setEnvVarsCallCount := 0
			fakeAdapter.SetEnvVarsMock = func() error {
				setEnvVarsCallCount++
				return fmt.Errorf("failed to set env vars")
			}
			err := cfgManager.SetupConfigFileLoader()
			assert.NotNil(t, err, "expected a failure")
			assert.Equal(t, "failed to set env vars", err.Error())
			assert.Equal(t, 1, getCfgFilePathCallCount)
			assert.Equal(t, 1, setCfgPathCallCount)
			assert.Equal(t, 1, loadCfgFromFileCallCount)
			assert.Equal(t, 1, setEnvVarsCallCount)
		})
	})
}

var SetupFileLoaderSuccessfulSetup = func(t *testing.T) {
	withContext(func(ctx common.Context) {
		withLogging(ctx, t, false, func(logger logger.Logger) {
			fakeAdapter := CreateFakeViperConfigAdapter()
			cfgManager := NewManager()
			cfgManager.adapter = fakeAdapter
			getCfgFilePathCallCount := 0
			cfgManager.getConfigFilePathFunc = func() (string, error) {
				getCfgFilePathCallCount++
				return "", nil
			}
			setCfgPathCallCount := 0
			fakeAdapter.SetConfigPathMock = func(path string) error {
				setCfgPathCallCount++
				return nil
			}
			loadCfgFromFileCallCount := 0
			fakeAdapter.LoadConfigFromFileMock = func() error {
				loadCfgFromFileCallCount++
				return nil
			}
			setEnvVarsCallCount := 0
			fakeAdapter.SetEnvVarsMock = func() error {
				setEnvVarsCallCount++
				return nil
			}
			err := cfgManager.SetupConfigFileLoader()
			assert.Nil(t, err, "expected to succeed")
			assert.Equal(t, 1, getCfgFilePathCallCount)
			assert.Equal(t, 1, setCfgPathCallCount)
			assert.Equal(t, 1, loadCfgFromFileCallCount)
			assert.Equal(t, 1, setEnvVarsCallCount)
		})
	})
}

var OverrideConfigUpdateAllConfigEntries = func(t *testing.T) {
	withContext(func(ctx common.Context) {
		withLogging(ctx, t, false, func(logger logger.Logger) {
			withConfig(ctx, GetDefaultTestConfigText(), func(cfg *AnchorConfig) {
				fakeAdapter := CreateFakeViperConfigAdapter()
				updateAllCallCount := 0
				fakeAdapter.UpdateAllMock = func(cfgToUpdate *AnchorConfig) error {
					updateAllCallCount++
					assert.Equal(t, cfgToUpdate, cfg)
					return nil
				}
				cfgManager := NewManager()
				cfgManager.adapter = fakeAdapter
				err := cfgManager.OverrideConfig(cfg)
				assert.Nil(t, err, "expected to succeed")
				assert.Equal(t, 1, updateAllCallCount)
			})
		})
	})
}

var OverrideConfigUpdateSingleConfigEntry = func(t *testing.T) {
	withContext(func(ctx common.Context) {
		withLogging(ctx, t, false, func(logger logger.Logger) {
			fakeAdapter := CreateFakeViperConfigAdapter()
			updateEntryCallCount := 0
			fakeAdapter.UpdateEntryMock = func(entryName string, value interface{}) error {
				updateEntryCallCount++
				assert.Equal(t, "test-entry", entryName)
				assert.Equal(t, "test-value", value)
				return nil
			}
			cfgManager := NewManager()
			cfgManager.adapter = fakeAdapter
			err := cfgManager.OverrideConfigEntry("test-entry", "test-value")
			assert.Nil(t, err, "expected to succeed")
			assert.Equal(t, 1, updateEntryCallCount)
		})
	})
}

var SwitchActiveConfigSwitchSuccessfully = func(t *testing.T) {
	withContext(func(ctx common.Context) {
		withLogging(ctx, t, false, func(logger logger.Logger) {
			withConfig(ctx, GetDefaultTestConfigText(), func(cfg *AnchorConfig) {
				cfgCtxName := "1st-anchorfiles"
				cfgManager := NewManager()
				err := cfgManager.SwitchActiveConfigContextByName(cfg, cfgCtxName)
				assert.Nil(t, err, "expected to succeed")
			})
		})
	})
}

var SwitchActiveConfigFailToGetConfigContext = func(t *testing.T) {
	withContext(func(ctx common.Context) {
		withLogging(ctx, t, false, func(logger logger.Logger) {
			withConfig(ctx, GetDefaultTestConfigText(), func(cfg *AnchorConfig) {
				cfgCtxName := "invalid-cfg-name"
				cfgManager := NewManager()
				err := cfgManager.SwitchActiveConfigContextByName(cfg, cfgCtxName)
				assert.NotNil(t, err, "expected to fail")
				assert.Contains(t, err.Error(), "could not identify config context")
			})
		})
	})
}

var CreateConfigObjectCreateSuccessfully = func(t *testing.T) {
	withContext(func(ctx common.Context) {
		withLogging(ctx, t, false, func(logger logger.Logger) {
			fakeAdapter := CreateFakeViperConfigAdapter()
			cfgManager := NewManager()
			cfgManager.adapter = fakeAdapter

			fakeAdapter.GetConfigByKeyMock = func(key string) string {
				if key == "author" {
					return "test-author"
				} else if key == "license" {
					return "test-license"
				}
				return ""
			}
			mergeConfigCallCount := 0
			fakeAdapter.AppendConfigMock = func(anchorConfig interface{}) error {
				cfg := anchorConfig.(*AnchorConfig)
				cfg.Config = &Config{}
				cfg.Config.Contexts = []*Context{}
				mergeConfigCallCount++
				return nil
			}
			validateCfgCallCount := 0
			cfgManager.validateConfigurationsFunc = func(anchorConfig *AnchorConfig) error {
				validateCfgCallCount++
				return nil
			}
			cfgManager.getDefaultRepoClonePathFunc = func(contextName string) (string, error) {
				return "", nil
			}
			result, err := cfgManager.CreateConfigObject()
			assert.Nil(t, err, "expected to succeed")
			assert.NotNil(t, result, "expected non-empty config object")
			assert.NotNil(t, "test-author", result.Author)
			assert.NotNil(t, "test-license", result.License)
			assert.Equal(t, 1, mergeConfigCallCount)
			assert.Equal(t, 1, validateCfgCallCount)
		})
	})
}

var CreateConfigObjectFailToMerge = func(t *testing.T) {
	withContext(func(ctx common.Context) {
		withLogging(ctx, t, false, func(logger logger.Logger) {
			fakeAdapter := CreateFakeViperConfigAdapter()
			cfgManager := NewManager()
			cfgManager.adapter = fakeAdapter

			fakeAdapter.GetConfigByKeyMock = func(key string) string {
				return ""
			}
			mergeConfigCallCount := 0
			fakeAdapter.AppendConfigMock = func(output interface{}) error {
				mergeConfigCallCount++
				return fmt.Errorf("failed to merge")
			}

			result, err := cfgManager.CreateConfigObject()
			assert.NotNil(t, err, "expected to fail")
			assert.NotNil(t, "failed to merge", err.Error())
			assert.Nil(t, result, "expected an empty config object")
			assert.Equal(t, 1, mergeConfigCallCount)
		})
	})
}

var CreateConfigObjectFailToValidate = func(t *testing.T) {
	withContext(func(ctx common.Context) {
		withLogging(ctx, t, false, func(logger logger.Logger) {
			fakeAdapter := CreateFakeViperConfigAdapter()
			cfgManager := NewManager()
			cfgManager.adapter = fakeAdapter

			fakeAdapter.GetConfigByKeyMock = func(key string) string {
				return ""
			}
			mergeConfigCallCount := 0
			fakeAdapter.AppendConfigMock = func(output interface{}) error {
				mergeConfigCallCount++
				return nil
			}
			validateCfgCallCount := 0
			cfgManager.validateConfigurationsFunc = func(anchorConfig *AnchorConfig) error {
				validateCfgCallCount++
				return fmt.Errorf("failed to validate")
			}

			result, err := cfgManager.CreateConfigObject()
			assert.NotNil(t, err, "expected to fail")
			assert.NotNil(t, "failed to validate", err.Error())
			assert.Nil(t, result, "expected an empty config object")
			assert.Equal(t, 1, mergeConfigCallCount)
			assert.Equal(t, 1, validateCfgCallCount)
		})
	})
}

var ConfigFilePathGetUserHomeDir = func(t *testing.T) {
	withContext(func(ctx common.Context) {
		withLogging(ctx, t, false, func(logger logger.Logger) {
			result, err := getConfigFilePath()
			assert.Nil(t, err, "expected to succeed")
			assert.NotNil(t, result, "expected non-empty config object")
			assert.Contains(t, result, ".config/anchor")
		})
	})
}

var DefaultRepoClonePathGetClonePath = func(t *testing.T) {
	withContext(func(ctx common.Context) {
		withLogging(ctx, t, false, func(logger logger.Logger) {
			cfgCtxName := "test-ctx"
			result, err := getDefaultRepoClonePath(cfgCtxName)
			assert.Nil(t, err, "expected to succeed")
			assert.NotNil(t, result, "expected non-empty config object")
			assert.Contains(t, result, fmt.Sprintf(".config/anchor/repositories/%s", cfgCtxName))
		})
	})
}

var ConfigTextUnmarshalFromYaml = func(t *testing.T) {
	withContext(func(ctx common.Context) {
		withLogging(ctx, t, false, func(logger logger.Logger) {
			validYamlText := `
config:
  currentContext: test-cfg-ctx
  contexts:
    - name: test-cfg-ctx
      context:
        repository: 
          local:
            path: /invalid/path
`
			result, err := YamlToConfigObj(validYamlText)
			assert.Nil(t, err, "expected to succeed")
			assert.NotNil(t, result, "expected non-empty config object")
			assert.Equal(t, "test-cfg-ctx", result.Config.CurrentContext)
		})
	})
}

var ConfigTextFailToUnmarshal = func(t *testing.T) {
	withContext(func(ctx common.Context) {
		withLogging(ctx, t, false, func(logger logger.Logger) {
			invalidYamlText := "@#$%!@#<invalid> yaml: -configuration"
			result, err := YamlToConfigObj(invalidYamlText)
			assert.NotNil(t, err, "expected to fail")
			assert.Nil(t, result, "expected empty config object")
		})
	})
}

var ConfigObjectToTextMissingObject = func(t *testing.T) {
	withContext(func(ctx common.Context) {
		withLogging(ctx, t, false, func(logger logger.Logger) {
			result, err := ConfigObjToYaml(nil)
			assert.NotNil(t, err, "expected to fail")
			assert.Empty(t, result, "expected empty config object")
		})
	})
}

//var ConfigObjectToTextFailToUnmarshalToYaml = func(t *testing.T) {
//	withContext(func(ctx common.Context) {
//		withLogging(ctx, t, false, func(logger logger.Logger) {
//			result, err := ConfigObjToYaml(nil)
//			assert.NotNil(t, err, "expected to fail")
//			assert.Empty(t, result, "expected empty config object")
//		})
//	})
//}

var ConfigObjectToTextUnmarshalToYaml = func(t *testing.T) {
	withContext(func(ctx common.Context) {
		withLogging(ctx, t, false, func(logger logger.Logger) {
			anchorCfg := &AnchorConfig{
				Config: &Config{},
			}
			anchorCfg.Author = "test-author"
			anchorCfg.License = "test-license"
			anchorCfg.Config.CurrentContext = "test-cfg-ctx"
			result, err := ConfigObjToYaml(anchorCfg)
			assert.Nil(t, err, "expected to fail")
			assert.Contains(t, result, "test-author")
			assert.Contains(t, result, "test-license")
			assert.Contains(t, result, "test-cfg-ctx")
		})
	})
}

var ContextEmtpyConfig = func(t *testing.T) {
	withContext(func(ctx common.Context) {
		withLogging(ctx, t, false, func(logger logger.Logger) {
			result := FromContext(ctx)
			assert.Nil(t, result, "expected not to have config in context")
		})
	})
}

var ContextRetrieveSuccessfully = func(t *testing.T) {
	withContext(func(ctx common.Context) {
		withLogging(ctx, t, false, func(logger logger.Logger) {
			withConfig(ctx, GetDefaultTestConfigText(), func(cfg *AnchorConfig) {
				result := FromContext(ctx)
				assert.NotNil(t, result, "expected not to have config in context")
				assert.Equal(t, cfg, result)
			})
		})
	})
}

var ContextSetConfigInContext = func(t *testing.T) {
	withContext(func(ctx common.Context) {
		withLogging(ctx, t, false, func(logger logger.Logger) {
			anchorCfg := &AnchorConfig{
				Config: &Config{},
			}
			anchorCfg.Author = "test-author"
			anchorCfg.License = "test-license"
			anchorCfg.Config.CurrentContext = "test-cfg-ctx"
			SetInContext(ctx, anchorCfg)
			result := FromContext(ctx)
			assert.NotNil(t, result, "expected to have config in context")
			assert.Contains(t, result.Author, "test-author")
			assert.Contains(t, result.License, "test-license")
			assert.Contains(t, result.Config.CurrentContext, "test-cfg-ctx")
		})
	})
}

var ConfigContextEmptyContexts = func(t *testing.T) {
	withContext(func(ctx common.Context) {
		withLogging(ctx, t, false, func(logger logger.Logger) {
			withConfig(ctx, GetDefaultTestConfigText(), func(cfg *AnchorConfig) {
				result := TryGetConfigContext(nil, "")
				assert.Nil(t, result, "expected to receive empty config context")
			})
		})
	})
}

var ConfigContextContextNotFound = func(t *testing.T) {
	withContext(func(ctx common.Context) {
		withLogging(ctx, t, false, func(logger logger.Logger) {
			withConfig(ctx, GetDefaultTestConfigText(), func(cfg *AnchorConfig) {
				cfgCtxName := "unknown-config-context"
				result := TryGetConfigContext(cfg.Config.Contexts, cfgCtxName)
				assert.Nil(t, result, "expected not to match config context")
			})
		})
	})
}

var ConfigContextFoundContext = func(t *testing.T) {
	withContext(func(ctx common.Context) {
		withLogging(ctx, t, false, func(logger logger.Logger) {
			withConfig(ctx, GetDefaultTestConfigText(), func(cfg *AnchorConfig) {
				cfgCtxName := "1st-anchorfiles"
				result := TryGetConfigContext(cfg.Config.Contexts, cfgCtxName)
				assert.NotNil(t, result, "expected to match config context")
			})
		})
	})
}

var ValidationsVerifyMandatoryConfigEntriesExist = func(t *testing.T) {
	anchorCfg := &AnchorConfig{
		Config: &Config{},
	}
	err := validateConfigurations(anchorCfg)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid configuration attribute")

	cfgYamlText := `
config:
  currentContext: test-cfg-ctx
  contexts:
    - name: test-cfg-ctx
      context:
`
	anchorCfg, _ = YamlToConfigObj(cfgYamlText)
	err = validateConfigurations(anchorCfg)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "context cannot be empty")

	cfgYamlText = `
config:
  currentContext: test-cfg-ctx
  contexts:
    - name: test-cfg-ctx
      context:
        repository:
`
	anchorCfg, _ = YamlToConfigObj(cfgYamlText)
	err = validateConfigurations(anchorCfg)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "context repository cannot be empty")

	cfgYamlText = `
config:
  currentContext: test-cfg-ctx
  contexts:
    - name: test-cfg-ctx
      context:
        repository:
          remote:
          local:
`
	anchorCfg, _ = YamlToConfigObj(cfgYamlText)
	err = validateConfigurations(anchorCfg)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "context repository must have valid remote/local attributes")

	cfgYamlText = `
config:
  currentContext: test-cfg-ctx
  contexts:
    - name: test-cfg-ctx
      context:
        repository: 
          local:
            path: /invalid/path
`
	anchorCfg, _ = YamlToConfigObj(cfgYamlText)
	err = validateConfigurations(anchorCfg)
	assert.Nil(t, err)
}

var PathGetConfigFilePath = func(t *testing.T) {
	cfgManager := NewManager()
	cfgManager.getConfigFilePathFunc = func() (string, error) {
		return "/some/path", nil
	}
	path, err := cfgManager.GetConfigFilePath()
	assert.Nil(t, err)
	assert.Equal(t, "/some/path", path)
}

var DefaultsDoNotOverrideOnExplicitValues = func(t *testing.T) {
	withContext(func(ctx common.Context) {
		withLogging(ctx, t, false, func(logger logger.Logger) {
			cfgYamlText := `
config:
  currentContext: test-cfg-ctx
  contexts:
    - name: test-cfg-ctx
      context:
        repository: 
          remote:
            branch: test-branch
            clonePath: /test/clone/path
`
			anchorCfg, _ := YamlToConfigObj(cfgYamlText)

			cfgManager := NewManager()
			cfgManager.getDefaultRepoClonePathFunc = func(contextName string) (string, error) {
				return "", nil
			}
			err := cfgManager.SetDefaultsPostCreation(anchorCfg)
			assert.Nil(t, err)
			context := TryGetConfigContext(anchorCfg.Config.Contexts, "test-cfg-ctx")
			assert.Equal(t, context.Context.Repository.Remote.Branch, "test-branch")
			assert.Equal(t, context.Context.Repository.Remote.ClonePath, "/test/clone/path")
		})
	})
}

var DefaultsSkipOnMissingRemoteConfig = func(t *testing.T) {
	withContext(func(ctx common.Context) {
		withLogging(ctx, t, false, func(logger logger.Logger) {
			cfgYamlText := `
config:
  currentContext: test-cfg-ctx
  contexts:
    - name: test-cfg-ctx
      context:
        repository: 
          remote:
`
			anchorCfg, _ := YamlToConfigObj(cfgYamlText)

			cfgManager := NewManager()
			err := cfgManager.SetDefaultsPostCreation(anchorCfg)
			assert.Nil(t, err)
		})
	})
}

var DefaultsFailOnCreationOfDefaultClonePathForConfigContext = func(t *testing.T) {
	withContext(func(ctx common.Context) {
		withLogging(ctx, t, false, func(logger logger.Logger) {
			cfgYamlText := `
config:
  currentContext: test-cfg-ctx
  contexts:
    - name: test-cfg-ctx
      context:
        repository: 
          remote:
            url: /remote/url
`
			anchorCfg, _ := YamlToConfigObj(cfgYamlText)

			cfgManager := NewManager()
			cfgManager.getDefaultRepoClonePathFunc = func(contextName string) (string, error) {
				return "", fmt.Errorf("failed to get default repo clone path")
			}
			err := cfgManager.SetDefaultsPostCreation(anchorCfg)
			assert.NotNil(t, err)
			assert.Equal(t, "failed to get default repo clone path", err.Error())
		})
	})
}

var DefaultsSetEmptyEntriesWithDefaultValuesOnAllContexts = func(t *testing.T) {
	withContext(func(ctx common.Context) {
		withLogging(ctx, t, false, func(logger logger.Logger) {
			cfgYamlText := `
config:
  currentContext: test-cfg-ctx-1
  contexts:
    - name: test-cfg-ctx-1
      context:
        repository: 
          remote:
            url: /remote/url/one
    - name: test-cfg-ctx-2
      context:
        repository: 
          remote:
            url: /remote/url/two
    - name: test-cfg-ctx-3
      context:
        repository: 
          remote:
            url: /remote/url/three
`
			anchorCfg, _ := YamlToConfigObj(cfgYamlText)

			cfgManager := NewManager()
			cfgManager.getDefaultRepoClonePathFunc = func(contextName string) (string, error) {
				return fmt.Sprintf("/default/clone/path/%s", contextName), nil
			}
			err := cfgManager.SetDefaultsPostCreation(anchorCfg)
			assert.Nil(t, err)
			context1 := TryGetConfigContext(anchorCfg.Config.Contexts, "test-cfg-ctx-1")
			assert.Equal(t, context1.Context.Repository.Remote.Branch, DefaultRemoteBranch)
			assert.Equal(t, context1.Context.Repository.Remote.ClonePath, "/default/clone/path/test-cfg-ctx-1")

			context2 := TryGetConfigContext(anchorCfg.Config.Contexts, "test-cfg-ctx-2")
			assert.Equal(t, context2.Context.Repository.Remote.Branch, DefaultRemoteBranch)
			assert.Equal(t, context2.Context.Repository.Remote.ClonePath, "/default/clone/path/test-cfg-ctx-2")

			context3 := TryGetConfigContext(anchorCfg.Config.Contexts, "test-cfg-ctx-3")
			assert.Equal(t, context3.Context.Repository.Remote.Branch, DefaultRemoteBranch)
			assert.Equal(t, context3.Context.Repository.Remote.ClonePath, "/default/clone/path/test-cfg-ctx-3")
		})
	})
}

// Duplicate these methods from the with.go test utility to prevent 'import cycle not allowed in test'
// config package is the only exception for this use case
func withContext(f func(ctx common.Context)) {
	// Every context must have a new registry instance
	ctx := common.EmptyAnchorContext(registry.New())
	f(ctx)
}

func withConfig(ctx common.Context, content string, f func(config *AnchorConfig)) {
	cfgManager := NewManager()
	if err := cfgManager.SetupConfigInMemoryLoader(content); err != nil {
		logger.Fatalf("Failed to create a fake config loader. error: %s", err)
	} else {
		cfg, err := cfgManager.CreateConfigObject()
		if err != nil {
			logger.Fatalf("Failed to create a fake config loader. error: %s", err)
		}
		SetInContext(ctx, cfg)
		// set current config context as the active config context
		_ = cfgManager.SwitchActiveConfigContextByName(cfg, cfg.Config.CurrentContext)
		ctx.Registry().Set(Identifier, cfgManager)
		f(cfg)
	}
}

func withLogging(ctx common.Context, t *testing.T, verbose bool, f func(logger logger.Logger)) {
	if fakeLogger, err := logger.CreateFakeTestingLogger(t, verbose); err != nil {
		println("Failed to create a fake testing logger. error: %s", err)
		os.Exit(1)
	} else {
		fakeLoggerManager := logger.CreateFakeLoggerManager()
		fakeLoggerManager.AppendStdoutLoggerMock = func(level string) (logger.Logger, error) {
			return fakeLogger, nil
		}
		fakeLoggerManager.AppendFileLoggerMock = func(level string) (logger.Logger, error) {
			return fakeLogger, nil
		}
		fakeLoggerManager.SetVerbosityLevelMock = func(level string) error {
			resolveAdapter := fakeLogger.(logger.LoggerLogrusAdapter)
			err = resolveAdapter.SetStdoutVerbosityLevel(level)
			if err != nil {
				return err
			}

			err = resolveAdapter.SetFileVerbosityLevel(level)
			if err != nil {
				return err
			}
			return nil
		}

		fakeLoggerManager.GetDefaultLoggerLogFilePathMock = func() (string, error) {
			return "/testing/logger/path", nil
		}
		fakeLoggerManager.SetActiveLoggerMock = func(log *logger.Logger) error {
			return nil
		}

		err = fakeLoggerManager.SetActiveLogger(&fakeLogger)
		if err != nil {
			return
		}
		ctx.Registry().Set(logger.Identifier, fakeLoggerManager)
		f(fakeLogger)
	}
}

func GetTestDataFolder() string {
	path, _ := ioutils.GetWorkingDirectory()
	repoRootPath := ioutils.GetRepositoryAbsoluteRootPath(path)
	if repoRootPath == "" {
		logger.Fatalf("failed to resolve the absolute path of the repository root.")
	}
	configFilePathTest := fmt.Sprintf("%s/test/data", repoRootPath)
	configFilePathTest = strings.TrimSuffix(configFilePathTest, "\n")
	return configFilePathTest
}

func GetTestConfigDirectoryPath() string {
	path, _ := ioutils.GetWorkingDirectory()
	repoRootPath := ioutils.GetRepositoryAbsoluteRootPath(path)
	if repoRootPath == "" {
		logger.Fatalf("failed to resolve the absolute path of the repository root.")
	}
	configFilePathTest := fmt.Sprintf("%s/test/data/config", repoRootPath)
	configFilePathTest = strings.TrimSuffix(configFilePathTest, "\n")
	return configFilePathTest
}

var AppendNewEmptyConfigContextSuccessfully = func(t *testing.T) {
	cfg := &AnchorConfig{
		Config: &Config{
			Contexts: []*Context{},
		},
	}
	createdCtx := AppendEmptyConfigContext(cfg, "new-cfg-ctx")
	assert.NotNil(t, createdCtx, "expected to succeed")
	assert.Equal(t, "new-cfg-ctx", createdCtx.Name)
	assert.Len(t, cfg.Config.Contexts, 1)
	assert.Equal(t, createdCtx, cfg.Config.Contexts[0])
}
