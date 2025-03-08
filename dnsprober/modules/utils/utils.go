package utils

import (
	"os"
	"strings"

	"golang.org/x/exp/rand"
)

func Deduplslice[T comparable](input []T) []T {
	if len(input) == 0 {
		return input
	}

	seen := make(map[T]bool)
	var unique []T

	for _, item := range input {
		if !seen[item] {
			seen[item] = true
			unique = append(unique, item)
		}
	}
	return unique
}

func GenerateRandomWords(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	word := make([]rune, n)
	for i := range word {
		word[i] = letters[rand.Intn(len(letters))]
	}
	return string(word)
}

func Set(urls []string) []string {
	sets := make(map[string]bool)
	var results []string

	for _, url := range urls {
		if _, ok := sets[url]; !ok {
			sets[url] = true
			results = append(results, url)
		}
	}
	return results
}

func IsStdin() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func SplitStrings(words string) []string {
	return strings.Split(words, ",")
}
