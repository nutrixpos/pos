// Package userio provides an interface for user interaction prompts, and a bubbletea implementation of such an interface.
package userio

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/elmawardy/nutrix/backend/common/logger"
)

// Prompter defines an interface for user interaction prompts.
type Prompter interface {
	// MultiChooseTree presents a tree of choices to the user and returns the selected elements.
	MultiChooseTree(msg string, choices []PromptTreeElement) (selected []PromptTreeElement, err error)

	// Confirmation prompts the user with a yes/no question and returns the result.
	Confirmation(msg string) (bool, error)
}

// BubbleTeaSeedablesPrompter is a Prompter that uses the BubbleTea library to present interactive prompts to the user.
//
// It implements the Prompter interface, and can be used to create prompts for the user to select multiple values from a tree.
type BubbleTeaSeedablesPrompter struct {
	cursor int

	// Logger is the logger to use for logging messages.
	Logger logger.ILogger

	// Message is the message to display to the user.
	Message string

	// TreeChoices is the tree of choices to present to the user.
	TreeChoices []PromptTreeElement

	// Breadcrump is the breadcrumb to display to the user.
	Breadcrump []string

	// FullySelectedTreeChar is the character to display when a tree element is fully selected.
	FullySelectedTreeChar string

	// SemiSelectedTreeChar is the character to display when a tree element is partially selected.
	SemiSelectedTreeChar string

	// UnselectedTreeChar is the character to display when a tree element is not selected.
	UnselectedTreeChar string

	// ElementsCount is the number of elements in the tree.
	ElementsCount int

	// IsTreeChoices is true if the prompt is a tree of choices.
	IsTreeChoices bool

	// IsConfirmation is true if the prompt is a yes/no confirmation.
	IsConfirmation bool

	// ConfirmationResult is the result of the confirmation prompt.
	ConfirmationResult bool

	// UserInputText is the text that the user has input.
	UserInputText string

	// isTerminating is true if the prompt should terminate.
	isTerminating bool
}

// PromptTreeElement represents an element in the tree of choices.
//
// Title is the title of the element.
// Selected is true if the element is selected.
// Level is the level of the element in the tree (0 is the root, 1 is a child of the root, etc).
// CounterIndex is the index of the element in the tree for use in the BubbleTea prompt.
// SubElements are the sub-elements of the element.
type PromptTreeElement struct {
	Title        string
	Selected     bool
	Level        int
	CounterIndex int
	SubElements  []PromptTreeElement
}

// Confirmation prompts the user with a yes/no question and returns the result.
//
// It sets the text of the prompt to the given message and then runs the BubbleTea
// program. When the program finishes, it returns the result of the confirmation
// prompt.
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

// MultiChooseTree presents a tree of choices to the user and returns the selected elements.
//
// It sets the text of the prompt to the given message and then runs the BubbleTea
// program. When the program finishes, it returns the selected elements.
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

// PropagateCounterIndexToTree takes a slice of PromptTreeElements and returns the total
// count of elements in the tree and the same slice with the CounterIndex field of each
// element set to the correct value. The CounterIndex of each element is the index of the
// element in the flattened tree.
//
// This function is used to set the CounterIndex of each element in the tree so that
// BubbleTea can use it to determine the index of the selected element.
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

// ToggleSelectedTreeElement takes a target index and a slice of PromptTreeElements
// and toggles the `Selected` field of the element at the target index and all
// of its children. If the target index is not found in the tree, it does nothing.
//
// It returns the modified tree and a boolean indicating whether the target index
// was found.
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

// Init is the BubbleTea initialization function. It's called when the program starts.
// In this case, there's no I/O to be done, so it simply returns `nil`.
func (m *BubbleTeaSeedablesPrompter) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

// Update processes incoming messages and updates the state of the BubbleTeaSeedablesPrompter.
// It handles key presses for quitting the program and delegates to specific update functions
// based on the current mode (tree selection or confirmation).
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

// UpdateConfirmation processes incoming messages and updates the state of the BubbleTeaSeedablesPrompter in confirmation mode.
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

// UpdateTreeSelection processes incoming messages and updates the state of the BubbleTeaSeedablesPrompter in tree selection mode.
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

// View renders the BubbleTeaSeedablesPrompter to a string.
func (m *BubbleTeaSeedablesPrompter) View() string {
	// The header

	if m.IsTreeChoices {
		return m.TreeSelectionView(0, m.Message, m.Breadcrump, m.TreeChoices, "\nPress q to return.\n")
	} else if m.IsConfirmation {
		return m.ConfirmationView()
	}

	return ""
}

// ConfirmationView renders a confirmation prompt to a string.
func (m *BubbleTeaSeedablesPrompter) ConfirmationView() string {
	output := m.Message

	output += " [Y/n] "
	output += m.UserInputText
	output += "\n"

	return output
}

// TreeSelectionView renders a tree selection prompt to a string.
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
