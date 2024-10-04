![Prick](images/logo.png)

Prick/Patch firewall rules in Azure resources:
- Storage Accounts
- Keyvaults
- Synapse Workspaces
- SQL Servers

By default use your current IP address, or you can provide a custom CIDR range.

Provide improved usability over `az` cli when adding mayny filewall rules.
example use cases:
- Applying local terraform without VPN
- Temporarily giving access to external people without permanent config.
- Allow easy patching of all firewall rules across an organisation in CI/CD.

Provide overview of all firewall rules.

Open UI:
```
prick list 
prick list --resource-type StorageAccount
prick list --resource-group rg1
prick list --resource-group rg1 --resource r1
prick list --resource-group rg1 --resource-type StorageAccount
```

```
prick poke --resource-type StorageAccount
prick poke --resource-group rg1
prick poke --resource-group rg1 --resource r1
prick poke --resource-group rg1 --resource-type StorageAccount
```

```
prick patch --resource-type StorageAccount
prick patch --resource-group rg1
prick patch --resource-group rg1 --resource r1
prick patch --resource-group rg1 --resource-type StorageAccount
```

## Config
coming soon
