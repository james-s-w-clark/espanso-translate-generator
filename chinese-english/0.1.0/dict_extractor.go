package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// thanks to https://stackoverflow.com/a/18479916/4261132 for read/write guidance
func main() {
	// ---------------- CONFIG START ----------------
	// https://github.com/pquentin/wiktionary-translations
	filePath, _ := filepath.Abs("./chinese-english/cedict_ts.u8.txt")
	l1 := "en"
	l2 := "zh"
	delimiter := "/"
	// ---------------- CONFIG END ----------------

	mostFrequentWords := getNMostFrequentWords(10_000)
	print(mostFrequentWords)

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
		pinyin := accentPinyinTones(
			strings.ToLower(
				strings.ReplaceAll(charsAndPinyin[1], "] ", "")))
		definition := parts[1]

		if !mostFrequentWords[simplified] {
			continue
		}

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

func getNMostFrequentWords(count int) map[string]bool {
	filePath, _ := filepath.Abs("./chinese-english/global_wordfreq.release_UTF-8.txt")
	file, err := os.Open(filePath)
	if err != nil {
		println(err)
	}
	defer file.Close()

	index := 0
	words := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		index++
		if index > count {
			break
		}
		line := scanner.Text()
		word := strings.Split(line, "\t")[0]
		words[word] = true
	}

	return words
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
		"  - trigger: \"%s:zh%s\"\n"+
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

func accentPinyinTones(pinyin string) string {
	syllables := strings.Split(pinyin, " ")
	var accented []string

	for _, syllable := range syllables {
		tone, _ := strconv.ParseInt(syllable[len(syllable)-1:], 0, 32)
		if tone == 0 {
			accented = append(accented, syllable) // e.g. syllable is "P", no tone mark
			continue
		}

		accented = append(accented, accentSyllable(syllable, tone))
	}

	return strings.Join(accented, " ")
}

func accentSyllable(syllable string, tone int64) string {
	tmp1 := strings.ReplaceAll(syllable, "u:", "ü")
	tmp2 := strings.ReplaceAll(tmp1, "U:", "Ü")
	tmp3 := tmp2[:len(tmp2)-1] // we can ignore the number at the end now that we've extracted it

	if strings.Index(syllable, "r") == 0 &&
		strings.LastIndex(syllable, "r") == 0 {
		return "er"
	}

	if strings.Contains(syllable, "iu") {
		return strings.ReplaceAll(tmp3, "u", toneMap["u"][tone-1])
	}

	for _, vowel := range vowels {
		if strings.Contains(syllable, vowel) {
			return strings.ReplaceAll(tmp3, vowel, toneMap[vowel][tone-1])
		}
	}

	return syllable
}

var toneMap = map[string][]string{
	"a": {"ā", "á", "ǎ", "à", "a"},
	"e": {"ē", "é", "ě", "è", "e"},
	"i": {"ī", "í", "ǐ", "ì", "i"},
	"o": {"ē", "é", "ě", "è", "e"},
	"u": {"ū", "ú", "ǔ", "ù", "u"},
	"ü": {"ǖ", "ǘ", "ǚ", "ǜ", "ü"},
}

var vowels = []string{"a", "o", "e", "i", "u", "ü"}
