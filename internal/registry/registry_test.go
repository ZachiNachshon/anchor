package registry

import (
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_RegistryShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "initialize only once",
			Func: InitializeOnlyOnce,
		},
		{
			Name: "set and get valid item",
			Func: SetAndGetValidItem,
		},
		{
			Name: "set and get invalid item",
			Func: SetAndGetInvalidItem,
		},
		{
			Name: "set and safe get valid item",
			Func: SetAndSafeGetValidItem,
		},
		{
			Name: "set and safe get invalid item",
			Func: SetAndSafeGetInvalidItem,
		},
	}
	harness.RunTests(t, tests)
}

var InitializeOnlyOnce = func(t *testing.T) {
	regFirst := Initialize()
	regSec := Initialize()
	assert.Equal(t, regFirst, regSec)
}

var SetAndGetValidItem = func(t *testing.T) {
	reg := New()
	testItems := []string{"one", "two", "three"}
	reg.Set("items", testItems)
	resolved := reg.Get("items")
	assert.NotNil(t, resolved)
	assert.Equal(t, 3, len(resolved.([]string)))
}

var SetAndGetInvalidItem = func(t *testing.T) {
	reg := New()
	testItems := []string{"one", "two", "three"}
	reg.Set("items", testItems)
	resolved := reg.Get("not-exist-items")
	assert.Nil(t, resolved)
}

var SetAndSafeGetValidItem = func(t *testing.T) {
	reg := New()
	testItems := []string{"one", "two", "three"}
	reg.Set("items", testItems)
	resolved, err := reg.SafeGet("items")
	assert.Nil(t, err, "expected to succeed")
	assert.NotNil(t, resolved)
	assert.Equal(t, 3, len(resolved.([]string)))
}

var SetAndSafeGetInvalidItem = func(t *testing.T) {
	reg := New()
	testItems := []string{"one", "two", "three"}
	reg.Set("items", testItems)
	resolved, err := reg.SafeGet("not-exist-items")
	assert.NotNil(t, err, "expected to fail")
	assert.Nil(t, resolved)
}
