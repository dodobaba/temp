package user

import (
	constants "shoppingzone/mylib/myconst"
	_ "shoppingzone/mylib/mylog" // :
	"shoppingzone/service"

	"github.com/gin-gonic/gin"
)

/*
//CreateUser :
func CreateUser() {
	t := User{
		StatusID:    1,
		DateOfBirth: time.Now().Format("2006-01-02"),
	}
	log.Println(t)
}
*/

// Router : user router
func Router(rg *gin.RouterGroup) {
	rg.POST("/signup", SignUp)
	rg.POST("/sendmobileverification", SendMobileVerification)
	rg.POST("/verificationmobilecode", VerificationMobileCode)
	rg.POST("/login", Login)
	rg.POST("/setsecure", service.Authorization(constants.ADMIN), SetSecure)
	rg.POST("/dropsecure", service.Authorization(constants.ADMIN), DropSecure)
	rg.POST("/changeStatus", service.Authorization(constants.ADMIN), ChangeStatus)
	rg.GET("/list", service.Authorization(constants.ADMIN), List)
}
