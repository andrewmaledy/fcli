package tui

import (
	"fmt"
	"strings"
)

// FormatSizeGB converts bytes to a human-readable GB string.
func FormatSizeGB(bytes int) string {
	gb := float64(bytes) / (1024 * 1024 * 1024)
	if gb >= 100 {
		return fmt.Sprintf("%.0f GB", gb)
	}
	return fmt.Sprintf("%.1f GB", gb)
}

// FormatRating formats a rating value with source-appropriate display.
func FormatRating(value float32, source string) string {
	if value == 0 {
		return "-"
	}
	switch source {
	case "imdb":
		return fmt.Sprintf("%.1f", value)
	case "tmdb":
		return fmt.Sprintf("%.0f%%", value*10)
	case "rt":
		return fmt.Sprintf("%d%%", int(value))
	default:
		return fmt.Sprintf("%.1f", value)
	}
}

// FormatGenres joins genres with commas, truncated to maxLen.
func FormatGenres(genres []string, maxLen int) string {
	joined := strings.Join(genres, ", ")
	if len(joined) > maxLen {
		return joined[:maxLen-3] + "..."
	}
	return joined
}

// FormatEpisodes renders "downloaded/total" episode count.
func FormatEpisodes(downloaded, total int) string {
	return fmt.Sprintf("%d/%d", downloaded, total)
}

// TruncateString truncates a string to maxLen with ellipsis.
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
