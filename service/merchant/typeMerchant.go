package merchant

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

//Merchant :
type Merchant struct {
	ID                                 bson.ObjectId `bson:"_id" form:"_id" json:"_id"`
	MerchantID                         int           `bson:"MerchantID" form:"MerchantID" json:"MerchantID"`
	MerchantTypeID                     int           `bson:"MerchantTypeID" form:"MerchantTypeID" json:"MerchantTypeID"`
	MerchantType                       string        `bson:"MerchantType" form:"MerchantType" json:"MerchantType"`
	MerchantNameInvariant              string        `bson:"MerchantNameInvariant" form:"MerchantNameInvariant" json:"MerchantNameInvariant"`
	MerchantCompanyName                string        `bson:"MerchantCompanyName" form:"MerchantCompanyName" json:"MerchantCompanyName"`
	BusinessRegistrationNo             string        `bson:"BusinessRegistrationNo" form:"BusinessRegistrationNo" json:"BusinessRegistrationNo"`
	MerchantSubdomain                  string        `bson:"MerchantSubdomain" form:"MerchantSubdomain" json:"MerchantSubdomain"`
	MerchantCode                       string        `bson:"MerchantCode" form:"MerchantCode" json:"MerchantCode"`
	MerchantDescInvariant              string        `bson:"MerchantDescInvariant" form:"MerchantDescInvariant" json:"MerchantDescInvariant"`
	HeaderLogoImage                    string        `bson:"HeaderLogoImage" form:"HeaderLogoImage" json:"HeaderLogoImage"`
	SmallLogoImage                     string        `bson:"SmallLogoImage" form:"SmallLogoImage" json:"SmallLogoImage"`
	LargeLogoImage                     string        `bson:"LargeLogoImage" form:"LargeLogoImage" json:"LargeLogoImage"`
	ProfileBannerImage                 string        `bson:"ProfileBannerImage" form:"ProfileBannerImage" json:"ProfileBannerImage"`
	ChatBackgroundImage                string        `bson:"ChatBackgroundImage" form:"ChatBackgroundImage" json:"ChatBackgroundImage"`
	IsListedMerchant                   string        `bson:"IsListedMerchant" form:"IsListedMerchant" json:"IsListedMerchant"`
	IsFeaturedMerchant                 string        `bson:"IsFeaturedMerchant" form:"IsFeaturedMerchant" json:"IsFeaturedMerchant"`
	IsRecommendedMerchant              string        `bson:"IsRecommendedMerchant" form:"IsRecommendedMerchant" json:"IsRecommendedMerchant"`
	IsSearchableMerchant               string        `bson:"IsSearchableMerchant" form:"IsSearchableMerchant" json:"IsSearchableMerchant"`
	IsAutoConfirmOrder                 string        `bson:"IsAutoConfirmOrder" form:"IsAutoConfirmOrder" json:"IsAutoConfirmOrder"`
	IsOrderConfirmNotify               string        `bson:"IsOrderConfirmNotify" form:"IsOrderConfirmNotify" json:"IsOrderConfirmNotify"`
	IsCrossBorder                      bool          `bson:"IsCrossBorder" form:"IsCrossBorder" json:"IsCrossBorder"`
	IsFeaturedRed                      string        `bson:"IsFeaturedRed" form:"IsFeaturedRed" json:"IsFeaturedRed"`
	IsFeaturedRedModify                string        `bson:"IsFeaturedRedModify" form:"IsFeaturedRedModify" json:"IsFeaturedRedModify"`
	IsFeaturedBlack                    string        `bson:"IsFeaturedBlack" form:"IsFeaturedBlack" json:"IsFeaturedBlack"`
	PriorityRed                        int           `bson:"PriorityRed" form:"PriorityRed" json:"PriorityRed"`
	PriorityBlack                      int           `bson:"PriorityBlack" form:"PriorityBlack" json:"PriorityBlack"`
	GeoCountryID                       int           `bson:"GeoCountryID" form:"GeoCountryID" json:"GeoCountryID"`
	GeoIDProvince                      int           `bson:"GeoIDProvince" form:"GeoIDProvince" json:"GeoIDProvince"`
	GeoIDCity                          int           `bson:"GeoIDCity" form:"GeoIDCity" json:"GeoIDCity"`
	District                           string        `bson:"District" form:"District" json:"District"`
	PostalCode                         string        `bson:"PostalCode" form:"PostalCode" json:"PostalCode"`
	Apartment                          string        `bson:"Apartment" form:"Apartment" json:"Apartment"`
	Floor                              string        `bson:"Floor" form:"Floor" json:"Floor"`
	BlockNo                            string        `bson:"BlockNo" form:"BlockNo" json:"BlockNo"`
	Building                           string        `bson:"Building" form:"Building" json:"Building"`
	StreetNo                           string        `bson:"StreetNo" form:"StreetNo" json:"StreetNo"`
	Street                             string        `bson:"Street" form:"Street" json:"Street"`
	StatusID                           int           `bson:"StatusID" form:"StatusID" json:"StatusID"`
	ShippingFee                        float32       `bson:"ShippingFee" form:"ShippingFee" json:"ShippingFee"`
	Priority                           int           `bson:"Priority" form:"Priority" json:"Priority"`
	DefaultShipFromInventoryLocationID int           `bson:"DefaultShipFromInventoryLocationID" form:"DefaultShipFromInventoryLocationID" json:"DefaultShipFromInventoryLocationID"`
	DefaultReturnToInventoryLocationID int           `bson:"DefaultReturnToInventoryLocationID" form:"DefaultReturnToInventoryLocationID" json:"DefaultReturnToInventoryLocationID"`
	FreeShippingThreshold              int           `bson:"FreeShippingThreshold" form:"FreeShippingThreshold" json:"FreeShippingThreshold"`
	FreeShippingFrom                   time.Time     `bson:"FreeShippingFrom" form:"FreeShippingFrom" json:"FreeShippingFrom"`
	FreeShippingTo                     time.Time     `bson:"FreeShippingTo" form:"FreeShippingTo" json:"FreeShippingTo"`
	LastStatus                         time.Time     `bson:"LastStatus" form:"LastStatus" json:"LastStatus"`
	LastCreated                        time.Time     `bson:"LastCreated" form:"LastCreated" json:"LastCreated"`
	LastModified                       time.Time     `bson:"LastModified" form:"LastModified" json:"LastModified"`
	Active                             bool          `bson:"Active" form:"Active" json:"Active"`
	CreateTime                         time.Time     `bson:"CreateTime" form:"CreateTime" json:"CreateTime"`
	ActiveTime                         time.Time     `bson:"ActiveTime" form:"ActiveTime" json:"ActiveTime"`
	AdminUser                          []string      `bson:"AdminUser" form:"AdminUser" json:"AdminUser"`
	ManagerUser                        []string      `bson:"ManagerUser" form:"ManagerUser" json:"ManagerUser"`
	User                               []string      `bson:"User" form:"User" json:"User"`
	MerchantKey                        string        `bson:"MerchantKey" form:"MerchantKey" json:"MerchantKey"`
}
