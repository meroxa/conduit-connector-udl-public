package destination

import "strings"

var DataModeValues = []string{"TEST", "REAL", "SIMULATED", "EXERCISE"}

var DataTypeValues = []string{"AIS", "ELSET"}

func SupportedStringValues(check string, supported []string) bool {
	for _, ds := range supported {
		if strings.ToUpper(strings.TrimSpace(check)) == ds {
			return true
		}
	}
	return false
}

func toTitleCase(str string) string {
	// Convert the entire string to lowercase
	str = strings.ToLower(str)

	// Use strings.Fields to obtain a slice of words
	words := strings.Fields(str)

	// Iterate over words and capitalize the first letter of each one
	for i, word := range words {
		words[i] = strings.Title(word)
	}

	// Rejoin words into a string
	return strings.Join(words, " ")
}

func containsString(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}
