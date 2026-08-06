package routes

import (
	regis "neura-go/src/handler"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/supabase-community/supabase-go"
)

func Setup(r *gin.Engine, db *supabase.Client, c *cache.Cache) {
	r.GET("/home", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "ini adalah Home"})
	})

	r.GET("etc/toram/regis", regis.RegisHandler(db, c))
	r.GET("etc/toram/trait", regis.TraitHandler(db, c))
}
