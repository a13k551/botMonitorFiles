package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var conf configuration

func init() {
	conf = getConf()
}

func main() {

	connectionString := fmt.Sprintf("%s:%s@/%s",
		conf.UserDB,
		conf.PassDB,
		conf.NameDB)

	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		log.Panic(err)
	}

	defer db.Close()

	bot, err := tgbotapi.NewBotAPI(conf.Token)
	if err != nil {
		log.Panic(err)
	}

	for {
		findedFiles, err := findFiles()
		if findedFiles == nil || err != nil {
			continue
		}

		findedFiles = filterFilesByDate(findedFiles)

		findedFiles = filterSendedFiles(findedFiles, db)

		sendfindedFiles(findedFiles, db, bot)

		time.Sleep(time.Second * 5)
	}
}

func getConf() configuration {

	conf := configuration{}

	data, err := os.ReadFile("./conf.json")

	if err != nil {
		log.Panic(err)
	}

	err = json.Unmarshal(data, &conf)

	if err != nil {
		log.Panic(err)
	}

	return conf
}

func findFiles() ([]string, error) {
	findString := fmt.Sprintf("%s%s", conf.Path, conf.Mask)
	findedFiles, err := filepath.Glob(findString)

	if err != nil {
		return nil, err
	}

	if findedFiles == nil {
		return nil, nil
	}

	return findedFiles, nil
}

func filterFilesByDate(files []string) []string {

	filtredFiles := []string{}

	for _, dir := range files {
		fileStat, err := os.Stat(dir)
		if err != nil {
			continue
		}

		fileTime := fileStat.ModTime()
		minDate, err := time.Parse("2/1/2006", conf.MinDate)
		if err != nil {
			continue
		}

		if minDate.Before(fileTime) {
			filtredFiles = append(filtredFiles, dir)
		}
	}

	return filtredFiles
}

func filterSendedFiles(files []string, db *sql.DB) []string {

	filtredFiles := []string{}

	for _, dir := range files {
		if !fileWasSended(db, dir) {
			filtredFiles = append(filtredFiles, dir)
		}
	}

	return filtredFiles
}

func fileWasSended(db *sql.DB, findedFilePath string) bool {
	fileWasSended := false

	// TO DO неоптимальное получение отправленных файлов
	selectSring := fmt.Sprintf("select %s from %s.%s where %s = ?",
		conf.FieldNameDB,
		conf.NameDB,
		conf.TableNameDB,
		conf.FieldNameDB)

	rows, err := db.Query(selectSring, findedFilePath)

	if err != nil {
		return fileWasSended
	}

	if rows.Next() {
		fileWasSended = true
	}

	return fileWasSended
}

func sendfindedFiles(findedFiles []string, db *sql.DB, bot *tgbotapi.BotAPI) {
	for _, findedFilePath := range findedFiles {

		fileBytes, err := os.ReadFile(findedFilePath)
		if err != nil {
			continue
		}

		telegrammFileBytes := tgbotapi.FileBytes{
			Name:  findedFilePath,
			Bytes: fileBytes,
		}

		Document := tgbotapi.NewDocument(conf.ChatID, telegrammFileBytes)
		msg, err := bot.Send(Document)

		if err != nil {
			continue
		}

		insertString := fmt.Sprintf("insert into %s.%s values (?)",
			conf.NameDB, conf.TableNameDB)
		result, err := db.Exec(insertString, findedFilePath)

		if err != nil {
			continue
		}

		fmt.Println(result, msg)
	}
}

type configuration struct {
	Token   string
	ChatID  int64
	Path    string
	Mask    string
	MinDate string
	UserDB  string
	PassDB  string
	// TO DO как можно создать БД с таблицами при старте приложения?
	NameDB      string
	TableNameDB string
	FieldNameDB string
}
