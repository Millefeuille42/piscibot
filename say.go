package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type levelNamePair struct {
	name  string
	level float64
}

func sendUser(session *discordgo.Session, message *discordgo.MessageCreate, user string) {
	if !Find(os.Args, user) {
		return
	}

	userDataParsed := UserInfoParsed{}

	fileData, err := ioutil.ReadFile(fmt.Sprintf("data/%s.json", user))
	checkError(err)
	err = json.Unmarshal(fileData, &userDataParsed)
	checkError(err)

	userMessage := fmt.Sprintf("<@%s>\n"+
		"```"+
		"\n\tEmail:                 %s"+
		"\n\tLocation:              %s"+
		"\n\tCorrection Points:     %d"+
		"\n\tNiveau:                %.2f"+
		"```",
		message.Author.ID,
		userDataParsed.Email,
		userDataParsed.Location,
		userDataParsed.CorrectionPoint,
		userDataParsed.Level,
	)

	_, err = session.ChannelMessageSend(message.ChannelID, userMessage)
	checkError(err)
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

func roadmap(session *discordgo.Session, message *discordgo.MessageCreate, status string) {
	if !Find([]string{"finished", "in_progress"}, status) {
		return
	}

	roadMessage := ""
	userList := os.Args
	projectList := make(map[string]string)
	re := regexp.MustCompile("[0-9]+")
	max := make(map[string]int)

	for _, user := range userList[1:] {
		userDataParsed := UserInfoParsed{}
		fileData, err := ioutil.ReadFile(fmt.Sprintf("data/%s.json", user))
		checkError(err)
		err = json.Unmarshal(fileData, &userDataParsed)
		checkError(err)

		max["Shell"] = 0
		max["Rush"] = 0
		max["C"] = 0
		max["Exam"] = 0

		for _, project := range userDataParsed.Projects {
			if project.ProjectStatus == status {
				cur, _ := strconv.Atoi(re.FindString(project.ProjectName))
				if cur > max[project.ProjectName[:strings.IndexByte(project.ProjectName, ':')]] || status == "in_progress" {
					if _, ok := projectList[project.ProjectName]; !ok {
						projectList[project.ProjectName] = "\n\t| " + user
					} else {
						projectList[project.ProjectName] = fmt.Sprintf("%s\n\t| %s", projectList[project.ProjectName], user)
					}
					max[project.ProjectName[:strings.IndexByte(project.ProjectName, ':')]] = cur
				}
			}
		}
	}

	for projectName, projectUsers := range projectList {
		roadMessage = fmt.Sprintf("%s\n\n%s%10s", roadMessage, projectName, projectUsers)
	}

	roadMessage = fmt.Sprintf("<@%s>, Roadmap for '%s'```%s ```", message.Author.ID, status, roadMessage)
	_, err := session.ChannelMessageSend(message.ChannelID, roadMessage)
	checkError(err)
}

func leaderboard(session *discordgo.Session, message *discordgo.MessageCreate) {
	var leadMessage = ""
	userList := os.Args
	userPair := make([]levelNamePair, 0)
	userDataParsed := UserInfoParsed{}

	for _, user := range userList[1:] {
		fileData, err := ioutil.ReadFile(fmt.Sprintf("data/%s.json", user))
		checkError(err)
		err = json.Unmarshal(fileData, &userDataParsed)
		checkError(err)
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
	checkError(err)
}

func sayHelp(session *discordgo.Session, message *discordgo.MessageCreate) {
	helpMessage := fmt.Sprintf("<@%s>`Read The Fucking Pin`", message.Author.ID)
	session.ChannelMessageSend(message.ChannelID, helpMessage)
}
