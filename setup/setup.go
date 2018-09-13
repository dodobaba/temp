package main

import (
	"log"
	"regexp"
	"shoppingzone/mylib/mydb"
	"shoppingzone/myutil"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var path = "../sql/setup"

func main() {
	var db mydb.MyDB
	var rs mydb.SQLResoult
	db.DBConn()
	defer db.Close()
	isExists := <-myutil.ExistPath(path)
	if isExists {
		startTime := time.Now()
		log.Println("Now start setup @ " + startTime.Format("2006-01-02 15:04:05.000"))
		files := myutil.ListFiles(path)
		readSQL := myutil.ReadFiles(files, path)
		runSQLResoult := db.RunSQLFile(readSQL)
		for s := range runSQLResoult {
			checkStatus := regexp.MustCompile(`Done|Error`).FindAllString(s, -1)
			checkFilename := regexp.MustCompile(`(?U)^.*\.sql`).FindAllString(s, -1)
			if checkStatus[0] == "Done" {
				l := "REPLACE INTO SqlRunner (SqlRunnerName,StatusId) VALUES (?,?)"
				rs = db.SQLExec(l, checkFilename[0], 1)
			}
			if checkStatus[0] == "Error" {
				checkError := regexp.MustCompile(`(?i:^xxxxx).*xxxxx`).FindAllString(s, -1)
				l := "REPLACE INTO SqlRunner (SqlRunnerName,StatusId,ErrorCount,ErrorText) VALUES (?,?,?,?)"
				rs = db.SQLExec(l, checkFilename[0], 1, 1, checkError[0])
			}
			log.Println(s, rs.Ids, rs.Err)
		}
		log.Println("Finish setup total use " + time.Since(startTime).String())
	} else {
		log.Println("The Path was not found")
	}
}
