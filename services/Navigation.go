package services

import (
	"github.com/zohari-tech/flowsim/models"
	"github.com/zohari-tech/flowsim/utils"
)

func Navigate(message *models.Message) (displayText string) {

	inputForScreen, start := message.GetLastScreen()
	if start {
		return message.GetScreen(utils.DEFAULT_SCREEN_LOCATION).Display()
	}
	next_index := inputForScreen.FormatBuild().NextPage(message.Content)
	return message.GetScreen(next_index).Display()

}
