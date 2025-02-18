package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fabiante/synotrain/app"
	"github.com/fabiante/synotrain/builtin"
	"math/rand/v2"
	"os"
	"strings"
)

type Model struct {
	// data contains the overall program data. This is a pointer value
	// to ensure no expansive deep copy is done when copying the model.
	data *app.Data

	solvedCount int

	learnModel app.LearnModel

	debug string
}

func NewModel(data *app.Data) Model {
	return Model{
		data:       data,
		learnModel: app.LearnModel{},
	}
}

func (m Model) startLearn() (Model, tea.Cmd) {
	// Select a random synonym group
	i := rand.IntN(len(m.data.Synonyms)) // TODO: Take into account if group was previously learned
	synonyms := m.data.Synonyms[i]
	m.learnModel = app.NewLearnModel(synonyms)
	return m, m.learnModel.Init()
}

func (m Model) Init() tea.Cmd {
	return tea.ClearScreen
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case app.LearnModelSolvedMsg:
		m.solvedCount += msg.SolvedCount
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

	// Pass on input to the learn model
	var cmd tea.Cmd
	m.learnModel, cmd = m.learnModel.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	var sb strings.Builder

	sb.WriteString("(ctrl+t = solve new synonym group) (ctrl+c = quit)\n\n")

	if m.learnModel.IsUnsolved() {
		sb.WriteString(m.learnModel.View())
	} else {
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
