package mylog

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"shoppingzone/conf"
	"strings"
	"time"
)

type mylog struct{}

var ip string

func init() {
	args := checkOSArgs()
	createLogFile(args.LogPath, args.HTTPPort)
}

//SetIP :
func SetIP(cip string) {
	ip = cip
}

func createLogFile(logPath string, server string) {
	t := time.Now().Format("20060102")
	isExists := <-existPath(logPath)
	if !isExists {
		if err := os.Mkdir(logPath, 0777); err != nil {
			log.Panic("can't make dir for log file! please check !")
		}
	}
	isExists = <-existPath(logPath + "/tar")
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
	log.SetOutput(outfile)
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
	files := listFiles(logPath)
	for f := range files {
		checkFilename := regexp.MustCompile(`.*\.log`).FindAllString(f, -1)
		if len(checkFilename) > 0 && f != filename {
			go tARfile(logPath, f)
		}
	}
}

//Ln :
func Ln(level string, service string, functaion string, data ...interface{}) {
	log.SetPrefix(level)
	log.Output(2, fmt.Sprintln(`[`, service, `] [`, functaion, `] [`, ip, `] `, data))
}

//Tf :
func Tf(level string, service string, functaion string, format string, args ...interface{}) {
	log.SetPrefix(level)
	log.Output(2, fmt.Sprintf(`[`+service+`] [`+functaion+`] [`+ip+`] `+format, args...))
}

type outOSArgs struct {
	HTTPPort string
	NetProt  string
	LogPath  string
}

func checkOSArgs() outOSArgs {
	var out outOSArgs
	var tp []string
	if len(os.Args) > 1 {
		for idx, v := range os.Args {
			if v == "-log" {
				out.LogPath = os.Args[idx+1]
				continue
			}
			tp = regexp.MustCompile(`^-log=`).FindAllString(v, -1)
			if len(tp) > 0 && tp[0] == "-log=" {
				out.LogPath = strings.Split(v, "=")[1]
				continue
			}
			if v == "-P" {
				out.HTTPPort = os.Args[idx+1]
				continue
			}
			tp = regexp.MustCompile(`^-P=`).FindAllString(v, -1)
			if len(tp) > 0 && tp[0] == "-P=" {
				out.HTTPPort = strings.Split(v, "=")[1]
				continue
			}
			if v == "-p" {
				out.NetProt = os.Args[idx+1]
				continue
			}
			tp = regexp.MustCompile(`^-p=`).FindAllString(v, -1)
			if len(tp) > 0 && tp[0] == "-p=" {
				out.NetProt = strings.Split(v, "=")[1]
				continue
			}
		}
	}
	if out.LogPath == "" {
		out.LogPath = conf.Config.LogFile
	}
	if out.HTTPPort == "" {
		out.HTTPPort = conf.Config.HTTPPort
	}
	if out.NetProt == "" {
		out.NetProt = conf.Config.NetPort
	}
	return out
}

func existPath(path string) <-chan bool {
	out := make(chan bool, 3)
	go func() {
		_, err := os.Stat(path)
		if err == nil {
			out <- true
		}
		if os.IsNotExist(err) {
			out <- false
			log.Println(err)
		}
		close(out)
	}()
	return out
}

func listFiles(path string) <-chan string {
	out := make(chan string, 3)
	go func() {
		dirList, err := ioutil.ReadDir(path)
		if err != nil {
			log.Println("List files in " + path + " was err")
			log.Println(err)
		}
		for _, v := range dirList {
			out <- v.Name()
		}
		close(out)
	}()
	return out
}

func tARfile(p string, f string) {
	tarCmd := exec.Command("tar", "-zcf", p+"/tar/"+f+".tar.gz", p+"/"+f, "--remove-files")
	tarCmd.Run()
	rmCmd := exec.Command("rm", "-fv", p+"/"+f)
	rmCmd.Run()
}
