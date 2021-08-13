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
	CreateConfigObject() (*AnchorConfig, error)

	GetConfigFilePath() (string, error)
	GetDefaultRepoClonePath(contextName string) (string, error)
}

type configManagerImpl struct {
	ConfigManager
	adapter ConfigViperAdapter
}

func NewManager() ConfigManager {
	return &configManagerImpl{
		adapter: NewAdapter(),
	}
}

func (cm *configManagerImpl) SetupConfigFileLoader() error {
	path, err := cm.GetConfigFilePath()
	if err != nil {
		return err
	}

	cm.adapter.SetConfigPath(path)

	err = cm.adapter.LoadConfigFromFile()
	if err != nil {
		return err
	}

	cm.adapter.SetEnvVars()
	return nil
}

func (cm *configManagerImpl) SetupConfigInMemoryLoader(yaml string) error {
	return cm.adapter.LoadConfigFromText(yaml)
}

func (cm *configManagerImpl) ListenOnConfigFileChanges(ctx common.Context) {
	cm.adapter.RegisterConfigChangesListener(func(e fsnotify.Event) {
		if err := cm.adapter.MergeConfig(ctx.Config()); err != nil {
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

func (cm *configManagerImpl) CreateConfigObject() (*AnchorConfig, error) {
	cfg := &AnchorConfig{
		Author:  cm.ReadConfig("author"),
		License: cm.ReadConfig("license"),
	}

	err := cm.adapter.MergeConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("Failed to merge configuration from file. error: %s \n", err)
	}

	err = validateConfigurations(cfg)
	if err != nil {
		return nil, err
	}

	setDefaultsPostCreation(cfg, cm.GetDefaultRepoClonePath)
	return cfg, nil
}

func (cm *configManagerImpl) GetConfigFilePath() (string, error) {
	if homeFolder, err := ioutils.GetUserHomeDirectory(); err != nil {
		logger.Errorf("failed to resolve home folder. err: %s", err.Error())
		return "", err
	} else {
		folderPath := fmt.Sprintf(defaultConfigFolderPathFormat, homeFolder)
		return fmt.Sprintf("%s/%s.%s", folderPath, defaultConfigFileName, defaultConfigFileType), nil
	}
}

func (cm *configManagerImpl) GetDefaultRepoClonePath(contextName string) (string, error) {
	if homeFolder, err := ioutils.GetUserHomeDirectory(); err != nil {
		logger.Errorf("failed to resolve home folder. err: %s", err.Error())
		return "", err
	} else {
		return fmt.Sprintf(defaultRepoClonePathFormat+"/%s", homeFolder, contextName), nil
	}
}

func YamlToConfigObj(yamlText string) AnchorConfig {
	cfg := AnchorConfig{}
	if err := converters.UnmarshalYamlToObj(yamlText, &cfg); err != nil {
		logger.Fatalf("Failed to generate config template. error: %s", err.Error())
	}
	return cfg
}

func ConfigObjToYaml(cfg *AnchorConfig) (string, error) {
	if out, err := converters.UnmarshalObjToYaml(cfg); err != nil {
		logger.Errorf("Failed to generate config template. error: %s", err.Error())
		return "", nil
	} else {
		return out, nil
	}
}

func FromContext(ctx common.Context) *AnchorConfig {
	return ctx.Config().(*AnchorConfig)
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

func setDefaultsPostCreation(anchorConfig *AnchorConfig, getRepoClonePathFunc func(contextName string) (string, error)) {
	cfg := anchorConfig.Config
	for _, ctx := range cfg.Contexts {
		if ctx.Context != nil {
			repo := ctx.Context.Repository
			if repo.Remote == nil {
				// Local must be set else validation would fail
				return
			}
			if repo.Remote.ClonePath == "" {
				clonePath, err := getRepoClonePathFunc(ctx.Name)
				if err != nil {
					logger.Fatal("failed to resolve default repo clone path")
				}
				repo.Remote.ClonePath = clonePath
			}

			if repo.Remote.Branch == "" {
				repo.Remote.Branch = DefaultRemoteBranch
			}
		}
	}
}
