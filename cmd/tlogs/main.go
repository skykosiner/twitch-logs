package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/skykosiner/twitch-logs/pkg/logs"
	"github.com/spf13/cobra"
)

type state int

const (
	dateSelectionState state = iota
	resultsListState
)

type model struct {
	cursor       int
	dates        []string
	channel      string
	username     string
	selectedDate string
	state        state
	viewport     viewport.Model
	ready        bool
}

func initialModel(dates []string, channel, username string) model {
	vp := viewport.New(0, 0)

	return model{
		dates:    dates,
		channel:  channel,
		username: username,
		state:    dateSelectionState,
		viewport: vp,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			return m, tea.Quit
		}

		if m.state == dateSelectionState {
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.dates)-1 {
					m.cursor++
				}
			case "enter":
				m.selectedDate = m.dates[m.cursor]
				parts := strings.Split(m.selectedDate, " ")
				if len(parts) >= 2 {
					logs := logs.GetLogs(m.channel, m.username, parts[0], parts[1])
					m.viewport.SetContent(strings.Join(logs.GetStringSlice(), "\n"))
					m.viewport.GotoTop()
					m.state = resultsListState
				}
			}
		} else {
			if msg.String() == "esc" {
				m.state = dateSelectionState
			}
		}

	case tea.WindowSizeMsg:
		headerHeight := 2
		footerHeight := 2
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}
	}

	if m.state == resultsListState {
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	var s strings.Builder
	if m.state == dateSelectionState {
		fmt.Fprintf(&s, "Logs for %s/%s\n\n", m.channel, m.username)
		s.WriteString("Select date:\n")
		for i, date := range m.dates {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			fmt.Fprintf(&s, "%s %s\n", cursor, date)
		}
		s.WriteString("\n(q to quit)")
	} else {
		fmt.Fprintf(&s, "--- Logs for %s (ESC to go back) ---\n", m.selectedDate)
		s.WriteString(m.viewport.View())
		s.WriteString("\n-------------------------------------------")
	}
	return s.String()
}

func main() {
	rootCmd := &cobra.Command{
		Short: "tlogs - Twitch Logs",
		Use:   "tlogs [channel] [username]",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				cmd.Help()
				return
			}
			channel := args[0]
			username := args[1]

			availableLogs := logs.GetLogsAvailable(channel, username)

			p := tea.NewProgram(
				initialModel(availableLogs.GetStringSlice(), channel, username),
				tea.WithAltScreen(),
			)

			if _, err := p.Run(); err != nil {
				fmt.Printf("Error: %v", err)
				os.Exit(1)
			}
		},
	}

	if err := rootCmd.Execute(); err != nil {
		slog.Error("Error executing command", "error", err)
		os.Exit(1)
	}
}
