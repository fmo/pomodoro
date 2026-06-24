package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"charm.land/bubbles/v2/progress"
	tea "charm.land/bubbletea/v2"
)

type model struct {
	progress progress.Model
	limit    int
	count    int
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q":
			m.Save()
			return m, tea.Quit
		}
	case tickMsg:
		if m.progress.Percent() == 1.0 || m.limit == 0 {
			m.Save()
			return m, tea.Quit
		}
		m.count--
		cmd := m.progress.IncrPercent(1.0 / float64(m.limit))
		return m, tea.Batch(cmd, tickCmd())
	case progress.FrameMsg:
		var cmd tea.Cmd
		m.progress, cmd = m.progress.Update(msg)
		return m, cmd

	}
	return m, nil
}

func (m model) View() tea.View {
	return tea.NewView(
		fmt.Sprintf("%s\nLeft: %d sec.", m.progress.View(), m.count),
	)
}

func (m model) Save() error {
	home, _ := os.UserHomeDir()
	fileToSave := filepath.Join(home, "Library", "Application Support", "pomodoro")

	err := os.MkdirAll(fileToSave, 0o700)
	if err != nil {
		return err
	}

	fullFileName := filepath.Join(fileToSave, "pomodoro.csv")

	var f *os.File

	f, err = os.OpenFile(fullFileName, os.O_APPEND|os.O_WRONLY, 0o700)
	if err != nil {
		f, err = os.Create(fullFileName)
		if err != nil {
			return err
		}
	}

	w := csv.NewWriter(f)
	err = w.Write([]string{time.Now().Format(time.RFC3339), strconv.Itoa(m.limit - m.count)})
	if err != nil {
		return err
	}
	w.Flush()

	return nil
}

func main() {
	duration := flag.String("duration", "3m", "timer duration")
	flag.Parse()

	d, err := time.ParseDuration(*duration)
	if err != nil {
		log.Fatal(err)
	}

	m := model{
		progress: progress.New(progress.WithDefaultBlend()),
		limit:    int(d.Seconds()),
		count:    int(d.Seconds()),
	}

	p := tea.NewProgram(m)
	p.Run()
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
