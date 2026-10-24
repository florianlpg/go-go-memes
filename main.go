package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Describe your workplace situation: ")
	sentence, err := reader.ReadString('\n')

	if err != nil && err != io.EOF {
		fmt.Println("Error reading input:", err)
		return
	}

	if strings.TrimSpace(sentence) == "" {
		fmt.Println("Please enter a sentence.")
		return
	}

	englishCategories, err := loadCategories("en")
	if err != nil {
		fmt.Println("Error loading English categories:", err)
		return
	}

	frenchCategories, err := loadCategories("fr")
	if err != nil {
		fmt.Println("Error loading French categories:", err)
		return
	}

	language := detectLanguage(
		sentence,
		englishCategories,
		frenchCategories,
	)

	fmt.Println("Detected language:", language)

	categories := englishCategories
	if language == "fr" {
		categories = frenchCategories
	}

	words := tokenize(sentence)
	scores := make(map[string]int)

	for _, category := range categories {
		scores[category.Name] = calculateScore(category, words)
	}

	bestCategory, bestScore := findBestScore(scores)

	if bestScore == 0 {
		fmt.Println("No matching category found.")
		return
	}

	fmt.Printf(
		"Best Category: %s (Score: %d)\n",
		bestCategory,
		bestScore,
	)
}
