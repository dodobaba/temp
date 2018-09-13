package myinnercommand

import (
	"shoppingzone/mylib/mydb"
)

//InnerUser : all user func in inner call in here
type InnerUser struct{}

var db mydb.MyDB

func init() {
	db.DBConn()
	//defer db.Close()
}

//CheckUserNameIsNull : check the user name in the database
func CheckUserNameIsNull(name string) <-chan bool {
	out := make(chan bool)
	go func() {
		//db.DBConn()
		//defer db.Close()
		l := `SELECT UserId FROM User WHERE UserName = ? LIMIT 1`
		rs := db.SQLQuery(l, name)
		if rs.Next() {
			out <- false // haven't user can be register
		} else {
			out <- true // have user as same name
		}
		close(out)
	}()
	return out
}

//CheckUserMobileIsNull : chack the user mobile number
func CheckUserMobileIsNull(mobile string) <-chan bool {
	out := make(chan bool)
	go func() {
		//db.DBConn()
		//defer db.Close()
		l := `SELECT UserId FROM User WHERE MobileNumber = ? LIMIT 1`
		rs := db.SQLQuery(l, mobile)
		if rs.Next() {
			out <- false // haven't user can be register
		} else {
			out <- true // have user as same name
		}
		close(out)
	}()
	return out
}
