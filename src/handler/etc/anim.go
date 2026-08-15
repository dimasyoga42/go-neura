package etc

import (
	"fmt"
	"net/http"
	"neura-go/src/helper"

	"github.com/gin-gonic/gin"
)


func Anime() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, err := helper.GetanimeHelper()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	}
}

func AnimeDetail() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		v := ctx.Query("path")
		data, err := helper.GetAnimeDetailHelper(v)
		fmt.Println(data)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	}
}
func AnimeEpisode() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		v := ctx.Query("path")
		data, err := helper.GetEpisodeHelper(v)
		fmt.Println(data)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	}
}
