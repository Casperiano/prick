![Prick](images/logo.png)

Prick is a CLI/TUI for patching or adding firewall rules to Azure resources, including:

- Storage Accounts
- Key Vaults
- Synapse Workspaces
- SQL Servers

## Features

Prick automatically uses your current IP address.
Simplifies the process compared to using the az CLI, especially when managing multiple firewall rules.

## Example Use Cases

- Temporarily granting external access without modifying permanent configurations.
- Seamlessly patching firewall rules across an organization in CI/CD pipelines.

## Installation

TODO

## Usage

To open the terminal user interface (TUI), run:
```bash
prick list
```

To see all other prick commands, run:

```bash
prick --help
```

## Planned Features

- Config file for default firewall rules
- Allow whitelisting of other IP addresses / CIDR ranges