package inventory

import (
	"shoppingzone/mylib/mylog"
	"shoppingzone/mylib/mymgo"
	"shoppingzone/myutil"
	"shoppingzone/service"
	"time"

	"github.com/gin-gonic/gin"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var validata myutil.Validata

/*
SignUpInventory :
*/
func SignUpInventory(c *gin.Context) {
	type qm struct {
		MerchantKey         string `bson:"MerchantKey" form:"MerchantKey" json:"MerchantKey"`
		MerchantCompanyName string `bson:"MerchantCompanyName" form:"MerchantCompanyName" json:"MerchantCompanyName"`
		InventoryKey        string `bson:"InventoryKey" form:"InventoryKey" json:"InventoryKey"`
		District            string `bson:"District" form:"District" json:"District"`
		PostalCode          string `bson:"PostalCode" form:"PostalCode" json:"PostalCode"`
		Apartment           string `bson:"Apartment" form:"Apartment" json:"Apartment"`
		Floor               string `bson:"Floor" form:"Floor" json:"Floor"`
		BlockNo             string `bson:"BlockNo" form:"BlockNo" json:"BlockNo"`
		Building            string `bson:"Building" form:"Building" json:"Building"`
		StreetNo            string `bson:"StreetNo" form:"StreetNo" json:"StreetNo"`
		Street              string `bson:"Street" form:"Street" json:"Street"`
		IsDefault           bool   `bson:"IsDefault" form:"IsDefault" json:"IsDefault"`
		Active              bool   `bson:"Active" form:"Active" json:"Active"`
	}
	var m qm
	contentType, cip, err := service.CheckContentType(c, &m)
	mylog.SetIP(cip)
	if err != nil {
		mylog.Tf("[Error]", "Inventory", "SignUpInventory", "%s", err.Error())
		c.JSON(200, gin.H{"err": err.Error()})
		return
	}
	go mylog.Tf("[Info]", "Inventory", "SignUpInventory", "Request : %+v", m)
	rxp := `^[0-9a-zA-Z\p{Han}][^\t\v\\\^\$\*\+\?\{\}\(\)\[\]\|]*$`
	v := []myutil.ValiStringType{
		{S: m.MerchantCompanyName, Rxp: rxp, Min: 1, Max: 255},
		{S: m.District, Rxp: rxp, Min: 1, Max: 255},
		{S: m.Apartment, Rxp: rxp, Min: 1, Max: 255},
		{S: m.Floor, Rxp: rxp, Min: 1, Max: 255},
		{S: m.BlockNo, Rxp: rxp, Min: 1, Max: 255},
		{S: m.Building, Rxp: rxp, Min: 1, Max: 255},
		{S: m.Street, Rxp: rxp, Min: 1, Max: 255},
		{S: m.StreetNo, Rxp: rxp, Min: 1, Max: 255},
		{S: m.MerchantKey, Rxp: rxp, Min: 64, Max: 64},
	}
	e := validata.ValiStrings(v...)
	for er := range e {
		if er != nil {
			mylog.Tf("[Error]", "Inventory", "SignUpInventory", "validata fail. %s", er.Error())
			c.JSON(200, gin.H{"Status": "Fail", "Message": "validata fail.", "err": er.Error()})
			return
		}
	}
	b := bson.M{
		"InventoryKey":        <-myutil.GetRandomString(64),
		"MerchantCompanyName": m.MerchantCompanyName,
		"MerchantKey":         m.MerchantKey,
		"District":            m.District,
		"PostalCode":          m.PostalCode,
		"Apartment":           m.Apartment,
		"Floor":               m.Floor,
		"BlockNo":             m.BlockNo,
		"Building":            m.Building,
		"StreetNo":            m.StreetNo,
		"Street":              m.Street,
		"IsDefault":           true,
		"Active":              true,
		"CreateTime":          time.Now(),
	}
	query := func(c *mgo.Collection) error {
		return c.Insert(b)
	}
	err = mymgo.Do("Invetory", query)
	if err != nil {
		mylog.Tf("[Error]", "Inventory", "SignUpInventory", "Fail to inster this Inventory into db. %s | %s", m.MerchantCompanyName, err.Error())
		c.JSON(200, gin.H{"Message": "Fail to inster this Inventory into db.", "ContentType": contentType, "Status": "Fail", "err": err.Error()})
		return
	}
	go mylog.Tf("[Info]", "Inventory", "SignUpInventory", "It's currect to set active status with this user. %+v", b)
	c.JSON(200, gin.H{"Message": "It's currect to sign up this inventory",
		"ContentType": contentType,
		"Status":      "Success",
		"Inventory":   b,
	})
}
