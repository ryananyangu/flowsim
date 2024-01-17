package services

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/zohari-tech/flowsim/database"
	"github.com/zohari-tech/flowsim/models"
	"github.com/zohari-tech/flowsim/utils"
)

const (
	AVAILABLE_ACCOUNT_NEXT = 2
)

func CheckAccountExists(sessionData map[string]interface{}, screen_ *models.Screen) (screen models.IScreen) {

	accounts := []struct {
		AccountNumber string `json:"external_id" binding:"required"`
		AccountStatus string `json:"account_status" binding:"required"`
	}{}
	routes := []models.Route{} // FIXME: Get this from details of the incomming screen
	phonenumber := sessionData["DialingNumber"]
	url := fmt.Sprintf("%s/api/v1/contact/list?phone=%s", os.Getenv("CRM_URL"), phonenumber)
	response, err := utils.Request("", map[string][]string{}, url, "GET", utils.RequestOptions{})
	if err != nil {
		_ = database.Db.Table("screens").Where("menu_id = ? AND location = ?", screen_.MenuID, screen_.Details["NextLocation"].(int)).First(&screen_).Error //NOTE: DEFAULT NEXT SCREEN
		// FIXME: Get the enter account number screen.
		return FormatBuild(sessionData, screen_)
	}

	if err = json.Unmarshal([]byte(response), &accounts); err != nil {
		logrus.Error(err)
		_ = database.Db.Table("screens").Where("menu_id = ? AND location = ?", screen_.MenuID, screen_.Details["NextLocation"].(int)).First(&screen_).Error //NOTE: DEFAULT NEXT SCREEN
		return FormatBuild(sessionData, screen_)
	}

	for _, account := range accounts {
		if strings.EqualFold(account.AccountStatus, "ACTIVE") {
			routes = append(routes, models.Route{
				Value:        account.AccountNumber,
				IsEnd:        false,
				NextLocation: AVAILABLE_ACCOUNT_NEXT,
			})

		}
	}
	if len(routes) <= 0 {
		_ = database.Db.Table("screens").Where("menu_id = ? AND location = ?", screen_.MenuID, screen_.Details["NextLocation"].(int)).First(&screen_).Error //NOTE: DEFAULT NEXT SCREEN
		// FIXME: Get the enter account number screen.
		return FormatBuild(sessionData, screen_)
	}

	return models.RouteScreen{
		CoreScreen: models.CoreScreen{
			Name:        screen_.Name,
			IsEnd:       screen_.IsEnd,
			BackEnabled: screen_.BackEnabled,
			Header:      screen_.Details["Header"].(string),
			ExitEnabled: screen_.ExitEnabled,
		},
		Routes: routes,
	}
}
