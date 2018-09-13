package main

import (
	"shoppingzone/service"
	"shoppingzone/service/inventory"
	"shoppingzone/service/merchant"
	"shoppingzone/service/user"

	"github.com/gin-gonic/gin"
)

func myrouter(r *gin.Engine, startService string) {
	s := startService
	r.POST("/upload", service.Uploadfiles)
	r.GET("/img/:hashkey", service.LoadImage)
	if s == "" || s == "user" {
		user.Router(r.Group("/user"))
	}
	if s == "" || s == "merchant" {
		merchant.Router(r.Group("/merchant"))
	}
	if s == "" || s == "inventory" {
		inventory.Router(r.Group("/inventory"))
	}
}
