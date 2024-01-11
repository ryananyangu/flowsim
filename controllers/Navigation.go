package controllers

import (
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zohari-tech/flowsim/models"
	"github.com/zohari-tech/flowsim/services"
	"github.com/zohari-tech/flowsim/utils"
)

func NavigateUSSD(ctx *gin.Context) {

	dial := struct {
		PhoneNumber string `form:"phoneNumber"`
		SessionID   string `form:"sessionId"`
		Content     string `form:"text"`
		ServiceCode string `form:"serviceCode"`
	}{}

	if err := ctx.ShouldBind(&dial); err != nil {
		return
	}
	msg := models.Message{
		Source:         dial.PhoneNumber,
		ConversationID: dial.SessionID,
		Destination:    dial.ServiceCode,
		Direction:      utils.INBOX,
		Content:        dial.Content,
	}

	content := services.Navigate(&msg)

	ctx.HTML(http.StatusOK, "ussd.hml", gin.H{"message": content})

}

func GetAvailableScreenTypes(ctx *gin.Context) {

	screenTypes := []string{
		string(utils.RAW_INPUT_SCREEN),
		string(utils.EXTERNAL_SCREEN),
		string(utils.LIST_SCREEN),
		string(utils.ROUTE_SCREEN),
	}
	ctx.HTML(http.StatusOK, "index.html", gin.H{"screen_types": screenTypes})
}

func DisplayScreen(ctx *gin.Context) {
	screen_type := ctx.Query("screen_type")
	ctx.HTML(http.StatusOK, "screen.html", gin.H{"ScreenType": screen_type})
}

func ListScreenItem(ctx *gin.Context) {
	jsonData, _ := io.ReadAll(ctx.Request.Body)
	log.Println(string(jsonData))
	// ctx.ShouldBind
	ctx.HTML(http.StatusOK, "listscreenitem.html", gin.H{"ls_option": "test"})
}
