package services

import (
	"github.com/zohari-tech/flowsim/models"
	"github.com/zohari-tech/flowsim/utils"
)

func Navigate(message *models.Message, country_code, networkCode string) (displayText string) {

	inputForScreen, start := message.GetLastScreen()
	if start {
		return message.GetScreen(utils.DEFAULT_SCREEN_LOCATION, country_code,networkCode).Display()
	}
	next_index := inputForScreen.FormatBuild(message.SessionData).NextPage(message.Content)
	return message.GetScreen(next_index, country_code,networkCode).Display()

}
