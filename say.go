package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type levelNamePair struct {
	name  string
	level float64
}

func sendInfo(session *discordgo.Session, message *discordgo.MessageCreate, arg []string) {
	userList, _ := getPisciList()

	if len(arg) <= 1 {
		user, err := getPisciPerID(message.Author.ID)
		if err != nil {
			_, err = session.ChannelMessageSend(message.ChannelID, "You are not registered / "+err.Error())
		}
		arg = append(arg, user)
	}

	for _, user := range arg[1:] {
		if !Find(userList, user) {
			return
		}

		userDataParsed := UserInfoParsed{}

		fileData, err := ioutil.ReadFile(fmt.Sprintf("./data/targets/%s.json", user))
		if err != nil {
			return
		}
		err = json.Unmarshal(fileData, &userDataParsed)
		if err != nil {
			return
		}

		userMessage := fmt.Sprintf("<@%s> Info about %s (%s)\n"+
			"```"+
			"\n\tLocation:              %s"+
			"\n\tCorrection Points:     %d"+
			"\n\tLevel:                 %.2f"+
			"```",
			message.Author.ID,
			userDataParsed.Login,
			userDataParsed.Gambler,
			userDataParsed.Location,
			userDataParsed.CorrectionPoint,
			userDataParsed.Level,
		)
		_, _ = session.ChannelMessageSend(message.ChannelID, userMessage)
	}
}

func sendUser(session *discordgo.Session, message *discordgo.MessageCreate, arg []string) {
	userList, _ := getPisciList()

	if len(arg) <= 1 {
		user, err := getPisciPerID(message.Author.ID)
		if err != nil {
			_, err = session.ChannelMessageSend(message.ChannelID, "You are not registered / "+err.Error())
		}
		arg = append(arg, user)
	}

	for _, user := range arg[1:] {
		if !Find(userList, user) {
			return
		}

		userDataParsed := UserInfoParsed{}

		fileData, err := ioutil.ReadFile(fmt.Sprintf("./data/targets/%s.json", user))
		if err != nil {
			return
		}
		err = json.Unmarshal(fileData, &userDataParsed)
		if err != nil {
			return
		}

		userMessage := fmt.Sprintf("<@%s> Info about %s (%s)\n"+
			"```"+
			"\n\tLocation:              %s"+
			"\n\tCorrection Points:     %d"+
			"\n\tLevel:                 %.2f"+
			"\n\tLatest Projects:       %s"+
			"\n\tCurrent Projects:      %s"+
			"```",
			message.Author.ID,
			userDataParsed.Login,
			userDataParsed.Gambler,
			userDataParsed.Location,
			userDataParsed.CorrectionPoint,
			userDataParsed.Level,
			getHighestProject(userDataParsed),
			getOngoingProject(userDataParsed),
		)
		_, _ = session.ChannelMessageSend(message.ChannelID, userMessage)
	}
}

func template(session *discordgo.Session, message *discordgo.MessageCreate, object string) {

	if !Find([]string{"lib", "bin"}, object) {
		return
	}

	file, err := ioutil.ReadFile(fmt.Sprintf("data/templates/%s", object))
	if err != nil {
		log.Print(err)
		return
	}

	_, err = session.ChannelFileSend(message.ChannelID, "Makefile_"+object, bytes.NewReader(file))
	logError(err)
}

func roadmapInP(session *discordgo.Session, message *discordgo.MessageCreate, status string) {

	roadMessage := ""
	userList, _ := getPisciList()
	projectList := make(map[string]string)
	for _, user := range userList {
		userDataParsed := UserInfoParsed{}
		fileData, err := ioutil.ReadFile(fmt.Sprintf("./data/targets/%s.json", user))
		if err != nil {
			continue
		}
		err = json.Unmarshal(fileData, &userDataParsed)
		if err != nil {
			continue
		}
		for _, project := range userDataParsed.Projects {
			if project.ProjectStatus == status {
				if _, ok := projectList[project.ProjectName]; !ok {
					projectList[project.ProjectName] = "\n\t| " + user
				} else {
					projectList[project.ProjectName] = fmt.Sprintf("%s\n\t| %s", projectList[project.ProjectName], user)
				}
			}
		}
	}
	for projectName, projectUsers := range projectList {
		roadMessage = fmt.Sprintf("%s\n\n%s%10s", roadMessage, projectName, projectUsers)
	}
	roadMessage = fmt.Sprintf("<@%s>, Roadmap for '%s'```%s ```", message.Author.ID, status, roadMessage)
	_, err := session.ChannelMessageSend(message.ChannelID, roadMessage)
	if err != nil {
		return
	}
}

