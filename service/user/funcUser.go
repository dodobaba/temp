package user

import (
	"shoppingzone/mylib/mylog"
	"shoppingzone/mylib/mymgo"
	sms "shoppingzone/mylib/mysms"
	"shoppingzone/myutil"
	"shoppingzone/service"
	"strconv"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/gin-gonic/gin"
)

var validata myutil.Validata

/*
SignUp :
@ 'Username'
@ 'Pwd'
@ 'Mobilecode'
@ 'Mobilenumber'
*/
func SignUp(c *gin.Context) {
	var user signupUser
	contentType, cip, err := service.CheckContentType(c, &user)
	mylog.SetIP(cip)
	if err != nil {
		mylog.Tf("[Error]", "User", "SignUp", "%s", err.Error())
		c.JSON(200, gin.H{"err": err.Error()})
	}
	go mylog.Tf("[Info]", "User", "SignUp", "Request : %+v", user)
	if len(user.Username) == 0 {
		user.Username = <-myutil.GetRandomString(10)
	}
	if _, err := validata.ValiString(user.Username, `^[a-zA-Z][^ ]*$`, 8, 16); err != nil {
		mylog.Tf("[Error]", "User", "SignUp", "validata username. %s  %s", user.Username, err.Error())
		c.JSON(200, gin.H{"err": err.Error()})
		return
	}
	if _, err := validata.ValiString(user.Mobilenumber, `^[0-9][^ a-zA-Z]*$`, 6, 16); err != nil {
		mylog.Tf("[Error]", "User", "SignUp", "validata mobile. %s  %s", user.Mobilenumber, err.Error())
		c.JSON(200, gin.H{"err": err.Error()})
		return
	}
	r := User{}
	b := bson.M{"$or": []bson.M{bson.M{"Username": user.Username}, bson.M{"Moblienumber": user.Mobilenumber}}}
	query := func(c *mgo.Collection) error {
		return c.Find(b).One(&r)
	}
	err = mymgo.Do("User", query)
	if err == nil {
		mylog.Tf("[Error]", "User", "SignUp", "Fail to signup with this name / mobilenumber,it always had.  %s", err.Error())
		c.JSON(200, gin.H{"err": "Fail to signup with this name / mobilenumber,it always had"})
		return
	}
	if err.Error() != "not found" {
		mylog.Tf("[Error]", "User", "SignUp", "Fail to find mongo query.  %s", err.Error())
		return
	}
	salt := <-myutil.GetRandomString(32)
	pwd := myutil.EnCrypto(user.Userpassword, salt)
	b = bson.M{"Username": user.Username, "Password": pwd, "Mobilecode": user.Mobilecode, "Mobilenumber": user.Mobilenumber, "Salt": salt, "Active": false, "CreateTime": time.Now()}
	query = func(c *mgo.Collection) error {
		return c.Insert(b)
	}
	err = mymgo.Do("User", query)
	if err != nil {
		mylog.Tf("[Error]", "User", "SignUp", "Fail register into db.  %s", err.Error())
		c.JSON(200, gin.H{"err": "Fail register into db"})
		return
	}
	go mylog.Tf("[Info]", "User", "SignUp", "Signup was okey. %+v", user)
	c.JSON(200, gin.H{"Message": "Signup was okey, please active this account with SMS code.",
		"ContentType": contentType,
		"UserName":    user.Username,
		"MobileCode":  user.Mobilecode,
		"Mobile":      user.Mobilenumber,
		"Status":      "Success",
		"UserStatus":  "InActive"})
}

