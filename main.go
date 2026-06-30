package main

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/fmo/pomodoro/cmd"
	"github.com/spf13/viper"
)

func main() {
	handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(handler)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		logger.Error("cant get the user directory", "err", err)
		os.Exit(1)
	}
	configPath := filepath.Join(homeDir, "Library", "Application Support", "pomodoro")

	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(configPath)

	configFile := filepath.Join(configPath, "config.yml")

	_, err = os.Create(configFile)
	if err != nil {
		logger.Error("cant create config", "err", err)
	}

	viper.Set("csv", "pomodoro.csv")

	err = viper.WriteConfig()
	if err != nil {
		logger.Error("cant write config", "err", err)
	}

	if err := viper.ReadInConfig(); err != nil {
		slog.Debug("cant read config file", "err", err)
		if err := viper.WriteConfig(); err != nil {
			slog.Error("cant create config file", "err", err)
			os.Exit(1)
		}
	}

	app := cmd.NewApp(logger, viper.GetViper())

	cmd.Execute(app)
}
