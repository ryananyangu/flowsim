package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zohari-tech/flowsim/controllers"
	"github.com/zohari-tech/flowsim/models"
	"github.com/zohari-tech/flowsim/services"
)

var Routes = map[string]map[string]gin.HandlerFunc{
	"/": {
		"GET": controllers.GetAvailableScreenTypes,
	},
	"/screenview": {
		"GET": controllers.DisplayScreen,
	},
	"/listscreen/item": {
		"POST": controllers.ListScreenItem,
	},
	"/ussd/navigate": {
		"POST": controllers.NavigateUSSD,
	},
}

var UssdRoutes map[string]models.ExternalHandlerFunc

func init() {
	UssdRoutes = map[string]models.ExternalHandlerFunc{
		"1": services.CheckAccountExists,
		"0": services.GetPaymentOptions,
		"2": services.SubmitPayment,
	}
}
