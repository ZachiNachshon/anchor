package prompter

import (
	"github.com/ZachiNachshon/anchor/test/harness"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_BellSkipperShould(t *testing.T) {
	tests := []harness.TestsHarness{
		{
			Name: "not ring the bell when bell character was sent",
			Func: NotRingTheBellWhenBellCharacterWasSent,
		},
		{
			Name: "pipe to stderr when bell character was not sent",
			Func: PipeToStderrWhenBellCharacterWasNotSent,
		},
	}
	harness.RunTests(t, tests)
}

var NotRingTheBellWhenBellCharacterWasSent = func(t *testing.T) {
	bs := newBellSkipper()
	bytesWritten, err := bs.Write([]byte{7})
	assert.Equal(t, 0, bytesWritten)
	assert.Nil(t, err)
}

var PipeToStderrWhenBellCharacterWasNotSent = func(t *testing.T) {
	tempFile := os.TempDir() + "bell_skipper.txt"
	dummyStderrFile, _ := os.Create(tempFile)
	defer dummyStderrFile.Close()
	bs := newBellSkipper()
	bs.stderr = dummyStderrFile
	bytesWritten, err := bs.Write([]byte{'a'})
	assert.Equal(t, 1, bytesWritten)
	assert.Nil(t, err)
	err = bs.Close()
	assert.Nil(t, err)
}
