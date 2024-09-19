package main

import (
	"fmt"
	"math"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/asrul/linux-command-gpt/gpt"
	"github.com/asrul/linux-command-gpt/reader"
	"github.com/atotto/clipboard"
)

const (
	HOST        = "https://api.openai.com/v1/"
	COMPLETIONS = "chat/completions"
	MODEL       = "gpt-4o-mini"
	PROMPT      = "Reply with linux command and nothing else. No need explanation. No need code blocks"

	// This file is created in the user's home directory
	// Example: /home/username/.openai_api_key
	API_KEY_FILE = ".openai_api_key"

	HELP = `

Usage: lcg [options]

  --help        -h  output usage information
  --version     -v  output the version number
  --file        -f  read command from file
  --update-key  -u  update the API key
  --delete-key  -d  delete the API key

Example Usage: lcg I want to extract linux-command-gpt.tar.gz file
Example Usage: lcg --file /path/to/file.json I want to print object questions with jq
  `

	VERSION        = "0.2.1"
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

func getCommand(gpt3 gpt.Gpt3, cmd string) (string, float64) {
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
	elapsed := time.Since(s).Seconds()
	elapsed = math.Round(elapsed*100) / 100

	if r == "" {
		return "", elapsed
	}
	return r, elapsed
}

func main() {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}

	args := os.Args
	cmd := ""
	file := ""
	if len(args) > 1 {
		start := 1
		if args[1] == "--file" || args[1] == "-f" {
			file = args[2]
			start = 3
		}
		cmd = strings.Join(args[start:], " ")
	}

	if file != "" {
		err := reader.FileToPrompt(&cmd, file)
		if err != nil {
			fmt.Println(err)
			return
		}
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
		Temperature:   0.01,
	}

	if h == CMD_UPDATE {
		gpt3.UpdateKey()
		return
	}

	if h == CMD_DELETE {
		gpt3.DeleteKey()
		return
	}

	c := "R"
	r := ""
	elapsed := 0.0
	for c == "R" || c == "r" {
		r, elapsed = getCommand(gpt3, cmd)
		c = "N"
		fmt.Printf("Completed in %v seconds\n\n", elapsed)
		fmt.Println(r)
		fmt.Print("\nDo you want to (c)opy, (r)egenerate, or take (N)o action on the command? (c/r/N): ")
		fmt.Scanln(&c)

		// No action
		if c == "N" || c == "n" {
			return
		}
	}

	if r == "" {
		return
	}

	// Copy to clipboard
	if c == "C" || c == "c" {
		clipboard.WriteAll(r)
		fmt.Println("\033[33mCopied to clipboard")
		return
	}
}
