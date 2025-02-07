package util

import (
	"bufio"
	"bytes"
	"log"
	"os"
)

func ReadFile(filepath string) (string, error) {
	fileBytes, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return string(fileBytes), nil
}

func ReadStdin() (string, error) {
	var buffer bytes.Buffer
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		buffer.WriteString(scanner.Text() + "\n")
	}

	return buffer.String(), scanner.Err()
}

func IsStdinPipe() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal("failed to check stdin?", err)
	}
	return (stat.Mode() & os.ModeCharDevice) == 0
}