/*
SendMobileVerification :
@ 'Mobilecode'
@ 'Mobilenumber'
*/
func SendMobileVerification(c *gin.Context) {
	var mv mobileVerification
	contentType, cip, err := service.CheckContentType(c, &mv)
	mylog.SetIP(cip)
	if err != nil {
		mylog.Tf("[Error]", "User", "SendMobileVerification", "%s", err.Error())
		c.JSON(200, gin.H{"err": err.Error()})
		return
	}
	go mylog.Tf("[Info]", "User", "SendMobileVerification", "Request : %+v", mv)
	if _, err := validata.ValiString(mv.Mobilenumber, `^[0-9][^ a-zA-Z]*$`, 6, 16); err != nil {
		mylog.Tf("[Error]", "User", "SendMobileVerification", "validata mobile. %s %s", mv.Mobilenumber, err.Error())
		c.JSON(200, gin.H{"err": err.Error()})
		return
	}
	rn := myutil.GetRandomRangeNumber(6)
	data := "[验证码]" + rn + ",请保管好验证码，不要向第三方或个人透露该验证码，请在5分钟内使用该验证码，过时作废，如需请再次发送，谢谢，祝您愉快！"
	go sms.SendSMS(mv.Mobilenumber, rn, data)
	if err != nil {
		mylog.Tf("[Error]", "User", "SendMobileVerification", "%s", err.Error())
		c.JSON(200, gin.H{"err": err.Error()})
		return
	}
	b := bson.M{"Mobile": mv.Mobilenumber, "VerificationCode": rn, "CreateTime": time.Now()}
	query := func(c *mgo.Collection) error {
		idx := mgo.Index{
			Key:         []string{"CreateTime"},
			ExpireAfter: time.Duration(time.Second * 300),
			Name:        "timer",
		}
		err := c.EnsureIndex(idx)
		if err != nil {
			mylog.Tf("[Error]", "User", "SendMobileVerification", "mgo index create error.  %s", err.Error())
		}
		return c.Insert(b)
	}
	err = mymgo.Do("VerificationCode", query)
	if err != nil {
		mylog.Tf("[Error]", "User", "SendMobileVerification", "Fail worten verification code into db. %s", err.Error())
		c.JSON(200, gin.H{"err": "Fail worten verification code into db"})
		return
	}
	go mylog.Tf("[Info]", "User", "SendMobileVerification", "Send mobile number was okey, please check verification code. %+v", mv)
	c.JSON(200, gin.H{"Message": "Send mobile number was okey, please check verification code.",
		"ContentType": contentType,
		"MobileCode":  mv.Mobilecode,
		"Mobile":      mv.Mobilenumber,
		"Status":      "Success",
		"SendStatus":  "Active"})
}

/*
VerificationMobileCode :
@ 'Mobilecode'
@ 'Mobilenumber'
@ 'VerificationCode'
*/
func VerificationMobileCode(c *gin.Context) {
	var vc verificationCode
	contentType, cip, err := service.CheckContentType(c, &vc)
	mylog.SetIP(cip)
	if err != nil {
		mylog.Tf("[Error]", "User", "VerificationMobileCode", "%s", err.Error())
		c.JSON(200, gin.H{"err": err.Error()})
		return
	}
	go mylog.Tf("[Info]", "User", "VerificationMobileCode", "Request : %+v ", vc)
	if _, err := validata.ValiString(vc.Mobilenumber, `^[0-9][^ a-zA-Z]*$`, 6, 16); err != nil {
		mylog.Tf("[Error]", "User", "VerificationMobileCode", "validata mobile. %s %s", vc.Mobilenumber, err.Error())
		c.JSON(200, gin.H{"err": err.Error()})
		return
	}
	if _, err := validata.ValiString(vc.VerificationCode, `^[0-9][^ a-zA-Z]*$`, 6, 6); err != nil {
		mylog.Tf("[Error]", "User", "VerificationMobileCode", "validata VCode. %s %s", vc.Mobilenumber, err.Error())
		c.JSON(200, gin.H{"err": err.Error()})
		return
	}
	if !sms.VeriticationSMSCode(vc.Mobilenumber, vc.VerificationCode) {
		mylog.Tf("[Error]", "User", "VerificationMobileCode", "Fail read verification code from db")
		c.JSON(200, gin.H{"Message": "Fail read verification code from db",
			"ContentType": contentType,
			"Status":      "Fail"})
		return
	}
	selector := bson.M{"Mobilenumber": vc.Mobilenumber, "Mobilecode": vc.Mobilecode}
	d := bson.M{"$set": bson.M{"Active": true, "ActiveTime": time.Now(), "UserKey": <-myutil.GetRandomString(64), "SecureRole": []string{"user"}}}
	query := func(c *mgo.Collection) error {
		return c.Update(selector, d)
	}
	err = mymgo.Do("User", query)
	if err != nil {
		mylog.Tf("[Error]", "User", "VerificationMobileCode", "%s", err.Error())
		c.JSON(200, gin.H{"err": err.Error()})
		return
	}
	go mylog.Tf("[Info]", "User", "VerificationMobileCode", "check the mobile verification code was okey, It's currect. %+v", vc)
	c.JSON(200, gin.H{"Message": "check the mobile verification code was okey, It's currect.",
		"ContentType": contentType,
		"Mobile":      vc.Mobilenumber,
		"Mobilecode":  vc.Mobilecode,
		"Status":      "Success"})
}

