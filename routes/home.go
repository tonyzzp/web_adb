package routes

import (
	"fmt"
	"net/http"
	"web_adb/adb"

	"github.com/gin-gonic/gin"
)

func Home(server *gin.Engine) {
	server.GET("/", func(ctx *gin.Context) {
		devices, e := adb.Adb.ListDevices()
		if e != nil {
			ctx.String(http.StatusOK, "ListDevices error:%s", e.Error())
			return
		}
		for _, device := range devices {
			var href = fmt.Sprintf("/device/%s/files", device.Serial)
			var text = fmt.Sprintf("%-20s: %s", device.Model, device.Serial)
			var s = fmt.Sprintf("<a href='%s'>%s</a>\n", href, text)
			ctx.Writer.WriteString(s)
		}
	})
}
