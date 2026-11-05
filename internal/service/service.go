package service

import (
	"fmt"
	"strings"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

// Convert determines the input format and converts Russian text or Morse code.
func Convert(data string) (string, error) {
	data = strings.TrimSpace(data)
	if data == "" {
		return "", fmt.Errorf("input is empty")
	}

	if strings.Trim(data, ".- \t\r\n") == "" {
		for _, code := range strings.Fields(data) {
			if morse.MorseToRune(code) == 0 {
				return "", fmt.Errorf("unknown Morse code %q", code)
			}
		}
		// Three spaces separate words; line breaks also separate words.
		data = strings.ReplaceAll(data, "\r\n", "\n")
		data = strings.NewReplacer("\r", "   ", "\n", "   ", "\t", " ").Replace(data)
		words := strings.Split(data, "   ")
		result := make([]string, 0, len(words))
		for _, word := range words {
			letters := strings.Fields(word)
			if len(letters) > 0 {
				result = append(result, morse.ToText(strings.Join(letters, " ")))
			}
		}
		return strings.Join(result, " "), nil
	}

	words := strings.Fields(data)
	for i, word := range words {
		for _, letter := range word {
			if morse.RuneToMorse(letter) == "" {
				return "", fmt.Errorf("unsupported character %q", letter)
			}
		}
		// The supplied converter skips spaces, so encode each word separately.
		words[i] = morse.ToMorse(word)
	}
	return strings.Join(words, "   "), nil
}
