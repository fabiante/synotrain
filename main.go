package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"strings"
)

type SynonymGroup = []string

type Data struct {
	Synonyms []SynonymGroup
}

func NewData() *Data {
	return &Data{
		Synonyms: make([][]string, 0),
	}
}

func ExampleSynonyms() SynonymGroup {
	return []string{"schön", "attraktiv", "bezaubernd", "charmant", "anziehend"}
}

type Model struct {
	// data contains the overall program data. This is a pointer value
	// to ensure no expansive deep copy is done when copying the model.
	data *Data

	// active is the currently active SynonymGroup - This may be nil, indicating that
	// the user does not currently train.
	active SynonymGroup

	// solved contains the synonyms from active which the user has correctly typed in.
	solved SynonymGroup

	startWord string

	inputHint string
	input     textinput.Model
}

func (m Model) isSolving() bool {
	return m.active != nil && len(m.active) > len(m.solved)
}

func (m Model) startSolve() Model {
	m.active = m.data.Synonyms[0] // FIXME: Improve selection
	m.solved = []string{}

	// Pick one random word from the synonyms as start word
	m.startWord = m.active[0] // FIXME: Not random
	m.active = m.active[1:]   // Remove start word from synonym group so it must not be typed in

	return m
}

func (m Model) isUnsolvedSynonym(s string) bool {
	// Check if s is a synonym
	for _, synonym := range m.active {
		if strings.EqualFold(s, synonym) {
			// Check if synonym is unsolved
			for _, solved := range m.solved {
				if strings.EqualFold(s, solved) {
					return false
				}
			}

			return true
		}
	}

	return false
}

func (m Model) solve(s string) Model {
	m.solved = append(m.solved, s)
	return m
}

func initialModel(data *Data) Model {
	ti := textinput.New()
	ti.Placeholder = "New word"
	ti.Focus()
	ti.CharLimit = 156

	return Model{
		data:  data,
		input: ti,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {
		// Confirm
		case "enter":
			if m.isSolving() {
				// Handle the input typed by the user
				input := m.input.Value()
				m.input.SetValue("")
				if m.isUnsolvedSynonym(input) {
					m.inputHint = "Correct"
					return m.solve(input), nil
				} else {
					m.inputHint = "Incorrect"
					return m, nil
				}
			}
		// Solve new synonym group
		case "ctrl+t":
			return m.startSolve(), textinput.Blink
		// Quit
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	// Pass on input to the text input
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	var sb strings.Builder

	sb.WriteString("(ctrl+t = solve new synonym group) (ctrl+c = quit)\n\n")

	if m.isSolving() {
		activeLen := len(m.active)
		sb.WriteString(fmt.Sprintf("Synonym group has %v synonyms - %v remaining\n", activeLen, activeLen-len(m.solved)))
		sb.WriteString("Find synonyms for: ")
		sb.WriteString(m.startWord)
		sb.WriteString("\n")

		sb.WriteString("You have found: ")
		sb.WriteString(strings.Join(m.solved, " "))
		sb.WriteString("\n\n")

		sb.WriteString(m.input.View())
		sb.WriteString("\n")
		sb.WriteString(m.inputHint)
		sb.WriteString("\n")
	}

	return sb.String()
}

func main() {
	data := NewData()
	data.Synonyms = append(data.Synonyms, ExampleSynonyms())
	p := tea.NewProgram(initialModel(data))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
