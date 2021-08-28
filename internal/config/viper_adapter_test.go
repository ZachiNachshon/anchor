package config

import (
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_ViperAdapterShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "set config path",
			Func: SetConfigPath,
		},
		{
			Name: "failed to read config from file",
			Func: FailedToReadConfigFromFile,
		},
		{
			Name: "read config from existing config file",
			Func: ReadConfigFromExistingConfigFile,
		},
		{
			Name: "fail to load config yaml text",
			Func: FailToLoadConfigYamlText,
		},
		{
			Name: "load config yaml text",
			Func: LoadConfigYamlText,
		},
		{
			Name: "register to config changes",
			Func: RegisterToConfigChanges,
		},
		{
			Name: "update all config: succeed",
			Func: UpdateAllConfigSucceed,
		},
		//{
		//	Name: "update all config: fail to marshal",
		//	Func: UpdateAllConfigFailToMarshal,
		//},
		//{
		//	Name: "update all config: fail to merge",
		//	Func: UpdateAllConfigFailToMerge,
		//},
		{
			Name: "updated single config: succeed",
			Func: UpdateSingleConfigSucceed,
		},
		//{
		//	Name: "updated single config: entry not set",
		//	Func: UpdateSingleConfigEntryNotSet,
		//},
		{
			Name: "get config entry to key",
			Func: GetConfigEntryByKey,
		},
		{
			Name: "set env vars",
			Func: SetEnvVars,
		},
		//{
		//	Name: "fail to append config",
		//	Func: FailToAppendConfig,
		//},
		{
			Name: "append config",
			Func: AppendConfig,
		},
	}
	harness.RunTests(t, tests)
}

var SetConfigPath = func(t *testing.T) {
	viperAdapter := NewAdapter()
	err := viperAdapter.SetConfigPath("/test/config/path")
	assert.Nil(t, err, "expected to succeed")
}

var FailedToReadConfigFromFile = func(t *testing.T) {
	testDataFolder := GetTestDataFolder()
	viperAdapter := NewAdapter()
	err := viperAdapter.SetConfigPath(testDataFolder + "/config_bad")
	err = viperAdapter.LoadConfigFromFile()
	assert.NotNil(t, err, "expected to succeed")
	assert.Contains(t, err.Error(), "could not read configuration from file")
}

var ReadConfigFromExistingConfigFile = func(t *testing.T) {
	testConfigFilePath := GetTestConfigDirectoryPath()
	viperAdapter := NewAdapter()
	err := viperAdapter.SetConfigPath(testConfigFilePath)
	err = viperAdapter.LoadConfigFromFile()
	assert.Nil(t, err, "expected to succeed")
}

var FailToLoadConfigYamlText = func(t *testing.T) {
	invalidYamlText := "@#$%!@#<invalid> yaml: -configuration"
	viperAdapter := NewAdapter()
	err := viperAdapter.LoadConfigFromText(invalidYamlText)
	assert.NotNil(t, err, "expected to fail")
	assert.Contains(t, err.Error(), "failed to read config from buffer")
}

var LoadConfigYamlText = func(t *testing.T) {
	validYamlText := `
config:
  currentContext: test-cfg-ctx
  contexts:
    - name: test-cfg-ctx
      context:
        repository: 
          local:
            path: /test/local/path
`
	viperAdapter := NewAdapter()
	err := viperAdapter.LoadConfigFromText(validYamlText)
	assert.Nil(t, err, "expected to succeed")
	assert.Equal(t, DefaultAuthor, viperAdapter.GetConfigByKey("author"))
	assert.Equal(t, DefaultLicense, viperAdapter.GetConfigByKey("license"))
}

var RegisterToConfigChanges = func(t *testing.T) {
	viperAdapter := NewAdapter()
	tempConfigFile := createTempConfigFile(t, viperAdapter)
	defer os.Remove(tempConfigFile)

	cfgChangeCallCount := 0
	viperAdapter.RegisterConfigChangesListener(func(e fsnotify.Event) {
		cfgChangeCallCount++
	})

	err := viperAdapter.UpdateEntry("currentContext", "updated-cfg-ctx-name")
	assert.Nil(t, err, "expected update entry to succeed")

	// TODO: need to solve the callback not triggering on test issue
	//assert.Equal(t, 1, cfgChangeCallCount)

	//yamlText, err := os.ReadFile(tempConfigFile)
	//cfgObj, err := YamlToConfigObj(string(yamlText))
	//assert.Nil(t, err, "expected to read config YAML successfully")
	//cfgObj.Author = "updated author"
	//cfgObj.Config = &Config{}
	//cfgObj.Config.CurrentContext = "updated-cfg-context"
	//updatedYamlText, err := ConfigObjToYaml(cfgObj)
	//err = os.WriteFile(tempConfigFile, []byte(updatedYamlText), 0)
	//assert.Nil(t, err, "expected write updated config to a temp file successfully")
}

var UpdateAllConfigSucceed = func(t *testing.T) {
	viperAdapter := NewAdapter()
	tempConfigFile := createTempConfigFile(t, viperAdapter)
	defer os.Remove(tempConfigFile)

	anchorCfg := &AnchorConfig{
		Config: &Config{},
	}
	anchorCfg.Author = "test-author"
	anchorCfg.License = "test-license"
	anchorCfg.Config.CurrentContext = "updated-test-cfg-ctx"
	err := viperAdapter.UpdateAll(anchorCfg)
	assert.Nil(t, err, "expected to succeed")
}

var UpdateSingleConfigSucceed = func(t *testing.T) {
	viperAdapter := NewAdapter()
	tempConfigFile := createTempConfigFile(t, viperAdapter)
	defer os.Remove(tempConfigFile)
	err := viperAdapter.UpdateEntry("currentContext", "updated-cfg-ctx-name")
	assert.Nil(t, err, "expected to succeed")
}

var GetConfigEntryByKey = func(t *testing.T) {
	viperAdapter := NewAdapter()
	tempConfigFile := createTempConfigFile(t, viperAdapter)
	defer os.Remove(tempConfigFile)
	authorValue := viperAdapter.GetConfigByKey("Author")
	LicenseValue := viperAdapter.GetConfigByKey("License")
	assert.Equal(t, authorValue, DefaultAuthor)
	assert.Equal(t, LicenseValue, DefaultLicense)
}

var SetEnvVars = func(t *testing.T) {
	viperAdapter := NewAdapter()
	err := viperAdapter.SetEnvVars()
	assert.Nil(t, err, "expected to succeed")
}

var AppendConfig = func(t *testing.T) {
	anchorCfg := &AnchorConfig{
		Author:  "test-author",
		License: "test-license",
	}

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
	viperAdapter := NewAdapter()
	err := viperAdapter.LoadConfigFromText(cfgYamlText)
	assert.Nil(t, err, "expected to load yaml config to memory")
	err = viperAdapter.AppendConfig(anchorCfg)
	assert.Nil(t, err)
	assert.NotNil(t, anchorCfg.Config)
	assert.Equal(t, "test-cfg-ctx", anchorCfg.Config.CurrentContext)
}

func createTempConfigFile(t *testing.T, adapter ConfigViperAdapter) string {
	tempDir := os.TempDir()
	tempConfigFile := tempDir + "/config.yaml"
	err := adapter.SetConfigPath(tempDir)
	assert.Nil(t, err, "expected to set temp config file path")

	err = adapter.LoadConfigFromFile()
	assert.Nil(t, err, "expected to load config from temp file")
	return tempConfigFile
}
