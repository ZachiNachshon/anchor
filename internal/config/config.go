package config

import (
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/common"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/converters"
	"github.com/ZachiNachshon/anchor/pkg/utils/ioutils"
	"github.com/fsnotify/fsnotify"
	"strings"
)

const (
	DefaultAuthor                 = "Zachi Nachshon <zachi.nachshon@gmail.com>"
	DefaultLicense                = "Apache"
	DefaultRemoteBranch           = "master"
	defaultConfigFileName         = "config"
	defaultConfigFileType         = "yaml"
	defaultConfigFolderPathFormat = "%s/.config/anchor"
	defaultRepoClonePathFormat    = "%s/.config/anchor/repositories"
)

const (
	Identifier string = "config-manager"
)

type ConfigManager interface {
	SetupConfigFileLoader() error
	SetupConfigInMemoryLoader(yaml string) error
	ListenOnConfigFileChanges(ctx common.Context)

	OverrideConfig(cfgToUpdate *AnchorConfig) error
	OverrideConfigEntry(entryName string, value interface{}) error
	ReadConfig(key string) string

	SwitchActiveConfigContextByName(cfg *AnchorConfig, cfgCtxName string) error
	CreateConfigObject(shouldValidateConfig bool) (*AnchorConfig, error)

	GetConfigFilePath() (string, error)
	SetDefaultsPostCreation(anchorConfig *AnchorConfig) error
}

type configManagerImpl struct {
	ConfigManager
	adapter ConfigViperAdapter

	getConfigFilePathFunc       func() (string, error)
	validateConfigurationsFunc  func(anchorConfig *AnchorConfig) error
	getDefaultRepoClonePathFunc func(contextName string) (string, error)
}

func NewManager() *configManagerImpl {
	return &configManagerImpl{
		adapter:                     NewAdapter(),
		getConfigFilePathFunc:       getConfigFilePath,
		validateConfigurationsFunc:  validateConfigurations,
		getDefaultRepoClonePathFunc: getDefaultRepoClonePath,
	}
}

func (cm *configManagerImpl) SetupConfigFileLoader() error {
	path, err := cm.getConfigFilePathFunc()
	if err != nil {
		return err
	}

	err = cm.adapter.SetConfigPath(path)
	if err != nil {
		return err
	}

	err = cm.adapter.LoadConfigFromFile()
	if err != nil {
		return err
	}

	err = cm.adapter.SetEnvVars()
	if err != nil {
		return err
	}
	return nil
}

func (cm *configManagerImpl) SetupConfigInMemoryLoader(yaml string) error {
	return cm.adapter.LoadConfigFromText(yaml)
}

func (cm *configManagerImpl) ListenOnConfigFileChanges(ctx common.Context) {
	cm.adapter.RegisterConfigChangesListener(func(e fsnotify.Event) {
		if err := cm.adapter.AppendConfig(ctx.Config()); err != nil {
			// TODO: is it really fatal?
			logger.Fatalf("Failed to reload in-memory configuration file after change was identified. error: %s", err.Error())
		} else {
			logger.Debugf("Config file changed, in-memory config state updated. event: %s", e.Name)
		}
	})
}

// OverrideConfig merge the configuration from disk with the in-memory configuration
func (cm *configManagerImpl) OverrideConfig(cfgToUpdate *AnchorConfig) error {
	return cm.adapter.UpdateAll(cfgToUpdate)
}

// OverrideConfigEntry allows to update delimited config values
// Note that updating dynamic config context values is not supported
func (cm *configManagerImpl) OverrideConfigEntry(entryName string, value interface{}) error {
	return cm.adapter.UpdateEntry(entryName, value)
}

func (cm *configManagerImpl) ReadConfig(key string) string {
	return cm.adapter.GetConfigByKey(key)
}

func (cm *configManagerImpl) SwitchActiveConfigContextByName(cfg *AnchorConfig, cfgCtxName string) error {
	if cfgCtx := TryGetConfigContext(cfg.Config.Contexts, cfgCtxName); cfgCtx == nil {
		return fmt.Errorf("could not identify config context. name: %s", cfgCtxName)
	} else {
		logger.Debugf("Loaded active config context. name: %s", cfgCtxName)
		cfg.Config.ActiveContext = cfgCtx
		return nil
	}
}

