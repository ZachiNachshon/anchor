package prompter

import (
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_BellSkipperShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "not ring the bell when bell character was sent",
			Func: NotRingTheBellWhenBellCharacterWasSent,
		},
		//{
		//	Name: "pipe to stderr when bell character was not sent",
		//	Func: PipetoStderrWhenBellCharacterWasNotSent,
		//},
	}
	harness.RunTests(t, tests)
}

var NotRingTheBellWhenBellCharacterWasSent = func(t *testing.T) {
	bs := &bellSkipper{}
	c, err := bs.Write([]byte{7})
	assert.Equal(t, 0, c)
	assert.Nil(t, err)
}

//var PipetoStderrWhenBellCharacterWasNotSent = func(t *testing.T) {
//	bs := &bellSkipper{}
//	c, err := bs.Write([]byte{'a'})
//	assert.NotEqual(t, 0, c)
//	assert.Nil(t, err)
//}
