package user

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

//User :
type User struct {
	UserKey               string        `bson:"UserKey" form:"UserKey" json:"UserKey"`
	UserName              string        `bson:"Username" form:"Username" json:"Username"`
	FirstName             string        `bson:"FirstName" form:"FirstName" json:"FirstName"`
	LastName              string        `bson:"LastName" form:"LastName" json:"LastName"`
	MIDdleName            string        `bson:"MIDdleName" form:"MIDdleName" json:"MIDdleName"`
	DisplayName           string        `bson:"DisplayName" form:"DisplayName" json:"DisplayName"`
	Email                 string        `bson:"Email" form:"Email" json:"Email"`
	MobileCode            string        `bson:"Mobilecode" form:"Mobilecode" json:"Mobilecode"`
	MobileNumber          string        `bson:"Mobilenumber" form:"Mobilenumber" json:"Mobilenumber"`
	LanguageID            int           `bson:"LanguageID" form:"LanguageID" json:"LanguageID"`
	TimeZoneID            int           `bson:"TimeZoneID" form:"TimeZoneID" json:"TimeZoneID"`
	ProfileImage          string        `bson:"ProfileImage" form:"ProfileImage" json:"ProfileImage"`
	ProfileAlternateImage string        `bson:"ProfileAlternateImage" form:"ProfileAlternateImage" json:"ProfileAlternateImage"`
	CoverImage            string        `bson:"CoverImage" form:"CoverImage" json:"CoverImage"`
	CoverAlternateImage   string        `bson:"CoverAlternateImage" form:"CoverAlternateImage" json:"CoverAlternateImage"`
	IsCurator             int           `bson:"IsCurator" form:"IsCurator" json:"IsCurator"`
	IsMerchant            int           `bson:"IsMerchant" form:"IsMerchant" json:"IsMerchant"`
	IsMm                  int           `bson:"IsMm" form:"IsMm" json:"IsMm"`
	IsFeatured            int           `bson:"IsFeatured" form:"IsFeatured" json:"IsFeatured"`
	InventoryLocationID   int           `bson:"InventoryLocationID" form:"InventoryLocationID" json:"InventoryLocationID"`
	DefaultUserAddressID  int           `bson:"DefaultUserAddressID" form:"DefaultUserAddressID" json:"DefaultUserAddressID"`
	UserIDentificationID  int           `bson:"UserIDentificationID" form:"UserIDentificationID" json:"UserIDentificationID"`
	ReferrerUserID        int           `bson:"ReferrerUserID" form:"ReferrerUserID" json:"ReferrerUserID"`
	Hash                  string        `bson:"Hash" form:"Hash" json:"Hash"`
	Salt                  string        `bson:"Salt" form:"Salt" json:"Salt"`
	StatusID              int           `bson:"StatusID" form:"StatusID" json:"StatusID"`
	StatusReasonCode      string        `bson:"StatusReasonCode" form:"StatusReasonCode" json:"StatusReasonCode"`
	GeoCountryID          int           `bson:"GeoCountryID" form:"GeoCountryID" json:"GeoCountryID"`
	GeoProvinceID         int           `bson:"GeoProvinceID" form:"GeoProvinceID" json:"GeoProvinceID"`
	GeoCityID             int           `bson:"GeoCityID" form:"GeoCityID" json:"GeoCityID"`
	Gender                string        `bson:"Gender" form:"Gender" json:"Gender"`
	DateOfBirth           string        `bson:"DateOfBirth" form:"DateOfBirth" json:"DateOfBirth"`
	UserDescription       string        `bson:"UserDescription" form:"UserDescription" json:"UserDescription"`
	SignupInviteCode      string        `bson:"SignupInviteCode" form:"SignupInviteCode" json:"SignupInviteCode"`
	ReferralInviteCode    string        `bson:"ReferralInviteCode" form:"ReferralInviteCode" json:"ReferralInviteCode"`
	Priority              int           `bson:"Priority" form:"Priority" json:"Priority"`
	CommissionRate        float32       `bson:"CommissionRate" form:"CommissionRate" json:"CommissionRate"`
	Password              string        `bson:"Password" form:"Password" json:"Password"`
	ID                    bson.ObjectId `bson:"_id" form:"_id" json:"_id"`
	Active                bool          `bson:"Active" form:"Active" json:"Active"`
	CreateTime            time.Time     `bson:"CreateTime" form:"CreateTime" json:"CreateTime"`
	ActiveTime            time.Time     `bson:"ActiveTime" form:"ActiveTime" json:"ActiveTime"`
	SecureLabel           []string      `bson:"SecureLabel" form:"SecureLabel" json:"SecureLabel"`
	SecureGroup           []string      `bson:"SecureGroup" form:"SecureGroup" json:"SecureGroup"`
	SecureRouter          []string      `bson:"SecureRouter" form:"SecureRouter" json:"SecureRouter"`
	SecureRole            []string      `bson:"SecureRole" form:"SecureRole" json:"SecureRole"`
	AdminMerchant         []string      `bson:"AdminMerchant" form:"AdminMerchant" json:"AdminMerchant"`
	ManagerMerchant       []string      `bson:"ManagerMerchant" form:"ManagerMerchant" json:"ManagerMerchant"`
	UserMerchant          []string      `bson:"UserMerchant" form:"UserMerchant" json:"UserMerchant"`
}

type signupUser struct {
	Username     string `form:"Username" json:"Username"`
	Userpassword string `form:"Pwd" json:"Pwd"`
	Mobilecode   string `form:"Mobilecode" json:"Mobilecode"`
	Mobilenumber string `form:"Mobilenumber" json:"Mobilenumber"`
}

type mobileVerification struct {
	Mobilecode   string `form:"Mobilecode" json:"Mobilecode"`
	Mobilenumber string `form:"Mobilenumber" json:"Mobilenumber"`
}

type verificationCode struct {
	Mobilecode       string `form:"Mobilecode" json:"Mobilecode"`
	Mobilenumber     string `form:"Mobilenumber" json:"Mobilenumber"`
	VerificationCode string `form:"VerificationCode" json:"VerificationCode"`
}

type loginUser struct {
	Signinstring   string `form:"Signinstring" json:"Signinstring"`
	Signinpassword string `form:"Signinpassword" json:"Signinpassword"`
}

type userSecure struct {
	SecureField []string `form:"SecureField" json:"SecureField" bson:"SecureField"`
	SecureRole  []string `form:"SecureRole" json:"SecureRole" bson:"SecureRole"`
	UserKey     string   `form:"UserKey" json:"UserKey" bson:"UserKey"`
	UserMobile  string   `form:"UserMobile" json:"UserMobile" bson:"UserMobile"`
}

type changeStatus struct {
	UserKey    string `form:"UserKey" json:"UserKey" bson:"UserKey"`
	UserMobile string `form:"UserMobile" json:"UserMobile" bson:"UserMobile"`
	Status     bool   `form:"Status" json:"Status" bson:"Status"`
}
