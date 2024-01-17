package services

import (
	"github.com/zohari-tech/flowsim/models"
)

type ExternalHandlerFunc func(map[string]interface{}, *models.Screen) models.IScreen

var UssdRoutes map[string]ExternalHandlerFunc

func init() {
	UssdRoutes = map[string]ExternalHandlerFunc{
		"1": CheckAccountExists,
		"0": GetPaymentOptions,
		"2": SubmitPayment,
	}

}
