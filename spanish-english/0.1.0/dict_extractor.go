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
	// https://github.com/open-dsl-dict/wiktionary-dict
	filePath, _ := filepath.Abs("./spanish-english/en-es-enwiktionary.txt")
	l1 := "en"
	l2 := "es"
	l1Index := 0
	l2Index := 1
	delimiter := " :: "
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
		if strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, delimiter)
		if len(parts) == 1 {
			continue
		}
		l1Word := parts[l1Index]
		if strings.Contains(l1Word, "{") {
			l1Word = strings.Split(l1Word, "{")[0]
		}
		l1Word = strings.TrimSpace(l1Word)

		l2Word := parts[l2Index]
		l2Word = strings.ReplaceAll(l2Word, "\"", "'")
		l2Word = strings.ReplaceAll(l2Word, "\\", "") // Espanso isn't happy with this character, \
		configLines = append(configLines, pairToEspanso(l1Word, l2Word, l2))
	}

	outputName := fmt.Sprintf("espanso-translate-%s-%s.yml", l1, l2)
	_ = writeLines(configLines, outputName)
}

func pairToEspanso(sourceWord string, targetWord string, targetLanguage string) string {
	return fmt.Sprintf(
		"  - trigger: \"%s:%s\"\n"+
			"    replace: \"%s\"",
		sourceWord, targetLanguage, targetWord)
}

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
