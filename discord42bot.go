package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func writeUsers(api Api42, session *discordgo.Session) Api42 {
	userList, err := getPisciList()
	userMap, err := getPisciMap()
	if err != nil {
		return api
	}

	for _, user := range userList {
		userData := UserInfo{}
		userDataParsed := UserInfoParsed{}
		var err error

		userData, api.Token, err = getUserInfo(user, api.Token, userData)
		if err != nil {
			logError(err)
			continue
		}
		fmt.Println(fmt.Sprintf("Request:\n\tGot raw data from %s", user))
		userDataParsed, err = processUserInfo(userData)
		if err != nil {
			logError(err)
			continue
		}
		fmt.Println("\tProcessed raw data")
		userDataParsed.Gambler = userMap[user]
		err = checkUserFile(user, userDataParsed, session)
		if err != nil {
			logError(err)
			continue
		}
		staticDataToDB(user)
		time.Sleep(3000 * time.Millisecond)
	}
	return api
}

func main() {
	api := Api42{}

	err := godotenv.Load("dev.env")
	checkError(err)
	fmt.Println("Started init")

	err = api.Token.getToken()
	checkError(err)
	fmt.Println("42 Token acquired")
	fmt.Println("Expires in:", api.Token.ExpiresIn)

	discordBot, err := discordgo.New("Bot " + os.Getenv("BOTTOKEN"))
	checkError(err)
	fmt.Println("Discord bot created")

	discordBot.AddHandler(messageHandler)

	err = discordBot.Open()
	checkError(err)
	fmt.Println("Discord Bot up and running")

	setupCloseHandler(discordBot)

	go func() {
		userList, err := getPisciList()
		if err != nil {
			return
		}
		for {
			for _, user := range userList {
				userDataToDB(user)
			}
			time.Sleep(6 * time.Hour)
		}
	}()

	for {
		api = writeUsers(api, discordBot)
	}
}

func setupCloseHandler(session *discordgo.Session) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		time.Sleep(2 * time.Second)
		_ = session.Close()
		os.Exit(0)
	}()
}
