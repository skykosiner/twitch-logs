package logs

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type AvailableLogs struct {
	AvailableLogs []struct {
		Year  string `json:"year"`
		Month string `json:"month"`
	} `json:"availableLogs"`
}

type Messages struct {
	Messages []struct {
		Text        string `json:"text"`
		DisplayName string `json:"displayName"`
		Timestamp   string `json:"timestamp"`
	} `json:"messages"`
}

func (a AvailableLogs) GetStringSlice() []string {
	var dates []string

	for _, d := range a.AvailableLogs {
		dates = append(dates, fmt.Sprintf("%s %s", d.Year, d.Month))

	}

	return dates
}

func (m Messages) GetStringSlice() []string {
	var messages []string

	for _, m := range m.Messages {
		messages = append(messages, fmt.Sprintf("%s %s %s", m.Timestamp, m.DisplayName, m.Text))
	}

	return messages
}

func GetLogsAvailable(channel, username string) AvailableLogs {
	var availableLogs AvailableLogs
	url := fmt.Sprintf("https://logs.zonian.dev/list?channel=%s&user=%s", channel, username)

	resp, err := http.Get(url)
	if err != nil {
		slog.Error("Error getting logs", "error", err)
		return availableLogs
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error getting logs", "error", err)
		return availableLogs
	}

	if err := json.Unmarshal(bytes, &availableLogs); err != nil {
		slog.Error("Error getting logs", "error", err)
		return availableLogs
	}

	return availableLogs
}

func GetLogs(channel, username, year, month string) Messages {
	var messages Messages
	url := fmt.Sprintf("https://logs.zonian.dev/channel/%s/user/%s/%s/%s?jsonBasic=1", channel, username, year, month)

	resp, err := http.Get(url)
	if err != nil {
		slog.Error("Error getting logs", "error", err)
		return messages
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error getting logs", "error", err)
		return messages
	}

	if err := json.Unmarshal(bytes, &messages); err != nil {
		slog.Error("Error getting logs", "error", err)
		return messages
	}

	return messages
}
