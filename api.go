package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

func startApi() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api/segbot/", apiTest).Methods("GET")
	router.HandleFunc("/api/segbot/{field}/{user}", apiUsers).Methods("GET")

	go func() {
		log.Println(http.ListenAndServe(":8080", router))
	}()
}

func apiTest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Printf("\nGET '%s'\n\tFrom: %s\n", r.URL.Path, r.RemoteAddr)
}

func apiUsers(w http.ResponseWriter, r *http.Request) {
	userData := UserInfoParsed{}
	vars := mux.Vars(r)

	_, _ = fmt.Printf("\nGET '%s'\n\tFrom: %s\n", r.URL.Path, r.RemoteAddr)

	fileData, err := ioutil.ReadFile(fmt.Sprintf("./data/%s.json", vars["user"]))
	if err != nil {
		_, _ = fmt.Fprintf(w, "Error: User Not Found")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = json.Unmarshal(fileData, &userData)
	if err != nil {
		_, _ = fmt.Fprintf(w, "Error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusFound)
	dataEncoder := json.NewEncoder(w)
	dataEncoder.SetIndent("", "\t")

	if vars["field"] == "all" {
		err = dataEncoder.Encode(userData)
		if err != nil {
			_, _ = fmt.Fprintf(w, "Error")
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		switch vars["field"] {
		case "Location":
			err = dataEncoder.Encode(userData.Location)
		case "Email":
			err = dataEncoder.Encode(userData.Email)
		case "Login":
			err = dataEncoder.Encode(userData.Login)
		case "Level":
			err = dataEncoder.Encode(userData.Level)
		case "Wallet":
			err = dataEncoder.Encode(userData.Wallet)
		case "CorrectionPoint":
			err = dataEncoder.Encode(userData.CorrectionPoint)
		case "Projects":
			err = dataEncoder.Encode(userData.Projects)
		default:
			_, _ = fmt.Fprintf(w, "Error: Bad Request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err != nil {
			_, _ = fmt.Fprintf(w, "Error")
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
