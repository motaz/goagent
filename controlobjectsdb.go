package main

import (
	"database/sql"
	"errors"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type controlObjectType struct {
	ID         int    `json:"id"`
	ObjectType string `json:"objecttype"`
	ObjectName string `json:"objectname"`
	Properties string `json:"properties"`
}

func checkControlObjectsTable(db *sql.DB) {

	checkTable(db, "controlobjects", `CREATE TABLE controlobjects (
				id INT NOT NULL AUTO_INCREMENT,
				ObjectType VARCHAR(25) NULL,
				ObjectName VARCHAR(45) NULL,
				Properties JSON NULL,
				PRIMARY KEY (id));
			  `)

}

func setControlObjectInDB(controlObject controlObjectType) (err error) {

	databasename := "simpletrunk"
	db, err := getMySQLConnection(databasename)
	if err == nil {
		checkDatabase(db, databasename)
		checkControlObjectsTable(db)
		err = insertOrUpdateControlObject(db, controlObject)

	}
	return
}

func insertOrUpdateControlObject(db *sql.DB, controlObject controlObjectType) (err error) {

	existControlObject, err := getControlObject(db, controlObject.ObjectType, controlObject.ObjectName)
	if err == nil {
		if existControlObject.ID > 0 {
			err = updateControlObject(db, existControlObject.ID, controlObject)
		} else {
			err = insertControlObject(db, controlObject)
		}
	}
	return
}

func insertControlObject(db *sql.DB, controlObject controlObjectType) (err error) {

	_, err = db.Exec("insert into controlobjects (ObjectType, ObjectName, Properties) "+
		"values (?, ?, ?)", controlObject.ObjectType, controlObject.ObjectName, controlObject.Properties)
	return
}

func updateControlObject(db *sql.DB, id int, controlObject controlObjectType) (err error) {

	_, err = db.Exec("update controlobjects set ObjectType=?, ObjectName=?, Properties=? \n "+
		"where id = ?", controlObject.ObjectType, controlObject.ObjectName, controlObject.Properties, id)
	return
}

func getControlObject(db *sql.DB, objectType string, objectName string) (controlObject controlObjectType, err error) {

	rows, err := db.Query("select id, objectType, objectName, Properties from controlobjects \n"+
		"where lower(objectType) = ? and lower(objectName) = ?",
		strings.ToLower(objectType), strings.ToLower(objectName))

	if err == nil {

		if rows.Next() {
			err = rows.Scan(&controlObject.ID, &controlObject.ObjectType, &controlObject.ObjectName,
				&controlObject.Properties)
		} else {
			controlObject.ID = -1
		}
	}
	return
}

func actualRemoveControlObject(db *sql.DB, objectType string, objectName string) (err error) {

	_, err = db.Exec("delete from controlobjects \n"+
		"where lower(objectType) = ? and lower(objectName) = ?",
		strings.ToLower(objectType), strings.ToLower(objectName))

	return
}

func removeControlObjectFromDB(objectType string, objectName string) error {

	databasename := "simpletrunk"
	db, err := getMySQLConnection(databasename)
	if err == nil {
		co, err := getControlObject(db, objectType, objectName)
		if err == nil && co.ID > 0 {

			err = actualRemoveControlObject(db, objectType, objectName)
		} else if err == nil && co.ID <= 0 {

			return errors.New("Object : " + objectName + "/" + objectType + ", does not exist")

		}

	}
	return err
}
