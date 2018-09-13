package mydb

import (
	"database/sql"
	"shoppingzone/conf"
	"shoppingzone/mylib/mylog"
	"shoppingzone/myutil"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql" //
)

//MyDB :
type MyDB struct {
	db *sql.DB
}

//SQLResoult :
type SQLResoult struct {
	Ids     int64
	Resoult sql.Result
	Err     error
}

//DBConn :
func (d *MyDB) DBConn() {
	var err error
	dbConnect := conf.DbConfig.Usr + ":" + conf.DbConfig.Pwd + "@" + conf.DbConfig.Protocol + "(" + conf.DbConfig.Host + ":" + strconv.Itoa(conf.DbConfig.Port) + ")/" + conf.DbConfig.DBname
	d.db, err = sql.Open("mysql", dbConnect)
	d.db.SetMaxIdleConns(conf.DbConfig.MaxIdleConns)
	d.db.SetMaxOpenConns(conf.DbConfig.MaxOpenConns)
	d.db.SetConnMaxLifetime(time.Duration(conf.DbConfig.SetConnMaxLifetime) * time.Second)
	if err != nil {
		mylog.Tf("[Error]", "MyDB", "DBconn", "DB connect is error! %s", err.Error())
	}
	mylog.Tf("[info]", "MyDB", "DBconnect", "%s", conf.DbConfig.Usr+"@"+conf.DbConfig.Host+" open "+conf.DbConfig.DBname+" is OK!")
}

//RunSQLFile :
func (d *MyDB) RunSQLFile(sqls <-chan myutil.ResoultReadFile) <-chan string {
	out := make(chan string, 3)
	go func() {
		for s := range sqls {
			startTime := time.Now()
			_, err := d.db.Exec(s.Context)
			if err != nil {
				mylog.Tf("[Error]", "MyDB", "RunSQLFile", "%s was run error. %s", s.FileName, err.Error())
				useTime := time.Since(startTime)
				out <- s.FileName + " Error xxxxx " + err.Error() + " xxxxx [" + useTime.String() + "]"
			}
			useTime := time.Since(startTime)
			out <- s.FileName + " Done [" + useTime.String() + "]"
		}
		close(out)
	}()
	return out
}

//SQLExec :
func (d *MyDB) SQLExec(query string, args ...interface{}) SQLResoult {
	rs, err := d.db.Exec(query, args...)
	if err != nil {
		mylog.Tf("[Error]", "MyDB", "SQLExec", "SQL run was error! %s", err.Error())
		return SQLResoult{0, nil, err}
	}
	inserid, err1 := rs.LastInsertId()
	rowsaffect, err2 := rs.RowsAffected()
	if err1 != nil || err2 != nil {
		mylog.Tf("[Error]", "MyDB", "SQLExec", "SQL run was error! %s %s", err1.Error(), err2.Error())
		return SQLResoult{0, nil, err}
	} else if inserid > 0 {
		return SQLResoult{inserid, rs, nil}
	}
	return SQLResoult{rowsaffect, rs, nil}
}

//SQLQuery :
func (d *MyDB) SQLQuery(query string, args ...interface{}) *sql.Rows {
	rs, err := d.db.Query(query, args...)
	if err != nil {
		mylog.Tf("[Error]", "MyDB", "SQLQuery", "SQL run was error! %s", err.Error())
	}
	return rs
}

//Close :
func (d *MyDB) Close() {
	d.db.Close()
}
