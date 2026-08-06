package main

import (
	lib "neura-go/src/lib"
	Regis "neura-go/src/routes"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

func main() {
	db := lib.Supabase()
	c := cache.New(5*time.Minute, 10*time.Minute)

	r := gin.Default()
	r.GET("/regis", Regis.RegisHandler(db, c))
	r.Run(":4000")
}
