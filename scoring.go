package main

import (
	"strings"
	"unicode"
)

func tokenize(sentence string) []string {
	return strings.FieldsFunc(
		strings.ToLower(sentence),
		func(r rune) bool {
			return unicode.IsSpace(r) ||
				unicode.IsPunct(r) ||
				unicode.IsSymbol(r)
		},
	)
}

func calculateScore(category Category, words []string) int {
	totalScore := 0

	normalized := " " + strings.Join(words, " ") + " "

	for keyword, score := range category.Keywords {
		cleanKeyword := strings.Join(tokenize(keyword), " ")

		if cleanKeyword == "" {
			continue
		}

		if strings.Contains(
			normalized,
			" "+cleanKeyword+" ",
		) {
			totalScore += score
		}
	}

	return totalScore
}

func findBestScore(scores map[string]int) (string, int) {
	bestCategory := ""
	bestScore := 0

	for category, score := range scores {
		if score > bestScore ||
			(score > 0 &&
				score == bestScore &&
				(bestCategory == "" || category < bestCategory)) {
			bestScore = score
			bestCategory = category
		}
	}

	return bestCategory, bestScore
}

func detectLanguage(
	sentence string,
	englishCategories []Category,
	frenchCategories []Category,
) string {

	words := tokenize(sentence)

	englishHints := map[string]int{
		"i": 4, "my": 3, "the": 3,
		"this": 3, "that": 2, "is": 3,
		"are": 2, "was": 2, "were": 2,
		"with": 3, "and": 2, "but": 2,
		"for": 2, "of": 2, "have": 2,
		"not": 3, "work": 3, "boss": 3,
		"meeting": 3, "today": 2,
		"another": 3, "why": 3,
		"what": 3, "how": 3, "it": 2,
	}

	frenchHints := map[string]int{
		"je": 4, "tu": 3, "il": 2,
		"elle": 2, "nous": 4, "vous": 4,
		"le": 3, "la": 3, "les": 3,
		"des": 3, "une": 3, "un": 2,
		"du": 3, "de": 2, "pour": 2,
		"avec": 3, "dans": 3, "est": 3,
		"suis": 4, "sont": 3, "pas": 4,
		"que": 2, "qui": 2, "mon": 3,
		"ma": 3, "mes": 3, "au": 3,
		"aux": 3, "ça": 4, "ai": 3,
		"encore": 3, "trop": 4,
		"réunion": 3, "boulot": 4,
		"bonjour": 4,
	}

	englishScore := 0
	frenchScore := 0

	for _, word := range words {
		englishScore += englishHints[word]
		frenchScore += frenchHints[word]

		if strings.ContainsAny(
			word,
			"àâäçéèêëîïôùûüÿœ",
		) {
			frenchScore += 4
		}
	}

	for _, category := range englishCategories {
		englishScore += calculateScore(category, words)
	}

	for _, category := range frenchCategories {
		frenchScore += calculateScore(category, words)
	}

	if frenchScore > englishScore {
		return "fr"
	}

	return "en"
}
