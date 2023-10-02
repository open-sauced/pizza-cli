package utils

import (
	"encoding/json"
	"strings"

	bubblesTable "github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v3"
)

func OutputJSON(entity interface{}) (string, error) {
	output, err := json.MarshalIndent(entity, "", "  ")
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func OutputYAML(entity interface{}) (string, error) {
	output, err := yaml.Marshal(entity)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(output), "\n"), nil
}

func OutputTable(rows []bubblesTable.Row, columns []bubblesTable.Column) string {
	styles := bubblesTable.Styles{
		Cell:     lipgloss.NewStyle().PaddingRight(1),
		Header:   lipgloss.NewStyle().Bold(true).PaddingRight(1),
		Selected: lipgloss.NewStyle(),
	}
	table := bubblesTable.New(
		bubblesTable.WithRows(rows),
		bubblesTable.WithColumns(columns),
		bubblesTable.WithHeight(len(rows)),
		bubblesTable.WithStyles(styles),
	)
	return table.View()
}

func GetMaxTableRowWidth(rows []bubblesTable.Row) int {
	var maxRowWidth int
	for i := range rows {
		rowWidth := len(rows[i][0])
		if rowWidth > maxRowWidth {
			maxRowWidth = rowWidth
		}
	}
	return maxRowWidth
}
