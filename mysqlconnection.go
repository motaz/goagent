package main

import (
	"database/sql"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func getMySQLConnection(databasename string) (db *sql.DB, err error) {

	cdrdbserver := getConfigValueLocal("cdrdbserver")
	cdruser := getConfigValueLocal("cdruser")
	cdrpass := getConfigValueLocal("cdrpass")

	connection := cdruser + ":" + cdrpass + "@tcp(" + cdrdbserver + ":3306)/" + databasename + "?charset=utf8"

	db, err = sql.Open("mysql", connection)

	return
}

func checkDatabase(db *sql.DB, databasename string) (err error) {

	_, err = db.Query("select now()")
	if err != nil {
		println(err.Error())
		if strings.Contains(err.Error(), "Unknown database") {
			newdb, err := getMySQLConnection("")
			_, err = newdb.Exec("CREATE DATABASE " + databasename +
				" CHARACTER SET utf8 COLLATE utf8_general_ci;")
			if err != nil {
				println("Error creating database: ", err.Error())
			}
		}
	}
	return
}

func checkTable(db *sql.DB, tablename string, tablescript string) (err error) {

	_, err = db.Query("select * from " + tablename + " limit 1")
	if err != nil {
		println(err.Error())
		if strings.Contains(err.Error(), "doesn't exist") {
			_, err = db.Exec(tablescript)
			if err != nil {
				println("Error creating table: ", err.Error())
			}
		}
	}
	return
}
