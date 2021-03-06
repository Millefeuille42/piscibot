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

func announceIsAccepted(login string, session *discordgo.Session) {
	message := fmt.Sprintf("%s Est pris a 42 !")
	fmt.Println(fmt.Sprintf("\t\tSending login for %s, on %s", login))
	_, err := session.ChannelMessageSend(os.Getenv("DISCORDCHANNEL"), message)
	logError(err)
}

func announceLocation(param string, newData, oldData UserInfoParsed, session *discordgo.Session) {
	fakeProject := Project{}

	fakeProject.ProjectName = "null"
	switch param {
	case "login":
		message := setVarsToMessage(phrasePicker("conf/login.txt"), fakeProject, newData, oldData)
		fmt.Println(fmt.Sprintf("\t\tSending login for %s, on %s", newData.Login, newData.Location))
		_, err := session.ChannelMessageSend(os.Getenv("DISCORDLOC"), message)
		logError(err)
	case "logout":
		message := setVarsToMessage(phrasePicker("conf/logout.txt"), fakeProject, newData, oldData)
		fmt.Println(fmt.Sprintf("\t\tSending logout for %s", newData.Login))
		_, err := session.ChannelMessageSend(os.Getenv("DISCORDLOC"), message)
		logError(err)
	case "newPos":
		message := setVarsToMessage(phrasePicker("conf/newPos.txt"), fakeProject, newData, oldData)
		fmt.Println(fmt.Sprintf("\t\tSending newPos for %s, from %s to %s", newData.Login, oldData.Location, newData.Location))
		_, err := session.ChannelMessageSend(os.Getenv("DISCORDLOC"), message)
		logError(err)
	}
}

func announceProject(param string, project Project, newData, oldData UserInfoParsed, session *discordgo.Session) {
	switch param {
	case "finished":
		message := setVarsToMessage(phrasePicker("conf/finished.txt"), project, newData, oldData)
		fmt.Println(fmt.Sprintf("\t\tSending finished for %s, on %s", newData.Login, project))
		_, err := session.ChannelMessageSend(os.Getenv("DISCORDPROK"), message)
		logError(err)
		_, err = session.ChannelMessageSend(os.Getenv("DISCORDLEADERBOARD"), sLeaderboard())
		logError(err)
	case "started":
		message := setVarsToMessage(phrasePicker("conf/started.txt"), project, newData, oldData)
		fmt.Println(fmt.Sprintf("\t\tSending started for %s, on %s", newData.Login, project))
		_, err := session.ChannelMessageSend(os.Getenv("DISCORDCHANNEL"), message)
		logError(err)
	}
}

func messageHandler(session *discordgo.Session, message *discordgo.MessageCreate) {

	botID, err := session.User("@me")
	if err != nil {
		logError(err)
		return
	}

	if botID.ID == message.Author.ID {
		return
	}

	switch message.Content {
	case "!leaderboard":
		leaderboard(session, message)
	case "!help":
		sayHelp(session, message)
	case "!whoIsStud":
		sayAccepted(session, message)
	case "!location":
		sayLocation(session, message)
	case "!haveAHumongousApparatus":
		haveAHumongousApparatus(session, message)
	case "!unregister":
		everyoneUnregister(session, message)
	case "!spectator":
		everyoneSpectator(session, message)

	default:
		switch {
		case strings.HasPrefix(message.Content, "!roadmap"):
			arg := strings.Split(message.Content, "-")
			if len(arg) > 1 {
				roadmap(session, message, arg[1])
			} else {
				roadmap(session, message, "in_progress")
			}

		case strings.HasPrefix(message.Content, "!project"):
			arg := strings.Split(message.Content, "-")
			if len(arg) > 1 {
				sayProject(session, message, arg[1])
			}
		case strings.HasPrefix(message.Content, "!register"):
			args := strings.Split(message.Content, "|")
			if len(args) == 4 {
				_ = registerUser(session, message, args)
			}

		case strings.HasPrefix(message.Content, "!user"):
			arg := strings.Split(message.Content, " ")
			sendUser(session, message, arg)
		case strings.HasPrefix(message.Content, "!isStud"):
			arg := strings.Split(message.Content, " ")
			isUserAccepted(session, message, arg)
		case strings.HasPrefix(message.Content, "!info"):
			arg := strings.Split(message.Content, " ")
			sendInfo(session, message, arg)
		}
	}
}
