// Package utils contains general purpose helper functions for Mindloop
package utils

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"os/exec"
	"time"

	"reflect"
	"strings"

	"github.com/snehmatic/mindloop/internal/log"

	"github.com/olekukonko/tablewriter"
)

var logger = log.Get()

var (
	green = "\033[32m"
	reset = "\033[0m"

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

// PrettyPrintBanner prints the Mindloop ASCII art banner
func PrettyPrintBanner() {
	greenBanner := fmt.Sprintf("%s%s%s", green, banner, reset)
	fmt.Println(greenBanner)
}

// PrettyPrint prints any data as a pretty JSON string
func PrettyPrint(x any) {
	b, err := json.MarshalIndent(x, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print(string(b))
}

// PrintTable renders a slice of structs as a terminal table
func PrintTable(data interface{}) {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		fmt.Println("Input must be a slice of structs")
		logger.Error("Input to PrintTable must be a slice of structs", nil)

		return
	}

	if v.Len() == 0 {
		fmt.Println("No records found.")
		logger.Info("len 0 of the provided data slice")

		return
	}

	first := v.Index(0)
	if first.Kind() != reflect.Struct {
		fmt.Println("Slice elements must be structs, type mismatch")
		logger.Error("Slice elements must be structs, type mismatch", nil)

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
	defer func(table *tablewriter.Table) {
		err := table.Close()
		if err != nil {
			logger.Error("Error closing table", err)
		}
	}(table)

	table.Header(headers)
	_ = table.Bulk(rows)
	_ = table.Render()
	logger.Info(fmt.Sprintf("Rendered table with %d records of type %s", v.Len(), first.Type()))
}

// PrintSuccessln prints a success message with a checkmark
func PrintSuccessln(a ...any) {
	if len(a) == 0 {
		_, _ = fmt.Fprintln(os.Stdout, greenTick)

		return
	}
	_, _ = fmt.Fprintln(os.Stdout, append([]any{greenTick}, a...)...)
}

// PrintSuccessf prints a formatted success message with a checkmark
func PrintSuccessf(format string, a ...any) {
	if len(a) == 0 {
		_, _ = fmt.Fprintf(os.Stdout, greenTick+" "+format, a...)

		return
	}
	_, _ = fmt.Fprintf(os.Stdout, greenTick+" "+format, a...)
}

// PrintRocketln prints a message with a rocket emoji
func PrintRocketln(a ...any) {
	if len(a) == 0 {
		_, _ = fmt.Fprintln(os.Stdout, rocket)

		return
	}
	_, _ = fmt.Fprintln(os.Stdout, append([]any{rocket}, a...)...)
}

// PrintRocketf prints a formatted message with a rocket emoji
func PrintRocketf(format string, a ...any) {
	if len(a) == 0 {
		_, _ = fmt.Fprintf(os.Stdout, rocket+" "+format, a...)

		return
	}
	_, _ = fmt.Fprintf(os.Stdout, rocket+" "+format, a...)
}

// PrintInfoln prints an informational message with an info symbol
func PrintInfoln(a ...any) {
	if len(a) == 0 {
		_, _ = fmt.Fprintln(os.Stdout, bulb)
		return
	}
	_, _ = fmt.Fprintln(os.Stdout, append([]any{bulb}, a...)...)
}

// PrintInfof prints a formatted informational message with an info symbol
func PrintInfof(format string, a ...any) {
	if len(a) == 0 {
		_, _ = fmt.Fprintf(os.Stdout, bulb+" "+format, a...)
		return
	}
	_, _ = fmt.Fprintf(os.Stdout, bulb+" "+format, a...)
}

// PrintLoadingln prints a message with a loading/gear symbol
func PrintLoadingln(a ...any) {
	if len(a) == 0 {
		_, _ = fmt.Fprintln(os.Stdout, timeSand)
		return
	}
	_, _ = fmt.Fprintln(os.Stdout, append([]any{timeSand}, a...)...)
}

// PrintLoadingf prints a formatted message with a loading/gear symbol
func PrintLoadingf(format string, a ...any) {
	if len(a) == 0 {
		_, _ = fmt.Fprintf(os.Stdout, timeSand+" "+format, a...)
		return
	}
	_, _ = fmt.Fprintf(os.Stdout, timeSand+" "+format, a...)
}

// PrintWarnln prints a warning message with a warning symbol
func PrintWarnln(a ...any) {
	if len(a) == 0 {
		_, _ = fmt.Fprintln(os.Stdout, warn)
		return
	}
	_, _ = fmt.Fprintln(os.Stdout, append([]any{warn}, a...)...)
}

// PrintWarnf prints a formatted warning message with a warning symbol
func PrintWarnf(format string, a ...any) {
	if len(a) == 0 {
		_, _ = fmt.Fprintf(os.Stdout, warn+" "+format, a...)
		return
	}
	_, _ = fmt.Fprintf(os.Stdout, warn+" "+format, a...)
}

// PrintErrorln prints an error message with a cross mark
func PrintErrorln(a ...any) {
	if len(a) == 0 {
		_, _ = fmt.Fprintln(os.Stdout, redCross)
		return
	}
	_, _ = fmt.Fprintln(os.Stdout, append([]any{redCross}, a...)...)
}

// PrintErrorf prints a formatted error message with a cross mark
func PrintErrorf(format string, a ...any) {
	if len(a) == 0 {
		_, _ = fmt.Fprintf(os.Stdout, redCross+" "+format, a...)
		return
	}
	_, _ = fmt.Fprintf(os.Stdout, redCross+" "+format, a...)
}

// WriteResponse writes a JSON response to the http.ResponseWriter
func WriteResponse(data interface{}, respWriter http.ResponseWriter, status int) {
	respWriter.Header().Set("content-type", "application/json; charset=utf-8")
	respWriter.WriteHeader(status)
	_ = json.NewEncoder(respWriter).Encode(data)
}

// GetEnvOrDie retrieves an environment variable or exits if not set
func GetEnvOrDie(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	logger.Fatal("failed to get environment variable", nil, log.Field{
		Key:   "key",
		Value: key,
	})
	return ""
}

// FileExists checks if a file exists at the given path
func FileExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

// FileWrite writes data to a file
func FileWrite(filename string, data []byte) error {
	if err := os.WriteFile(filename, data, 0644); err != nil {
		logger.Error("failed to write file", err,
			log.Field{
				Key:   "file",
				Value: filename,
			},
		)
		return err
	}
	logger.Info("file written successfully",
		log.Field{
			Key:   "file",
			Value: filename,
		},
	)
	return nil
}

// FileRead reads data from a file
func FileRead(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		logger.Error("failed to read file", err, log.Field{
			Key:   "file",
			Value: filename,
		})
		return nil, err
	}
	logger.Info("file read successfully", log.Field{
		Key:   "file",
		Value: filename,
	})
	return data, nil
}

// FileDelete deletes a file at the given path
func FileDelete(filename string) error {
	if err := os.Remove(filename); err != nil {
		logger.Error("failed to delete file", err, log.Field{
			Key:   "file",
			Value: filename,
		})
		return err
	}
	logger.Info("file deleted successfully", log.Field{
		Key:   "file",
		Value: filename,
	})
	return nil
}

// CaptureJournalWithEditor launches the user's default editor to capture journal content
func CaptureJournalWithEditor() (string, error) {
	header := "# Mindloop Journal\n# Write your thoughts below. Lines starting with # will be ignored.\n\n"
	return CaptureWithEditor("mindloop_journal_*.md", header, "")
}

// CaptureWithEditor launches the user's default editor to capture content
func CaptureWithEditor(filenamePattern, header, initialContent string) (string, error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	tmpFile, err := os.CreateTemp("", filenamePattern)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			logger.Error("failed to delete temp file", err, log.Field{
				Key:   "file",
				Value: tmpFile.Name(),
			})
		}
	}()

	if header != "" {
		_, _ = tmpFile.WriteString(header)
	}
	if initialContent != "" {
		_, _ = tmpFile.WriteString(initialContent)
	}
	if err := tmpFile.Close(); err != nil {
		return "", err
	}

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

// GetDateRange returns start and end time for common periods (daily, weekly, yearly)
func GetDateRange(period string) (time.Time, time.Time) {
	now := time.Now()
	switch period {
	case "yearly", "year":
		return time.Date(now.Year()-1, now.Month(), now.Day(), 0, 0, 0, 0, now.Location()), now
	case "weekly", "week":
		end := time.Now()
		start := end.AddDate(0, 0, -7)
		start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
		return start, end
	case "daily", "day":
		return now.Add(-24 * time.Hour), now
	default:
		return now.Add(-24 * time.Hour), now
	}
}

// FormatMinutes converts float64 minutes into a human-readable string like "1hr 2min"
func FormatMinutes(minutes float64) string {
	totalMinutes := int(math.Round(minutes))
	hours := totalMinutes / 60
	mins := totalMinutes % 60

	switch {
	case hours > 0 && mins > 0:
		return fmt.Sprintf("%dhr %dmin", hours, mins)
	case hours > 0:
		return fmt.Sprintf("%dhr", hours)
	default:
		return fmt.Sprintf("%dmin", mins)
	}
}
