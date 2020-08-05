package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"strings"
)

func setVarsToMessage(phrase string, project Project, newData, oldData UserInfoParsed) string {
	phrase = strings.Replace(phrase, "#{userName}", newData.Login, -1)
	if project.ProjectName != "null" {
		phrase = strings.Replace(phrase, "#{project}", fmt.Sprintf("%s", project.ProjectName), -1)
		phrase = strings.Replace(phrase, "#{mark}", fmt.Sprintf("%d", project.ProjectMark), -1)
	}
	phrase = strings.Replace(phrase, "#{oldLocation}", oldData.Location, -1)
	phrase = strings.Replace(phrase, "#{proverb}", phrasePicker("conf/proverbs.txt"), -1)
	phrase = strings.Replace(phrase, "#{newLocation}", newData.Location, -1)
	phrase = strings.Replace(phrase, "#{oldLevel}", fmt.Sprintf("%.2f", oldData.Level), -1)
	phrase = strings.Replace(phrase, "#{newLevel}", fmt.Sprintf("%.2f", newData.Level), -1)

	return phrase
}

func announceLocation(param string, newData, oldData UserInfoParsed, session *discordgo.Session) {
	fakeProject := Project{}

	fakeProject.ProjectName = "null"
	switch param {
	case "login":
		message := setVarsToMessage(phrasePicker("conf/login.txt"), fakeProject, newData, oldData)
		fmt.Println(fmt.Sprintf("\t\tSending login for %s, on %s", newData.Login, newData.Location))
		_, err := session.ChannelMessageSend(os.Getenv("DISCORDLOC"), message)
		checkError(err)
	case "logout":
		message := setVarsToMessage(phrasePicker("conf/logout.txt"), fakeProject, newData, oldData)
		fmt.Println(fmt.Sprintf("\t\tSending logout for %s", newData.Login))
		_, err := session.ChannelMessageSend(os.Getenv("DISCORDLOC"), message)
		checkError(err)
	case "newPos":
		message := setVarsToMessage(phrasePicker("conf/newPos.txt"), fakeProject, newData, oldData)
		fmt.Println(fmt.Sprintf("\t\tSending newPos for %s, from %s to %s", newData.Login, oldData.Location, newData.Location))
		_, err := session.ChannelMessageSend(os.Getenv("DISCORDLOC"), message)
		checkError(err)
	}
}

func announceProject(param string, project Project, newData, oldData UserInfoParsed, session *discordgo.Session) {
	switch param {
	case "finished":
		message := setVarsToMessage(phrasePicker("conf/finished.txt"), project, newData, oldData)
		fmt.Println(fmt.Sprintf("\t\tSending finished for %s, on %s", newData.Login, project))
		_, err := session.ChannelMessageSend(os.Getenv("DISCORDPROK"), message)
		checkError(err)
	case "started":
		message := setVarsToMessage(phrasePicker("conf/started.txt"), project, newData, oldData)
		fmt.Println(fmt.Sprintf("\t\tSending started for %s, on %s", newData.Login, project))
		_, err := session.ChannelMessageSend(os.Getenv("DISCORDCHANNEL"), message)
		checkError(err)
	}
}

func messageHandler(session *discordgo.Session, message *discordgo.MessageCreate) {

	botID, err := session.User("@me")
	checkError(err)

	if botID.ID == message.Author.ID {
		return
	}

	if message.Content == "!leaderboard" {
		leaderboard(session, message)
	}

	if strings.HasPrefix(message.Content, "!roadmap") {
		arg := strings.Split(message.Content, "-")
		if len(arg) > 1 {
			roadmap(session, message, arg[1])
		} else {
			roadmap(session, message, "in_progress")
		}
	}

	if strings.HasPrefix(message.Content, "!template") {
		arg := strings.Split(message.Content, "-")
		if len(arg) > 1 {
			template(session, message, arg[1])
		} else {
			template(session, message, "bin")
		}
	}

	if strings.HasPrefix(message.Content, "!user") {
		arg := strings.Split(message.Content, " ")
		if len(arg) > 1 {
			sendUser(session, message, arg[1])
		}
	}

	if strings.HasPrefix(message.Content, "!help") {
		sayHelp(session, message)
	}

	if strings.HasPrefix(message.Content, "!project") {
		arg := strings.Split(message.Content, "-")
		if len(arg) > 1 {
			sayProject(session, message, arg[1])
		}
	}

	if strings.HasPrefix(message.Content, "!location") {
		sayLocation(session, message)
	}

	if strings.HasPrefix(message.Content, "!register") {
		args := strings.Split(message.Content, "-")
		if len(args) == 4 {
			_ = registerUser(session, message, args)
		}
	}
}
