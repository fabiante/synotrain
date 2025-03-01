package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fabiante/synotrain/app"
	"github.com/fabiante/synotrain/builtin"
)

type Model struct {
	// data contains the overall program data. This is a pointer value
	// to ensure no expansive deep copy is done when copying the model.
	data *app.Data

	solvedCount int

	childModel tea.Model

	debug string
}

func NewModel(data *app.Data) Model {
	return Model{
		data: data,
	}
}

func (m Model) startLearn() (Model, tea.Cmd) {
	// Select a random synonym group
	i := rand.IntN(len(m.data.Synonyms)) // TODO: Take into account if group was previously learned
	synonyms := m.data.Synonyms[i]

	// Create a child model
	m.childModel = app.NewLearnModel(synonyms)
	return m, m.childModel.Init()
}

func (m Model) Init() tea.Cmd {
	return tea.ClearScreen
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case app.LearnModelSolvedMsg:
		m.solvedCount += msg.SolvedCount
		m.childModel = nil
	case tea.KeyMsg:
		switch msg.String() {
		// Solve new synonym group
		case "ctrl+t":
			return m.startLearn()
		// Quit
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	// Pass on input to the child model
	if m.childModel != nil {
		var cmd tea.Cmd
		m.childModel, cmd = m.childModel.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	var sb strings.Builder

	// Show child model or main view
	if m.childModel != nil {
		sb.WriteString(m.childModel.View())
	} else {
		sb.WriteString("(ctrl+t = solve new synonym group) (ctrl+c = quit)\n\n")

		if m.solvedCount > 0 {
			sb.WriteString(fmt.Sprintf("Today you have learned %v synonyms\n\n", m.solvedCount))
		} else {
			sb.WriteString("Let's learn some synonyms!\n\n")
		}
	}

	if m.debug != "" {
		sb.WriteString(m.debug)
		sb.WriteString("\n")
	}

	return sb.String()
}

func main() {
	data := builtin.Data()
	p := tea.NewProgram(NewModel(data))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
