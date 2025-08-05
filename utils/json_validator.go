package utils

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

func IsValidURL(urlStr string) error {
	if urlStr == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	parsed, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	if parsed.Scheme == "" {
		return fmt.Errorf("URL must include scheme (http or https)")
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("URL scheme must be http or https")
	}

	if parsed.Host == "" {
		return fmt.Errorf("URL must include host")
	}

	return nil
}

func IsValidJSON(data string) error {
	data = strings.TrimSpace(data)
	if data == "" {
		return fmt.Errorf("response body is empty")
	}

	var js interface{}
	if err := json.Unmarshal([]byte(data), &js); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}

	return nil
}

func IsJSONContentType(contentType string) bool {
	contentType = strings.ToLower(strings.TrimSpace(contentType))
	return strings.Contains(contentType, "application/json") ||
		strings.Contains(contentType, "text/json")
}