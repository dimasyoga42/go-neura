package neru

import (
	"encoding/json"
	"net/http"
	"neura-go/src/config"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/supabase-community/supabase-go"
)

func XtalHandler(db *supabase.Client, c *cache.Cache) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		name := ctx.Query("name")
		cacheKey := "xtal:" + name

		if cached, found := c.Get(cacheKey); found {
			elapsed := time.Since(start)
			ctx.JSON(http.StatusOK, config.ApiRespon{
				Data: cached,
				Time: float64(elapsed.Microseconds()) / 1000,
			})
			return
		}

		data, _, err := db.
			From("xtal").
			Select("*", "", false).
			Ilike("name", "%"+name+"%").
			Execute()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var res []config.Xtal
		if err := json.Unmarshal(data, &res); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Set(cacheKey, res, cache.DefaultExpiration)

		elapsed := time.Since(start)
		ctx.JSON(http.StatusOK, config.ApiRespon{
			Data: res,
			Time: float64(elapsed.Microseconds()) / 1000,
		})
	}
}
