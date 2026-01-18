package model

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/risbern21/pokedex/internal/commands"
	"github.com/risbern21/pokedex/internal/styles"
)

const (
	normalView = iota
	inputView
)

type model struct {
	state       uint8
	config      *commands.Config
	focused     bool
	ready       bool
	width       int
	height      int
	textinput   textinput.Model
	inputLabel  string
	inputCaller string
	spinner     spinner.Model
	list        list.Model
	pokedexList list.Model
	viewport    viewport.Model
}

type item struct {
	title string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return i.title }

func New() model {
	ti := textinput.New()
	ti.Focus()
	ti.Width = 30

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styles.SpinnerStyle

	items := []list.Item{
		item{title: "map"},
		item{title: "map back"},
		item{title: "explore"},
		item{title: "catch"},
		item{title: "inspect"},
		item{title: "pokedex"},
		item{title: "help"},
		item{title: "exit"},
	}

	d := list.NewDefaultDelegate()
	d.Styles.SelectedTitle = styles.ListItemSelectedStyle
	d.Styles.SelectedDesc = styles.ListItemSelectedStyle

	l := list.New(items, d, 40, 40)
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)

	pl := list.New(nil, d, 60, 40)
	pl.Title = styles.PokedexListTitleStyle.Render("Your Pokemon")
	pl.SetShowStatusBar(false)

	return model{
		config:      commands.Newconfig(),
		focused:     true,
		width:       100,
		height:      30,
		textinput:   ti,
		spinner:     s,
		list:        l,
		pokedexList: pl,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() string {
	var s string

	switch m.state {
	case normalView:
		var (
			listView     string
			viewportView string
		)
		if m.focused {
			listView = styles.ListFocusedStyle.Render(m.list.View())
		} else {
			listView = styles.ListUnFocusedStyle.Render(m.list.View())
		}

		viewportView = fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
		s += lipgloss.JoinHorizontal(lipgloss.Top, listView, viewportView)
	case inputView:
		s += fmt.Sprintf("%s\n%s\n\n", m.inputLabel, m.textinput.View())
	}

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, s)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	if !m.ready {
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.state == normalView {
		if m.focused {
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
		}

		if !m.focused {
			m.viewport, cmd = m.viewport.Update(msg)
			cmds = append(cmds, cmd)
		}
	}
	m.textinput, cmd = m.textinput.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch m.state {
		case normalView:
			switch key {
			case "ctrl+c", "q", "esc":
				return m, tea.Quit
			case "enter":
				i := m.list.SelectedItem().(item)
				switch i.title {
				case "map", "map back":
					content, err := commands.Commands[i.title].Callback(m.config)
					if err != nil {
						m.viewport.SetContent(err.Error())
						return m, tea.Batch(cmds...)
					}
					m.viewport.SetContent(content)
				case "explore":
					m.inputLabel = "Enter the area name"
					m.textinput.SetValue("")
					m.textinput.Focus()
					m.inputCaller = "explore"
					m.state = inputView
				case "catch", "inspect":
					m.inputLabel = "Enter the pokemons name"
					m.textinput.SetValue("")
					m.textinput.Focus()
					switch i.title {
					case "catch":
						m.inputCaller = "catch"
					case "inspect":
						m.inputCaller = "inspect"
					}
					m.state = inputView
				case "pokedex":
					content, err := commands.Commands["pokedex"].Callback(m.config)
					if err != nil {
						m.viewport.SetContent(err.Error())
						return m, tea.Batch(cmds...)
					}
					m.viewport.SetContent(content)

				case "help":
					var s string

					for _, c := range commands.Commands {
						s += fmt.Sprintf("%s:\n%s\n\n", c.Title, c.Description)
					}

					m.viewport.SetContent(s)
				case "exit":
					return m, tea.Quit
				}
			}
		case inputView:
			switch key {
			case "enter":
				var content string
				var err error

				switch m.inputCaller {
				case "explore":
					locationName := m.textinput.Value()
					m.config.LocationName = &locationName

					content, err = commands.Commands["explore"].Callback(m.config)
				case "catch":
					pokemonName := m.textinput.Value()
					m.config.PokemonName = &pokemonName

					content, err = commands.Commands["catch"].Callback(m.config)
				case "inspect":
					pokemonName := m.textinput.Value()
					m.config.PokemonName = &pokemonName

					content, err = commands.Commands["inspect"].Callback(m.config)
				}

				if err != nil {
					m.viewport.SetContent(err.Error())
					return m, tea.Batch(cmds...)
				}

				m.viewport.SetContent(content)
				m.state = normalView
			}
		}
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		m.width = msg.Width
		m.height = msg.Height

		if !m.ready {
			m.viewport = viewport.New(m.width, m.height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.SetContent("TUI pokedex")
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) headerView() string {
	title := styles.TitleStyle.Render("Pokedex")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m model) footerView() string {
	info := styles.InfoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}
