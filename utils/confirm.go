package utils

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type confirmModel struct {
	prompt      string
	yamlPreview string
	confirmed   bool
	quitting    bool
	yesActive   bool // true if 'Yes' is selected, false if 'No' is selected
}

func (m confirmModel) Init() tea.Cmd {
	return nil
}

func (m confirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			m.confirmed = false
			return m, tea.Quit
		case "left", "h", "tab":
			m.yesActive = !m.yesActive
			return m, nil
		case "right", "l":
			m.yesActive = !m.yesActive
			return m, nil
		case "y", "Y":
			m.confirmed = true
			m.quitting = true
			return m, tea.Quit
		case "n", "N", "q":
			m.confirmed = false
			m.quitting = true
			return m, tea.Quit
		case "enter":
			m.confirmed = m.yesActive
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m confirmModel) View() string {
	if m.quitting {
		return ""
	}

	// Stylings
	var yesStyle, noStyle lipgloss.Style
	if m.yesActive {
		yesStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#00FF00")).
			Bold(true).
			Padding(0, 2)
		noStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#333333")).
			Padding(0, 2)
	} else {
		yesStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#333333")).
			Padding(0, 2)
		noStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#FF0000")).
			Bold(true).
			Padding(0, 2)
	}

	promptStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FAFAFA"))

	buttons := lipgloss.JoinHorizontal(lipgloss.Top,
		yesStyle.Render("Yes"),
		"  ",
		noStyle.Render("No"),
	)

	var previewBox string
	if m.yamlPreview != "" {
		previewStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#00FFFF")).
			Padding(1, 2).
			MarginBottom(1)
		previewBox = previewStyle.Render(m.yamlPreview) + "\n"
	}

	return fmt.Sprintf("\n%s%s\n\n%s\n\n", previewBox, promptStyle.Render(m.prompt), buttons)
}

func AskConfirm(prompt string, yamlPreview string) bool {
	p := tea.NewProgram(confirmModel{
		prompt:      prompt,
		yamlPreview: yamlPreview,
		yesActive:   true, // Default to Yes
	})
	m, err := p.Run()
	if err != nil {
		fmt.Printf("Error running prompt: %v\n", err)
		return false
	}
	if cm, ok := m.(confirmModel); ok {
		return cm.confirmed
	}
	return false
}