/*
Login :
@ 'Signinstring'
@ 'Signinpassword'
*/
func Login(c *gin.Context) {
	var lu loginUser
	contentType, cip, err := service.CheckContentType(c, &lu)
	mylog.SetIP(cip)
	if err != nil {
		mylog.Tf("[Error]", "User", "Login", "%s", err.Error())
		c.JSON(200, gin.H{"err": err.Error()})
		return
	}
	go mylog.Tf("[Info]", "User", "Login", "Request : %+v ", lu)
	if _, err := validata.ValiString(lu.Signinstring, `^[a-zA-Z0-9][^ ]*$`, 8, 16); err != nil {
		mylog.Tf("[Error]", "User", "Login", "validata signinstring. %s %s", lu.Signinstring, err.Error())
		c.JSON(200, gin.H{"err": err.Error(), "filed": 1})
		return
	}
	if _, err := validata.ValiString(lu.Signinpassword, `^[a-zA-Z0-9][^ ]*$`, 6, 255); err != nil {
		mylog.Tf("[Error]", "User", "Login", "validata signinpassword. %s %s", lu.Signinpassword, err.Error())
		c.JSON(200, gin.H{"err": err.Error(), "filed": 2})
		return
	}
	r := User{}
	if sms.VeriticationSMSCode(lu.Signinstring, lu.Signinpassword) {
		b := bson.M{"Moblienumber": lu.Signinstring, "Active": true}
		query := func(c *mgo.Collection) error {
			return c.Find(b).One(&r)
		}
		err = mymgo.Do("User", query)
		if err != nil {
			mylog.Tf("[Error]", "User", "Login", "This mobilenumber not signup,please register. %s", err.Error())
			c.JSON(200, gin.H{"Message": "This mobilenumber not signup,please register",
				"ContentType": contentType,
				"err":         err.Error(),
				"Status":      "Fail"})
			return
		}
		cacheData := myutil.TakeHashWord(r.UserKey)
		token := <-myutil.EnCryptoToken(r.UserKey)
		cached := <-service.Cache(r.UserKey, cacheData, 3600*72)
		mylog.Tf("[Info]", "User", "Login", "used mobile and veritication code to login. %+v %s %+v", lu, token, r)
		c.JSON(200, gin.H{"Message": "used mobile and veritication code to login.",
			"ContentType": contentType,
			"Status":      "Success",
			"Cache":       cached,
			"Token":       token,
			"User":        r})
		return
	}
	query := func(c *mgo.Collection) error {
		b := bson.M{"Active": true, "$or": []bson.M{bson.M{"Username": lu.Signinstring}, bson.M{"Mobilenumber": lu.Signinstring}}}
		return c.Find(b).One(&r)
	}
	err = mymgo.Do("User", query)
	if err != nil {
		mylog.Tf("[Error]", "User", "Login", "Fail to Login with this name / mobilenumber,please check. %s", err.Error())
		c.JSON(200, gin.H{"Message": "Fail to Login with this name / mobilenumber,please check",
			"ContentType": contentType,
			"err":         err.Error(),
			"Status":      "Fail"})
		return
	}
	salt := r.Salt
	if myutil.EnCrypto(lu.Signinpassword, salt) != r.Password {
		mylog.Tf("[Error]", "User", "Login", "Fail to Login with this password,please check. %s", err.Error())
		c.JSON(200, gin.H{"Message": "Fail to Login with this password,please check",
			"ContentType": contentType,
			"Status":      "Fail"})
		return
	}
	cacheData := myutil.TakeHashWord(r.UserKey)
	token := <-myutil.EnCryptoToken(r.UserKey)
	cached := <-service.Cache(r.UserKey, cacheData, 3600*72)
	go mylog.Tf("[Info]", "User", "Login", "It's currect. Login by this user. %+v %s %+v", lu, token, r)
	c.JSON(200, gin.H{"Message": "It's currect. Login by this user",
		"ContentType": contentType,
		"Status":      "Success",
		"Cache":       cached,
		"Token":       token,
		"User":        r})
}

