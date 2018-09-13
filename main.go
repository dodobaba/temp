package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"shoppingzone/conf"
	"shoppingzone/mylib/mylog"
	"shoppingzone/mylib/mynet"
	"shoppingzone/myutil"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	port         = flag.String("p", ":30000", "net server port")
	httport      = flag.String("P", ":3000", "http service port")
	logPath      = flag.String("log", conf.Config.LogFile, "log file path")
	ginlog       = flag.Bool("ginlogfile", false, "access output gin log to log file")
	startService = flag.String("S", "", "S=[service' name] can start service alone")
)

var r *gin.Engine

func init() {
	mylog.Tf("[Info]", "APP", "Server", "this is number of http port %s", *httport)
	flag.Parse()
	if *ginlog {
		ginLogFile := logFile(*logPath, *httport)
		gin.DefaultWriter = io.MultiWriter(ginLogFile)
	}
	r = gin.Default()
	myrouter(r, *startService)
}

func main() {
	r.Run(*httport)
	var server mynet.MyNet
	l, err := net.Listen("tcp", *port)
	if err != nil {
		mylog.Tf("[Error]", "APP", "Server", "failed to listen: %s", err.Error())
	}
	defer l.Close()
	for {
		c, err := l.Accept()
		if err != nil {
			mylog.Tf("[Error]", "APP", "Server", "Failed to accept: %s", err.Error())
		}
		mylog.Tf("[Info]", "APP", "Server", "Accept from  %s", c.RemoteAddr().String())
		defer c.Close()
		go server.ConnHandel(c)
	}

}

func logFile(logPath string, server string) *os.File {
	t := time.Now().Format("20060102")
	isExists := <-myutil.ExistPath(logPath)
	if !isExists {
		if err := os.Mkdir(logPath, 0777); err != nil {
			log.Panic("can't make dir for log file! please check !")
		}
	}
	isExists = <-myutil.ExistPath(logPath + "/tar")
	if !isExists {
		if err := os.Mkdir(logPath+"/tar", 0777); err != nil {
			log.Panic("can't make tar dir for log file! please check !")
		}
	}
	s := strings.Replace(server, ":", "_", 1)
	filename := t + s + ".log"
	logfile := logPath + "/" + filename
	outfile, err := os.OpenFile(logfile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		panic("log file open failed")
	}
	return outfile
}
