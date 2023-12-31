package tui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error

type state int

const (
	loading state = iota
	quitting
	stopped
	focusRight
)

var activeStyle = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder(), true).
	BorderForeground(lipgloss.Color("205"))

var inactiveStyle = activeStyle.Copy().
	Border(lipgloss.HiddenBorder(), true)

type model struct {
	err     error
	spinner spinner.Model
	state   state
	height  int
	width   int
}

var quitKeys = key.NewBinding(
	key.WithKeys("q", "esc", "ctrl+c"),
	key.WithHelp("", "press q to quit"),
)

var stopKeys = key.NewBinding(
	key.WithKeys(" "),
	key.WithHelp("", "press space or enter to stop"),
)

var focusKeys = key.NewBinding(
	key.WithKeys("tab"),
	key.WithHelp("", "press tab to focus"),
)

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{spinner: s, state: loading}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		return m, nil

	case tea.KeyMsg:
		if key.Matches(msg, quitKeys) {
			m.state = quitting
			return m, tea.Quit
		}
		if key.Matches(msg, stopKeys) {
			if m.state == stopped {
				m.state = loading
				return m, m.spinner.Tick
			}
			if m.state == loading {
				m.state = stopped
			}
			return m, nil
		}
		if key.Matches(msg, focusKeys) {
			if m.state == focusRight {
				m.state = loading
				return m, m.spinner.Tick
			}
			m.state = focusRight
		}
		return m, nil
	case errMsg:
		m.err = msg
		return m, nil

	default:
		var cmd tea.Cmd
		if m.state == loading {
			m.spinner, cmd = m.spinner.Update(msg)
		}
		return m, cmd
	}
}

func (m model) View() string {
	x, y := activeStyle.GetFrameSize()
	height := m.height - y
	width := (m.width - x) / 2
	if m.err != nil {
		return m.err.Error()
	}
	spinnerView := m.spinner.View()
	leftContent := fmt.Sprintf(
		"\n\n   %s Loading forever... %s\n\n",
		spinnerView,
		quitKeys.Help().Desc,
	)
	rightContent := fmt.Sprintf(
		"\n\n    %d, %d",
		m.height,
		m.width,
	)
	if m.state == quitting {
		return leftContent + "\n"
	}
	var left, right string
	if m.state == focusRight {
		left = inactiveStyle.Width(width).Height(height).Render(leftContent)
		right = activeStyle.Width(width).Height(height).Render(rightContent)
	} else {
		left = activeStyle.Width(width).Height(height).Render(leftContent)
		right = inactiveStyle.Width(width).Height(height).Render(rightContent)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}

func Run() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
