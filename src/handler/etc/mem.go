package etc

import (
	"image"
	_ "image/jpeg"
	"image/png"
	"io"
	"net/http"
	"strconv"

	"neura-go/src/helper"

	"github.com/gin-gonic/gin"
	"path/filepath"
  "runtime"
)
var font = func() string {
	_, thisFile, _, _ := runtime.Caller(0)
	// thisFile = .../neura-go/src/handler/etc/mem.go
	root := filepath.Join(filepath.Dir(thisFile), "..", "..", "..")
	return filepath.Join(root, "Geologica_Auto-ExtraBold.ttf")
}()
func CreateMem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		url := ctx.Query("url")
		top := ctx.Query("top")
		bottom := ctx.Query("bottom")

		if url == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "url is required"})
			return
		}

		size, err := strconv.ParseFloat(ctx.Query("size"), 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid size parameter"})
			return
		}

		req, err := http.NewRequestWithContext(ctx.Request.Context(), http.MethodGet, url, nil)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid url"})
			return
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"error": "failed to fetch image"})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			ctx.JSON(http.StatusBadGateway, gin.H{"error": "image url returned non-200 status"})
			return
		}

		img, _, err := image.Decode(io.LimitReader(resp.Body, 10<<20))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to decode image"})
			return
		}

		result, err := helper.DrawText(img, font, size, top, bottom)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}


		ctx.Writer.Header().Set("Content-Type", "image/png")
		ctx.Status(http.StatusOK)
		if err := png.Encode(ctx.Writer, result); err != nil {
			// header/status sudah terkirim, jadi cuma bisa log di sini
			return
		}
	}
}
