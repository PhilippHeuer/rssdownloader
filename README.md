# RSS File Downloader / CLI

A configurable cli to download files from a feed.

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
    # rules that must match
    rules:
      - type: regex
        value: ".*"
    # exclude items that matched the rules above
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
