package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// thanks to https://stackoverflow.com/a/18479916/4261132 for read/write guidance
func main() {
	// https://github.com/pquentin/wiktionary-translations
	filePath := "frwiktionary-20140612-euradicfmt.txt"
	l1 := "en"
	l2 := "fr"

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
		parts := strings.Split(line, ";")
		l1Word := parts[3] // en
		l2Word := parts[0] // fr
		configLines = append(configLines, pairToEspanso(l1Word, l2Word, l2))
		configLines = append(configLines, pairToEspanso(l2Word, l1Word, l1))
	}

	writeLines(configLines, "espanso-translate-en-fr.yml")
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
