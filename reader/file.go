package reader

import (
	"bufio"
	"os"
)

func FileToPrompt(cmd *string, filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	*cmd = *cmd + "\nFile path: " + filePath + "\n"
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		*cmd = *cmd + "\n" + line
	}
	return nil
}
