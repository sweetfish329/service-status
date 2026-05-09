# Minimal Status Page

A lightning-fast, minimal server status page designed for Cloudflare Workers.
Built with:
- **Go (WebAssembly)** using `syumai/workers` for high-performance concurrent fetching.
- **Alpine.js** for lightweight frontend reactivity.
- **Materialize CSS** for modern, clean UI.

## Configuration

Servers are configured via the `SERVERS_JSON` environment variable in `wrangler.toml` or via Cloudflare Dashboard.

Format:
```json
[
  {"name": "Website", "type": "http", "url": "https://example.com"},
  {"name": "Palworld", "type": "palworld", "url": "http://ip:8212/v1/api/info", "auth": "admin:password"}
]
```

## Local Development

```bash
# Build the project
make build

# Run locally using wrangler
npx wrangler dev
```

## Deployment

```bash
npx wrangler deploy
```
