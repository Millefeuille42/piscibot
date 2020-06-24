package main

type Project struct {
	ProjectName   string
	ProjectStatus string
	ProjectMark   int
}

type UserInfoParsed struct {
	Login           string
	Email           string
	Location        string
	CorrectionPoint int
	Level           float64
	Projects        map[string]Project
}

func processUserInfo(userData UserInfo) (UserInfoParsed, error) {

	project := Project{}
	userDataParsed := UserInfoParsed{}

	userDataParsed.Login = userData.Login
	userDataParsed.Email = userData.Email
	userDataParsed.CorrectionPoint = userData.CorrectionPoint

	userDataParsed.Location = userData.Location
	if userData.Location == "" {
		userDataParsed.Location = "null"
	}

	for _, cursus := range userData.CursusUsers {
		userDataParsed.Level = cursus.Level
	}

	userDataParsed.Projects = make(map[string]Project)

	for _, projectRaw := range userData.ProjectsUsers {
		project.ProjectName = projectRaw.Project.Name
		project.ProjectStatus = projectRaw.Status
		if projectRaw.FinalMark == nil {
			projectRaw.FinalMark = 0
		}
		project.ProjectMark = (projectRaw.FinalMark).(int)
		userDataParsed.Projects[projectRaw.Project.Name] = project
	}

	return userDataParsed, nil
}
