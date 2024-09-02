# fcli

`fcli` is a command-line interface (CLI) tool for managing media and tv shows data through the Radarr, Overseer, and Sonarr APIs. It provides functionalities to quickly query and delete content that is taking up the most disk space. This tool is very early stages and has very limited functionality.

## Features

- **Manage Movies:**
  - Retrieve top movies (default to top 10, specify with --limit flag)
  - Delete movies from Radarr.
  - Delete related requests from Overseer (or Jellyseer)

- **Manage TV Series:**
  - Retrieve top shows/series (default to top 10, specify with --limit flag).
  - Select and delete entire series or specific seasons.
  - If partial deletion (only a specific season is deleted), then updated sonarr to not track that particular season.

- **Configuration:**
  - Supports configuration via `.fcli-config` file (should reside in home directory)
  - Environment variable overrides if preferred.

## Installation

### Building from Source

1. Clone the repository:

   ```bash
   git clone https://github.com/andrewmaledy/fcli.git
   cd fcli
   go build -o fcli main.go
   sudo mv fcli /usr/local/bin/```

### Example Configuration File

Example `.fcli-config`

```
radarr:
  base_url: "http://localhost:7878"
  api_key: "your-radarr-api-key"

overseer:
  base_url: "http://localhost:5055"
  api_key: "your-overseer-api-key"

sonarr:
  base_url: "http://localhost:8989"
  api_key: "your-sonarr-api-key"
```
