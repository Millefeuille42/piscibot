package main

import "strings"

type Project struct {
	ProjectName   string
	ProjectStatus string
	ProjectMark   int
}

type UserInfoParsed struct {
	Gambler         string
	Login           string
	Email           string
	Location        string
	CorrectionPoint int
	Level           float64
	Projects        map[string]Project
	IsIn            bool
}

func processUserInfo(userData UserInfo) (UserInfoParsed, error) {

	project := Project{}
	userDataParsed := UserInfoParsed{}

	userDataParsed.IsIn = false
	userDataParsed.Login = userData.Login
	userDataParsed.Email = userData.Email
	userDataParsed.CorrectionPoint = userData.CorrectionPoint

	userDataParsed.Location = userData.Location
	if userData.Location == "" {
		userDataParsed.Location = "null"
	}

	for _, cursus := range userData.CursusUsers {
		userDataParsed.Level = cursus.Level
		if cursus.CursusID == 21 {
			userDataParsed.IsIn = true
		}
	}

	userDataParsed.Projects = make(map[string]Project)

	for _, projectRaw := range userData.ProjectsUsers {
		project.ProjectName = strings.Replace(projectRaw.Project.Name, "C Piscine ", "", -1)
		project.ProjectStatus = projectRaw.Status
		if projectRaw.FinalMark == nil {
			projectRaw.FinalMark = 0
		}

		switch projectRaw.FinalMark.(type) {
		case int:
			project.ProjectMark = projectRaw.FinalMark.(int)
		case float64:
			project.ProjectMark = int(projectRaw.FinalMark.(float64))
		default:
			projectRaw.FinalMark = 0
		}

		userDataParsed.Projects[project.ProjectName] = project
	}

	return userDataParsed, nil
}
