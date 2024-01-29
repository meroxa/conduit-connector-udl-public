package destination

import (
	"reflect"
	"strings"
)

var DataModeValues = []string{"TEST", "REAL", "SIMULATED", "EXERCISE"}

var DataTypeValues = []string{"AIS", "ELSET", "EPHEMERIS"}

func SupportedStringValues(check string, supported []string) bool {
	for _, ds := range supported {
		if strings.ToUpper(strings.TrimSpace(check)) == ds {
			return true
		}
	}
	return false
}

func allZero(slice []float64) bool {
	for _, v := range slice {
		if v != 0 {
			return false
		}
	}
	return true
}

func replaceUnderscoresWithSpaces(s string) string {
	return strings.ReplaceAll(s, "_", " ")
}

func replaceUnderscoresInStruct(v interface{}) {
	val := reflect.ValueOf(v).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		switch field.Kind() {
		case reflect.String:
			// Replace underscores with spaces and set the new value
			newValue := strings.ReplaceAll(field.String(), "_", " ")
			if field.CanSet() {
				field.SetString(newValue)
			}
		case reflect.Struct:
			// Recursively process nested structs
			replaceUnderscoresInStruct(field.Addr().Interface())
			// Optionally handle other container types like slices, arrays, maps, etc.
		}
	}
}
