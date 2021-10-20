package config

import "github.com/ZachiNachshon/anchor/internal/common"

var CreateFakeConfigManager = func() *fakeConfigManager {
	return &fakeConfigManager{}
}

type fakeConfigManager struct {
	ConfigManager
	SetupConfigFileLoaderMock           func() error
	SetupConfigInMemoryLoaderMock       func(yaml string) error
	ListenOnConfigFileChangesMock       func(ctx common.Context)
	OverrideConfigMock                  func(cfgToUpdate *AnchorConfig) error
	OverrideConfigEntryMock             func(entryName string, value interface{}) error
	ReadConfigMock                      func(key string) string
	SwitchActiveConfigContextByNameMock func(cfg *AnchorConfig, cfgCtxName string) error
	CreateConfigObjectMock              func(shouldValidateConfig bool) (*AnchorConfig, error)
	GetConfigFilePathMock               func() (string, error)
	SetDefaultsPostCreationMock         func(anchorConfig *AnchorConfig) error
}

func (cm *fakeConfigManager) SetupConfigFileLoader() error {
	return cm.SetupConfigFileLoaderMock()
}

func (cm *fakeConfigManager) SetupConfigInMemoryLoader(yaml string) error {
	return cm.SetupConfigInMemoryLoaderMock(yaml)
}

func (cm *fakeConfigManager) ListenOnConfigFileChanges(ctx common.Context) {
	cm.ListenOnConfigFileChangesMock(ctx)
}

func (cm *fakeConfigManager) OverrideConfig(cfgToUpdate *AnchorConfig) error {
	return cm.OverrideConfigMock(cfgToUpdate)
}

func (cm *fakeConfigManager) OverrideConfigEntry(entryName string, value interface{}) error {
	return cm.OverrideConfigEntryMock(entryName, value)
}

func (cm *fakeConfigManager) ReadConfig(key string) string {
	return cm.ReadConfigMock(key)
}

func (cm *fakeConfigManager) SwitchActiveConfigContextByName(cfg *AnchorConfig, cfgCtxName string) error {
	return cm.SwitchActiveConfigContextByNameMock(cfg, cfgCtxName)
}

func (cm *fakeConfigManager) CreateConfigObject(shouldValidateConfig bool) (*AnchorConfig, error) {
	return cm.CreateConfigObjectMock(shouldValidateConfig)
}

func (cm *fakeConfigManager) GetConfigFilePath() (string, error) {
	return cm.GetConfigFilePathMock()
}

func (cm *fakeConfigManager) SetDefaultsPostCreation(anchorConfig *AnchorConfig) error {
	return cm.SetDefaultsPostCreationMock(anchorConfig)
}
