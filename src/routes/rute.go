package routes

import (
	"neura-go/src/handler/etc"
	neru "neura-go/src/handler/toram"
	"neura-go/src/helper"
	"neura-go/src/middleware"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/supabase-community/supabase-go"
)

func Setup(r *gin.Engine, db *supabase.Client, c *cache.Cache) {
	r.Use(middleware.RateLimiter())
	r.GET("/home", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "ini adalah Home"})
	})

	r.GET("etc/toram/regis", neru.RegisHandler(db, c))
	r.GET("etc/toram/trait", neru.TraitHandler(db, c))
	r.GET("etc/toram/xtal", neru.XtalHandler(db, c))
	r.GET("etc/thumbnail", etc.TumbnailView(db))
	r.GET("etc/mem", etc.CreateMem())
	r.GET("etc/dye", helper.DyeView())
	r.GET("etc/waifu", etc.WaifuGacha(db))
}
