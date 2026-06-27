package utils

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type osOption struct {
	label string
	value string
}

type multiSelectOSModel struct {
	options      []osOption
	cursor       int
	selected     map[int]bool
	quitting     bool
	confirmed    bool
	errorMessage string
}

func initialMultiSelectOSModel() multiSelectOSModel {
	return multiSelectOSModel{
		options: []osOption{
			{label: "Windows", value: "windows"},
			{label: "macOS (Darwin)", value: "mac"},
			{label: "Linux", value: "linux"},
		},
		selected: map[int]bool{
			0: true, // Default to select Windows
		},
	}
}

func (m multiSelectOSModel) Init() tea.Cmd {
	return nil
}

func (m multiSelectOSModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			m.confirmed = false
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
			m.errorMessage = ""
			return m, nil
		case "down", "j", "tab":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
			m.errorMessage = ""
			return m, nil
		case " ":
			m.selected[m.cursor] = !m.selected[m.cursor]
			m.errorMessage = ""
			return m, nil
		case "enter":
			hasSelection := false
			for _, sel := range m.selected {
				if sel {
					hasSelection = true
					break
				}
			}
			if !hasSelection {
				m.errorMessage = "Please select at least one target OS."
				return m, nil
			}
			m.confirmed = true
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m multiSelectOSModel) View() string {
	if m.quitting {
		return ""
	}

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FAFAFA"))
	cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Bold(true)
	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF")).Bold(true)
	unselectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Italic(true)

	var s string
	s += "\n" + titleStyle.Render("Select target OSes:") + "\n\n"

	for i, opt := range m.options {
		cursor := " "
		if m.cursor == i {
			cursor = cursorStyle.Render(">")
		}

		checked := "[ ]"
		if m.selected[i] {
			checked = selectedStyle.Render("[x]")
		}

		label := opt.label
		if m.cursor == i {
			label = cursorStyle.Render(label)
		} else if m.selected[i] {
			label = selectedStyle.Render(label)
		} else {
			label = unselectedStyle.Render(label)
		}

		s += fmt.Sprintf("%s %s %s\n", cursor, checked, label)
	}

	if m.errorMessage != "" {
		s += "\n" + errorStyle.Render(m.errorMessage) + "\n"
	}

	s += "\n(Space to toggle, Enter to confirm, Esc to quit)\n\n"
	return s
}

func PromptTargetOSes() ([]string, error) {
	p := tea.NewProgram(initialMultiSelectOSModel())
	m, err := p.Run()
	if err != nil {
		return nil, err
	}
	if sm, ok := m.(multiSelectOSModel); ok {
		if !sm.confirmed {
			return nil, fmt.Errorf("initialization cancelled")
		}
		var selected []string
		for i, opt := range sm.options {
			if sm.selected[i] {
				selected = append(selected, opt.value)
			}
		}
		return selected, nil
	}
	return nil, fmt.Errorf("invalid model type")
}
