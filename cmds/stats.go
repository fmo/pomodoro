package cmds

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func NewStatsCmd(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "stats",
		Short: "Pomodoro stats",
		Run: func(cmd *cobra.Command, args []string) {
			total, err := app.pomodoroManager.TotalTime("today")
			if err != nil {
				app.logger.Error("cant get total time", "err", err)
				os.Exit(1)
			}

			fmt.Printf("%s\n", total)
			os.Exit(0)
		},
	}
}
