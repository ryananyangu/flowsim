package services

import (
	"fmt"
	"time"

	"github.com/zohari-tech/flowsim/database"
	"github.com/zohari-tech/flowsim/models"
	"github.com/zohari-tech/flowsim/utils"
)



func Navigate(message *models.Message, country_code, networkCode string) (displayText string) {

	inputForScreen, start := message.GetLastScreen()
	if start {
		return GetScreen(utils.DEFAULT_SCREEN_LOCATION, message, country_code, networkCode).Display()
	}
	next_index := FormatBuild(message.SessionData, &inputForScreen).NextPage(message.Content)
	return GetScreen(next_index, message, country_code, networkCode).Display()

}

func FormatBuild(sessionData map[string]interface{}, scrn *models.Screen) models.IScreen {
	core := models.CoreScreen{
		Name:        scrn.Name,
		Header:      scrn.Details["Header"].(string),
		IsEnd:       scrn.IsEnd,
		BackEnabled: scrn.BackEnabled,
		ExitEnabled: scrn.ExitEnabled,
	}
	switch scrn.ScreenType {
	case utils.EXTERNAL_SCREEN:

		// logrus.Info(scrn)
		// FIXME: Handle when id not in routes map
		externalHandlerFunc := UssdRoutes[fmt.Sprintf("%d", scrn.ID)]
		return externalHandlerFunc(sessionData, scrn)
	case utils.LIST_SCREEN:
		return models.ListScreen{
			CoreScreen:   core,
			NextLocation: uint(scrn.Details["NextLocation"].(int)),
			Options:      scrn.Details["Options"].([]string),
		}
	case utils.RAW_INPUT_SCREEN:
		return models.RawInputScreen{
			CoreScreen:   core,
			NextLocation: uint(scrn.Details["NextLocation"].(float64)),
		}
	case utils.ROUTE_SCREEN:
		return models.RouteScreen{
			CoreScreen: core,
			Routes:     scrn.Details["Routes"].([]models.Route),
		}
	}
	return nil
}

func GetScreen(location int, msg *models.Message, countr_code, networkCode string) (screen models.IScreen) {
	menu := models.Menu{}
	screens := []models.Screen{}
	nextscreen := models.Screen{}
	// FIXME: Handle these errors
	_ = database.Db.Table("menus").Where("shortcode = ? AND country_code = ? AND telco = ?", msg.Destination, countr_code, networkCode).First(&menu).Error

	_ = database.Db.Table("screens").Where("menu_id", menu.ID).Order("created_at ASC").Find(&screens).Error

	for _, screen_ := range screens {
		if location == int(screen_.Location) {
			nextscreen = screens[location]
			break
		}
	}

	screen = FormatBuild(msg.SessionData, &nextscreen)

	// NOTE: This is where we set the screen to be set
	utils.CacheInstance.Set(msg.ConversationID, nextscreen, time.Minute)
	return
}



