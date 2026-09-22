# Notification Channels

`internal/notifications` (`Dispatcher`, `dispatcher.go`/`channels.go`/`retry.go`) fans an alert,
finding, or digest out to every enabled channel in `notification_config` (`GET`/`PUT
/api/notifications/config`, `NotificationChannel` in `docs/openapi.yaml`). Every channel below
shares the same delivery-tracking and bounded-retry machinery
(`notification_deliveries`, `internal/notifications/retry.go`): a failed send is retried on the
schedule in `retrySchedule` (1m, 5m, 15m, 30m, 1h) up to `maxDeliveryAttempts` (5) before it stays
`failed` for good.

## Webhook signing (HMAC-SHA256)

The generic `webhook` channel — and `discord`/`slack`, which use the same HTTP transport — can be
configured with a shared secret. When set, every delivery carries:

- `X-WiseLabz-Timestamp` — Unix seconds at send time.
- `X-WiseLabz-Signature` — `sha256=<hex>`, where the digest is
  `HMAC-SHA256(secret, timestamp + "." + raw_request_body)`.

Verify a delivery by recomputing the same HMAC over the exact bytes received (before any
re-serialization) and comparing constant-time. There's no replay window enforced server-side;
receivers that care should reject stale timestamps themselves.

The secret is submitted once as `config.secret` (write-only) and stored encrypted at rest
(`WISELABZ_ENCRYPTION_KEY`); reads return `config.secretSet: true` instead of the value. Omitting
`secret` on an update keeps the previously stored one; submitting an empty string clears it. An
unset secret sends the request unsigned — the API accepts a signed webhook only if the receiver
requires one.

## Channel reference

Every `config` field below lives under `NotificationChannel.config` for that channel's entry in
`notification_config.channels`.

| Channel    | `type`     | Config fields                                              | Notes |
|------------|------------|--------------------------------------------------------------|-------|
| Generic webhook | `webhook`  | `url`, `secret` (optional)                             | JSON body `{"title","message"}`, HMAC-signed if `secret` is set. |
| Discord    | `discord`  | `url`, `secret` (optional)                                 | Discord incoming-webhook URL. Body is an embed: `{"embeds":[{"title","description"}]}`. |
| Slack      | `slack`    | `url`, `secret` (optional)                                 | Slack incoming-webhook URL. Body: `{"text":"*title*\nmessage"}`. |
| ntfy       | `ntfy`     | `url` (server, default `https://ntfy.sh`), `topic`, `priority` (optional), `tags` (optional) | Publishes to `<url>/<topic>` per [ntfy's HTTP API](https://docs.ntfy.sh/publish/); title/priority/tags ride as `X-Title`/`X-Priority`/`X-Tags` headers, body is the plain-text message. |
| Telegram   | `telegram` | `chatId`, `secret` (bot token)                              | Posts to the Bot API's `sendMessage` (`https://api.telegram.org/bot<token>/sendMessage`) with Markdown formatting. The bot token rides in the same write-only `secret` field as webhook signing secrets — it's just as sensitive, so it reuses the same encrypted-at-rest storage. |
| SMTP       | `smtp`     | `host`, `port` (default 587), `username` (optional), `secret` (password, optional), `from`, `to` (comma-separated) | Sends a plain-text email with opportunistic STARTTLS. Dials through the same SSRF-guarded dialer as the HTTP channels, since `host` is admin-configured. |
| In-app     | `in_app`   | — | Always delivered first; not configurable, never fails. |

All HTTP-based channels (webhook, Discord, Slack, ntfy, Telegram) share one guarded transport
(`internal/notifications/webhook.go`): a 10s timeout, a dialer that refuses loopback/link-local/
unspecified/multicast destinations, and no redirect following. Response bodies are read up to 64
KiB and discarded — never stored or surfaced.

## Adding a channel type

A channel type is a `channelSender` function (`internal/notifications/channels.go`) registered in
the `channelSenders` map, plus (for a type intended to be user-configurable) an entry in
`NotificationChannelType` in `docs/openapi.yaml`. The dispatcher, delivery recording, and retry
loop dispatch by channel type through that map and need no changes.
