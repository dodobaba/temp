package inventory

import (
	_ "shoppingzone/mylib/mylog" // :

	"github.com/gin-gonic/gin"
)

// Router : user router
func Router(rg *gin.RouterGroup) {
	rg.POST("/signupinventory" /*service.Authorization("admin"),*/, SignUpInventory)
}
