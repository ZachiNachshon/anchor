package status

import "github.com/ZachiNachshon/anchor/internal/common"

var createFakeOrchestrator = func() *fakeOrchestrator {
	return &fakeOrchestrator{}
}

type fakeOrchestrator struct {
	orchestrator
	bannerMock  func()
	prepareMock func(ctx common.Context) error
	runMock     func(ctx common.Context) error
}

func (o *fakeOrchestrator) banner() {
	o.bannerMock()
}

func (o *fakeOrchestrator) prepare(ctx common.Context) error {
	return o.prepareMock(ctx)
}

func (o *fakeOrchestrator) run(ctx common.Context) error {
	return o.runMock(ctx)
}
