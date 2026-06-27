package utils

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type inputModel struct {
	textInput textinput.Model
	err       error
	pkgName   string
	quitting  bool
}

func initialInputModel() inputModel {
	ti := textinput.New()
	ti.Placeholder = "package-name (e.g. goforge)"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 30

	return inputModel{
		textInput: ti,
		err:       nil,
	}
}

func (m inputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			val := m.textInput.Value()
			if val != "" {
				m.pkgName = val
				m.quitting = true
				return m, tea.Quit
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			m.pkgName = ""
			return m, tea.Quit
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m inputModel) View() string {
	if m.quitting {
		return ""
	}

	// Border style around the input field
	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#00FF00")).
		Padding(0, 1)

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FAFAFA"))

	return fmt.Sprintf(
		"\n%s\n\n%s\n\n(press Enter to confirm, Esc to quit)\n\n",
		titleStyle.Render("Enter package name for initialization:"),
		inputStyle.Render(m.textInput.View()),
	)
}

func PromptPackageName() (string, error) {
	p := tea.NewProgram(initialInputModel())
	m, err := p.Run()
	if err != nil {
		return "", err
	}
	if im, ok := m.(inputModel); ok {
		if im.pkgName == "" {
			return "", fmt.Errorf("initialization cancelled")
		}
		return im.pkgName, nil
	}
	return "", fmt.Errorf("invalid model type")
}
