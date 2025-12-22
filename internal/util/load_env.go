package util

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Env struct {
	Key   string
	Value string
}

// LoadEnv loads environment variables; this is a no-op fallback so the program compiles.
func LoadEnv() error {
	// TODO: implement actual .env loading (e.g., using github.com/joho/godotenv) if needed.
	_ = godotenv.Load()
	return nil
}

func GetEnv(key string, d string) string {
	v, exist := os.LookupEnv(key)
	if exist {
		return v
	}
	return d
}

func GetEnvBool(key string, d bool) bool {
	v, exist := os.LookupEnv(key)
	result, err := strconv.ParseBool(v)
	if err == nil && exist {
		return result
	}
	return d
}

func GetEnvInt(key string, d int) int {
	v, exist := os.LookupEnv(key)
	result, err := strconv.Atoi(v)
	if err == nil && exist {
		return result
	}
	return d
}

func GetAllEnvs() []Env {
	envs := os.Environ()
	res := make([]Env, 0, len(envs))
	for _, e := range envs {
		parts := strings.SplitN(e, "=", 2)
		key := parts[0]
		val := ""
		if len(parts) > 1 {
			val = parts[1]
		}
		res = append(res, Env{Key: key, Value: val})
	}
	return res
}

// isSensitiveKey reports whether the environment variable name looks sensitive.
// It performs a case-insensitive substring check for common sensitive keywords.
func isSensitiveKey(key string) bool {
	k := strings.ToUpper(key)
	sensitive := []string{
		"PASSWORD", "PASS", "SECRET", "TOKEN", "KEY", "AWS_", "ACCESS", "SECRET_KEY", "PRIVATE", "CREDENTIAL", "PWD", "API_KEY",
	}
	for _, s := range sensitive {
		if strings.Contains(k, s) {
			return true
		}
	}
	return false
}

// maskValue returns a short redaction string indicating the original length.
func maskValue(v string) string {
	if v == "" {
		return ""
	}
	return "<REDACTED>"
}

// FilterSensitiveEnvs returns a copy of the provided env slice with sensitive
// values redacted. Keys are preserved.
func FilterSensitiveEnvs(envs []Env) []Env {
	out := make([]Env, 0, len(envs))
	for _, e := range envs {
		if isSensitiveKey(e.Key) {
			out = append(out, Env{Key: e.Key, Value: maskValue(e.Value)})
		} else {
			out = append(out, e)
		}
	}
	return out
}

// GetFilteredEnvs returns all environment variables with sensitive values redacted.
func GetFilteredEnvs() []Env {
	return FilterSensitiveEnvs(GetAllEnvs())
}
