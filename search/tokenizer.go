package search

import (
	"strings"
	"unicode"

	snowballeng "github.com/kljensen/snowball/english"
)

// analyze performs a series of operations on the given text to prepare it for indexing.
// It first tokenizes the text into words.
// Then it converts all words to lowercase.
// After that, it removes all stop words from the words.
// Finally, it applies stemming to the words.
// It returns a slice of the processed words.
//
// Parameters:
//
//	text: The string to analyze.
//
// Returns:
//
//	[]string: A slice of processed words.
func analyze(text string) []string {
	tokens := tokenize(text)
	tokens = lowercaseFilter(tokens)
	tokens = stopWordFilter(tokens)
	tokens = stemFilter(tokens)
	return tokens
}

// tokenize splits the given text into words.
// It uses the strings.FieldsFunc function with a function that returns true for any rune that is not a letter or a number.
// This means that it splits the text at any character that is not a letter or a number.
// It returns a slice of words.
//
// Parameters:
//
//	text: The string to tokenize.
//
// Returns:
//
//	[]string: A slice of words.
func tokenize(text string) []string {
	// split the text into words
	// return a slice of words
	return strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

// lowercaseFilter converts all tokens in the given slice to lowercase.
// It iterates over the tokens and for each one, it converts it to lowercase.
// It returns a new slice with the lowercase tokens.
//
// Parameters:
//
//	tokens: A slice of string tokens to convert to lowercase.
//
// Returns:
//
//	[]string: A slice of lowercase tokens.
func lowercaseFilter(tokens []string) []string {
	// convert all tokens to lowercase
	// return a slice of lowercase tokens
	r := make([]string, len(tokens)) // make a slice of the same length as tokens
	for i, token := range tokens {
		r[i] = strings.ToLower(token)

	}
	return r
}

// stopWordFilter removes all stop words from the given slice of tokens.
// It uses a map of stop words for the English language.
// It returns a new slice with the tokens that are not stop words.
//
// Parameters:
//
//	tokens: A slice of string tokens to filter.
//
// Returns:
//
//	[]string: A slice of tokens without stop words.
func stopWordFilter(tokens []string) []string {
	// remove all stop words from the tokens
	// return a slice of tokens without stop words
	var stopWords = map[string]struct{}{
		"a":    {},
		"an":   {},
		"and":  {},
		"the":  {},
		"of":   {},
		"be":   {},
		"to":   {},
		"it":   {},
		"that": {},
		"have": {},
		"for":  {},
		"not":  {},
		"on":   {},
		"with": {},
		"as":   {},
		"do":   {},
		"at":   {},
		"this": {},
		"but":  {},
		"by":   {},
	}

	r := make([]string, 0, len(tokens)) // make a slice of length 0 and capacity of len(tokens)
	for _, token := range tokens {
		if _, ok := stopWords[token]; !ok { // if token is not a stop word
			r = append(r, token)
		}
	}

	return r
}

// stemFilter applies stemming to each token in the given slice.
// It uses the snowballeng.Stem function to stem each token.
// It returns a new slice with the stemmed tokens.
//
// Parameters:
//
//	tokens: A slice of string tokens to stem.
//
// Returns:
//
//	[]string: A slice of stemmed tokens.
func stemFilter(tokens []string) []string {
	// stem each token
	// return a slice of stemmed tokens
	r := make([]string, len(tokens))

	for i, token := range tokens {
		r[i] = snowballeng.Stem(token, false)
	}
	return r
}
