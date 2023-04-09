package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"

	"github.com/asrul/linux-command-gpt/gpt"
)

const (
	HOST        = "https://api.openai.com/v1/"
	COMPLETIONS = "chat/completions"
	MODEL       = "gpt-3.5-turbo"
	PROMPT      = "I want you to reply with linux command and nothing else. Do not write explanations."

	// This file is created in the user's home directory
	// Example: /home/username/.openai_api_key
	API_KEY_FILE = ".openai_api_key"

	HELP = `

Usage: lcg [options]

  --help         output usage information
  --version      output the version number
  --update-key   update the API key
  --delete-key   delete the API key

  `

	VERSION        = "0.1.0"
	CMD_HELP       = 100
	CMD_VERSION    = 101
	CMD_UPDATE     = 102
	CMD_DELETE     = 103
	CMD_COMPLETION = 110
)

func handleCommand(cmd string) int {
	if cmd == "" || cmd == "--help" || cmd == "-h" {
		return CMD_HELP
	}
	if cmd == "--version" || cmd == "-v" {
		return CMD_VERSION
	}
	if cmd == "--update-key" || cmd == "-u" {
		return CMD_UPDATE
	}
	if cmd == "--delete-key" || cmd == "-d" {
		return CMD_DELETE
	}
	return CMD_COMPLETION
}

func main() {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}

	args := os.Args
	cmd := ""
	if len(args) > 1 {
		cmd = strings.Join(args[1:], " ")
	}
	h := handleCommand(cmd)

	if h == CMD_HELP {
		fmt.Println(HELP)
		return
	}

	if h == CMD_VERSION {
		fmt.Println(VERSION)
		return
	}

	gpt3 := gpt.Gpt3{
		CompletionUrl: HOST + COMPLETIONS,
		Model:         MODEL,
		Prompt:        PROMPT,
		HomeDir:       currentUser.HomeDir,
		ApiKeyFile:    API_KEY_FILE,
	}

	if h == CMD_UPDATE {
		gpt3.UpdateKey()
		return
	}

	if h == CMD_DELETE {
		gpt3.DeleteKey()
		return
	}

	gpt3.InitKey()
	s := time.Now()
	done := make(chan bool)
	go func() {
		loadingChars := []rune{'-', '\\', '|', '/'}
		i := 0
		for {
			select {
			case <-done:
				fmt.Printf("\r")
				return
			default:
				fmt.Printf("\rLoading %c", loadingChars[i])
				i = (i + 1) % len(loadingChars)
				time.Sleep(30 * time.Millisecond)
			}
		}
	}()

	r := gpt3.Completions(cmd)
	done <- true
	if r == "" {
		return
	}

	c := "Y"
	elapsed := time.Since(s).Seconds()
	elapsed = math.Round(elapsed*100) / 100
	fmt.Printf("Completed in %v seconds\n\n", elapsed)
	fmt.Println(r)
	fmt.Print("\nAre you sure you want to execute the command? (Y/n): ")
	fmt.Scanln(&c)
	if c != "Y" && c != "y" {
		return
	}

	cmsplit := strings.Split(r, " ")
	cm := exec.Command(cmsplit[0], cmsplit[1:]...)
	out, err := cm.Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(out))
}
