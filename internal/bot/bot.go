package bot

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/a13k551/botMonitorFiles/internal/config"
	"github.com/a13k551/botMonitorFiles/internal/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartBot() {

	conf := config.GetConf()

	db.CreateBase(conf)
	db.CreateTables(conf)

	connection, err := db.CreateConnection(conf)

	if err != nil {
		log.Panic(err)
	}

	defer connection.Close()

	bot, err := tgbotapi.NewBotAPI(conf["token"].(string))

	if err != nil {
		log.Panic(err)
	}

	for {

		if !connection.IsAlive() {
			connection, err = db.CreateConnection(conf)
			if err != nil {
				log.Panic(err)
			}
		}

		findedFiles, err := findFiles(conf)

		if findedFiles == nil || err != nil {
			time.Sleep(time.Second * 5)
			continue
		}

		for _, filepath := range findedFiles {
			if !validDate(filepath, conf) {
				continue
			}
			if db.FileWasSended(connection, filepath, conf) {
				continue
			}
			err = sendFileToChat(filepath, bot, conf)
			if err == nil {
				db.WriteFilepathToDB(filepath, connection, conf)
			}
		}

		time.Sleep(time.Second * 5)
	}
}

func findFiles(conf map[string]interface{}) ([]string, error) {
	findString := fmt.Sprintf("%s%s", conf["path"], conf["mask"])
	findedFiles, err := filepath.Glob(findString)

	if err != nil {
		return nil, err
	}

	if findedFiles == nil {
		return nil, nil
	}

	return findedFiles, nil
}

func validDate(filepath string, conf map[string]interface{}) bool {

	fileStat, err := os.Stat(filepath)
	if err != nil {
		return false
	}

	fileTime := fileStat.ModTime()
	minDate, err := time.Parse("2/1/2006", conf["mindate"].(string))
	if err != nil {
		return false
	}

	if minDate.Before(fileTime) {
		return true
	} else {
		return false
	}
}

func sendFileToChat(findedFilePath string, bot *tgbotapi.BotAPI, conf map[string]interface{}) error {

	fileBytes, err := os.ReadFile(findedFilePath)
	if err != nil {
		log.Println(err)
		return err
	}

	telegrammFileBytes := tgbotapi.FileBytes{
		Name:  findedFilePath,
		Bytes: fileBytes,
	}

	chatidint := conf["chatid"].(int)
	chatidint64 := int64(chatidint)

	Document := tgbotapi.NewDocument(chatidint64, telegrammFileBytes)
	_, err = bot.Send(Document)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
