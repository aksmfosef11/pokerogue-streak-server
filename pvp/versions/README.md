# Version profiles

`PVP_VERSION_REGISTRY=pvp/versions.json` loads the supported version map. Without this setting, the existing `PVP_RULES_DIR` and `PVP_AUDIT_URL` configuration remains supported.

Each entry names its exact four-part game version, exported `rules.json` and `rules.manifest.json` directory, audit policy, and native worker URL. Relative directories resolve from the registry file. The startup loader checks the version, rule hash, schema, and required rule tables. An unregistered client version is excluded from PvP while ordinary uploads retain the existing soft qualification behavior.

The deployment registry currently enables only 1.12.2.0. Older rules and sources are retained in the repository for compatibility work, but their routes and workers are not enabled in the current deployment. Before retiring an engine on any existing installation, ensure its active runs, battles, and pending audit jobs have been handled.

When adding a future version, use its native client and corresponding rules and audit worker. A copied or relabeled rules export does not establish native compatibility. Preserve a previously active version for any runs or audits that still depend on it; test save migration and PvP compatibility before retiring it.

The original 1.12.1.3 account/team tables and seed domains retain their identities. Additional versions use isolated state, teams, qualification, and seed domains. Existing enhanced system data first imported into another version still uses operator baseline review; approval is scoped to the selected version. Use `pvp-admin review -version VERSION USERNAME` and the existing reasoned approval command against a matching system save.

All profiles share the server's audit concurrency and 64 in-memory pending/running queue slots. Pending audits are durable and recovered when workers become available. There is no per-account daily registration cap; historical attempts remain recorded. Operator review opens the registry without starting background audit recovery.
