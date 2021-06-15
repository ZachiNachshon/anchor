package printer

import "github.com/ZachiNachshon/anchor/models"

var CreateFakePrinter = func() *fakePrinter {
	return &fakePrinter{}
}

type fakePrinter struct {
	Printer
	PrintApplicationsMock func(apps []*models.AppContent)
}

func (l *fakePrinter) PrintApplications(apps []*models.AppContent) {
	l.PrintApplicationsMock(apps)
}
