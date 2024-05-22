# RSS File Downloader / CLI

[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/philippheuer/rssdownloader/badge)](https://securityscorecards.dev/viewer/?uri=github.com/philippheuer/rssdownloader)

> A flexible cli to download files from rss feeds, with support for filtering and templating.

## Download

Download the binary from the [GitHub Releases](https://github.com/PhilippHeuer/rssdownloader/releases).

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

> The download command will create a feed-state.json in the output directory to track the timestamp of the last download, only downloading newly added files.

## Usage

```bash
rssdownloader download --config feeds.yaml
```

## License

Released under the [MIT license](./LICENSE).