func roadmap(session *discordgo.Session, message *discordgo.MessageCreate, status string) {
	if !Find([]string{"finished", "in_progress"}, status) {
		return
	}

	if status == "in_progress" {
		roadmapInP(session, message, status)
		return
	}

	roadMessage := ""
	userList, _ := getPisciList()
	projectList := make(map[string]string)
	re := regexp.MustCompile("[0-9]+")

	for _, user := range userList {
		max := make(map[string]int)
		maxP := make(map[string]Project)
		userDataParsed := UserInfoParsed{}
		fileData, err := ioutil.ReadFile(fmt.Sprintf("./data/targets/%s.json", user))
		if err != nil {
			continue
		}
		err = json.Unmarshal(fileData, &userDataParsed)
		if err != nil {
			continue
		}
		fmt.Printf(user + "\n")
		for _, project := range userDataParsed.Projects {
			if project.ProjectStatus == status && project.ProjectMark != 0 {
				cur := 0
				pName := project.ProjectName
				if strings.Contains(project.ProjectName, " ") {
					pName = project.ProjectName[:strings.IndexByte(project.ProjectName, ' ')]
					cur, _ = strconv.Atoi(re.FindString(project.ProjectName))
				}
				if _, ok := max[pName]; !ok {
					max[pName] = 0
				}
				if cur >= max[pName] {
					maxP[pName] = project
					max[pName] = cur
				}
			}
		}

		for _, project := range maxP {
			if _, ok := projectList[project.ProjectName]; !ok {
				projectList[project.ProjectName] = "\n\t| " + user
			} else {
				projectList[project.ProjectName] = fmt.Sprintf("%s\n\t| %s", projectList[project.ProjectName], user)
			}
		}
	}

	for projectName, projectUsers := range projectList {
		roadMessage = fmt.Sprintf("%s\n\n%s%10s", roadMessage, projectName, projectUsers)
	}

	roadMessage = fmt.Sprintf("<@%s>, Roadmap for '%s'```%s ```", message.Author.ID, status, roadMessage)
	_, err := session.ChannelMessageSend(message.ChannelID, roadMessage)
	if err != nil {
		return
	}
}

func leaderboard(session *discordgo.Session, message *discordgo.MessageCreate) {
	var leadMessage = ""
	userList, _ := getPisciList()
	userPair := make([]levelNamePair, 0)
	userDataParsed := UserInfoParsed{}

	for _, user := range userList {
		fileData, err := ioutil.ReadFile(fmt.Sprintf("./data/targets/%s.json", user))
		if err != nil {
			continue
		}
		err = json.Unmarshal(fileData, &userDataParsed)
		if err != nil {
			continue
		}
		userPair = append(userPair, levelNamePair{userDataParsed.Login, userDataParsed.Level})
	}

	sort.Slice(userPair, func(i, j int) bool {
		return userPair[i].level > userPair[j].level
	})

	for i, user := range userPair[:len(userList)-1] {
		leadMessage = fmt.Sprintf("%s\n%2d: %-15s%.2f", leadMessage, i+1, user.name, user.level)
	}

	leadMessage = fmt.Sprintf("<@%s>```%s```", message.Author.ID, leadMessage)
	_, err := session.ChannelMessageSend(message.ChannelID, leadMessage)
	if err != nil {
		return
	}
}

func sayProject(session *discordgo.Session, message *discordgo.MessageCreate, project string) {
	users, _ := getPisciList()
	prMessage := ""
	prList := make(map[string]int)

	for _, user := range users {
		userDataParsed := UserInfoParsed{}
		fileData, err := ioutil.ReadFile(fmt.Sprintf("./data/targets/%s.json", user))
		if err != nil {
			continue
		}
		err = json.Unmarshal(fileData, &userDataParsed)
		if err != nil {
			continue
		}
		for _, userProject := range userDataParsed.Projects {
			if userProject.ProjectName == project && userProject.ProjectStatus == "finished" {
				prList[user] = userProject.ProjectMark
			}
		}
	}

	for _, index := range rankMapStringInt(prList) {
		prMessage = fmt.Sprintf("%s\n%-15s%d", prMessage, index, prList[index])
	}

	if prMessage == "" {
		prMessage = "Perhaps the archives are incomplete..."
	} else {
		prMessage = fmt.Sprintf("<@%s>, Grades for %s```%s ```", message.Author.ID, project, prMessage)
	}
	_, err := session.ChannelMessageSend(message.ChannelID, prMessage)
	if err != nil {
		return
	}
}

func sayLocation(session *discordgo.Session, message *discordgo.MessageCreate) {
	users, _ := getPisciList()
	locMessage := ""
	nullList := make([]string, 0)

	for _, user := range users {
		userDataParsed := UserInfoParsed{}
		fileData, err := ioutil.ReadFile(fmt.Sprintf("./data/targets/%s.json", user))
		if err != nil {
			continue
		}
		err = json.Unmarshal(fileData, &userDataParsed)
		if err != nil {
			continue
		}

		if userDataParsed.Location != "null" {
			locMessage = fmt.Sprintf("%s\n%-15s%s", locMessage, user, userDataParsed.Location)
		} else {
			nullList = append(nullList, user)
		}
	}

	for _, user := range nullList {
		locMessage = fmt.Sprintf("%s\n%-15sOffline", locMessage, user)
	}
	locMessage = fmt.Sprintf("<@%s>```%s ```", message.Author.ID, locMessage)
	_, err := session.ChannelMessageSend(message.ChannelID, locMessage)
	if err != nil {
		return
	}
}

func sayHelp(session *discordgo.Session, message *discordgo.MessageCreate) {
	helpMessage := fmt.Sprintf("<@%s>`Read The Fucking Pin`", message.Author.ID)
	_, _ = session.ChannelMessageSend(message.ChannelID, helpMessage)
}

func sayAccepted(session *discordgo.Session, message *discordgo.MessageCreate) {

	users, _ := getPisciList()
	accMessage := ""

	for _, user := range users {
		userDataParsed := UserInfoParsed{}
		fileData, err := ioutil.ReadFile(fmt.Sprintf("./data/targets/%s.json", user))
		if err != nil {
			continue
		}
		err = json.Unmarshal(fileData, &userDataParsed)
		if err != nil {
			continue
		}
		if userDataParsed.IsIn {
			accMessage = fmt.Sprintf("%s\n%-15s Accepted", accMessage, user)
		}
	}
	if accMessage == "" {
		accMessage = "Personne n'est pris"
	} else {
		accMessage = fmt.Sprintf("<@%s>```%s ```", message.Author.ID, accMessage)
	}
	_, _ = session.ChannelMessageSend(message.ChannelID, accMessage)
}
