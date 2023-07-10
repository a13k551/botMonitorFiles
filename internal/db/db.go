package db

import (
	"fmt"
	"log"

	"github.com/jackc/pgx"
)

func CreateBase(conf map[string]interface{}) {

	userdb := conf["userdb"].(string)
	passdb := conf["passdb"].(string)
	namedb := conf["namedb"].(string)
	portdbint := conf["port"].(int)
	portdb := uint16(portdbint)

	ConnConfig := pgx.ConnConfig{User: userdb, Password: passdb, Port: portdb, Database: namedb}
	db, err := pgx.Connect(ConnConfig)

	if err != nil {
		fmt.Println(err)
		log.Panic(err)
	}

	defer db.Close()

	rows, err := db.Query("SELECT datname FROM pg_database where datname = $1", namedb)
	if err != nil {
		log.Panic(err)
	}

	if rows.Next() {
		return
	}

	createDbString := fmt.Sprintf("CREATE DATABASE %s", namedb)
	_, err = db.Exec(createDbString)

	if err != nil {
		panic(err)
	}
}

func CreateTables(conf map[string]interface{}) {

	userdb := conf["userdb"].(string)
	passdb := conf["passdb"].(string)
	namedb := conf["namedb"].(string)
	portdbint := conf["port"].(int)
	portdb := uint16(portdbint)

	ConnConfig := pgx.ConnConfig{User: userdb, Password: passdb, Port: portdb, Database: namedb}
	db, err := pgx.Connect(ConnConfig)

	if err != nil {
		log.Panic(err)
	}

	defer db.Close()

	createTableString := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s varchar(200))",
		conf["tablenamedb"], conf["fieldnamedb"])
	_, err = db.Exec(createTableString)
	if err != nil {
		panic(err)
	}
}

func CreateConnection(conf map[string]interface{}) (*pgx.Conn, error) {
	userdb := conf["userdb"].(string)
	passdb := conf["passdb"].(string)
	namedb := conf["namedb"].(string)
	portdbint := conf["port"].(int)
	portdb := uint16(portdbint)

	ConnConfig := pgx.ConnConfig{User: userdb, Password: passdb, Port: portdb, Database: namedb}
	db, err := pgx.Connect(ConnConfig)

	return db, err
}

func WriteFilepathToDB(filePath string, db *pgx.Conn, conf map[string]interface{}) {

	insertString := fmt.Sprintf("insert into %s values ($1)", conf["tablenamedb"])
	_, err := db.Exec(insertString, filePath)

	if err != nil {
		log.Println(err)
	}
}

func FileWasSended(db *pgx.Conn, findedFilePath string, conf map[string]interface{}) bool {
	fileWasSended := false

	selectSring := fmt.Sprintf("select %s from %s where %s = $1",
		conf["fieldnamedb"],
		conf["tablenamedb"],
		conf["fieldnamedb"])

	rows, err := db.Query(selectSring, findedFilePath)

	if err != nil {
		log.Println(err)
		return false
	}
	defer rows.Close()

	if rows.Next() {
		fileWasSended = true
	}

	return fileWasSended
}
