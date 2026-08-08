package etc

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"neura-go/src/config"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/supabase-community/supabase-go"
)

func WaifuGacha(db *supabase.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res, _, err := db.From("waifu").Select("*", "", false).Execute()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, config.Error{
				Message: "Gagal mengambil data dari database",
				Status:  false,
			})
			return
		}

		var data []config.Waifu

		if err := json.Unmarshal(res, &data); err != nil {
			ctx.JSON(http.StatusInternalServerError, config.Error{
				Message: "Gagal parser data dari database",
				Status:  false,
			})
			return
		}

		if len(data) == 0 {
			ctx.JSON(http.StatusNotFound, config.Error{
				Message: "Tidak ada data dalam database",
				Status:  false,
			})
			return
		}

		random := data[rand.Intn(len(data))]

		link := strings.TrimSpace(random.Link)
		if link == "" {
			ctx.JSON(http.StatusBadRequest, config.Error{
				Message: "Link tidak valid",
				Status:  false,
			})
			return
		}

		resp, err := http.Get(link)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, config.Error{
				Message: "Gagal mengambil gambar",
				Status:  false,
			})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
			ctx.JSON(http.StatusBadGateway, config.Error{
				Message: "Gagal mengambil gambar dari sumber",
				Status:  false,
			})
			return
		}

		ctx.DataFromReader(
			http.StatusOK,
			resp.ContentLength,
			resp.Header.Get("Content-Type"),
			resp.Body,
			nil,
		)
	}
}
