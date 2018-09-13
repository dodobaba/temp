package mysms

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"shoppingzone/conf"
	"shoppingzone/mylib/mymgo"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

/*
SendSMS :
@m #mobilenumber
@c #the code of veritication
@d #the data of messages
*/
func SendSMS(m string, vc string, d string) error {
	ed := url.QueryEscape(d)
	l := conf.SMSConfig.URL + "?un=" + conf.SMSConfig.Un + "&pw=" + conf.SMSConfig.Pw + "&da=" + m + "&sm=" + ed + "&rd=1&tf=3&dc=15"
	resp, err := http.Get(l)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	query := func(c *mgo.Collection) error {
		selector := bson.M{"Mobile": m, "VerificationCode": vc}
		b := bson.M{"$set": bson.M{"SmsExpRes": string(body), "LastModify": time.Now()}}
		return c.Update(selector, b)
	}
	err = mymgo.Do("VerificationCode", query)
	if err != nil {
		return err
	}
	return nil
}

//VeriticationSMSCode :
func VeriticationSMSCode(m string, vc string) bool {
	r := responesVerificationCode{}
	b := bson.M{"Mobile": m, "VerificationCode": vc}
	query := func(c *mgo.Collection) error {
		return c.Find(b).One(&r)
	}
	err := mymgo.Do("VerificationCode", query)
	if err != nil {
		return false
	}
	return true
}

type responesVerificationCode struct {
	Mobile           string `bson:"Mobile"`
	VerificationCode string `bson:"VerificationCode"`
	SmsExpRes        string `bson:"SmsExpRes"`
}
