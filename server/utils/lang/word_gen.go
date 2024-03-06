package lang

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
)

const (
	LanguageEnglish = "en"
	LanguageRussian = "ru"
)

func GenerateWords(language string, wordPoolCount int) (wordPool []string, err error) {

	jsonFile, err := os.Open("lang/" + language + ".json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var jsonWordsMap map[string][]string
	json.Unmarshal([]byte(byteValue), &jsonWordsMap)

	if len(jsonWordsMap["words"]) < wordPoolCount {
		return nil, errors.New(fmt.Sprintf("not enough words, expected at least: %d, got: %d", wordPoolCount, len(jsonWordsMap["words"])))
	}

	jsonWords := jsonWordsMap["words"]

	wordPool = GetRandomWords(jsonWords, wordPoolCount)

	return wordPool, nil
}

func GetRandomWords(wordPool []string, wordCount int) []string {
	var winWords []string
	wordPoolTemp := make([]string, len(wordPool)) // Used to store unused win words for loop
	copy(wordPoolTemp, wordPool)
	for i := 0; i < wordCount; i++ {
		randomIndex := rand.Intn(len(wordPoolTemp))
		winWords = append(winWords, wordPoolTemp[randomIndex])

		wordPoolTemp[randomIndex] = wordPoolTemp[len(wordPoolTemp)-1]
		wordPoolTemp = wordPoolTemp[:len(wordPoolTemp)-1]
	}

	return winWords
}
