package test

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	db, _ = sql.Open("mysql", "root:123456@tcp(10.210.215.245:3306)/test?charset=utf8")
	db.SetMaxOpenConns(0)
	db.SetMaxIdleConns(50)
	db.Ping()
}

func TestMysql() {
	rows, err := db.Query("SELECT * FROM mlog limit 1")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}

	record := make(map[string]string)
	for rows.Next() {
		//将行数据保存到record字典
		err = rows.Scan(scanArgs...)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = string(col.([]byte))
			}
		}
	}
	fmt.Println(record)
}