/*
SetSecure :
@UserKey # (2 option 1)
@UserMobile # (2 option 1)
@SecureField #set secure fields arrary just 4 options ["SecureLabel" "SecureGroup" "SecureRouter" "SecureRole"]
@SecureRole #set secure role arrary
*/
func SetSecure(c *gin.Context) {
	var us userSecure
	contentType, cip, err := service.CheckContentType(c, &us)
	mylog.SetIP(cip)
	if err != nil {
		mylog.Tf("[Error]", "User", "SetSecure", "%s", err.Error())
		c.JSON(200, gin.H{"err": err.Error()})
		return
	}
	go mylog.Tf("[Info]", "User", "SetSecure", "Request : %+v", us)
	if len(us.SecureField) != len(us.SecureRole) {
		mylog.Tf("[Error]", "User", "SetSecure", "Fail to set user to these authraztaions,please check.")
		c.JSON(200, gin.H{"Message": "Fail to set user to these authraztaions,please check",
			"ContentType": contentType,
			"Status":      "Fail"})
		return
	}
	if len(us.UserKey) == 0 && len(us.UserMobile) == 0 {
		mylog.Tf("[Error]", "User", "SetSecure", "Fail to set userkey or user mobile ,it's necessary.")
		c.JSON(200, gin.H{"Message": "Fail to set userkey or user mobile ,it's necessary",
			"ContentType": contentType,
			"Status":      "Fail"})
		return
	}
	if us.UserMobile != "" {
		if _, err := validata.ValiString(us.UserMobile, `^[0-9][^ a-zA-Z]*$`, 6, 16); err != nil {
			mylog.Tf("[Error]", "User", "SetSecure", "validata mobile. %s %s", us.UserMobile, err.Error())
			c.JSON(200, gin.H{"err": err.Error()})
			return
		}
	}
	if us.UserKey != "" {
		if _, err := validata.ValiString(us.UserKey, `^[a-zA-Z0-9]*$`, 32, 32); err != nil {
			mylog.Tf("[Error]", "User", "SetSecure", "validata userkey. %s %s", us.UserKey, err.Error())
			c.JSON(200, gin.H{"err": err.Error()})
			return
		}
	}
	var sl, sg, sr, srl []string
	for idx, v := range us.SecureField {
		switch v {
		case "SecureLabel":
			sl = append(sl, us.SecureRole[idx])
			break
		case "SecureGroup":
			sg = append(sg, us.SecureRole[idx])
			break
		case "SecureRouter":
			sr = append(sr, us.SecureRole[idx])
			break
		case "SecureRole":
			srl = append(srl, us.SecureRole[idx])
			break
		default:
			mylog.Tf("[Error]", "User", "SetSecure", "cann't set "+v+". havn't this role. %s", err.Error())
			break
		}
	}
	query := func(c *mgo.Collection) error {
		selector := bson.M{"$or": []bson.M{bson.M{"Mobilenumber": us.UserMobile}, bson.M{"UserKey": us.UserKey}}}
		b := bson.M{"$addToSet": bson.M{"SecureLabel": bson.M{"$each": sl}, "SecureGroup": bson.M{"$each": sg}, "SecureRouter": bson.M{"$each": sr}, "SecureRole": bson.M{"$each": srl}}}
		return c.Update(selector, b)
	}
	err = mymgo.Do("User", query)
	if err != nil {
		mylog.Tf("[Error]", "User", "SetSecure", "Fail to set this user as some anthraztion. %s", err.Error())
		c.JSON(200, gin.H{"Message": "Fail to set this user as some anthraztion",
			"ContentType": contentType,
			"err":         err.Error(),
			"Status":      "Fail"})
		return
	}
	go mylog.Tf("[Info]", "User", "SetSecure", "It's currect set anthrazion to this user. %+v", us)
	c.JSON(200, gin.H{"Message": "It's currect set anthrazion to this user",
		"ContentType": contentType,
		"Status":      "Success"})
}

