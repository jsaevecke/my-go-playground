package main

import (
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)

// Bubbletea is based on the Elm Architecture and is a Go library for building terminal applications. It is easier to think about
// MVC (Model-View-Controller) in the context of Bubbletea.

// Model: The state of the application. It can be a simple data structure or a complex one.
// View: A function that takes the model and returns a string to be displayed in the terminal.
// Update: A function that takes the model and a message and returns a new model. It is responsible for updating the state of the application.
// Message: A message that describes an event that has occurred in the application. It can be a user input, a timer event, etc.

type Model struct {
	fp           filepicker.Model
	selectedFile string
	quitting     bool
	err          error
}

func (m Model) Init() tea.Cmd {
	return m.fp.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.fp, cmd = m.fp.Update(msg)

	if ok, file := m.fp.DidSelectFile(msg); ok {
		m.selectedFile = file
	}

	return m, cmd
}

func (m Model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	var s strings.Builder
	if m.err != nil {
		s.WriteString("Error: " + m.err.Error() + "\n")
	} else if m.selectedFile == "" {
		s.WriteString("No file selected.\n")
	} else {
		s.WriteString("Selected file: " + m.selectedFile + "\n")
	}

	s.WriteString(m.fp.View())

	return s.String()
}

func main() {
	// Set up logging to a file for debugging purposes because the TUI occupies the stdout and stderr streams.
	// 'tail -f debug.log' in another terminal to see the logs in real time.
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fp := filepicker.New()
	fp.AllowedTypes = []string{"sql", "csv"}
	fp.CurrentDirectory = workingDir

	p := tea.NewProgram(Model{
		fp: fp,
	}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
