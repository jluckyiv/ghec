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
	starting state = iota
	quitting
	ready
)

var (
	listStyle = lipgloss.NewStyle().
			Margin(1, 2)
	containerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(lipgloss.Color("205"))
)

var escKey = key.NewBinding(
	key.WithKeys("esc"),
)

var levelKeys = key.NewBinding(
	key.WithKeys("1", "2", "3", "4", "5", "6", "7", "8", "9"),
	key.WithHelp("1–9", "card lvl"),
)

var previousEnhancementKeys = key.NewBinding(
	key.WithKeys("p", "P"),
	key.WithHelp("p/P", "prev enh"),
)

var targetKeys = key.NewBinding(
	key.WithKeys("+", "=", "-"),
	key.WithHelp("+/-", "cur tgts"),
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
	// modifiers are values that affect the enhancement cost.
	modifiers struct {
		// level is the card level, which affects the enhancement cost.
		level ghec.Level
		// prev is the number of previous enhancements on the card, which affects the
		// enhancement cost.
		prev ghec.PreviousEnhancements
		// targets is the current number of targets on the card.
		// Multiple targets double the enhancement cost and adding a hex applies a
		// formula based on the number of targets.
		targets int
	}
	// state is the current state of the UI.
	state state
	// width and height are the current terminal dimensions.
	width  int
	height int
}

func initialModel() model {
	// Set the initial state.
	state := starting
	// Set the list items and the data map.
	data, items := enhancementsData()
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
	// Set the model from the data.
	m := model{state: state, data: data, list: l}
	// Set default values for level, targets, and previous enhancements.
	return m.resetModifiers()
}

func enhancementsData() (map[string]ghec.BaseEnhancement, []list.Item) {
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
	return data, items
}

func (m model) level() ghec.Level {
	return m.modifiers.level
}

func (m model) targets() int {
	return m.modifiers.targets
}

func (m model) prev() ghec.PreviousEnhancements {
	return m.modifiers.prev
}

func (m model) title() string {
	title := fmt.Sprintf(
		"Level: %1d, Targets: %2d, Previous: %1d",
		m.level(), m.targets(), m.prev(),
	)
	cost, err := m.cost()
	if err != nil {
		return title
	}
	return fmt.Sprintf("%s, Cost: %3d", title, cost)
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
	return ghec.NewEnhancement(be,
		ghec.OptionWithLevel(m.level()),
		ghec.OptionWithMultipleTarget(m.targets()),
		ghec.OptionWithPreviousEnhancements(m.prev()),
	).Cost()
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

	case errMsg:
		m.err = msg
		return m, nil

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		m.state = ready
		return m, nil

	case tea.KeyMsg:
		if key.Matches(msg, escKey) && !m.list.IsFiltered() {
			// If the list is filtered, don't quit the app.
			// Instead, reset the model and return so the list is not updated.
			if m.level() == 1 && m.targets() == 1 && m.prev() == 0 {
				m.state = quitting
				return m, tea.Quit
			}
			return m.resetModifiers(), nil
		}

		if key.Matches(msg, levelKeys) {
			m = m.setCardLevel(msg)
		}
		if key.Matches(msg, previousEnhancementKeys) {
			m = m.setPreviousEnhancements(msg)
		}
		if key.Matches(msg, targetKeys) {
			m = m.setCurrentTargets(msg)
		}
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) setCurrentTargets(msg tea.KeyMsg) model {
	if msg.String() == "+" || msg.String() == "=" {
		m.modifiers.targets = m.modifiers.targets + 1
	}
	if msg.String() == "-" && m.modifiers.targets > 1 {
		m.modifiers.targets = m.modifiers.targets - 1
	}
	return m
}

func (m model) setPreviousEnhancements(msg tea.KeyMsg) model {
	if msg.String() == "P" {
		m.modifiers.prev = ghec.DecrementPrevious(m.modifiers.prev)
	}
	if msg.String() == "p" {
		m.modifiers.prev = ghec.IncrementPrevious(m.modifiers.prev)
	}
	return m
}

func (m model) setCardLevel(msg tea.KeyMsg) model {
	if l, err := strconv.Atoi(msg.String()); err == nil {
		m.modifiers.level = ghec.Level(l)
	}
	return m
}

func (m model) resetModifiers() model {
	m.modifiers.level = 1
	m.modifiers.targets = 1
	m.modifiers.prev = 0
	return m
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	if m.state == quitting {
		return "\n  Quitting"
	}
	// Reduce the container size by its borders.
	borderW, borderH := containerStyle.GetFrameSize()
	containerW := max((m.width-borderW)/2, 80)
	containerH := m.height - borderH

	// Reduce the list size by its margins and padding,
	// plus the container borders.
	frameH, frameW := listStyle.GetFrameSize()
	listW := containerW - frameH
	listH := containerH - frameW
	m.list.SetWidth(listW)
	m.list.SetHeight(listH)

	// Set the contents of the list.
	m.list.Title = m.title()
	content := listStyle.Render(m.list.View())
	return containerStyle.
		Width(containerW).
		Height(containerH).
		Render(content)
}

func Run() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