/*
DropSecure :
@UserKey # (2 option 1)
@UserMobile # (2 option 1)
@SecureField #drop secure fields arrary just 4 options ["SecureLabel" "SecureGroup" "SecureRouter" "SecureRole"]
@SecureRole #drop secure role arrary
*/
func DropSecure(c *gin.Context) {
	var us userSecure
	contentType, cip, err := service.CheckContentType(c, &us)
	mylog.SetIP(cip)
	if err != nil {
		mylog.Tf("[Error]", "User", "DropSecure", "%s", err.Error())
		c.JSON(200, gin.H{"err": err.Error()})
		return
	}
	go mylog.Tf("[Info]", "User", "DropSecure", "Request : %+v", us)
	if len(us.SecureField) != len(us.SecureRole) {
		mylog.Tf("[Error]", "User", "DropSecure", "Fail to set user to these authraztaions,please check")
		c.JSON(200, gin.H{"Message": "Fail to set user to these authraztaions,please check",
			"ContentType": contentType,
			"Status":      "Fail"})
		return
	}
	if len(us.UserKey) == 0 && len(us.UserMobile) == 0 {
		mylog.Tf("[Error]", "User", "DropSecure", "Fail to set userkey or user mobile ,it's necessary")
		c.JSON(200, gin.H{"Message": "Fail to set userkey or user mobile ,it's necessary",
			"ContentType": contentType,
			"Status":      "Fail"})
		return
	}
	if us.UserMobile != "" {
		if _, err := validata.ValiString(us.UserMobile, `^[0-9][^ a-zA-Z]*$`, 6, 16); err != nil {
			mylog.Tf("[Error]", "User", "DropSecure", "validata mobile. %s %s", us.UserMobile, err.Error())
			c.JSON(200, gin.H{"err": err.Error()})
			return
		}
	}
	if us.UserKey != "" {
		if _, err := validata.ValiString(us.UserKey, `^[a-zA-Z0-9]*$`, 32, 32); err != nil {
			mylog.Tf("[Error]", "User", "DropSecure", "validata userkey. %s %s", us.UserKey, err.Error())
			c.JSON(200, gin.H{"err": err.Error()})
			return
		}
	}
	var sl, sg, sr, srl []string
	for idx, v := range us.SecureField {
		switch v {
		case "SecureLabel":
			sl = append(sl, us.SecureRole[idx])
			break
		case "SecureGroup":
			sg = append(sg, us.SecureRole[idx])
			break
		case "SecureRouter":
			sr = append(sr, us.SecureRole[idx])
			break
		case "SecureRole":
			srl = append(srl, us.SecureRole[idx])
			break
		default:
			mylog.Tf("[Error]", "User", "DropSecure", "cann't set "+v+". havn't this role. %s", err.Error())
			break
		}
	}
	query := func(c *mgo.Collection) error {
		selector := bson.M{"$or": []bson.M{bson.M{"Mobilenumber": us.UserMobile}, bson.M{"UserKey": us.UserKey}}}
		b := bson.M{"$pull": bson.M{"SecureLabel": bson.M{"$in": sl}, "SecureGroup": bson.M{"$in": sg}, "SecureRouter": bson.M{"$in": sr}, "SecureRole": bson.M{"$in": srl}}}
		return c.Update(selector, b)
	}
	err = mymgo.Do("User", query)
	if err != nil {
		mylog.Tf("[Error]", "User", "DropSecure", "Fail to drop some anthraztion from this user. %s", err.Error())
		c.JSON(200, gin.H{"Message": "Fail to drop some anthraztion from this user",
			"ContentType": contentType,
			"err":         err.Error(),
			"Status":      "Fail"})
		return
	}
	go mylog.Tf("[Info]", "User", "DropSecure", "It's currect drop anthrazion from this user. %+v", us)
	c.JSON(200, gin.H{"Message": "It's currect drop anthrazion from this user",
		"ContentType": contentType,
		"Status":      "Success"})
}

