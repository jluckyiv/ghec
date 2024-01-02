package tui

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jluckyiv/ghec"
)

type errMsg error

type state int

const (
	loading state = iota
	quitting
	focusRight
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

var activeStyle = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder(), true).
	BorderForeground(lipgloss.Color("205"))

var inactiveStyle = activeStyle.Copy().
	Border(lipgloss.HiddenBorder(), true)

var focusKeys = key.NewBinding(
	key.WithKeys("tab"),
	key.WithHelp("focus", "press tab to focus"),
)

var levelKeys = key.NewBinding(
	key.WithKeys("1", "2", "3", "4", "5", "6", "7", "8", "9"),
	key.WithHelp("1–9", "card level"),
)

var previousEnhancementKeys = key.NewBinding(
	key.WithKeys("p", "P"),
	key.WithHelp("p/P", "num prev enh"),
)

var targetKeys = key.NewBinding(
	key.WithKeys("+", "=", "-"),
	key.WithHelp("+/-", "num targets"),
)

type model struct {
	list             list.Model
	err              error
	baseEnhancements []ghec.BaseEnhancement
	state            state
	height           int
	width            int
	level            ghec.Level
	prev             ghec.PreviousEnhancements
	targets          int
}

func initialModel() model {
	identity := func(be ghec.BaseEnhancement) ghec.BaseEnhancement { return be }
	baseEnhancements := ghec.List(identity)
	items := ghec.List(newItem)
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			levelKeys,
			previousEnhancementKeys,
			targetKeys,
		}
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			levelKeys,
			previousEnhancementKeys,
			targetKeys,
		}
	}
	m := model{
		list: l, state: loading, level: 1, prev: 0, targets: 1,
		baseEnhancements: baseEnhancements,
	}
	return m
}

func (m model) Title() string {
	return fmt.Sprintf("Level: %d, Targets: %d, Previous: %d, Cost: %d",
		m.level, m.targets, m.prev, m.cost())
}

func (m model) cost() ghec.Cost {
	be := m.baseEnhancements[m.list.Index()]
	cost, _ := ghec.NewEnhancement(be).
		WithLevel(m.level).
		WithPreviousEnhancements(m.prev).
		WithMultipleTarget(m.targets).
		Cost()
	return cost
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		return m, nil

	case errMsg:
		m.err = msg
		return m, nil

	case tea.KeyMsg:
		if key.Matches(msg, levelKeys) {
			if l, err := strconv.Atoi(msg.String()); err == nil {
				m.level = ghec.Level(l)
			}
		}
		if key.Matches(msg, previousEnhancementKeys) {
			if msg.String() == "P" {
				m.prev = ghec.DecrementPrevious(m.prev)
			}
			if msg.String() == "p" {
				m.prev = ghec.IncrementPrevious(m.prev)
			}
		}
		if key.Matches(msg, targetKeys) {
			if msg.String() == "+" || msg.String() == "=" {
				m.targets = m.targets + 1
			}
			if msg.String() == "-" && m.targets > 1 {
				m.targets = m.targets - 1
			}
		}
		if key.Matches(msg, focusKeys) {
			if m.state == focusRight {
				m.state = loading
			} else {
				m.state = focusRight
			}
		}
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	x, y := activeStyle.GetFrameSize()
	frameWidth := (m.width - x) / 2
	frameHeight := m.height - y
	x2, y2 := docStyle.GetFrameSize()
	width := frameWidth - x2
	height := frameHeight - y2
	m.list.SetSize(width, height)
	m.list.Title = m.Title()
	leftContent := docStyle.Width(width).Height(height).Render(m.list.View())
	rightContent := fmt.Sprintf(
		`
        m.height: %d, m.width: %d
     frameHeight: %d, frameWidth: %d
          height: %d,   width: %d
         m.level: %d
          m.prev: %d
            Cost: %d
    `,
		m.height,
		m.width,
		frameHeight,
		frameWidth,
		height,
		width,
		m.level,
		m.prev,
		m.cost(),
	)
	if m.state == quitting {
		return leftContent + "\nQuitting"
	}
	var left, right string
	if m.state == focusRight {
		left = inactiveStyle.Width(frameWidth).Height(frameHeight).Render(leftContent)
		right = activeStyle.Width(frameWidth).Height(frameHeight).Render(rightContent)
	} else {
		left = activeStyle.Width(frameWidth).Height(frameHeight).Render(leftContent)
		right = inactiveStyle.Width(frameWidth).Height(frameHeight).Render(rightContent)
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
