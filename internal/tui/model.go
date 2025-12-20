package tui

import (
	"lazycurl/internal/curl"
	"lazycurl/internal/load"
	"lazycurl/internal/model"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Pane represents a focusable area of the TUI.
type Pane int

const (
	PaneRequests Pane = iota
	PaneEditor
	PaneResponse
)

// EditorField represents a field in the editor pane.
type EditorField int

const (
	FieldMethod EditorField = iota
	FieldURL
	FieldTabs    // Focus on the tab bar
	FieldContent // Focus on the content (Body or Headers list)
)

// EditorTab represents the active tab in the editor.
type EditorTab int

const (
	TabHeaders EditorTab = iota
	TabBody
	TabLoad
)

// InputPair represents a key-value input pair (for headers).
type InputPair struct {
	Key   textinput.Model
	Value textinput.Model
}

type LoadConfig struct {
	Concurrency textinput.Model
	Duration    textinput.Model
}

type LoadState struct {
	IsRunning bool
	Stats     *load.Stats
	Sub       chan load.StatsMsg // Active subscription
}

// Model represents the state of the TUI.
type Model struct {
	// Global State
	ActivePane Pane
	KeyMap     KeyMap
	Help       help.Model
	Width      int
	Height     int
	Executor   *curl.Executor

	// Requests Pane State
	Requests       []model.Request
	SelectedReqIdx int

	// Editor Pane State
	ActiveEditorTab  EditorTab
	EditorInputs     []textinput.Model // Method, URL
	HeaderInputs     []InputPair       // Headers
	EditorBody       textarea.Model    // Body
	LoadConfig       LoadConfig        // Load Config
	FocusedField     EditorField
	FocusedHeaderIdx int  // Index of the header being edited
	FocusedHeaderKey bool // True if editing Key, False if Value
	IsEditing        bool // True if user is typing in a field

	// Load Test State
	LoadState LoadState

	// Response Pane State
	Response *model.Response
}

// NewModel creates the initial model.
func NewModel() Model {
	// Initialize inputs
	methodInput := textinput.New()
	methodInput.Placeholder = "GET"
	methodInput.CharLimit = 10

	urlInput := textinput.New()
	urlInput.Placeholder = "https://example.com"
	urlInput.Width = 50

	bodyInput := textarea.New()
	bodyInput.Placeholder = "Request Body (JSON)..."
	bodyInput.SetHeight(10)
	bodyInput.ShowLineNumbers = false

	// Load Inputs
	concInput := textinput.New()
	concInput.Placeholder = "10"
	concInput.SetValue("10")
	concInput.CharLimit = 5

	durInput := textinput.New()
	durInput.Placeholder = "5s"
	durInput.SetValue("5s")
	durInput.CharLimit = 5

	// Initial default request
	defaultReq := model.NewRequest()

	m := Model{
		ActivePane:       PaneRequests,
		KeyMap:           DefaultKeyMap(),
		Help:             help.New(),
		Executor:         curl.NewExecutor(),
		Requests:         []model.Request{defaultReq},
		SelectedReqIdx:   0,
		ActiveEditorTab:  TabBody, // Default to Body
		EditorInputs:     []textinput.Model{methodInput, urlInput},
		EditorBody:       bodyInput,
		LoadConfig:       LoadConfig{Concurrency: concInput, Duration: durInput},
		FocusedField:     FieldMethod,
		IsEditing:        false,
		FocusedHeaderIdx: 0,
		HeaderInputs:     []InputPair{},
		LoadState:        LoadState{IsRunning: false},
	}

	m.SyncEditorToRequest()
	return m
}

// SyncEditorToRequest populates editor fields from the currently selected request.
func (m *Model) SyncEditorToRequest() {
	if len(m.Requests) == 0 {
		return
	}
	req := m.Requests[m.SelectedReqIdx]
	m.EditorInputs[0].SetValue(req.Method)
	m.EditorInputs[1].SetValue(req.URL)
	m.EditorBody.SetValue(req.Body)

	// Sync Headers
	m.HeaderInputs = []InputPair{}
	for k, v := range req.Headers {
		kInput := textinput.New()
		kInput.Placeholder = "Header"
		kInput.SetValue(k)

		vInput := textinput.New()
		vInput.Placeholder = "Value"
		vInput.SetValue(v)

		m.HeaderInputs = append(m.HeaderInputs, InputPair{Key: kInput, Value: vInput})
	}
	// Always have at least one empty row for new headers if empty
	if len(m.HeaderInputs) == 0 {
		kInput := textinput.New()
		kInput.Placeholder = "Header"
		vInput := textinput.New()
		vInput.Placeholder = "Value"
		m.HeaderInputs = append(m.HeaderInputs, InputPair{Key: kInput, Value: vInput})
	}
}

// SyncRequestToEditor updates the selected request from editor fields.
func (m *Model) SyncRequestToEditor() {
	if len(m.Requests) == 0 {
		return
	}
	req := &m.Requests[m.SelectedReqIdx]
	req.Method = m.EditorInputs[0].Value()
	req.URL = m.EditorInputs[1].Value()
	req.Body = m.EditorBody.Value()

	// Sync Headers
	req.Headers = make(map[string]string)
	for _, pair := range m.HeaderInputs {
		k := pair.Key.Value()
		v := pair.Value.Value()
		if k != "" {
			req.Headers[k] = v
		}
	}
}

// Init initializes the model.
func (m Model) Init() tea.Cmd {
	return nil
}
