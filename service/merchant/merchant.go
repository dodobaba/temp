package merchant

import (
	constants "shoppingzone/mylib/myconst"
	_ "shoppingzone/mylib/mylog" // :
	"shoppingzone/service"

	"github.com/gin-gonic/gin"
)

// Router : user router
func Router(rg *gin.RouterGroup) {
	rg.POST("/signupmerchant" /*service.Authorization("admin"),*/, SignUpMerchant)
	rg.POST("/setupmerchantimage" /*service.Authorization("admin"),*/, SetupMerchantImage)
	rg.POST("/setupmerchantmanager" /*service.Authorization("admin"),*/, SetupMerchantManager)
	rg.POST("/changestatus", service.Authorization(constants.MERCHANTADMIN), ChangeStatus)
	rg.GET("/list", service.Authorization(constants.MERCHANTADMIN), List)
}
