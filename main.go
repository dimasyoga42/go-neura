package main

import (
	lib "neura-go/src/lib"
	"neura-go/src/middleware"
	"neura-go/src/routes"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

func main() {
	go middleware.CleanupVisitors()
	db := lib.Supabase()
	c := cache.New(5*time.Minute, 10*time.Minute)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.Setup(r, db, c)

	r.Run(":4000")
}
