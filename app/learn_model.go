package app

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

type LearnModelSolvedMsg struct {
	SolvedCount int
}

type LearnModel struct {
	// active is the currently active SynonymGroup
	active SynonymGroup

	// solved contains the synonyms from active which the user has correctly typed in.
	solved SynonymGroup

	startWord string

	inputHint string
	input     textinput.Model
}

// NewLearnModel creates a new LearnModel using the given SynonymGroup.
//
// From those synonyms, a single word is chosen as start word, show to the user
// as initial hint for what synonyms are being searched.
func NewLearnModel(synonyms SynonymGroup) LearnModel {
	startWord := synonyms[0]

	ti := textinput.New()
	ti.Placeholder = "Enter a synonym for: " + startWord
	ti.Focus()
	ti.CharLimit = 156

	return LearnModel{
		active:    synonyms[1:],
		solved:    make(SynonymGroup, 0),
		startWord: startWord, // TODO: Choose random start word - fix the slicing of synonyms above too
		inputHint: "",
		input:     ti,
	}
}

func (m LearnModel) IsUnsolved() bool {
	return len(m.active) > len(m.solved)
}

func (m LearnModel) IsSolved() bool {
	return len(m.active) == len(m.solved)
}

func (m LearnModel) isUnsolvedSynonym(s string) bool {
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

func (m LearnModel) solve(s string) (LearnModel, tea.Cmd) {
	m.inputHint = "Correct"
	m.solved = append(m.solved, s)
	if m.IsSolved() {
		return m, tea.Cmd(func() tea.Msg {
			return LearnModelSolvedMsg{
				SolvedCount: len(m.solved),
			}
		})
	} else {
		return m, nil
	}
}

func (m LearnModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m LearnModel) Update(msg tea.Msg) (LearnModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// Confirm
		case "enter":
			if m.IsUnsolved() {
				// Handle the input typed by the user
				input := m.input.Value()
				m.input.SetValue("")
				if m.isUnsolvedSynonym(input) {
					return m.solve(input)
				} else {
					m.inputHint = "Incorrect"
					return m, nil
				}
			}
		}
	}

	// Pass on input to the text input
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	return m, cmd
}

func (m LearnModel) View() string {
	var sb strings.Builder

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

	return sb.String()
}
