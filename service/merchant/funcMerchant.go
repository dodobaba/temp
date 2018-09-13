package merchant

import (
	"shoppingzone/mylib/mylog"
	"shoppingzone/mylib/mymgo"
	"shoppingzone/myutil"
	"shoppingzone/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var validata myutil.Validata

/*
SignUpMerchant :
@MerchantCompanyName 1-255
@MerchantNameInvariant 1-255
@MerchantDescInvariant 1-1024
@MerchantCode 1-255
@MerchantType 1-255
@BusinessRegistrationNo 1-255
@District 1-255
@Apartment 1-255
@Floor 1-255
@BlockNo 1-255
@Building 1-255
@Street 1-255
@StreetNo 1-255
@IsCrossBorder bool
*/
func SignUpMerchant(c *gin.Context) {
	var m Merchant
	contentType, cip, err := service.CheckContentType(c, &m)
	mylog.SetIP(cip)
	if err != nil {
		mylog.Tf("[Error]", "Merchant", "SignUpMerchant", "%s", err.Error())
		c.JSON(200, gin.H{"err": err.Error()})
		return
	}
	go mylog.Tf("[Info]", "Merchant", "SignUpMerchant", "Request : %+v", m)
	rxp := `^[0-9a-zA-Z\p{Han}][^\t\v\\\^\$\*\+\?\{\}\(\)\[\]\|]*$`
	v := []myutil.ValiStringType{
		{S: m.MerchantCompanyName, Rxp: rxp, Min: 1, Max: 255},
		{S: m.MerchantNameInvariant, Rxp: rxp, Min: 1, Max: 255},
		{S: m.MerchantDescInvariant, Rxp: rxp, Min: 1, Max: 1024},
		{S: m.MerchantCode, Rxp: rxp, Min: 1, Max: 8},
		{S: m.MerchantType, Rxp: rxp, Min: 1, Max: 255},
		{S: m.BusinessRegistrationNo, Rxp: rxp, Min: 1, Max: 255},
		{S: m.District, Rxp: rxp, Min: 1, Max: 255},
		{S: m.Apartment, Rxp: rxp, Min: 1, Max: 255},
		{S: m.Floor, Rxp: rxp, Min: 1, Max: 255},
		{S: m.BlockNo, Rxp: rxp, Min: 1, Max: 255},
		{S: m.Building, Rxp: rxp, Min: 1, Max: 255},
		{S: m.Street, Rxp: rxp, Min: 1, Max: 255},
		{S: m.StreetNo, Rxp: rxp, Min: 1, Max: 255},
	}
	e := validata.ValiStrings(v...)
	for er := range e {
		if er != nil {
			mylog.Tf("[Error]", "Merchant", "SignUpMerchant", "validata fail. %s", er.Error())
			c.JSON(200, gin.H{"err": er.Error()})
			return
		}
	}
	b := bson.M{
		"MerchantCompanyName":    m.MerchantCompanyName,
		"MerchantNameInvariant":  m.MerchantDescInvariant,
		"MerchantDescInvariant":  m.MerchantDescInvariant,
		"MerchantCode":           m.MerchantCode,
		"MerchantType":           m.MerchantType,
		"BusinessRegistrationNo": m.BusinessRegistrationNo,
		"District":               m.District,
		"Apartment":              m.Apartment,
		"Floor":                  m.Floor,
		"BlockNo":                m.BlockNo,
		"Building":               m.Building,
		"Street":                 m.Street,
		"StreetNo":               m.StreetNo,
		"IsCrossBorder":          false,
		"Active":                 false,
		"CreateTime":             time.Now(),
		"ActiveTime":             time.Now(),
		"MerchantKey":            <-myutil.GetRandomString(64),
	}
	query := func(c *mgo.Collection) error {
		return c.Insert(b)
	}
	err = mymgo.Do("Merchant", query)
	if err != nil {
		mylog.Tf("[Error]", "Merchant", "SignUpMerchant", "Fail to inster this Merchant into db. %s | %s", m.MerchantCompanyName, err.Error())
		c.JSON(200, gin.H{"Message": "Fail to inster this Merchant into db.", "ContentType": contentType, "Status": "Fail", "err": err.Error()})
		return
	}
	go mylog.Tf("[Info]", "Merchant", "SignUpMerchant", "Signup was okey. %+v", b)
	c.JSON(200, gin.H{"Message": "Signup was okey.",
		"ContentType":    contentType,
		"Status":         "Success",
		"MerchantStatus": "Active",
		"MerhantInfo":    b,
	})
}

/*
SetupMerchantImage :
@MerchantKey
@MerchantCompanyName
@LargeLogoImage
@HeaderLogoImage
@ProfileBannerImage
@SmallLogoImage
*/
func SetupMerchantImage(c *gin.Context) {
	var m Merchant
	contentType, cip, err := service.CheckContentType(c, &m)
	mylog.SetIP(cip)
	if err != nil {
		mylog.Tf("[Error]", "Merchant", "SetupMerchantImage", "%s", err.Error())
		c.JSON(200, gin.H{"err": err.Error()})
		return
	}
	go mylog.Tf("[Info]", "Merchant", "SetupMerchantImage", "Request : %+v", m)
	rxp := `^[0-9a-zA-Z\p{Han}][^\t\v\\\^\$\*\+\?\{\}\(\)\[\]\|]*$`
	v := []myutil.ValiStringType{
		{S: m.MerchantKey, Rxp: rxp, Min: 64, Max: 64},
		{S: m.MerchantCompanyName, Rxp: rxp, Min: 1, Max: 255},
		{S: m.LargeLogoImage, Rxp: rxp, Min: 64, Max: 64},
		{S: m.HeaderLogoImage, Rxp: rxp, Min: 64, Max: 64},
		{S: m.ProfileBannerImage, Rxp: rxp, Min: 64, Max: 64},
		{S: m.SmallLogoImage, Rxp: rxp, Min: 64, Max: 64},
	}
	e := validata.ValiStrings(v...)
	for er := range e {
		if er != nil {
			mylog.Tf("[Error]", "Merchant", "SetupMerchantImage", "validata fail. %s", er.Error())
			c.JSON(200, gin.H{"err": er.Error()})
			return
		}
	}
	selection := bson.M{"MerchantKey": m.MerchantKey, "MerchantCompanyName": m.MerchantCompanyName}
	b := bson.M{"$set": bson.M{"LargeLogoImage": m.LargeLogoImage, "HeaderLogoImage": m.HeaderLogoImage, "ProfileBannerImage": m.ProfileBannerImage, "SmallLogoImage": m.SmallLogoImage, "LastModifTime": time.Now()}}
	query := func(c *mgo.Collection) error {
		return c.Update(selection, b)
	}
	err = mymgo.Do("Merchant", query)
	if err != nil {
		mylog.Tf("[Error]", "Merchant", "SetupMerchantImage", "Fail to update Merchant images into db. %s | %s", m.MerchantCompanyName, err.Error())
		c.JSON(200, gin.H{"Message": "Fail to update Merchant images into db.", "ContentType": contentType, "Status": "Fail", "err": err.Error()})
		return
	}
	go mylog.Tf("[Info]", "Merchant", "SetupMerchantImage", "Setup images was okey. %+v", b)
	c.JSON(200, gin.H{"Message": "Signup was okey.",
		"ContentType":    contentType,
		"Status":         "Success",
		"MerchantStatus": m.Active,
		"MerhantInfo":    b,
	})

}

/*
SetupMerchantManager :
@MerchantKey
@MerchantCompanyName
@AdminUser
@ManagerUser
@User
*/
func SetupMerchantManager(c *gin.Context) {
	var mm, finder Merchant
	contentType, cip, err := service.CheckContentType(c, &mm)
	mylog.SetIP(cip)
	if err != nil {
		mylog.Tf("[Error]", "Merchant", "SetupMerchantManager", "%s", err.Error())
		c.JSON(200, gin.H{"err": err.Error()})
		return
	}
	go mylog.Tf("[Info]", "Merchant", "SetupMerchantManager", "Request : %+v", mm)
	_, err = validata.ValiString(mm.MerchantKey, `^[0-9a-zA-Z\p{Han}][^\t\v\\\^\$\*\+\?\{\}\(\)\[\]\|]*$`, 64, 64)
	if err != nil {
		mylog.Tf("[Error]", "Merchant", "SetupMerchantManager", "validata fail. %s", err.Error())
		c.JSON(200, gin.H{"err": err.Error()})
		return
	}
	if len(mm.AdminUser) <= 0 || len(mm.User) <= 0 || len(mm.ManagerUser) <= 0 {
		mylog.Tf("[Error]", "Merchant", "SetupMerchantManager", "validata fail to must have set one user ")
		c.JSON(200, gin.H{"Messge": "validata fail to must have set one user ", "Status": "Fail", "ContentType": contentType})
		return
	}
	selection := bson.M{"MerchantKey": mm.MerchantKey, "MerchantCompanyName": mm.MerchantCompanyName}
	find := func(c *mgo.Collection) error {
		return c.Find(selection).One(&finder)
	}
	err = mymgo.Do("Merchant", find)
	if err != nil {
		mylog.Tf("[Error]", "Merchant", "SetupMerchantManager", "This Merchant not have in db. %s | %s", mm.MerchantCompanyName, err.Error())
		c.JSON(200, gin.H{"Message": "This Merchant not have in db.", "ContentType": contentType, "Status": "Fail", "err": err.Error()})
		return
	}
	b := bson.M{"$addToSet": bson.M{"AdminUser": bson.M{"$each": mm.AdminUser}, "ManagerUser": bson.M{"$each": mm.ManagerUser}, "User": bson.M{"$each": mm.User}}}
	query := func(c *mgo.Collection) error {
		return c.Update(selection, b)
	}
	err = mymgo.Do("Merchant", query)
	if err != nil {
		mylog.Tf("[Error]", "Merchant", "SetupMerchantManager", "Fail to update Merchant manager into db. %s | %s", mm.MerchantCompanyName, err.Error())
		c.JSON(200, gin.H{"Message": "Fail to update Merchant manager into db.", "ContentType": contentType, "Status": "Fail", "err": err.Error()})
		return
	}
	go func() {
		fa, ra, _ := myutil.CompareArrary(finder.AdminUser, mm.AdminUser)
		fm, rm, _ := myutil.CompareArrary(finder.ManagerUser, mm.ManagerUser)
		fu, ru, _ := myutil.CompareArrary(finder.User, mm.User)
		if len(fa) > 0 {
			s1 := bson.M{"UserKey": bson.M{"$in": fa}}
			b := bson.M{"$pull": bson.M{"AdminMerchant": mm.MerchantKey}}
			query := func(c *mgo.Collection) error {
				return c.Update(s1, b)
			}
			if err := mymgo.Do("User", query); err != nil {
				mylog.Tf("[Error]", "Merchant", "SetupMerchantManager", "Fail to Remove Merchant Admin from User. %s | %s", mm.MerchantCompanyName, err.Error())
			} else {
				mylog.Tf("[Info]", "Merchant", "SetupMerchantManager", "Success to Remove Merchant Admin from User. %s | %s", mm.MerchantCompanyName, fa)
			}
		}
		if len(fm) > 0 {
			s1 := bson.M{"UserKey": bson.M{"$in": fm}}
			b := bson.M{"$pull": bson.M{"ManagerMerchant": mm.MerchantKey}}
			query := func(c *mgo.Collection) error {
				return c.Update(s1, b)
			}
			if err := mymgo.Do("User", query); err != nil {
				mylog.Tf("[Error]", "Merchant", "SetupMerchantManager", "Fail to Remove Merchant Manager from User. %s | %s", mm.MerchantCompanyName, err.Error())
			} else {
				mylog.Tf("[Info]", "Merchant", "SetupMerchantManager", "Success to Remove Merchant Manager from User. %s | %s", mm.MerchantCompanyName, fm)
			}
		}
		if len(fu) > 0 {
			s1 := bson.M{"UserKey": bson.M{"$in": fu}}
			b := bson.M{"$pull": bson.M{"UserMerchant": mm.MerchantKey}}
			query := func(c *mgo.Collection) error {
				return c.Update(s1, b)
			}
			if err := mymgo.Do("User", query); err != nil {
				mylog.Tf("[Error]", "Merchant", "SetupMerchantManager", "Fail to Remove Merchant User from User. %s | %s", mm.MerchantCompanyName, err.Error())
			} else {
				mylog.Tf("[Info]", "Merchant", "SetupMerchantManager", "Success to Remove Merchant User from User. %s | %s", mm.MerchantCompanyName, fu)
			}
		}
		if len(ra) > 0 {
			s1 := bson.M{"UserKey": bson.M{"$in": ra}}
			b := bson.M{"$addToSet": bson.M{"AdminMerchant": mm.MerchantKey}}
			query := func(c *mgo.Collection) error {
				return c.Update(s1, b)
			}
			if err := mymgo.Do("User", query); err != nil {
				mylog.Tf("[Error]", "Merchant", "SetupMerchantManager", "Fail to Add Merchant Admin into User. %s | %s", mm.MerchantCompanyName, err.Error())
			} else {
				mylog.Tf("[Info]", "Merchant", "SetupMerchantManager", "Success to Add Merchant Admin into User. %s | %s", mm.MerchantCompanyName, ra)
			}
		}
		if len(rm) > 0 {
			s1 := bson.M{"UserKey": bson.M{"$in": rm}}
			b := bson.M{"$addToSet": bson.M{"ManagerMerchant": mm.MerchantKey}}
			query := func(c *mgo.Collection) error {
				return c.Update(s1, b)
			}
			if err := mymgo.Do("User", query); err != nil {
				mylog.Tf("[Error]", "Merchant", "SetupMerchantManager", "Fail to Add Merchant Manager into User. %s | %s", mm.MerchantCompanyName, err.Error())
			} else {
				mylog.Tf("[Info]", "Merchant", "SetupMerchantManager", "Success to Add Merchant Merchant into User. %s | %s", mm.MerchantCompanyName, rm)
			}
		}
		if len(ru) > 0 {
			s1 := bson.M{"UserKey": bson.M{"$in": ru}}
			b := bson.M{"$addToSet": bson.M{"UserMerchant": mm.MerchantKey}}
			query := func(c *mgo.Collection) error {
				return c.Update(s1, b)
			}
			if err := mymgo.Do("User", query); err != nil {
				mylog.Tf("[Error]", "Merchant", "SetupMerchantManager", "Fail to Add Merchant User into User. %s | %s", mm.MerchantCompanyName, err.Error())
			} else {
				mylog.Tf("[Info]", "Merchant", "SetupMerchantManager", "Success to Add Merchant User into User. %s | %s", mm.MerchantCompanyName, ru)
			}
		}
	}()
	go mylog.Tf("[Info]", "Merchant", "SetupMerchantManager", "Signup was okey. %+v", b)
	c.JSON(200, gin.H{"Message": "Signup was okey.",
		"ContentType":    contentType,
		"Status":         `Success`,
		"MerchantStatus": "Active",
	})
}

/*
ChangeStatus :
@MerchantKey
@MerchantCompanyName
*/
func ChangeStatus(c *gin.Context) {
	var m Merchant
	contentType, cip, err := service.CheckContentType(c, &m)
	mylog.SetIP(cip)
	if err != nil {
		mylog.Tf("[Error]", "Merchant", "ChangeStatus", "%s", err.Error())
		c.JSON(200, gin.H{"err": err.Error()})
		return
	}
	go mylog.Tf("[Info]", "Merchant", "ChangeStatus", "Request : %+v", m)
	_, err = validata.ValiString(m.MerchantKey, `^[0-9a-zA-Z\p{Han}][^\t\v\\\^\$\*\+\?\{\}\(\)\[\]\|]*$`, 64, 64)
	if err != nil {
		mylog.Tf("[Error]", "Merchant", "ChangeStatus", "validata fail. %s", err.Error())
		c.JSON(200, gin.H{"Message": "Fail to viladata.", "ContentType": contentType, "Status": "Fail", "err": err.Error()})
		return
	}
	selection := bson.M{"MerchantKey": m.MerchantKey, "MerchantCompanyName": m.MerchantCompanyName}
	find := func(c *mgo.Collection) error {
		return c.Find(selection).One(&m)
	}
	if err := mymgo.Do("Merchant", find); err != nil {
		mylog.Tf("[Error]", "Merchant", "ChangeStatus", "Feil to find merchant. %s", err.Error())
		c.JSON(200, gin.H{"Message": "Fail to find this merchant from db.", "ContentType": contentType, "Status": "Fail", "err": err.Error()})
		return
	}
	b := bson.M{"$set": bson.M{"Active": !m.Active}}
	query := func(c *mgo.Collection) error {
		return c.Update(selection, b)
	}
	err = mymgo.Do("Merchant", query)
	if err != nil {
		mylog.Tf("[Error]", "Merchant", "ChangeStatus", "Fail to change merchant status. %s", err.Error())
		c.JSON(200, gin.H{"Message": "Fail to change Merchant status.", "ContentType": contentType, "Status": "Fail", "err": err.Error()})
		return
	}
	go mylog.Tf("[Info]", "Merchant", "ChangeStatus", "Success to change merchant status okey. %+v", b)
	c.JSON(200, gin.H{"Message": "Signup was okey.",
		"ContentType":    contentType,
		"Status":         `Success`,
		"MerchantStatus": !m.Active,
	})
}

//List :
func List(c *gin.Context) {
	contentType := c.Request.Header.Get("Content-Type")
	ip := c.ClientIP()
	mylog.SetIP(ip)
	p, _ := strconv.Atoi(c.Query("page"))
	s, _ := strconv.Atoi(c.Query("pagesize"))
	type listmerchant struct {
		MerchantCompanyName string        `bson:"MerchantCompanyName" form:"MerchantCompanyName" json:"MerchantCompanyName"`
		SmallLogoImage      string        `bson:"SmallLogoImage" form:"SmallLogoImage" json:"SmallLogoImage"`
		MerchantKey         string        `bson:"MerchantKey" form:"MerchantKey" json:"MerchantKey"`
		ID                  bson.ObjectId `bson:"_id" form:"_id" json:"_id"`
		Active              bool          `bson:"Active" form:"Active" json:"Active"`
		CreateTime          time.Time     `bson:"CreateTime" form:"CreateTime" json:"CreateTime"`
	}
	var merchants []listmerchant
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{}).Sort("CreateTime").Skip(p * s).Limit(s).All(&merchants)
	}
	if err := mymgo.Do("Merchant", query); err != nil {
		mylog.Tf("[Error]", "Merchant", "List", "Fail to list merchants. %s", err.Error())
		c.JSON(200, gin.H{"Message": "Fail to list merchants",
			"ContentType": contentType,
			"err":         err.Error(),
			"Status":      "Fail"})
		return
	}
	go mylog.Tf("[Info]", "Merchant", "List", "List merchants currect.")
	c.JSON(200, gin.H{"Message": "List merchants currect.",
		"ContentType": contentType,
		"Status":      "Success",
		"Merchants":   merchants,
	})
}
