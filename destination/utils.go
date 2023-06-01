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
