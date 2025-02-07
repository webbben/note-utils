package cleanup

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
)

type TestCase struct {
	Name           string   `json:"name"`
	InputFile      string   `json:"input"`
	ExpectedOutput []string `json:"output"`
}

// go test -run ^TestCleanNoteFast$ github.com/webbben/note-utils/pkg/note-cleanup
func TestCleanNoteFast(t *testing.T) {
	runTests(t, func(s string) (string, error) {
		return CleanNoteWithOpts(s, CleanNoteOpts{Fast: true})
	})
}

// go test -run ^TestCleanNoteCOT$ github.com/webbben/note-utils/pkg/note-cleanup
func TestCleanNoteCOT(t *testing.T) {
	runTests(t, func(s string) (string, error) {
		return CleanNoteCOT(s)
	})
}

func runTests(t *testing.T, funcToTest func(string) (string, error)) {
	var testCases []TestCase
	rawBytes, err := os.ReadFile("tests.json")
	if err != nil {
		log.Fatal("failed to load tests.json")
	}

	err = json.Unmarshal(rawBytes, &testCases)
	if err != nil {
		log.Fatal("failed to unmarshall tests.json data")
	}

	sumAccuracy := float32(0)

	for _, testCase := range testCases {
		filename := testCase.InputFile
		fileBytes, err := os.ReadFile(fmt.Sprintf("tests/%s", filename))
		if err != nil {
			log.Fatal(err)
		}
		fileString := string(fileBytes)

		output, err := funcToTest(fileString)
		if err != nil {
			log.Fatal(err)
		}

		// format string to prepare for comparisons
		output = strings.ToLower(output)
		output = strings.ReplaceAll(output, "*", "")

		fail := 0
		// check contents of output to see if the integrity of original note is intact
		for _, expOut := range testCase.ExpectedOutput {
			if !strings.Contains(output, strings.ToLower(expOut)) {
				log.Printf("Expected string not found: \"%s\"\n", expOut)
				fail++
			}
		}

		total := len(testCase.ExpectedOutput)
		accuracy := float32(total-fail) / float32(total)

		if fail > 0 {
			log.Println(testCase.Name)
			log.Printf("Accuracy: %v (%v/%v)\n", accuracy, total-fail, total)
			log.Printf("\n%s\n", output)
			t.Fail()
		} else {
			log.Println(testCase.Name, "Pass!")
		}

		sumAccuracy += accuracy
	}

	log.Printf("\n\nAverage accuracy: %v\n", sumAccuracy/float32(len(testCases)))
}
