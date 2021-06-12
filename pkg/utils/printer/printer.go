package printer

import (
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/models"
)

type printerImpl struct {
	Printer
}

func (p *printerImpl) PrintApplications(apps []*models.AppContent) {
	logger.Info("------ Applications ------")
	for _, app := range apps {
		logger.Infof("Name: %s", app)
	}
}
