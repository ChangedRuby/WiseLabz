# WiseLabz

**Automated documentation for your homelab.**

WiseLabz is a self-hosted, open-source service that automatically generates and maintains
documentation for your home lab. Register your services, Proxmox, Docker/Portainer,
pfSense/OPNsense, and more. Connect their APIs, and WiseLabz fetches live data to produce
structured, readable documentation from configurable templates. A diff engine detects
changes since the last sync and either triggers an AI-assisted section update or surfaces
an alert so nothing drifts undocumented. A web dashboard gives you a live view of your
services alongside the generated docs.

<p>
  <a href="https://github.com/WiseLabz/WiseLabz/actions/workflows/ci.yml"><img src="https://github.com/WiseLabz/WiseLabz/actions/workflows/ci.yml/badge.svg" alt="CI status"></a>
  <a href="https://github.com/WiseLabz/WiseLabz/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-Apache%202.0-blue.svg" alt="License"></a>
  <a href="https://github.com/WiseLabz/WiseLabz/releases/latest"><img src="https://img.shields.io/github/v/release/WiseLabz/WiseLabz?label=latest" alt="Latest release"></a>
  <a href="https://github.com/WiseLabz/WiseLabz/pkgs/container/wiselabz"><img src="https://img.shields.io/badge/container-ghcr.io-blue" alt="Docker image"></a>
</p>

---

## Features

- **Multi-service connectors** — Built-in support for Proxmox, Docker/Portainer, and
  pfSense/OPNsense. Pull live configuration and state from each service via their native
  APIs.
- **Template-driven documentation** — Define what your docs look like with configurable
  templates. Control which sections appear, how data is formatted, and what detail level
  you get.
- **Change-aware diff engine** — WiseLabz compares each sync against the previous snapshot.
  When something changes, it either updates the relevant doc section automatically or
  flags it for review.
- **AI-assisted updates** — When a service change is detected, WiseLabz can draft the
  updated documentation section for you, keeping manual effort near zero.
- **Web dashboard** — Browse your generated documentation and monitor live service data
  side by side in one place.
- **Community connectors** — The connector interface is open. Anyone can write and
  contribute a connector for a service they rely on.
- **Single binary, single compose file** — Deploy with `docker compose up`. SQLite for
  small labs, PostgreSQL when you need it.

## Quick start

```bash
# Clone the repo
git clone https://github.com/WiseLabz/WiseLabz.git
cd WiseLabz

# Copy and edit the configuration
cp config.example.yaml config.yaml
# Edit config.yaml — add at least one service connection

# Start everything
docker compose up -d
```

Open `http://localhost:8080` and register your first service. A `config.example.yaml`
file is included in the repository root with commented examples for every supported
service.

## Configuration

WiseLabz reads settings from a `config.yaml` file at the repository root. Every key in
the YAML file can be overridden with an environment variable prefixed with `WISELABZ_`.
Environment variables always take precedence over the file, so you can keep secrets out of
your config file and tune settings per deployment without editing YAML.

```yaml
# config.yaml
server:
  port: 8080

database:
  driver: sqlite                # or "postgres"
  dsn: ./data/wiselabz.db

services:
  - name: home-proxmox
    type: proxmox
    url: https://192.168.1.10:8006
    # token_id and token_secret are read from env:
    #   WISELABZ_SERVICES_0_TOKEN_ID
    #   WISELABZ_SERVICES_0_TOKEN_SECRET
```

For a full list of configuration keys, see [deploy/config.example.yaml](deploy/config.example.yaml).
Run `server config schema` to print a JSON Schema generated from the server's `Config`
struct (for example, `server config schema > config.schema.json`). Each scalar key has
an `x-env` annotation naming its `WISELABZ_*` override. OIDC providers (`auth.oidc`)
are file-only collections, marked `x-config-file-only`; their item fields have no env
annotation. The schema describes key names and types, includes secret field names,
and never loads or prints configured values. It does not replace semantic validation:
use `server config validate` to check deployment settings and `server config print
--redacted` to inspect them with secrets masked.

## Supported services

| Service            | Status               |
|--------------------|----------------------|
| Proxmox VE         | Built-in             |
| Docker             | Built-in             |
| pfSense / OPNsense | Built-in             |
| Netbird            | Built-in             |
| Pi-hole (v5/v6)    | Built-in             |
| Cloudflare         | Built-in             |
| Traefik            | Built-in             |
| Portainer          | Built-in             |
| AdGuard Home       | Built-in             |
| UniFi              | Built-in             |
| Everything else    | Community connectors |

New service connectors are community-driven. If the service you run isn't here yet,
[request it][connector-request] or build it yourself, the connector interface is
straightforward and well-documented.

## Contributing

Contributions are welcome. Whether it's a bug report, a new connector, or a docs fix,
start with [`CONTRIBUTING.md`](CONTRIBUTING.md) for the full process, local setup
instructions, and branch conventions.

## License

WiseLabz is licensed under the Apache License 2.0. See [`LICENSE`](LICENSE) for the full
text.

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md).
Participation means you agree to its terms. Report concerns to the project maintainers
via GitHub's private reporting mechanisms.

[connector-request]: https://github.com/WiseLabz/WiseLabz/issues/new?template=connector_request.yml
