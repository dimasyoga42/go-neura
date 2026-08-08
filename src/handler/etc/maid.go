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

func MaidGacha (db *supabase.Client) gin.HandlerFunc {
	return  func(ctx *gin.Context) {
		var res []config.Maid
		data, _, err := db.From("maid").Select("*", "", false).Execute();
		if err != nil {
			ctx.JSON(http.StatusNotFound, config.Error{
				Message: "terjadi kesalahan pada server",
				Status: false,
			})
		}
		if err := json.Unmarshal(data, &res); err != nil {
			ctx.JSON(http.StatusBadRequest, config.Error{
				Message: "gagal dalam parser data",
				Status:  false,
			})
		}

		if len(res) == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "data kosong",
			})
			return
		}

		rans := res[rand.Intn(len(res))]
		if strings.TrimSpace(rans.Image) == "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "link thumbnail tidak valid",
			})
			return
		}
		resp, err := http.Get(rans.Image)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "gagal mengambil gambar",
			})
			return
		}
		defer resp.Body.Close()

		ctx.DataFromReader(
			http.StatusOK,
			resp.ContentLength,
			resp.Header.Get("Content-Type"),
			resp.Body,
			nil,
		)
	}
}
