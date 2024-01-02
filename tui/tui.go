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
	err error
	// list holds a list of items and a delegate for rendering the list.
	list list.Model
	// data holds the base enhancements for the list.
	// Because the list.Model can't return the base enhancement, the
	// data map is used to look up the base enhancement from the list item.
	// The key is the list item's FilterValue(), because that's the only
	// method in the list.Item interface.
	data map[string]ghec.BaseEnhancement
	// level is the card level, which affects the enhancement cost.
	level ghec.Level
	// prev is the number of previous enhancements on the card, which affects the
	// enhancement cost.
	prev ghec.PreviousEnhancements
	// targets is the current number of targets on the card.
	// Multiple targets double the enhancement cost and adding a hex applies a
	// formula based on the number of targets.
	targets int
	// state is the current state of the UI.
	state state
	// width and height are the current terminal dimensions.
	width  int
	height int
}

func initialModel() model {
	// Get a temporary list of the base enhancements.
	baseEnhancements := ghec.BaseEnhancements()
	// Use a map instead of a slice to look up the base enhancement.
	// Don't use a slice because the Index() method of the list.Model
	// returns the index of the selected item from the visible items,
	// not the index from the original list.
	data := make(map[string]ghec.BaseEnhancement)
	items := make([]list.Item, len(baseEnhancements))
	// Assign the base enhancements to the list items and the data map
	// from the same loop.
	for i, be := range baseEnhancements {
		item := newItem(be)
		items[i] = item
		data[item.FilterValue()] = baseEnhancements[i]
	}
	// Create the list.Model.
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
		list:    l,
		data:    data,
		level:   1,
		prev:    0,
		targets: 1,
		state:   loading,
	}
	return m
}

func (m model) Title() string {
	title := fmt.Sprintf(
		"Level: %d, Targets: %d, Previous: %d",
		m.level, m.targets, m.prev,
	)
	cost, err := m.cost()
	if err != nil {
		return title
	}
	return fmt.Sprintf("%s, Cost: %d", title, cost)
}

func (m model) cost() (ghec.Cost, error) {
	selected := m.list.SelectedItem()
	if selected == nil {
		return ghec.Cost(0), fmt.Errorf("no base enhancement selected")
	}
	be, ok := m.data[selected.FilterValue()]
	if !ok {
		return ghec.Cost(0), fmt.Errorf("base enhancement not found")
	}
	return ghec.NewEnhancement(be).
		WithLevel(m.level).
		WithPreviousEnhancements(m.prev).
		WithMultipleTarget(m.targets).
		Cost()
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {

	case errMsg:
		m.err = msg
		return m, nil

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
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
	w, h := activeStyle.GetFrameSize()
	frameWidth := (m.width - w) / 2
	frameHeight := m.height - h
	w2, h2 := docStyle.GetFrameSize()
	width := frameWidth - w2
	height := frameHeight - h2
	m.list.SetWidth(width)
	m.list.SetHeight(height)
	m.list.Title = m.Title()
	leftContent := docStyle.
		Width(width).
		Height(height).
		Render(m.list.View())
	cost, _ := m.cost()
	rightContent := fmt.Sprintf(
		`
        m.height: %d, m.width: %d
     frameHeight: %d, frameWidth: %d
          height: %d,   width: %d
         m.level: %d
          m.prev: %d
         index(): %d
    selectedItem: %s
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
		m.list.Index(),
		m.list.SelectedItem(),
		cost,
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
