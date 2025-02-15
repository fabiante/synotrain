package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fabiante/synotrain/app"
	"os"
	"strings"
)

type Data struct {
	Synonyms []app.SynonymGroup
}

func NewData() *Data {
	return &Data{
		Synonyms: make([][]string, 0),
	}
}

func ExampleSynonyms() app.SynonymGroup {
	return []string{"schön", "attraktiv", "bezaubernd", "charmant", "anziehend"}
}

type Model struct {
	// data contains the overall program data. This is a pointer value
	// to ensure no expansive deep copy is done when copying the model.
	data *Data

	learnModel app.LearnModel
}

func NewModel(data *Data) Model {
	return Model{
		data:       data,
		learnModel: app.LearnModel{},
	}
}

func (m Model) startLearn() (Model, tea.Cmd) {
	synonyms := m.data.Synonyms[0] // TODO: Use random
	m.learnModel = app.NewLearnModel(synonyms)
	return m, m.learnModel.Init()
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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

	if m.learnModel.IsSolving() {
		sb.WriteString(m.learnModel.View())
	}

	return sb.String()
}

func main() {
	data := NewData()
	data.Synonyms = append(data.Synonyms, ExampleSynonyms())
	p := tea.NewProgram(NewModel(data))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
