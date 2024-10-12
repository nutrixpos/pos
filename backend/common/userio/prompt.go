package userio

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/elmawardy/nutrix/common/logger"
)

type Prompter interface {
	MultiChooseTree(msg string, choices []PromptTreeElement) (selected []PromptTreeElement, err error)
	Confirmation(msg string) (bool, error)
}

type BubbleTeaSeedablesPrompter struct {
	cursor int

	Logger                logger.ILogger
	Message               string
	TreeChoices           []PromptTreeElement
	Breadcrump            []string
	FullySelectedTreeChar string // the character to display when a tree element is fully selected
	SemiSelectedTreeChar  string // the character to display when a tree element is partially selected
	UnselectedTreeChar    string // the character to display when a tree element is not selected
	ElementsCount         int

	IsTreeChoices      bool
	IsConfirmation     bool
	ConfirmationResult bool
	UserInputText      string
	isTerminating      bool
}

type PromptTreeElement struct {
	Title        string
	Selected     bool
	Level        int
	CounterIndex int
	SubElements  []PromptTreeElement
}

func (m *BubbleTeaSeedablesPrompter) Confirmation(msg string) (result bool, err error) {

	m.IsConfirmation = true
	m.IsTreeChoices = false
	m.Message = msg
	m.UserInputText = ""

	p := tea.NewProgram(m)
	if _, err = p.Run(); err != nil {
		m.Logger.Error("Alas, there's been an error in bubbletea prompt: %v", err)
		os.Exit(1)
	}

	return m.ConfirmationResult, err
}

func (m *BubbleTeaSeedablesPrompter) MultiChooseTree(msg string, choices []PromptTreeElement) (selected []PromptTreeElement, err error) {

	m.IsConfirmation = false
	m.IsTreeChoices = true
	m.ConfirmationResult = true
	m.Message = msg
	m.cursor = 0
	m.TreeChoices = choices

	m.ElementsCount, m.TreeChoices = m.PropagateCounterIndexToTree(0, choices)

	m.FullySelectedTreeChar = "X"
	m.SemiSelectedTreeChar = "─"
	m.UnselectedTreeChar = " "

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		m.Logger.Error("Alas, there's been an error in bubbletea prompt: %v", err)
		os.Exit(1)
	}

	if m.isTerminating {
		os.Exit(0)
	}

	return m.TreeChoices, err
}

func (m *BubbleTeaSeedablesPrompter) PropagateCounterIndexToTree(total int, tree []PromptTreeElement) (int, []PromptTreeElement) {

	counter := 0

	for index, element := range tree {
		tree[index].CounterIndex = total + index
		counter++

		if len(element.SubElements) > 0 {
			var new_total int
			new_total, tree[index].SubElements = m.PropagateCounterIndexToTree(total+counter, element.SubElements)
			total += (new_total - counter)
		}
	}

	total += counter

	m.Logger.Info(fmt.Sprintf("total: %v", total))
	return total, tree
}

func ToggleSelectedTreeElement(targetindex int, tree []PromptTreeElement) (result []PromptTreeElement, found bool) {

	found = false

	for index, element := range tree {
		if tree[index].CounterIndex == targetindex {
			tree[index].Selected = !tree[index].Selected
			for subElementIndex := range tree[index].SubElements {
				tree[index].SubElements[subElementIndex].Selected = tree[index].Selected
			}
			return tree, true
		}

		tree[index].SubElements, found = ToggleSelectedTreeElement(targetindex, element.SubElements)
		if found {
			break
		}
	}

	return tree, found
}

func (m *BubbleTeaSeedablesPrompter) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m *BubbleTeaSeedablesPrompter) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "ctrl+c" {
			m.isTerminating = true
			return m, tea.Quit
		}

		if k == "q" || k == "esc" {
			m.isTerminating = false
			return m, tea.Quit
		}
	}

	if m.IsTreeChoices {
		return UpdateTreeSelection(m, msg)
	} else if m.IsConfirmation {
		return m.UpdateConfirmation(msg)
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m *BubbleTeaSeedablesPrompter) UpdateConfirmation(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		case "enter":
			if m.UserInputText == "y" || m.UserInputText == "Y" {
				m.ConfirmationResult = true
				return m, tea.Quit
			} else if m.UserInputText == "n" || m.UserInputText == "N" {
				m.ConfirmationResult = false
				return m, tea.Quit
			} else if m.UserInputText == "" {
				m.ConfirmationResult = true
				return m, tea.Quit
			} else {
				m.UserInputText = ""
			}

		case "backspace":
			if len(m.UserInputText) > 0 {
				m.UserInputText = m.UserInputText[:len(m.UserInputText)-1]
			}
		default:
			m.UserInputText += msg.String()

		}
	}

	return m, nil
}

func UpdateTreeSelection(m *BubbleTeaSeedablesPrompter, msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c":
			m.isTerminating = true
			return m, tea.Quit

		case "q":
			m.isTerminating = false
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < m.ElementsCount-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case " ", "enter":
			m.TreeChoices, _ = ToggleSelectedTreeElement(m.cursor, m.TreeChoices)

		}
	}

	return m, nil
}

func (m *BubbleTeaSeedablesPrompter) View() string {
	// The header

	if m.IsTreeChoices {
		return m.TreeSelectionView(0, m.Message, m.Breadcrump, m.TreeChoices, "\nPress q to return.\n")
	} else if m.IsConfirmation {
		return m.ConfirmationView()
	}

	return ""
}

func (m *BubbleTeaSeedablesPrompter) ConfirmationView() string {
	output := m.Message

	output += " [Y/n] "
	output += m.UserInputText
	output += "\n"

	return output
}

func (m *BubbleTeaSeedablesPrompter) TreeSelectionView(level int, message string, breadcrump []string, choices []PromptTreeElement, finishMessage string) string {

	output := message

	for _, choice := range choices {

		cursor := " " // no cursor
		if m.cursor == choice.CounterIndex {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if choice.Selected {
			checked = "x" // selected!
		}

		for range level {
			output += "   "
		}

		output += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice.Title)
		output += m.TreeSelectionView(level+1, "", breadcrump, choice.SubElements, "")
	}

	output += finishMessage
	return output

}
