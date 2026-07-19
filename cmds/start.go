package cmds

import (
	"log"
	"time"

	"charm.land/bubbles/v2/progress"
	tea "charm.land/bubbletea/v2"
	"github.com/spf13/cobra"
)

func NewStartCmd(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {
			duration, _ := cmd.Flags().GetString("duration")
			d, err := time.ParseDuration(duration)
			if err != nil {
				log.Fatal(err)
			}

			m := model{
				app:      app,
				progress: progress.New(progress.WithDefaultBlend()),
				limit:    int(d.Seconds()),
				count:    int(d.Seconds()),
			}

			p := tea.NewProgram(m)
			p.Run()
		},
	}
}
