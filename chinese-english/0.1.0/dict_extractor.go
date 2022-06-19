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
	filePath, _ := filepath.Abs("./chinese-english/cedict_ts.u8.txt")
	l1 := "en"
	l2 := "zh"
	//l1Index := 3
	//l2Index := 0
	delimiter := "/"
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
		charsAndPinyin := strings.Split(parts[0], "[") // "\\["
		chars := strings.Split(charsAndPinyin[0], " ")
		traditional := chars[0]
		simplified := chars[1]
		pinyin := strings.ToLower(
			strings.ReplaceAll(charsAndPinyin[1], "] ", ""))
		definition := parts[1]

		configLines = append(configLines, chineseToEnglish(traditional, pinyin, definition))
		if traditional != simplified { // don't write duplicate triggers
			configLines = append(configLines, chineseToEnglish(simplified, pinyin, definition))
		}

		// ignore lengthy English definitions (anything with a space)
		// it's unlikely users will be typing long definition sentences
		// TODO optimise and break the trigger pattern a bit? trad/simp same -> match on just :zh ?
		if !strings.Contains(definition, " ") {
			configLines = append(configLines, englishToChinese(definition, traditional, pinyin, "t"))
			configLines = append(configLines, englishToChinese(definition, simplified, pinyin, "s"))
		}
	}

	outputName := fmt.Sprintf("espanso-translate-%s-%s.yml", l1, l2)
	writeLines(configLines, outputName)
}

func chineseToEnglish(character string, pinyin string, definiton string) string {
	return fmt.Sprintf(
		"  - trigger: \"%s:en\"\n"+
			"    replace: \"{%s}(%s)[%s]\"",
		character,
		character, pinyin, definiton)
}

func englishToChinese(definiton string, character string, pinyin string, simpOrTrad string) string {
	return fmt.Sprintf(
		"  - trigger: \"%s:zh-%s\"\n"+
			"    replace: \"{%s}(%s)[%s]\"",
		definiton, simpOrTrad,
		character, pinyin, definiton)
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
