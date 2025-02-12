package util

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
)

func ReadFileOrStdin(filepath string) (string, error) {
	if filepath != "" {
		return ReadFile(filepath)
	}

	return ReadStdin()
}

func ReadFile(filepath string) (string, error) {
	fileBytes, err := os.ReadFile(filepath)
	if err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}
	return string(fileBytes), nil
}

func ReadStdin() (string, error) {
	var buffer bytes.Buffer
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		buffer.WriteString(scanner.Text() + "\n")
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading from stdin: %w", err)
	}

	return buffer.String(), nil
}

func IsStdinPipe() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal("failed to check stdin?", err)
	}
	return (stat.Mode() & os.ModeCharDevice) == 0
}
