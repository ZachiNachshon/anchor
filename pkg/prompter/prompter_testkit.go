package prompter

import "github.com/ZachiNachshon/anchor/models"

func GetAppByName(appsArr []*models.AppContent, name string) *models.AppContent {
	for _, v := range appsArr {
		if v.Name == name {
			return v
		}
	}
	return nil
}

func GetInstructionItemById(instructions *models.Instructions, id string) *models.PromptItem {
	for _, v := range instructions.Items {
		if v.Id == id {
			return v
		}
	}
	return nil
}
