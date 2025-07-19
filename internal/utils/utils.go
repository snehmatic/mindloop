package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
)

func PrettyPrint(x any) {
	b, err := json.MarshalIndent(x, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print(string(b))
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
	log.Fatal().Str("key", key).Msg("failed to get environment variable")
	return ""
}
