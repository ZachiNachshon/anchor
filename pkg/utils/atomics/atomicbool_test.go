package atomics

import (
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_AtomicBoolShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "set and get true",
			Func: SetAndGetTrue,
		},
		{
			Name: "set and get false",
			Func: SetAndGetFalse,
		},
		{
			Name: "set and get default false",
			Func: SetAndGetDefaultFalse,
		},
	}
	harness.RunTests(t, tests)
}

var SetAndGetTrue = func(t *testing.T) {
	a := &AtomicBool{}
	a.Set(true)
	val := a.Get()
	assert.True(t, val)
}

var SetAndGetFalse = func(t *testing.T) {
	a := &AtomicBool{}
	a.Set(false)
	val := a.Get()
	assert.False(t, val)
}

var SetAndGetDefaultFalse = func(t *testing.T) {
	a := &AtomicBool{}
	val := a.Get()
	assert.False(t, val)
}
