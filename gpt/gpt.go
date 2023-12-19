package gpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Gpt3 struct {
	CompletionUrl string
	Prompt        string
	Model         string
	HomeDir       string
	ApiKeyFile    string
	ApiKey        string
}

type Chat struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Gpt3Request struct {
	Model    string `json:"model"`
	Messages []Chat `json:"messages"`
}

type Gpt3Response struct {
	Choices []struct {
		Message Chat `json:"message"`
	} `json:"choices"`
}

func (gpt3 *Gpt3) deleteApiKey() {
	filePath := gpt3.HomeDir + string(filepath.Separator) + gpt3.ApiKeyFile
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return
	}
	err := os.Remove(filePath)
	if err != nil {
		panic(err)
	}
}

func (gpt3 *Gpt3) updateApiKey(apiKey string) {
	filePath := gpt3.HomeDir + string(filepath.Separator) + gpt3.ApiKeyFile
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	apiKey = strings.TrimSpace(apiKey)
	_, err = file.WriteString(apiKey)
	if err != nil {
		panic(err)
	}
	gpt3.ApiKey = apiKey
}

func (gpt3 *Gpt3) storeApiKey(apiKey string) {
	if apiKey == "" {
		return
	}
	filePath := gpt3.HomeDir + string(filepath.Separator) + gpt3.ApiKeyFile
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	apiKey = strings.TrimSpace(apiKey)
	_, err = file.WriteString(apiKey)
	if err != nil {
		panic(err)
	}
	gpt3.ApiKey = apiKey
}

func (gpt3 *Gpt3) loadApiKey() bool {
	dirSeparator := string(filepath.Separator)
	apiKeyFile := gpt3.HomeDir + dirSeparator + gpt3.ApiKeyFile
	if _, err := os.Stat(apiKeyFile); os.IsNotExist(err) {
		return false
	}
	apiKey, err := os.ReadFile(apiKeyFile)
	if err != nil {
		return false
	}
	gpt3.ApiKey = string(apiKey)

	return true
}

func (gpt3 *Gpt3) UpdateKey() {
	var apiKey string
	fmt.Print("OpenAI API Key: ")
	fmt.Scanln(&apiKey)
	gpt3.updateApiKey(apiKey)
}

func (gpt3 *Gpt3) DeleteKey() {
	var c string
	fmt.Print("Are you sure you want to delete the API key? (y/N): ")
	fmt.Scanln(&c)
	if c == "Y" || c == "y" {
		gpt3.deleteApiKey()
	}
}

func (gpt3 *Gpt3) InitKey() {
	load := gpt3.loadApiKey()
	if load {
		return
	}
	var apiKey string
	fmt.Print("OpenAI API Key: ")
	fmt.Scanln(&apiKey)
	gpt3.storeApiKey(apiKey)
}

func (gpt3 *Gpt3) Completions(ask string) string {
	req, err := http.NewRequest("POST", gpt3.CompletionUrl, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(gpt3.ApiKey))

	messages := []Chat{
		{"system", gpt3.Prompt},
		{"user", ask},
	}
	payload := Gpt3Request{
		gpt3.Model,
		messages,
	}
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	req.Body = io.NopCloser(bytes.NewBuffer(payloadJson))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println(string(body))
		return ""
	}

	var res Gpt3Response
	err = json.Unmarshal(body, &res)
	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(res.Choices[0].Message.Content)
}
