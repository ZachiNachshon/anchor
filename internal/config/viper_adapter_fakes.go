package config

import (
	"github.com/fsnotify/fsnotify"
)

var CreateFakeViperConfigAdapter = func() *fakeViperConfigAdapter {
	return &fakeViperConfigAdapter{}
}

type fakeViperConfigAdapter struct {
	ConfigViperAdapter
	SetConfigPathMock                 func(path string)
	LoadConfigFromFileMock            func() error
	LoadConfigFromTextMock            func(yaml string) error
	RegisterConfigChangesListenerMock func(callback func(e fsnotify.Event))
	UpdateAllMock                     func(cfgToUpdate *AnchorConfig) error
	UpdateEntryMock                   func(entryName string, value interface{}) error
	GetConfigByKeyMock                func(key string) string
	SetDefaultsMock                   func()
	SetEnvVarsMock                    func()
	MergeConfigMock                   func(output interface{}) error
	flushToNewConfigFileMock          func() error
	flushMock                         func() error
}

func (ca *fakeViperConfigAdapter) SetConfigPath(path string) {
	ca.SetConfigPathMock(path)
}

func (ca *fakeViperConfigAdapter) LoadConfigFromFile() error {
	return ca.LoadConfigFromFileMock()
}

func (ca *fakeViperConfigAdapter) LoadConfigFromText(yaml string) error {
	return ca.LoadConfigFromTextMock(yaml)
}

func (ca *fakeViperConfigAdapter) RegisterConfigChangesListener(callback func(e fsnotify.Event)) {
	ca.RegisterConfigChangesListenerMock(callback)
}

func (ca *fakeViperConfigAdapter) UpdateAll(cfgToUpdate *AnchorConfig) error {
	return ca.UpdateAllMock(cfgToUpdate)
}

func (ca *fakeViperConfigAdapter) UpdateEntry(entryName string, value interface{}) error {
	return ca.UpdateEntryMock(entryName, value)
}

func (ca *fakeViperConfigAdapter) GetConfigByKey(key string) string {
	return ca.GetConfigByKeyMock(key)
}

func (ca *fakeViperConfigAdapter) SetDefaults() {
	ca.SetDefaultsMock()
}

func (ca *fakeViperConfigAdapter) SetEnvVars() {
	ca.SetEnvVarsMock()
}

func (ca *fakeViperConfigAdapter) MergeConfig(output interface{}) error {
	return ca.MergeConfigMock(output)
}

func (ca *fakeViperConfigAdapter) flushToNewConfigFile() error {
	return ca.flushToNewConfigFileMock()
}

func (ca *fakeViperConfigAdapter) flush() error {
	return ca.flushMock()
}
