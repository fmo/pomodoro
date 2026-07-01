package main

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/fmo/pomodoro/cmd"
	"github.com/spf13/viper"
)

func main() {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})
	logger := slog.New(handler)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		logger.Error("cant get the user directory", "err", err)
		os.Exit(1)
	}
	configPath := filepath.Join(homeDir, "Library", "Application Support", "pomodoro")

	viper.WithLogger(logger)
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(configPath)

	if err := os.MkdirAll(configPath, 0o700); err != nil {
		logger.Error("cant create directories", "err", err)
		os.Exit(1)
	}

	configFile := filepath.Join(configPath, "config.yml")

	_, err = os.Open(configFile)
	if err != nil {
		_, err = os.Create(configFile)
		if err != nil {
			logger.Error("cant create config", "err", err)
		}
	}
	if os.Getenv("env") != "" {
		viper.Set("env", "dev")
	}
	viper.Set("csv", "pomodoro.csv")

	err = viper.WriteConfig()
	if err != nil {
		logger.Error("cant write config", "err", err)
	}

	if err := viper.ReadInConfig(); err != nil {
		logger.Error("cant read config file", "err", err)
		os.Exit(1)
	}

	app := cmd.NewApp(logger, viper.GetViper())

	cmd.Execute(app)
}
