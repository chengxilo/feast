package file

import (
	"feast/logger"
	"feast/types"
	"feast/ui"
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go.uber.org/zap"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strconv"
	"time"
)

type detail struct {
	name         string
	dateModified time.Time
	type_        string
	size         int64
	mode         string
}

func (f detail) toTableRow() table.Row {
	return table.Row{
		f.name,
		f.dateModified.Format("2006-01-02 15:04:05"),
		f.type_,
		strconv.FormatInt(f.size, 10),
		f.mode,
	}
}

func (m *Model) getFiles() ([]detail, error) {
	dir, err := os.ReadDir(m.path)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir: %e", err)
	}
	details := make([]detail, 0, len(dir))
	log.Debug("get files", zap.String("path", m.path), zap.String("detail", fmt.Sprint(details)))
	for _, file := range dir {
		info, err := file.Info()
		if err != nil {
			return nil, fmt.Errorf("failed to read info: %e", err)
		}
		details = append(details, detail{
			name:         file.Name(),
			dateModified: info.ModTime(),
			type_: func() string {
				if info.IsDir() {
					return "directory"
				} else {
					return "file"
				}
			}(),
			size: info.Size(),
			mode: info.Mode().String(),
		})
	}
	return details, nil
}

var log = logger.Logger.With(zap.String("FileModel", "file"))

type Model struct {
	// file path
	path string
	// right arrow target
	rightArrowTargetList []string
	table                table.Model
	height               int
	width                int
}

func (m *Model) SetWidth(width int) {
	m.width = width
}

func (m *Model) SetHeight(height int) {
	m.height = height
}

func NewModel() ui.Model {
	m := &Model{}
	initPath, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("failed to find user home directory")
	}

	m.path = initPath
	columns := []table.Column{
		{Title: "Name", Width: 25},
		{Title: "Date Modified", Width: 10},
		{Title: "type_", Width: 10},
		{Title: "Size", Width: 10},
		{Title: "Mode", Width: 10},
	}

	var rows []table.Row

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(9),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)
	m.table = t
	return m
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	log.Debug("update", zap.Any("msg", msg))
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "enter":
			if m.table.SelectedRow() == nil {
				return m, nil
			}
			// select a file to explore
			absPath := path.Join(m.path, m.table.SelectedRow()[0])
			f, err := os.Stat(absPath)
			if err != nil {
				log.Error("Failed to get file when select", zap.String("path", absPath), zap.Error(err))
				return m, nil
			}
			if f.IsDir() {
				m.rightArrowTargetList = nil
				// set path to target dir
				m.path = absPath
			}
		case "left":
			m.rightArrowTargetList = slices.Concat([]string{m.path}, m.rightArrowTargetList)
			m.path = filepath.Dir(m.path)
			m.table.SetCursor(0)
		case "right":
			log.Debug("Right arrow target list", zap.String("path", m.path), zap.Strings("rightArrowTargetList", m.rightArrowTargetList))
			if len(m.rightArrowTargetList) != 0 {
				m.path = m.rightArrowTargetList[0]
				m.rightArrowTargetList = m.rightArrowTargetList[1:]
			}
			m.table.SetCursor(0)
		case "q":
			return m, types.RouteCmd("/")
		}
	}
	details, err := m.getFiles()
	if err != nil {
		log.Error("Failed to get file details", zap.Error(err))
		return m, nil
	}

	rows := make([]table.Row, len(details))
	for i, detail := range details {
		rows[i] = detail.toTableRow()
	}

	m.table.SetRows(rows)
	m.table.SetHeight(m.height)

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	s := lipgloss.NewStyle().Width(m.width).Height(m.height)
	tableView := m.table.View()
	return s.Render(tableView)
}
