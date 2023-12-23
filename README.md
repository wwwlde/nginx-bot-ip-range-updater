# Nginx Bot IP Range Updater

This project fetches IP ranges for various search engine bots (like Googlebot, Bingbot, etc.) and generates an Nginx configuration file to handle these IP ranges. It allows for easy management of Nginx configurations based on updated bot IP ranges.

## Features

- Fetch IP ranges for various bots from specified JSON URLs.
- Generate an Nginx configuration file using a customizable template.
- Easy configuration management through a YAML file.

## Configuration

### YAML Configuration File

The application uses a YAML configuration file to specify the output file path, the template for the Nginx configuration, and the URLs for the IP ranges of different bots.

Example [`example.yaml`](example.yaml):

```yaml
file: /etc/nginx/conf.d/bots.conf
template: |
  geo $searchbot {
      default 0;
      {{- range .Prefixes }}
      {{- if .Ipv4Prefix }}
      {{ .Ipv4Prefix }} 1;
      {{- end }}
      {{- if .Ipv6Prefix }}
      {{ .Ipv6Prefix }} 1;
      {{- end }}
      {{- end }}
  }

google:
  url: https://developers.google.com/static/search/apis/ipranges/googlebot.json
bing:
  url: https://www.bing.com/toolbox/bingbot.json
yandex:
  url: https://example.com/yandexbot-ranges.json
yahoo:
  url: https://example.com/yahoobot-ranges.json
```

Replace the URLs with the actual JSON URLs for bot IP ranges and set the correct file path for the Nginx configuration.

### Command-Line Usage

Use the `--config` flag to specify the path to the YAML configuration file.

```bash
./nginx-bot-ip-range-updater --config /path/to/your/config.yaml
```
