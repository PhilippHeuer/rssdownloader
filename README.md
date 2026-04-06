# RSS File Downloader / CLI

[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/PhilippHeuer/rssdownloader/badge)](https://securityscorecards.dev/viewer/?uri=github.com/PhilippHeuer/rssdownloader)

> A flexible cli to download files from rss feeds, with support for filtering and templating.

## Download

Download the binary from the [GitHub Releases](https://github.com/PhilippHeuer/rssdownloader/releases).

## Features

- **RSS/Atom Feed Parsing** - Supports standard RSS 2.0 and Atom feeds
- **Filtering** - Include/exclude based on complex regex rules
- **Filename Templating** - Customizable filenames with placeholders
- **State Tracking** - Remembers last download timestamp, only fetches new items

## Configuration

```yaml
feeds:
  - name: my-feed
    enabled: true
    output: /target-dir
    url: https://example.com/my-feed
    # use item title as filename
    template: "{title}"
    # rules, if set items must match at least one rule
    rules:
      - type: regex
        value: ".*"
    # exclude, all items matching a rule will be excluded (always has precedence over rules)
    exclude:
      - type: regex
        value: "prefix.*"
```

### Configuration Options

| Field      | Required | Description                             |
|------------|----------|-----------------------------------------|
| `name`     | Yes      | Unique feed identifier                  |
| `url`      | Yes      | RSS/Atom feed URL                       |
| `output`   | Yes      | Output directory (auto-created)         |
| `template` | Yes      | Filename template                       |
| `enabled`  | No       | Enable/disable feed (default: true)     |
| `rules`    | No       | Include rules (regex)                   |
| `exclude`  | No       | Exclude rules (regex, takes precedence) |

## Usage

```bash
# Download all enabled feeds
rssdownloader download --config feeds.yaml

# Validate configuration
rssdownloader validate --config feeds.yaml

# List configured feeds
rssdownloader list --config feeds.yaml
```

## License

Released under the [MIT license](./LICENSE).
