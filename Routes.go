package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zohari-tech/flowsim/controllers"
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
}
