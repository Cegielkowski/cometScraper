package utils

import (
	"bytes"
	"cometScraper/entity"
	"context"
	"encoding/json"
	"log"
)

func CompactJSON(data []byte) string {
	var js map[string]interface{}
	if json.Unmarshal(data, &js) != nil {
		return string(data)
	}

	result := new(bytes.Buffer)
	if err := json.Compact(result, data); err != nil {
		log.Println(err)
	}
	return result.String()
}

// GetReqID get request id from echo context
func GetReqID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if reqID, ok := ctx.Value(entity.RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}
