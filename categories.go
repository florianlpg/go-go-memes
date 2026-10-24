package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"
)

var categoryFiles embed.FS

type Category struct {
	ID          string         `json:"id"`
	Language    string         `json:"language"`
	Name        string         `json:"name"`
	Keywords    map[string]int `json:"keywords"`
	SearchTerms []string       `json:"search_terms"`
}

func loadCategoryFromFile(path string) (Category, error) {
	data, err := categoryFiles.ReadFile(path)
	if err != nil {
		return Category{}, err
	}

	var category Category

	if err := json.Unmarshal(data, &category); err != nil {
		return Category{}, fmt.Errorf(
			"invalid JSON in %s: %w",
			path,
			err,
		)
	}

	return category, nil
}

func loadCategories(language string) ([]Category, error) {
	if language != "en" && language != "fr" {
		return nil, fmt.Errorf(
			"unsupported language: %s",
			language,
		)
	}

	directory := "categories/" + language

	entries, err := categoryFiles.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	categories := make([]Category, 0)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		path := directory + "/" + entry.Name()

		category, err := loadCategoryFromFile(path)
		if err != nil {
			return nil, err
		}

		if category.Language != language {
			return nil, fmt.Errorf(
				"language mismatch in %s",
				path,
			)
		}

		categories = append(categories, category)
	}

	if len(categories) == 0 {
		return nil, fmt.Errorf(
			"no categories found for %s",
			language,
		)
	}

	return categories, nil
}
