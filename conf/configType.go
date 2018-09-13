package conf

import "github.com/gin-contrib/sessions"

type myConfig struct {
	PwdSrec         string
	CryptoKey       string
	NsqAddress      []string
	NsqLooupAddress []string
	NsqPoolThresds  int
	LogFile         string
	HTTPPort        string
	NetPort         string
	NetMaxConns     int
}

type redisConfig struct {
	Host        string
	Port        int
	Protocol    string
	Poolconf    redisPoolConfig
	Sessionconf sessions.Options
}

type dbConfig struct {
	Host               string
	Port               int
	Usr                string
	Pwd                string
	Protocol           string
	MaxIdleConns       int
	MaxOpenConns       int
	SetConnMaxLifetime int
	DBname             string
}

type redisPoolConfig struct {
	MaxIdle     int
	MaxActive   int
	IdleTimeout int64
}

type sMS struct {
	URL string
	Un  string
	Pw  string
}

type weChatConfig struct {
	appID        string
	appToken     string
	appsecret    string
	mchID        string
	appIDH5      string
	appsecretH5  string
	appIDApp     string
	appsecretApp string
	notifyURL    string
	paykey       string
	redirctURL   string
	pfx          string
}

type uploadConfig struct {
	Uploadpath   string
	Maxsize      int64
	Allowtype    []string
	Allowheaders []string
	Maxfiles     int
}
