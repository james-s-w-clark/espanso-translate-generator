package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// thanks to https://stackoverflow.com/a/18479916/4261132 for read/write guidance
func main() {
	// ---------------- CONFIG START ----------------
	// https://github.com/pquentin/wiktionary-translations
	filePath, _ := filepath.Abs("./french-english/frwiktionary-20140612-euradicfmt.txt")
	l1 := "en"
	l2 := "fr"
	l1Index := 3
	l2Index := 0
	delimiter := ";"
	// ---------------- CONFIG END ----------------

	file, err := os.Open(filePath)
	if err != nil {
		println(err)
	}
	defer file.Close()

	var configLines []string
	configLines = append(configLines, "matches:")
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, delimiter)
		l1Word := parts[l1Index]
		l2Word := parts[l2Index]
		configLines = append(configLines, pairToEspanso(l1Word, l2Word, l2))
		configLines = append(configLines, pairToEspanso(l2Word, l1Word, l1))
	}

	outputName := fmt.Sprintf("espanso-translate-%s-%s.yml", l1, l2)
	writeLines(configLines, outputName)
}

func pairToEspanso(sourceWord string, targetWord string, targetLanguage string) string {
	return fmt.Sprintf(
		"  - trigger: \"%s:%s\"\n"+
			"    replace: \"%s\"",
		sourceWord, targetLanguage, targetWord)
}

// writeLines writes the lines to the given file.
func writeLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}
