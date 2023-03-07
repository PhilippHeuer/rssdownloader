# RSS File Downloader / CLI

A configurable cli to download files from a feed.

## Download



## Configuration

```yaml
feeds:
  - name: my-feed
    enabled: true
    output: /target-dir
    url: https://example.com/my-feed
    rules:
      - type: regex
        value: ".*"
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