func (cm *configManagerImpl) CreateConfigObject(shouldValidateConfig bool) (*AnchorConfig, error) {
	cfg := &AnchorConfig{
		Author:  cm.ReadConfig("author"),
		License: cm.ReadConfig("license"),
	}

	err := cm.adapter.AppendConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("Failed to merge configuration from file. error: %s \n", err)
	}

	if shouldValidateConfig {
		err = cm.validateConfigurationsFunc(cfg)
		if err != nil {
			return nil, err
		}
	}

	err = cm.SetDefaultsPostCreation(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (cm *configManagerImpl) GetConfigFilePath() (string, error) {
	return cm.getConfigFilePathFunc()
}

func (cm *configManagerImpl) SetDefaultsPostCreation(anchorConfig *AnchorConfig) error {
	cfg := anchorConfig.Config
	for _, ctx := range cfg.Contexts {
		if ctx.Context != nil {
			repo := ctx.Context.Repository
			if repo.Remote == nil {
				// Local must be set else validation would fail
				continue
			}
			if len(repo.Remote.ClonePath) == 0 {
				clonePath, err := cm.getDefaultRepoClonePathFunc(ctx.Name)
				if err != nil {
					logger.Error("failed to resolve default repo clone path")
					return err
				}
				repo.Remote.ClonePath = clonePath
			}

			if len(repo.Remote.Branch) == 0 {
				repo.Remote.Branch = DefaultRemoteBranch
			}
		}
	}
	return nil
}

func getDefaultRepoClonePath(contextName string) (string, error) {
	if homeFolder, err := ioutils.GetUserHomeDirectory(); err != nil {
		logger.Errorf("failed to resolve home folder. err: %s", err.Error())
		return "", err
	} else {
		return fmt.Sprintf(defaultRepoClonePathFormat+"/%s", homeFolder, contextName), nil
	}
}

func getConfigFilePath() (string, error) {
	if homeFolder, err := ioutils.GetUserHomeDirectory(); err != nil {
		logger.Errorf("failed to resolve home folder. err: %s", err.Error())
		return "", err
	} else {
		return fmt.Sprintf(defaultConfigFolderPathFormat, homeFolder), nil
	}
}

func YamlToConfigObj(yamlText string) (*AnchorConfig, error) {
	cfg := &AnchorConfig{}
	if err := converters.UnmarshalYamlToObj(yamlText, &cfg); err != nil {
		logger.Errorf("Failed to generate config template. error: %s", err.Error())
		return nil, err
	}
	return cfg, nil
}

func ConfigObjToYaml(cfg *AnchorConfig) (string, error) {
	if cfg == nil {
		return "", fmt.Errorf("cannot unmarshal nil config object to string")
	}
	if out, err := converters.UnmarshalObjToYaml(cfg); err != nil {
		logger.Errorf("Failed to generate config template. error: %s", err.Error())
		return "", err
	} else {
		return out, nil
	}
}

func FromContext(ctx common.Context) *AnchorConfig {
	if ctx.Config() != nil {
		return ctx.Config().(*AnchorConfig)
	}
	return nil
}

func SetInContext(ctx common.Context, config *AnchorConfig) {
	ctx.(common.ConfigSetter).SetConfig(config)
}

func TryGetConfigContext(contexts []*Context, cfgCtxName string) *Context {
	if contexts == nil {
		return nil
	}
	for _, cfgCtx := range contexts {
		if strings.EqualFold(cfgCtx.Name, cfgCtxName) {
			return cfgCtx
		}
	}
	return nil
}

func AppendEmptyConfigContext(cfg *AnchorConfig, name string) *Context {
	cfgCtx := emptyContext()
	cfgCtx.Name = name
	cfg.Config.Contexts = append(cfg.Config.Contexts, cfgCtx)
	return cfgCtx
}

func validateConfigurations(anchorConfig *AnchorConfig) error {
	cfg := anchorConfig.Config
	if cfg.Contexts == nil || len(cfg.Contexts) == 0 {
		return fmt.Errorf("invalid configuration attribute. name: contexts")
	}

	for _, ctx := range cfg.Contexts {
		if ctx.Context == nil {
			return fmt.Errorf("invalid configuration attribute. context cannot be empty")
		}

		if ctx.Context.Repository == nil {
			return fmt.Errorf("invalid configuration attribute. context repository cannot be empty")
		}

		if ctx.Context.Repository.Remote == nil &&
			ctx.Context.Repository.Local == nil {
			return fmt.Errorf("invalid configuration attribute. context repository must have valid remote/local attributes")
		}
	}
	return nil
}
