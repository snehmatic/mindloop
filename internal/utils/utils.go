package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"reflect"
	"strings"

	"github.com/snehmatic/mindloop/internal/log"

	"github.com/olekukonko/tablewriter"
)

var logger = log.Get()

func PrettyPrint(x any) {
	b, err := json.MarshalIndent(x, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print(string(b))
}

func PrintTable(data interface{}) {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		fmt.Println("Input must be a slice of structs")
		logger.Error().Msg("Input to PrintTable must be a slice of structs")
		return
	}

	if v.Len() == 0 {
		fmt.Println("No records found.")
		logger.Info().Msg("len 0 of the provided data slice")
		return
	}

	first := v.Index(0)
	if first.Kind() != reflect.Struct {
		fmt.Println("Slice elements must be structs, type mismatch")
		logger.Error().Msg("Slice elements must be structs, type mismatch")
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
	logger.Info().Msgf("Rendered table with %d records of type %s", v.Len(), first.Type())
}

func WriteResponse(data interface{}, respWriter http.ResponseWriter, status int) {
	respWriter.Header().Set("content-type", "application/json; charset=utf-8")
	respWriter.WriteHeader(status)
	json.NewEncoder(respWriter).Encode(data)
}

func GetEnvOrDie(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	logger.Fatal().Str("key", key).Msg("failed to get environment variable")
	return ""
}

func FileExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
func FileWrite(filename string, data []byte) error {
	if err := os.WriteFile(filename, data, 0644); err != nil {
		logger.Error().Err(err).Str("file", filename).Msg("failed to write file")
		return err
	}
	logger.Info().Str("file", filename).Msg("file written successfully")
	return nil
}
func FileRead(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		logger.Error().Err(err).Str("file", filename).Msg("failed to read file")
		return nil, err
	}
	logger.Info().Str("file", filename).Msg("file read successfully")
	return data, nil
}
func FileDelete(filename string) error {
	if err := os.Remove(filename); err != nil {
		logger.Error().Err(err).Str("file", filename).Msg("failed to delete file")
		return err
	}
	logger.Info().Str("file", filename).Msg("file deleted successfully")
	return nil
}
