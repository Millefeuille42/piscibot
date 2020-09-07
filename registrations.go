package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"os"
	"strings"
)

type userRegistrationFile struct {
	Login  string
	Pisci  string
	Mail   string
	UserID string
}

func formatRegistrationArgs(args *[]string) {
	(*args)[1] = strings.TrimSpace((*args)[1])
	(*args)[2] = strings.TrimSpace((*args)[2])
	(*args)[3] = strings.TrimSpace((*args)[3])
	(*args)[1] = strings.ToLower((*args)[1])
	(*args)[2] = strings.ToLower((*args)[2])
}

func checkRegistrationError(session *discordgo.Session, message *discordgo.MessageCreate, args []string) error {
	_, loginExist := os.Stat(fmt.Sprintf("./data/registrations/%s.json", message.Author.ID))
	_, pisciExist := os.Stat(fmt.Sprintf("./data/targets/%s.json", args[2]))
	if loginExist == nil {
		_, _ = session.ChannelMessageSend(message.ChannelID,
			fmt.Sprintf("<@%s> You are already registered!", message.Author.ID))
		return os.ErrExist
	}
	if pisciExist == nil {
		_, _ = session.ChannelMessageSend(message.ChannelID,
			fmt.Sprintf("<@%s> %s is already taken!", message.Author.ID, args[2]))
		return os.ErrExist
	}
	return nil
}

func createRegistrationFile(message *discordgo.MessageCreate, args []string) error {
	userRegistration := userRegistrationFile{
		Login:  args[1],
		Pisci:  args[2],
		Mail:   args[3],
		UserID: message.Author.ID,
	}

	registrationJson, err := json.MarshalIndent(userRegistration, "", "\t")
	if err != nil {
		logError(err)
		return err
	}
	file, err := os.Create(fmt.Sprintf("./data/registrations/%s.json", message.Author.ID))
	if err != nil {
		logError(err)
		return err
	}
	defer file.Close()
	err = ioutil.WriteFile(fmt.Sprintf("./data/registrations/%s.json", message.Author.ID),
		registrationJson, 0644)
	if err != nil {
		logError(err)
		return err
	}
	return nil
}

func registerUser(session *discordgo.Session, message *discordgo.MessageCreate, args []string) error {
	formatRegistrationArgs(&args)
	if err := checkRegistrationError(session, message, args); err != nil {
		return err
	}
	if err := createRegistrationFile(message, args); err != nil {
		_, _ = session.ChannelMessageSend(message.ChannelID, err.Error())
		logError(err)
		return err
	}
	f, err := os.OpenFile("./data/registrations/userlist.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		_, _ = session.ChannelMessageSend(message.ChannelID, err.Error())
		logError(err)
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(fmt.Sprintf("%s\n", message.Author.ID)); err != nil {
		_, _ = session.ChannelMessageSend(message.ChannelID, err.Error())
		logError(err)
		return err
	}
	_ = session.GuildMemberRoleAdd(message.GuildID, message.Author.ID, os.Getenv("DISCORDREGISTEREDROLE"))
	_ = session.GuildMemberRoleRemove(message.GuildID, message.Author.ID, os.Getenv("DISCORDUNREGISTEREDROLE"))
	_, _ = session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s> Is now registered", message.Author.ID))
	_ = session.ChannelMessagePin(message.ChannelID, message.ID)
	return nil
}

func getPisciList() ([]string, error) {
	UserRegistration := userRegistrationFile{}
	var userList []string

	lines, _ := parseFileToLines("./data/registrations/userlist.txt")
	for _, line := range lines {
		if line == "" {
			continue
		}
		fileData, err := ioutil.ReadFile(fmt.Sprintf("./data/registrations/%s.json", line))
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(fileData, &UserRegistration)
		if err != nil {
			return nil, err
		}
		userList = append(userList, UserRegistration.Pisci)
	}
	return userList, nil
}

func getPisciMap() (map[string]string, error) {
	UserRegistration := userRegistrationFile{}
	userMap := make(map[string]string)

	lines, _ := parseFileToLines("./data/registrations/userlist.txt")
	for _, line := range lines {
		if line == "" {
			continue
		}
		fileData, err := ioutil.ReadFile(fmt.Sprintf("./data/registrations/%s.json", line))
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(fileData, &UserRegistration)
		if err != nil {
			return nil, err
		}
		userMap[UserRegistration.Pisci] = UserRegistration.Login
	}
	return userMap, nil
}

func getPisciPerID(userID string) (string, error) {
	UserRegistration := userRegistrationFile{}

	fileData, err := ioutil.ReadFile(fmt.Sprintf("./data/registrations/%s.json", userID))
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(fileData, &UserRegistration)
	if err != nil {
		return "", err
	}
	return UserRegistration.Pisci, nil
}
