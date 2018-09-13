package conf

import (
	"github.com/gin-contrib/sessions"
)

// Config been use in won config
var Config = myConfig{
	PwdSrec:         "iloveekain1000years",
	CryptoKey:       "100009ea99999ee9a99f99e9e9120811",
	NsqAddress:      []string{"10.51.20.240:4150", "10.51.20.240:4250"},
	NsqLooupAddress: []string{"10.51.20.240:4161", "10.51.20.240:4261"},
	NsqPoolThresds:  10,
	LogFile:         "logs",
	HTTPPort:        ":3000",
	NetPort:         ":30000",
	NetMaxConns:     3,
}

// RedisConf been use only for reids
var RedisConf = redisConfig{
	Host:     "127.0.0.1:6379",
	Port:     6379,
	Protocol: "tcp",
	Poolconf: redisPoolConfig{
		MaxIdle:     3,
		MaxActive:   10,
		IdleTimeout: 240,
	},
	Sessionconf: sessions.Options{
		MaxAge:   300,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	},
}

//DbConfig been used for mysql
var DbConfig = dbConfig{
	Host:               "127.0.0.1",
	Port:               3306,
	Usr:                "root",
	Pwd:                "pa33w0rd",
	Protocol:           "tcp",
	MaxIdleConns:       20,
	MaxOpenConns:       20,
	SetConnMaxLifetime: 3,
	DBname:             "shoppingzone",
}

//SMSConfig configuration
var SMSConfig = sMS{
	URL: "http://222.73.117.138:7891/mt",
	Un:  "N18521779092",
	Pw:  "263310",
}

//WeChatConfig configuration
var WeChatConfig = weChatConfig{
	appID:        "wx38bc9f7435db5d60",
	appToken:     "mymmtest",
	appsecret:    "3082995df6fce9dde26bbb914ea811db",
	mchID:        "1486317882",
	appIDH5:      "wxdc7f3a20643e6f31",
	appsecretH5:  "e5e520cd85d5ffb4caafedd2d8d20dca",
	appIDApp:     "wxdc7f3a20643e6f31",
	appsecretApp: "e5e520cd85d5ffb4caafedd2d8d20dca",
	notifyURL:    "https://dev.mymm.com/wxpay/notify",
	paykey:       "p1Sqsd8n2sn8FkxvORd2NyosR5k4yoIx",
	redirctURL:   "https://beta-m.mymm.com/#/login",
	pfx:          "apiclient_cert.p12",
}

//Uploadconfig :
var Uploadconfig = uploadConfig{
	Uploadpath: "uploads",
	Maxsize:    20 * 1024 * 1024,
	Maxfiles:   5,
	Allowtype: []string{
		"jpg",
		"pdf",
		"csv",
		"png",
		"bmp",
		"xls",
		"xlsx",
		"txt",
		"doc",
		"docx",
		"gif",
		"zip",
	},
	Allowheaders: []string{
		"application/pdf",
		"application/vnd.ms-excel",
		"application/msword",
		"application/zip",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"image/jpeg",
		"image/png",
		"image/gif",
		"image/bmp",
		"text/csv",
		"text/plain",
	},
}
