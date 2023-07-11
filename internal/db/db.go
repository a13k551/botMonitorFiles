package db

import (
	"fmt"
	"log"

	"github.com/a13k551/botMonitorFiles/internal/config"
	"github.com/jackc/pgx"
)

func CreateBase(conf config.Config) {

	ConnConfig := pgx.ConnConfig{User: conf.UserDB, Password: conf.PassDB,
		Port: conf.Port, Database: conf.NameDB}
	db, err := pgx.Connect(ConnConfig)

	if err != nil {
		fmt.Println(err)
		log.Panic(err)
	}

	defer db.Close()

	rows, err := db.Query("SELECT datname FROM pg_database where datname = $1", conf.NameDB)
	if err != nil {
		log.Panic(err)
	}

	if rows.Next() {
		return
	}

	createDbString := fmt.Sprintf("CREATE DATABASE %s", conf.NameDB)
	_, err = db.Exec(createDbString)

	if err != nil {
		panic(err)
	}
}

func CreateTables(conf config.Config) {

	ConnConfig := pgx.ConnConfig{User: conf.UserDB, Password: conf.PassDB,
		Port: conf.Port, Database: conf.NameDB}
	db, err := pgx.Connect(ConnConfig)

	if err != nil {
		log.Panic(err)
	}

	defer db.Close()

	createTableString := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s varchar(200))",
		conf.TableNamedb, conf.FieldnameDB)
	_, err = db.Exec(createTableString)
	if err != nil {
		panic(err)
	}
}

func CreateConnection(conf config.Config) (*pgx.Conn, error) {

	ConnConfig := pgx.ConnConfig{User: conf.UserDB, Password: conf.PassDB,
		Port: conf.Port, Database: conf.NameDB}
	db, err := pgx.Connect(ConnConfig)

	return db, err
}

func WriteFilepathToDB(filePath string, db *pgx.Conn, conf config.Config) {

	insertString := fmt.Sprintf("insert into %s values ($1)", conf.TableNamedb)
	_, err := db.Exec(insertString, filePath)

	if err != nil {
		log.Println(err)
	}
}

func FileWasSended(db *pgx.Conn, findedFilePath string, conf config.Config) bool {
	fileWasSended := false

	selectSring := fmt.Sprintf("select %s from %s where %s = $1",
		conf.FieldnameDB,
		conf.TableNamedb,
		conf.FieldnameDB)

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
