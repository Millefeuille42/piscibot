package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"os"
)

type userRegistrationFile struct {
	Login  string
	Pisci  string
	Mail   string
	UserID string
}

func checkRegistrationError(session *discordgo.Session, message *discordgo.MessageCreate, args []string) error {
	_, loginExist := os.Stat(fmt.Sprintf("./data/registrations/%s.json", args[1]))
	_, pisciExist := os.Stat(fmt.Sprintf("./data/targets/%s.json", args[2]))
	if loginExist == nil {
		_, _ = session.ChannelMessageSend(message.ChannelID,
			fmt.Sprintf("<@%s> %s is already registered!", message.Author.ID, message.Author.ID))
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
		return err
	}
	file, err := os.Create(fmt.Sprintf("./data/registrations/%s.json", message.Author.ID))
	if err != nil {
		return err
	}
	defer file.Close()
	err = ioutil.WriteFile(fmt.Sprintf("./data/registrations/%s.json", message.Author.ID),
		registrationJson, 0644)
	if err != nil {
		return err
	}
	return nil
}

func registerUser(session *discordgo.Session, message *discordgo.MessageCreate, args []string) error {
	if err := checkRegistrationError(session, message, args); err != nil {
		return err
	}
	if err := createRegistrationFile(message, args); err != nil {
		_, _ = session.ChannelMessageSend(message.ChannelID, err.Error())
		return err
	}
	f, err := os.OpenFile("./data/registrations/userlist.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		_, _ = session.ChannelMessageSend(message.ChannelID, err.Error())
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(fmt.Sprintf("%s\n", message.Author.ID)); err != nil {
		_, _ = session.ChannelMessageSend(message.ChannelID, err.Error())
		return err
	}
	return nil
}

func getPisciList() ([]string, error) {
	UserRegistration := userRegistrationFile{}
	var userList []string

	file, err := os.Open("./data/registrations/userlist.txt")
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		fileData, err := ioutil.ReadFile(fmt.Sprintf("./data/registrations/%s.json", scanner.Text()))
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(fileData, &UserRegistration)
		if err != nil {
			return nil, err
		}
		userList = append(userList, UserRegistration.Pisci)
	}
	_ = file.Close()
	return userList, nil
}
