package tui

import "fmt"

// RenderConnectionError returns a styled error box for API connection failures.
func RenderConnectionError(serviceName string, endpoint string, err error) string {
	title := ErrorBoxTitleStyle.Render("CONNECTION ERROR")

	body := fmt.Sprintf(
		"%s\n\n"+
			"Failed to connect to %s API\n\n"+
			"Endpoint:  %s\n"+
			"Error:     %s\n\n"+
			"Please check:\n"+
			"  - %s is running and accessible\n"+
			"  - API endpoint URL is correct\n"+
			"  - API key is valid\n"+
			"  - Network connectivity",
		title, serviceName,
		endpoint,
		TruncateString(err.Error(), 50),
		serviceName,
	)

	return ErrorBoxStyle.Render(body)
}