/*
ChangeStatus :
@UserKey # (2 option 1)
@UserMobile # (2 option 1)
@Status #ture or false
*/
func ChangeStatus(c *gin.Context) {
	var cs changeStatus
	s := false
	contentType, cip, err := service.CheckContentType(c, &cs)
	mylog.SetIP(cip)
	if err != nil {
		mylog.Tf("[Error]", "User", "ChangeStatus", "%s", err.Error())
		c.JSON(200, gin.H{"err": err.Error()})
		return
	}
	go mylog.Tf("[Info]", "User", "ChangeStatus", "Request : %+v", cs)
	if len(cs.UserKey) == 0 && len(cs.UserMobile) == 0 {
		mylog.Tf("[Error]", "User", "ChangeStatus", "Fail to set userkey or user mobile ,it's necessary")
		c.JSON(200, gin.H{"Message": "Fail to set userkey or user mobile ,it's necessary",
			"ContentType": contentType,
			"Status":      "Fail"})
		return
	}
	if cs.UserMobile != "" {
		if _, err := validata.ValiString(cs.UserMobile, `^[0-9][^ a-zA-Z]*$`, 6, 16); err != nil {
			mylog.Tf("[Error]", "User", "ChangeStatus", "validata mobile. %s %s", cs.UserMobile, err.Error())
			c.JSON(200, gin.H{"err": err.Error()})
			return
		}
	}
	if cs.UserKey != "" {
		if _, err := validata.ValiString(cs.UserKey, `^[a-zA-Z0-9]*$`, 32, 32); err != nil {
			mylog.Tf("[Error]", "User", "ChangeStatus", "validata userkey. %s %s", cs.UserKey, err.Error())
			c.JSON(200, gin.H{"err": err.Error()})
			return
		}
	}
	if cs.Status {
		s = true
	}
	query := func(c *mgo.Collection) error {
		selector := bson.M{"$or": []bson.M{bson.M{"Mobilenumber": cs.UserMobile}, bson.M{"UserKey": cs.UserKey}}}
		b := bson.M{"$set": bson.M{"Active": s}}
		return c.Update(selector, b)
	}
	err = mymgo.Do("User", query)
	if err != nil {
		mylog.Tf("[Error]", "User", "ChangeStatus", "Fail to set this user active status. %s", err.Error())
		c.JSON(200, gin.H{"Message": "Fail to set this user active status",
			"ContentType": contentType,
			"err":         err.Error(),
			"Status":      "Fail"})
		return
	}
	go mylog.Tf("[Info]", "User", "ChangeStatus", "It's currect to set active status with this user. %+v", cs)
	c.JSON(200, gin.H{"Message": "It's currect to set active status with this user",
		"ContentType": contentType,
		"Status":      "Success"})
}

//List :
func List(c *gin.Context) {
	contentType := c.Request.Header.Get("Content-Type")
	ip := c.ClientIP()
	mylog.SetIP(ip)
	p, _ := strconv.Atoi(c.Query("page"))
	s, _ := strconv.Atoi(c.Query("pagesize"))
	type listuser struct {
		UserKey     string        `bson:"UserKey" form:"UserKey" json:"UserKey"`
		DisplayName string        `bson:"DisplayName" form:"DisplayName" json:"DisplayName"`
		ID          bson.ObjectId `bson:"_id" form:"_id" json:"_id"`
		Active      bool          `bson:"Active" form:"Active" json:"Active"`
		CreateTime  time.Time     `bson:"CreateTime" form:"CreateTime" json:"CreateTime"`
	}
	var users []listuser
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{}).Sort("CreateTime").Skip(p * s).Limit(s).All(&users)
	}
	if err := mymgo.Do("User", query); err != nil {
		mylog.Tf("[Error]", "User", "List", "Fail to list users. %s", err.Error())
		c.JSON(200, gin.H{"Message": "Fail to list users",
			"ContentType": contentType,
			"err":         err.Error(),
			"Status":      "Fail"})
		return
	}
	go mylog.Tf("[Info]", "User", "List", "List users currect.")
	c.JSON(200, gin.H{"Message": "List users currect.",
		"ContentType": contentType,
		"Status":      "Success",
		"Users":       users,
	})
}
