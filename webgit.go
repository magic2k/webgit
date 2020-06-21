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
	"time"
)

var gitdir string

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

func handlerRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, this is git utility\nPossible urls:\n/status\n/log\n/pull\n/rollback\n/branch?<branch_name>")
}

func handlerStatus(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("git", "-C", gitdir, "status")
	output, err := cmd.Output()
	//	log.Printf("Output: %s", output)
	fmt.Fprintf(w, "Done\n, %s", output)

	if err != nil {
		log.Printf("Command finished with error: %v", err)
	}
}

func handlerPull(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("git", "-C", gitdir, "pull")
	output, err := cmd.Output()
	fmt.Fprintf(w, "Done\n, %s", output)

	if err != nil {
		log.Printf("Command finished with error: %v", err)
	}
}

func handlerRollback(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("git", "-C", gitdir, "reset", "--hard", "HEAD^")
	output, err := cmd.Output()
	fmt.Fprintf(w, "Done\n, %s", output)

	if err != nil {
		log.Printf("Command finished with error: %v", err)
	}
}

func handlerLog(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("git", "-C", gitdir, "log", "--oneline")
	output, err := cmd.Output()
	fmt.Fprintf(w, "Done\n, %s", output)

	if err != nil {
		log.Printf("Command finished with error: %v", err)
	}
}

func handlerBranch(w http.ResponseWriter, r *http.Request) {
	branch := r.URL.RawQuery
	cmd := exec.Command("git", "-C", gitdir, "checkout", branch)
	output, err := cmd.Output()
	fmt.Fprintf(w, "Done, %s", output)

	if err != nil {
		log.Printf("Command finished with error: %v", err)
	}
}

func main() {
	configMap := readConfig("properties.conf")

	port := configMap["port"]
	gitdir = configMap["gitdir"]
	username := configMap["username"]
	password := configMap["password"]
	msg := "Please enter your credentials for this site"

	siteMux := http.NewServeMux()
	siteMux.HandleFunc("/", handlerRoot)
	siteMux.HandleFunc("/status", handlerStatus)
	siteMux.HandleFunc("/pull", handlerPull)
	siteMux.HandleFunc("/rollback", handlerRollback)
	siteMux.HandleFunc("/log", handlerLog)
	siteMux.HandleFunc("/branch", handlerBranch)

	siteHandler := accessLogMiddleware(siteMux)
	siteHandler = basicAuth(siteHandler, username, password, msg)
	siteHandler = errorMiddleware(siteHandler)
	http.ListenAndServe(port, siteHandler)
}

func accessLogMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		fmt.Printf("[%s] %s, %s %s, %s\n",
			r.Method, r.RemoteAddr, r.URL.Path, time.Since(start), time.Now())
	})
}

func errorMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// var err error
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recovered from panic: %v", err)
			}
		}()
		h.ServeHTTP(w, r)
	})
}

func basicAuth(h http.Handler, username, password, realm string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user, pass, ok := r.BasicAuth()

		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised.\n"))
			return
		}

		h.ServeHTTP(w, r)
	})
}
