package etc

import (
	"bytes"
	"context"
	"encoding/json"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"github.com/supabase-community/supabase-go"
)

type Thumbnail struct {
	ID   int    `json:"id"`
	Link string `json:"link"`
}

const (
	fullScreenWidth  uint = 1080
	fullScreenHeight uint = 1920
)

var thumbnailHTTPClient = &http.Client{
	Timeout: 15 * time.Second,
}

func TumbnailView(db *supabase.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, _, err := db.
			From("tumb").
			Select("*", "", false).
			Execute()
		if err != nil {
			log.Println("gagal mengambil data tumb:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "gagal mengambil data thumbnail",
			})
			return
		}

		var res []Thumbnail
		if err := json.Unmarshal(data, &res); err != nil {
			log.Println("gagal parsing data tumb:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "gagal memproses data thumbnail",
			})
			return
		}

		if len(res) == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "data kosong",
			})
			return
		}

		rans := res[rand.Intn(len(res))]
		if strings.TrimSpace(rans.Link) == "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "link thumbnail tidak valid",
			})
			return
		}

		reqCtx, cancel := context.WithTimeout(ctx.Request.Context(), 15*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, rans.Link, nil)
		if err != nil {
			log.Println("gagal membuat request:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "gagal memproses link thumbnail",
			})
			return
		}
		req.Header.Set("User-Agent", "Mozilla/5.0")

		resp, err := thumbnailHTTPClient.Do(req)
		if err != nil {
			log.Println("gagal mengambil gambar:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "gagal mengambil gambar",
			})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Println("status tidak ok dari sumber gambar:", resp.StatusCode, rans.Link)
			ctx.JSON(http.StatusBadGateway, gin.H{
				"error": "gagal mengambil gambar dari sumber",
			})
			return
		}

		contentType := resp.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "image/") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "link bukan gambar",
			})
			return
		}

		srcImg, format, err := image.Decode(resp.Body)
		if err != nil {
			log.Println("gagal decode gambar:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "gagal memproses gambar",
			})
			return
		}

		resizedImg := resize.Thumbnail(fullScreenWidth, fullScreenHeight, srcImg, resize.Lanczos3)

		buf := new(bytes.Buffer)
		outContentType := contentType

		switch format {
		case "jpeg":
			err = jpeg.Encode(buf, resizedImg, &jpeg.Options{Quality: 90})
			outContentType = "image/jpeg"
		case "png":
			err = png.Encode(buf, resizedImg)
			outContentType = "image/png"
		case "gif":
			err = gif.Encode(buf, resizedImg, nil)
			outContentType = "image/gif"
		default:
			err = jpeg.Encode(buf, resizedImg, &jpeg.Options{Quality: 100})
			outContentType = "image/jpeg"
		}

		if err != nil {
			log.Println("gagal encode gambar hasil resize:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "gagal memproses gambar",
			})
			return
		}

		ctx.Data(http.StatusOK, outContentType, buf.Bytes())
	}
}
