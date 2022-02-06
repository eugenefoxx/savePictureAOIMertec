package config

import "os"

type FilesConfig struct {
	SourceFile    string
	HomeSavingDir string
}

type HTTPConfig struct {
	Port         string
	StartURL     string
	IndexHTML    string
	IndexExecute string
	WorkDirWEB   string
}

type Config struct {
	Files FilesConfig
	HTTP  HTTPConfig
}

func GetConfig() *Config {
	return &Config{}
}

// New returns a new Config struct
func New() *Config {
	return &Config{
		Files: FilesConfig{
			SourceFile:    getEnv("SourceFile", ""),
			HomeSavingDir: getEnv("HomeSavingDIR", ""),
		},
		HTTP: HTTPConfig{
			Port:         getEnv("Port", ""),
			StartURL:     getEnv("StartURL", ""),
			IndexHTML:    getEnv("IndexHTML", ""),
			IndexExecute: getEnv("IndexExecute", ""),
			WorkDirWEB:   getEnv("WorkDirWEB", ""),
		},
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
