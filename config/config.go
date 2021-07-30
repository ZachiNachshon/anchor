package config

import (
	"bytes"
	"fmt"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/ioutils"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"strings"
)

const (
	DefaultAuthor                     = "Zachi Nachshon <zachi.nachshon@gmail.com>"
	DefaultLicense                    = "Apache"
	DefaultRemoteBranch               = "master"
	defaultConfigFileName             = "config"
	defaultConfigFileType             = "yaml"
	defaultConfigFolderPathFormat     = "%s/.config/anchor"
	defaultRepoClonePathFormat        = "%s/.config/anchor/repositories"
	defaultLoggerLogFilePathFormat    = "%s/.config/anchor/anchor.log"
	defaultScriptOutputFilePathFormat = "%s/.config/anchor/scripts-output.log"
)

func GetConfigFilePath() (string, error) {
	if homeFolder, err := ioutils.GetUserHomeFolder(); err != nil {
		logger.Errorf("failed to resolve home folder. err: %s", err.Error())
		return "", err
	} else {
		folderPath := fmt.Sprintf(defaultConfigFolderPathFormat, homeFolder)
		return fmt.Sprintf("%s/%s.%s", folderPath, defaultConfigFileName, defaultConfigFileType), nil
	}
}

func GetDefaultRepoClonePath(contextName string) (string, error) {
	if homeFolder, err := ioutils.GetUserHomeFolder(); err != nil {
		logger.Errorf("failed to resolve home folder. err: %s", err.Error())
		return "", err
	} else {
		return fmt.Sprintf(defaultRepoClonePathFormat+"/%s", homeFolder, contextName), nil
	}
}

func GetDefaultLoggerLogFilePath() (string, error) {
	if homeFolder, err := ioutils.GetUserHomeFolder(); err != nil {
		logger.Errorf("failed to resolve home folder. err: %s", err.Error())
		return "", err
	} else {
		return fmt.Sprintf(defaultLoggerLogFilePathFormat, homeFolder), nil
	}
}

func GetDefaultScriptOutputLogFilePath() (string, error) {
	if homeFolder, err := ioutils.GetUserHomeFolder(); err != nil {
		logger.Errorf("failed to resolve home folder. err: %s", err.Error())
		return "", err
	} else {
		return fmt.Sprintf(defaultScriptOutputFilePathFormat, homeFolder), nil
	}
}

func FromContext(ctx common.Context) AnchorConfig {
	return ctx.Config().(AnchorConfig)
}

func SetInContext(ctx common.Context, config AnchorConfig) {
	ctx.(common.ConfigSetter).SetConfig(config)
}

func LoadActiveConfigByName(cfg *AnchorConfig, cfgCtxName string) error {
	loadedActiveCfgCtx := false
	for _, cfgCtx := range cfg.Config.Contexts {
		if strings.EqualFold(cfgCtx.Name, cfgCtxName) {
			logger.Debugf("Loaded active config context. name: %s", cfgCtxName)
			cfg.Config.ActiveContext = cfgCtx
			loadedActiveCfgCtx = true
		}
	}

	if !loadedActiveCfgCtx {
		return fmt.Errorf("could not identify config context. name: %s", cfgCtxName)
	}
	return nil
}

func initConfigPath() error {
	if homeFolder, err := ioutils.GetUserHomeFolder(); err != nil {
		return err
	} else {
		viper.SetConfigName(defaultConfigFileName)
		viper.SetConfigType(defaultConfigFileType)
		viper.AddConfigPath(fmt.Sprintf(defaultConfigFolderPathFormat, homeFolder)) // path to look for the config file in
		//viper.AddConfigPath(".")                      		// optionally look for config in the working directory
		return nil
	}
}

func setDefaults() {
	viper.SetDefault("author", DefaultAuthor)
	viper.SetDefault("license", DefaultLicense)
}

func createConfigFileWithDefaults() {
	err := viper.SafeWriteConfig() // Write defaults
	if err != nil {
		logger.Errorf("Could not create config file with defaults: %s \n", err)
	}
}

func listenOnConfigFileChanges() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// Suggest to git fetch the repository
		logger.Infof("Config file changed:", e.Name)
	})
}

func createConfigObject() (*AnchorConfig, error) {
	var config Config
	if err := viper.UnmarshalKey("config", &config); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal configuration file. error: %s \n", err)
	}

	err := validateConfigurations(&config)
	if err != nil {
		return nil, err
	}

	setDefaultsPostCreation(&config)

	return &AnchorConfig{
		Config:  &config,
		Author:  viper.GetString("author"),
		License: viper.GetString("license"),
	}, nil
}

func validateConfigurations(cfg *Config) error {
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

func setDefaultsPostCreation(cfg *Config) {
	for _, ctx := range cfg.Contexts {
		if ctx.Context != nil {
			repository := ctx.Context.Repository
			if repository.Remote == nil {
				// Local must be set else validation would fail
				return
			}
			if repository.Remote.ClonePath == "" {
				clonePath, err := GetDefaultRepoClonePath(ctx.Name)
				if err != nil {
					logger.Fatal("failed to resolve default repo clone path")
				}
				repository.Remote.ClonePath = clonePath
			}

			if repository.Remote.Branch == "" {
				repository.Remote.Branch = DefaultRemoteBranch
			}
		}
	}
}

var ViperConfigInMemoryLoader = func(yaml string) (*AnchorConfig, error) {
	viper.SetConfigType("yaml")
	setDefaults()

	if err := viper.ReadConfig(bytes.NewBuffer([]byte(yaml))); err != nil {
		logger.Errorf("Failed to read config from buffer. error: %s", err)
		return nil, err
	}

	return createConfigObject()
}

var ViperConfigFileLoader = func() (*AnchorConfig, error) {
	err := initConfigPath()
	if err != nil {
		logger.Fatalf("Failed to initialize anchor configuration path. error: %s", err.Error())
		return nil, err
	}

	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {
		// Handle errors reading the config file
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			setDefaults()
			createConfigFileWithDefaults()
			listenOnConfigFileChanges()
		} else {
			logger.Errorf("Config file was found but an error occurred. error: %s", err)
			return nil, err
		}
	}

	// Every viper.Get request auto checks for ANCHOR_<flag-name> before reading from config file
	viper.SetEnvPrefix("ANCHOR")
	viper.AutomaticEnv()

	return createConfigObject()
}
