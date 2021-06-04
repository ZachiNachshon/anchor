package config

import (
	"bytes"
	"fmt"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/ioutils"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const (
	DefaultAuthor  = "Zachi Nachshon <zachi.nachshon@gmail.com>"
	DefaultLicense = "Apache"
)

func FromContext(ctx common.Context) AnchorConfig {
	return ctx.Config().(AnchorConfig)
}

func SetInContext(ctx common.Context, config AnchorConfig) {
	ctx.(common.ConfigSetter).SetConfig(config)
}

type AnchorConfig struct {
	Config  Config
	Author  string
	License string
}

type Config struct {
	RepositoryFiles RepositoryFiles `yaml:"repositoryFiles"`
}

type RepositoryFiles struct {
	Remote Remote `yaml:"remote"`
	Local  Local  `yaml:"local"`
}

type Remote struct {
	Url       string `yaml:"url"`
	Revision  string `yaml:"revision"`
	Branch    string `yaml:"branch"`
	LocalPath string `yaml:"localPath"`
}

type Local struct {
	Path string `yaml:"path"`
}

func initConfigPath() {
	viper.SetConfigName("anchorConfig")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("${HOME}/.config/anchor") // path to look for the config file in
	//viper.AddConfigPath(".")                      // optionally look for config in the working directory
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

func createConfigObject() *AnchorConfig {
	var config Config
	if err := viper.UnmarshalKey("config", &config); err != nil {
		logger.Fatal(fmt.Sprintf("Failed to unmarshal configuration file. error: %s \n", err))
	}

	return &AnchorConfig{
		Config:  config,
		Author:  viper.GetString("author"),
		License: viper.GetString("license"),
	}
}

var ViperConfigInMemoryLoader = func(yaml string) (*AnchorConfig, error) {
	viper.SetConfigType("yaml")
	setDefaults()

	if err := viper.ReadConfig(bytes.NewBuffer([]byte(yaml))); err != nil {
		logger.Errorf("Failed to read config from buffer. error: %s", err)
		return nil, err
	}

	return createConfigObject(), nil
}

var ViperConfigFileLoader = func() (*AnchorConfig, error) {
	initConfigPath()

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

	return createConfigObject(), nil
}

var ResolveAnchorfilesPathFromConfig = func(anchorConfig AnchorConfig) (string, error) {
	// Checks if repositoryFiles config attribute is empty
	if anchorConfig.Config.RepositoryFiles == (RepositoryFiles{}) {
		return "", fmt.Errorf("missing required config value. name: repositoryFiles")
	}

	if localPath, localRepoErr := tryResolveFromLocalPath(anchorConfig.Config.RepositoryFiles.Local); localRepoErr != nil {

		if remotePath, remoteRepoErr := tryResolveFromRemoteRepo(anchorConfig.Config.RepositoryFiles.Remote); remoteRepoErr == nil && remotePath != "" {

			logger.Infof("Using cloned anchorfiles remote repository. path: %s", remotePath)
			return remotePath, nil
		}

	} else if localPath != "" {
		logger.Infof("Using local anchorfiles repository. path: %s", localPath)
		return localPath, nil
	}

	return "", fmt.Errorf("could not resolve anchorfiles local repo path or git tracked repo path")
}

func tryResolveFromLocalPath(local Local) (string, error) {
	pathToUse := local.Path
	if len(pathToUse) > 0 {
		if !ioutils.IsValidPath(pathToUse) {
			return "", fmt.Errorf("local anchorfiles repository path is invalid. path: %s", pathToUse)
		} else {
			return pathToUse, nil
		}
	}
	return "", nil
}

func tryResolveFromRemoteRepo(remote Remote) (string, error) {
	pathToUse := remote.Url

	// TODO: resolve from remote repo in here...

	if len(pathToUse) > 0 {

		if !ioutils.IsValidPath(pathToUse) {
			return "", fmt.Errorf("remote anchorfiles cloned repository path is invalid. path: %s", pathToUse)
		} else {
			return pathToUse, nil
		}
	}
	return "", nil
}
