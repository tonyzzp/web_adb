package routes

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"web_adb/adb"

	"github.com/gin-gonic/gin"
	goadb "github.com/zach-klippenstein/goadb"
)

type fileReq struct {
	Serial string `uri:"serial" binding:"required"`
	Path   string `uri:"path" binding:"required"`
}

func serveFile(ctx *gin.Context, device *goadb.Device, path string) {
	entry, e := device.Stat(path)
	if e != nil {
		ctx.String(http.StatusInternalServerError, "Stat error: "+e.Error())
		return
	}
	r, e := device.OpenRead(path)
	if e != nil {
		ctx.String(http.StatusInternalServerError, "OpenRead error: "+e.Error())
		return
	}
	var ext = filepath.Ext(path)
	var t = mime.TypeByExtension(ext)
	if t == "" {
		t = "application/octstream"
	}
	ctx.Header("Content-Length", strconv.Itoa(int(entry.Size)))
	ctx.Header("Content-Type", t)
	ctx.Stream(func(w io.Writer) bool {
		io.Copy(w, r)
		return false
	})
}

func Files(server *gin.Engine) {
	server.GET("/device/:serial/files/*path", func(ctx *gin.Context) {
		var fileReq = &fileReq{}
		if ctx.ShouldBindUri(fileReq) != nil {
			ctx.String(http.StatusBadRequest, "需要serial/path")
			return
		}
		var device = adb.DeviceBySerial(fileReq.Serial)
		entry, e := device.Stat(fileReq.Path)
		if e != nil {
			ctx.String(http.StatusInternalServerError, "Stat error: "+e.Error())
			return
		}
		if entry.Mode.IsRegular() {
			serveFile(ctx, device, fileReq.Path)
			return
		}
		entries, e := device.ListDirEntries(fileReq.Path)
		if e != nil {
			ctx.String(http.StatusInternalServerError, "ListDirEntries error: "+e.Error())
			return
		}
		names, e := entries.ReadAll()
		if e != nil {
			ctx.String(http.StatusInternalServerError, "ListDirEntries error: "+e.Error())
			return
		}
		for _, entry := range names {
			var href = entry.Name
			var text = entry.Name
			var target = "_blank"
			if entry.Mode.IsDir() || entry.Mode&os.ModeSymlink > 0 {
				text = text + "/"
				href = href + "/"
				target = "_self"
			}
			ctx.Writer.WriteString(fmt.Sprintf("<a href='%s' target='%s'>%s</a><br/>\n", href, target, text))
		}
	})
}
