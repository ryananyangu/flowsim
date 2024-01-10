package main

import (
	"github.com/gin-gonic/gin"
)

// SetupRouter specifiy routes
func SetupRouter() *gin.Engine {
	// gin initialization and middleware configuration
	r := gin.Default()
	r.Use(gin.Recovery())

	for path, handlers := range Routes {
		for method, handler := range handlers {
			switch method {
			case "GET":
				r.GET(path, handler)

			case "POST":
				r.POST(path, handler)

			case "PUT":
				r.PUT(path, handler)

			case "PATCH":
				r.PATCH(path, handler)

			case "DELETE":
				r.DELETE(path, handler)
			}
		}
	}
	// r.Static("/assets", "./assets")

	r.LoadHTMLGlob("templates/*")
	return r
}

func main() {

	r := SetupRouter()
	r.Run("0.0.0.0:8080")

}
