package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

func main() {
	port := ":8080"
	gitdir := "/var/www/html/CustomCode"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, this is git utility\nPossible urls:\n/status\n/log\n/pull\n/rollback\n/branch?<branch_name>")
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		cmd := exec.Command("git", "-C", gitdir, "status")
		output, err := cmd.Output()
		//	log.Printf("Output: %s", output)
		fmt.Fprintf(w, "Done\n, %s", output)

		if err != nil {
			log.Printf("Command finished with error: %v", err)
		}
	})

	http.HandleFunc("/pull", func(w http.ResponseWriter, r *http.Request) {
		cmd := exec.Command("git", "-C", gitdir, "pull")
		output, err := cmd.Output()
		fmt.Fprintf(w, "Done\n, %s", output)

		if err != nil {
			log.Printf("Command finished with error: %v", err)
		}
	})

	http.HandleFunc("/rollback", func(w http.ResponseWriter, r *http.Request) {
		cmd := exec.Command("git", "-C", gitdir, "reset", "--hard", "HEAD^")
		output, err := cmd.Output()
		fmt.Fprintf(w, "Done\n, %s", output)

		if err != nil {
			log.Printf("Command finished with error: %v", err)
		}
	})

	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		cmd := exec.Command("git", "-C", gitdir, "log", "--oneline")
		output, err := cmd.Output()
		fmt.Fprintf(w, "Done\n, %s", output)

		if err != nil {
			log.Printf("Command finished with error: %v", err)
		}
	})

	http.HandleFunc("/branch", func(w http.ResponseWriter, r *http.Request) {
		branch := r.URL.RawQuery
		cmd := exec.Command("git", "-C", gitdir, "checkout", branch)
		output, err := cmd.Output()
		fmt.Fprintf(w, "Done, %s", output)

		if err != nil {
			log.Printf("Command finished with error: %v", err)
		}
	})

	http.ListenAndServe(port, nil)
}
