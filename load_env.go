package main

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
