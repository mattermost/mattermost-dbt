package tui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mattermost/mattermost-dbt/internal/store"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type monitorModel struct {
	nodeStores []*store.PgedgeNodeStore
	table      *table.Table
	timer      timer.Model
	timeout    time.Duration
	keymap     keymap
	help       help.Model
	quitting   bool
	logger     log.FieldLogger
}

type keymap struct {
	quit key.Binding
}

func (m monitorModel) Init() tea.Cmd {
	return m.timer.Init()
}

func (m monitorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case timer.TimeoutMsg:
		// Refresh the table view and restart timer.
		table, err := buildMonitorTable(m.nodeStores)
		if err != nil {
			m.logger.WithError(err).Error("Error encountered during monitoring")
			return m, tea.Quit
		}
		m.table = table
		m.timer.Timeout = m.timeout
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			m.quitting = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m monitorModel) helpView() string {
	return "\n" + m.help.ShortHelpView([]key.Binding{m.keymap.quit})
}

func (m monitorModel) View() string {
	s := m.table.String()
	s += "\n"
	s += "\n"
	s += "Refreshing in " + m.timer.View()
	s += m.helpView()

	return s
}

func StartMonitoring(timeout time.Duration, nodeStores []*store.PgedgeNodeStore, logger log.FieldLogger) {
	table, err := buildMonitorTable(nodeStores)
	if err != nil {
		logger.WithError(err).Error("Failed to build monitoring table")
		os.Exit(1)
	}

	m := monitorModel{
		nodeStores: nodeStores,
		table:      table,
		timer:      timer.NewWithInterval(timeout, time.Second),
		timeout:    timeout,
		keymap: keymap{
			quit: key.NewBinding(
				key.WithKeys("q", "ctrl+c"),
				key.WithHelp("q", "quit"),
			),
		},
		help:   help.New(),
		logger: logger,
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		logger.WithError(err).Error("Error encountered during monitoring")
		os.Exit(1)
	}
}

func buildMonitorTable(nodeStores []*store.PgedgeNodeStore) (*table.Table, error) {
	data := [][]string{}
	for _, nodeStore := range nodeStores {
		start := time.Now()

		version, err := nodeStore.Store.GetSpockVersion()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get spock version")
		}

		versionNum, err := nodeStore.Store.GetSpockVersionNum()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get spock version num")
		}

		status, err := nodeStore.Store.GetSpockReplicationStatus()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get spock replication status")
		}

		lag, err := nodeStore.Store.GetSpockLag()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get spock lag")
		}

		connections, err := nodeStore.Store.GetConnectionCount()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get connection count")
		}

		row := []string{
			nodeStore.Node.Name,
			fmt.Sprintf("%s (%s)", version.Version, versionNum.VersionNum),
			fmt.Sprintf("%d milliseconds", time.Since(start).Milliseconds()),
			fmt.Sprintf("%s [%s]", status.SubscriptionName, status.Status),
			fmt.Sprintf("%s [%s]", lag.CommitLSN, lag.ReplicationLag),
			fmt.Sprintf("%d", connections),
		}

		data = append(data, row)
	}

	re := lipgloss.NewRenderer(os.Stdout)
	baseStyle := re.NewStyle().Padding(0, 1)
	headerStyle := baseStyle.Foreground(lipgloss.Color("252")).Bold(true)
	headers := []string{"Node", "Spock", "Compute Time", "Replication", "Lag", "DB Conns"}

	capitalizeHeaders := func(data []string) []string {
		for i := range data {
			data[i] = strings.ToUpper(data[i])
		}
		return data
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(re.NewStyle().Foreground(lipgloss.Color("238"))).
		Headers(capitalizeHeaders(headers)...).
		Rows(data...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				return headerStyle
			}
			return baseStyle.Foreground(lipgloss.Color("252"))

		})

	return t, nil
}
