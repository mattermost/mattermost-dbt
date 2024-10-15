# Mattermost Database Tooling

This repository contains tooling for reviewing Mattermost databases with complex configurations.

### Building

Before building the `mmdbt` CLI, ensure that you have Go installed and are using a version that is at least as new as the one [referenced here](https://github.com/mattermost/mattermost-dbt/blob/main/go.mod#L3).

To build and install `mmdbt` locally:

```bash
make install
```

To build `mmdbt` for Linux/AMD64:

```bash
make build
```

The `mmdbt` CLI binary will be placed in the root of the `mattermost-dbt` directory. You can export env var `ARCH` to `arm64` or `arm` to build for arm architectures.
