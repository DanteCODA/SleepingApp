package repos

import (
	"context"
	"strings"
	"time"
)

// createContext create a new context with timeout
func createContext(ctx context.Context, t uint64) (context.Context, context.CancelFunc) {
	timeout := time.Duration(t) * time.Millisecond
	return context.WithTimeout(ctx, timeout*time.Millisecond)
}

// stringsToUpperCase converts slide of string to slides of lower case string
func stringsToUpperCase(strs []string) (upperStrings []string, err error) {
	for _, str := range strs {
		upperString := strings.ToUpper(str)
		upperStrings = append(upperStrings, upperString)
	}

	return upperStrings, nil
}
