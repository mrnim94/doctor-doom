package helper

import "os"

// GetEnvOrDefault returns the value of the environment variable represented by the key.
// If the key doesn't exist, it returns the defaultVal.
func GetEnvOrDefault(key string, defaultVal string) string {
	val, exists := os.LookupEnv(key)
	if !exists {
		return defaultVal
	}
	return val
}
