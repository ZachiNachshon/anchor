package harness

import (
	"testing"
)

type TestsHarness struct {
	Name string
	Func func(t *testing.T)
}

func RunTests(t *testing.T, tests []TestsHarness) {
	for _, tt := range tests {
		t.Run(tt.Name, tt.Func)
	}
}
