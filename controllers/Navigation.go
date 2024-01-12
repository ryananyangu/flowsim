package controllers

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nyaruka/phonenumbers"
	"github.com/sirupsen/logrus"

	"github.com/zohari-tech/flowsim/models"
	"github.com/zohari-tech/flowsim/services"
	"github.com/zohari-tech/flowsim/utils"
)

func NavigateUSSD(ctx *gin.Context) {

	dial := struct {
		PhoneNumber string `form:"phoneNumber" binding:"required"`
		SessionID   string `form:"sessionId" binding:"required"`
		Content     string `form:"text" binding:"-"`
		ServiceCode string `form:"serviceCode" binding:"required"`
		NetworkCode string `form:"networkCode" binding:"required"`
	}{}

	if err := ctx.ShouldBind(&dial); err != nil {
		logrus.Error(err)
		ctx.HTML(http.StatusOK, "ussd.html", gin.H{"message": "system error"})
		return
	}
	parsed, err := phonenumbers.Parse(dial.PhoneNumber, "KE")
	if err != nil {
		logrus.Error(err)
		ctx.HTML(http.StatusOK, "ussd.html", gin.H{"message": "system error"})
		return
	}
	msg := models.Message{
		Source:         dial.PhoneNumber,
		ConversationID: dial.SessionID,
		Destination:    dial.ServiceCode,
		Direction:      utils.INBOX,
		Content:        dial.Content,
	}

	content := services.Navigate(&msg, fmt.Sprintf("%d", *parsed.CountryCode), dial.NetworkCode)

	ctx.HTML(http.StatusOK, "ussd.html", gin.H{"message": content})

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
	ctx.HTML(http.StatusOK, "listscreenitem.html", gin.H{"ls_option": "test"})
}
