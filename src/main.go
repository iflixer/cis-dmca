package main

import (
	"dmca/serv"
	"dmca/telegram"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("START")

	log.Println("runtime.GOMAXPROCS:", runtime.GOMAXPROCS(0))

	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Cant load .env: ", err)
	}

	dleApiURL := os.Getenv("DLE_API_MASS_ALIAS_URL")
	dleBasicLogin := os.Getenv("DLE_BASIC_LOGIN")
	telegramReportGroupID, _ := strconv.ParseInt(os.Getenv("TELEGRAM_GROUP_ID"), 10, 64)

	telegramApiToken := os.Getenv("TELEGRAM_APITOKEN")
	if os.Getenv("TELEGRAM_APITOKEN_FILE") != "" {
		telegramApiToken_, err := os.ReadFile(os.Getenv("TELEGRAM_APITOKEN_FILE"))
		if err != nil {
			log.Fatal(err)
		}
		telegramApiToken = strings.TrimSpace(string(telegramApiToken_))
	}

	dleBasicPassword := os.Getenv("DLE_BASIC_PASSWORD")
	if os.Getenv("DLE_BASIC_PASSWORD_FILE") != "" {
		dleBasicPassword_, err := os.ReadFile(os.Getenv("DLE_BASIC_PASSWORD_FILE"))
		if err != nil {
			log.Fatal(err)
		}
		dleBasicPassword = strings.TrimSpace(string(dleBasicPassword_))
	}

	telegramService, err := telegram.NewService(telegramApiToken)
	if err != nil {
		log.Fatal(err)
	}

	telegramService.Send(telegramReportGroupID, fmt.Sprintf("dmca started"))

	httpService, err := serv.NewService("8099", telegramService, dleApiURL, telegramReportGroupID, dleBasicLogin, dleBasicPassword)
	if err != nil {
		log.Fatal(err)
	}
	httpService.Run()
}
