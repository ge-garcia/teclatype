package source

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

type WordsSource struct {
	words []string
	count uint
}

func NewWordsSource(filename string, count uint) WordsSource {
	words, err := readWords(filename)
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	return WordsSource{
		words: words,
		count: count,
	}
}

func (ws WordsSource) Generate() string {
	selectedWords := make([]string, ws.count)

	for i := range selectedWords {
		selectedWords[i] = ws.words[rand.Intn(len(ws.words))]
	}

	return strings.Join(selectedWords, " ")
}

func readWords(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return nil, err
	}

	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return words, nil
}
