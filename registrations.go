package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type userRegistrationFile struct {
	User  string
	Pisci string
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
