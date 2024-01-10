package controllers

import (
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zohari-tech/flowsim/utils"
)

// func GetMenu(shortcode string) (menu Menu) {

// 	ls := &ListScreen{
// 		CoreScreen: CoreScreen{
// 			Header: "Please select youre bank selection",
// 			Name:   "HomePage",
// 		},
// 		NextLocation: 1,
// 		Options:      []string{"KCB", "NCBA", "COOP"},
// 	}
// 	rs := &RouteScreen{
// 		CoreScreen: CoreScreen{
// 			Name:   "RoutedPage",
// 			Header: "Show routed selection",
// 		},
// 		Routes: []Route{{Value: "Selected", NextLocation: 2}},
// 	}
// 	rsi := &RawInputScreen{
// 		CoreScreen: CoreScreen{
// 			Header: "Please enter your name",
// 			Name:   "RawInputTest",
// 		},
// 		NextLocation: 3,
// 	}

// 	es := &ExternalScreen{
// 		CoreScreen: CoreScreen{
// 			Header: "Please enter your name",
// 			Name:   "RawInputTest",
// 		},
// 	}

// 	es.SetDisplay("1.YTUTY0\n2.KJ8789")
// 	es.SetLocation(-1, []Route{{Value: "Account number", NextLocation: 4}})
// 	menu = append(menu, ls, rs, rsi, es)
// 	return
// }

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
