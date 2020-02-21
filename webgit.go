package main

import (
	"bufio"
	"crypto/subtle"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

func readConfig(confName string) map[string]string {
	m := make(map[string]string)

	confFile, err := ioutil.ReadFile(confName)
	if err != nil {
		log.Fatal(err)
		panic("Conf file cannot be read")
	}

	scanner := bufio.NewScanner(strings.NewReader(string(confFile)))
	for scanner.Scan() {
		prop := strings.Split(scanner.Text(), " ")
		m[prop[0]] = prop[1]
	}
	fmt.Println(m)
	return m
}

func main() {
	configMap := readConfig("properties.conf")

	port := configMap["port"]
	gitdir := configMap["gitdir"]
	username := configMap["username"]
	password := configMap["password"]
	msg := "Please enter your username and password for this site"

	handlerRoot := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, this is git utility\nPossible urls:\n/status\n/log\n/pull\n/rollback\n/branch?<branch_name>")
	}

	handlerStatus := func(w http.ResponseWriter, r *http.Request) {
		cmd := exec.Command("git", "-C", gitdir, "status")
		output, err := cmd.Output()
		//	log.Printf("Output: %s", output)
		fmt.Fprintf(w, "Done\n, %s", output)

		if err != nil {
			log.Printf("Command finished with error: %v", err)
		}
	}

	handlerPull := func(w http.ResponseWriter, r *http.Request) {
		cmd := exec.Command("git", "-C", gitdir, "pull")
		output, err := cmd.Output()
		fmt.Fprintf(w, "Done\n, %s", output)

		if err != nil {
			log.Printf("Command finished with error: %v", err)
		}
	}

	handlerRollback := func(w http.ResponseWriter, r *http.Request) {
		cmd := exec.Command("git", "-C", gitdir, "reset", "--hard", "HEAD^")
		output, err := cmd.Output()
		fmt.Fprintf(w, "Done\n, %s", output)

		if err != nil {
			log.Printf("Command finished with error: %v", err)
		}
	}

	handlerLog := func(w http.ResponseWriter, r *http.Request) {
		cmd := exec.Command("git", "-C", gitdir, "log", "--oneline")
		output, err := cmd.Output()
		fmt.Fprintf(w, "Done\n, %s", output)

		if err != nil {
			log.Printf("Command finished with error: %v", err)
		}
	}

	handlerBranch := func(w http.ResponseWriter, r *http.Request) {
		branch := r.URL.RawQuery
		cmd := exec.Command("git", "-C", gitdir, "checkout", branch)
		output, err := cmd.Output()
		fmt.Fprintf(w, "Done, %s", output)

		if err != nil {
			log.Printf("Command finished with error: %v", err)
		}
	}

	http.HandleFunc("/", basicAuth(handlerRoot, username, password, msg))
	http.HandleFunc("/status", basicAuth(handlerStatus, username, password, msg))
	http.HandleFunc("/pull", basicAuth(handlerPull, username, password, msg))
	http.HandleFunc("/rollback", basicAuth(handlerRollback, username, password, msg))
	http.HandleFunc("/log", basicAuth(handlerLog, username, password, msg))
	http.HandleFunc("/branch", basicAuth(handlerBranch, username, password, msg))

	http.ListenAndServe(port, nil)
}

func basicAuth(handler http.HandlerFunc, username, password, realm string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user, pass, ok := r.BasicAuth()

		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised.\n"))
			return
		}

		handler(w, r)
	}
}
