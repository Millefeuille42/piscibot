package main

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func logError(err error) {
	if err != nil {
		log.Print(err)
	}
}

func Find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func rankMapStringInt(values map[string]int) []string {
	type kv struct {
		Key   string
		Value int
	}

	var ss []kv

	for k, v := range values {
		ss = append(ss, kv{k, v})
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})
	ranked := make([]string, len(values))
	for i, kv := range ss {
		ranked[i] = kv.Key
	}
	return ranked
}

func getHighestProject(userDataParsed UserInfoParsed) string {
	re := regexp.MustCompile("[0-9]+")
	max := make(map[string]int)
	maxP := make(map[string]Project)
	var highestList = ""

	for _, project := range userDataParsed.Projects {
		if project.ProjectStatus == "finished" && project.ProjectMark != 0 {
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
		highestList = fmt.Sprintf("%s\n\t\t                   %-15s%d", highestList, project.ProjectName, project.ProjectMark)
	}
	return highestList
}

func getOngoingProject(userDataParsed UserInfoParsed) string {
	var prList = ""

	for _, project := range userDataParsed.Projects {
		prList = fmt.Sprintf("%s\n\t\t                   %s", prList, project.ProjectName)
	}
	return prList
}
