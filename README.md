# anything-to-ntfy

anything-to-ntfy aims to provide a variety of *bridge* endpoints to convert from
common webhook formats to a notification/alert to a ntfy.sh compatible instance
(note that ntfy.sh can be self hosted).

## Build status

![OCI Imagebuilding](https://github.com/0robustus1/anything-to-ntfy/actions/workflows/image.yml/badge.svg)

## CLI parameters

```
Usage: anything-to-ntfy --ntfy-token=STRING

Flags:
  -h, --help                            Show context-sensitive help.
      --ntfy-token=STRING               Token to use to communicate with ntfy instance ($NTFY_TOKEN)
      --ntfy-default-instance=STRING    Which ntfy instance to use by default ($NTFY_DEFAULT_INSTANCE)
      --ntfy-default-topic=STRING       Which ntfy topic to use by default ($NTFY_DEFAULT_TOPIC)
      --listen-host=STRING              Which host to listen on, should be an address. Defaults to empty string which is
                                        equivalent to 0.0.0.0 ($LISTEN_HOST)
      --listen-port=5000                Which port to listen on ($LISTEN_PORT).
```

## Endpoints

### General parameters

Generally an anything-to-ntfy instance should be run with `NTFY_TOKEN`,
`NTFY_DEFAULT_INSTANCE` and `NTFY_DEFAULT_TOPIC` environment variables set.

All three parameters can also be set as part of the HTTP requests via the following parameters:

* `?ntfyToken=`, `?ntfyInstance=`, and `?ntfyTopic` query parameters as part of the URL
* `X-NTFY-Token`, `X-NTYF-INSTANCE`, and `X-NTFY-Topic` HTTP request headers

### Supported endpoints

#### Grafana Alerting Webhooks

* `/grafana/webhook/:topic` (where `:topic` is an optional parameter with the name of an ntfy.sh topic)

#### Slack Incoming Webhook

* `/slack/incoming_webhook/:topic` (where `:topic` is an optional parameter with the name of an ntfy.sh topic)
