package ui

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/rs/zerolog"
)

// Interface defines the UI operations
type Interface interface {
	ShowSuccess(message string)
	ShowError(message string)
	ShowInfo(message string)
	ShowWarning(message string)
	ShowRocket(message string)
	ShowLoading(message string)
	ShowTable(data interface{})
	ShowBanner()
	CaptureJournalWithEditor() (string, error)
	ConfirmAction(message string) bool
}

// CLIInterface implements Interface for command-line interface
type CLIInterface struct {
	logger zerolog.Logger
}

// NewCLIInterface creates a new CLI interface
func NewCLIInterface(logger zerolog.Logger) Interface {
	return &CLIInterface{
		logger: logger,
	}
}

var (
	green          = "\033[32m"
	red            = "\033[31m"
	reset          = "\033[0m"
	tick           = "✓"
	cross          = "✖"
	loading        = "↻"
	greenTickSmall = fmt.Sprintf("%s%s%s", green, tick, reset)
	redCrossSmall  = fmt.Sprintf("%s%s%s", red, cross, reset)

	// emojis
	greenTick = "✅"
	redCross  = "❌"
	warn      = "⚠️"
	rocket    = "🚀"
	timeSand  = "⏳"
	bulb      = "💡"

	// banner
	banner = `
 _  _  __  __ _  ____  __     __    __  ____ 
( \/ )(  )(  ( \(    \(  )   /  \  /  \(  _ \
/ \/ \ )( /    / ) D (/ (_/\(  O )(  O )) __/
\_)(_/(__)\_)__)(____/\____/ \__/  \__/(__)  
`
)

func (ui *CLIInterface) ShowSuccess(message string) {
	fmt.Fprintln(os.Stdout, greenTick, message)
}

func (ui *CLIInterface) ShowError(message string) {
	fmt.Fprintln(os.Stdout, redCross, message)
}

func (ui *CLIInterface) ShowInfo(message string) {
	fmt.Fprintln(os.Stdout, bulb, message)
}

func (ui *CLIInterface) ShowWarning(message string) {
	fmt.Fprintln(os.Stdout, warn, message)
}

func (ui *CLIInterface) ShowRocket(message string) {
	fmt.Fprintln(os.Stdout, rocket, message)
}

func (ui *CLIInterface) ShowLoading(message string) {
	fmt.Fprintln(os.Stdout, timeSand, message)
}

func (ui *CLIInterface) ShowBanner() {
	greenBanner := fmt.Sprintf("%s%s%s", green, banner, reset)
	fmt.Println(greenBanner)
}

func (ui *CLIInterface) ShowTable(data interface{}) {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		fmt.Println("Input must be a slice of structs")
		ui.logger.Error().Msg("Input to ShowTable must be a slice of structs")
		return
	}

	if v.Len() == 0 {
		fmt.Println("No records found.")
		ui.logger.Info().Msg("len 0 of the provided data slice")
		return
	}

	first := v.Index(0)
	if first.Kind() != reflect.Struct {
		fmt.Println("Slice elements must be structs, type mismatch")
		ui.logger.Error().Msg("Slice elements must be structs, type mismatch")
		return
	}

	// Extract headers
	var headers []string
	t := first.Type()
	for i := 0; i < t.NumField(); i++ {
		headers = append(headers, strings.ToUpper(t.Field(i).Name))
	}

	// Extract data
	var rows [][]string
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		var row []string
		for j := 0; j < elem.NumField(); j++ {
			val := elem.Field(j)
			row = append(row, fmt.Sprintf("%v", val.Interface()))
		}
		rows = append(rows, row)
	}

	// Print in table format
	table := tablewriter.NewWriter(os.Stdout)
	table.Header(headers)
	table.Bulk(rows)
	table.Render()
	ui.logger.Info().Msgf("Rendered table with %d records of type %s", v.Len(), first.Type())
}

func (ui *CLIInterface) CaptureJournalWithEditor() (string, error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	tmpFile, err := os.CreateTemp("", "mindloop_journal_*.md")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name())

	header := "# Mindloop Journal\n# Write your thoughts below. Lines starting with # will be ignored.\n\n"
	tmpFile.WriteString(header)
	tmpFile.Close()

	cmd := exec.Command(editor, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}

	data, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(data), "\n")
	var content strings.Builder
	for _, line := range lines {
		if !strings.HasPrefix(line, "#") {
			content.WriteString(line + "\n")
		}
	}

	return strings.TrimSpace(content.String()), nil
}

func (ui *CLIInterface) ConfirmAction(message string) bool {
	fmt.Printf("%s %s (type 'yes' to confirm): ", bulb, message)
	var confirmation string
	fmt.Scanln(&confirmation)
	return confirmation == "yes"
}
