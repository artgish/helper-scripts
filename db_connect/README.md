# db_connect

> Note. AI generated readme.

A unified command-line interface for managing and connecting to multiple database types with reusable connection configurations.

## Description

`db_connect` is a Bash-based database connection manager that simplifies the process of connecting to various databases by storing connection credentials in a centralized configuration file. Instead of remembering complex connection strings and credentials for each database, you can store them once and connect with a simple command.

The tool provides an interactive configuration experience with fuzzy search for database selection and supports multiple database types including ClickHouse, MySQL, PostgreSQL, and MongoDB.

## Features

- **Multi-database support**: ClickHouse, MySQL, PostgreSQL, and MongoDB
- **Encrypted credential storage**: Configuration at `~/.config/db_connect.yaml` is encrypted with [sops](https://github.com/getsops/sops) using your GPG key; decryption is transparent and handled by `gpg-agent`
- **Interactive configuration**: Fuzzy search interface for selecting database types
- **Bash completion**: Tab completion for commands and stored configurations
- **Enhanced CLI clients**: Uses modern, feature-rich database clients (mycli, pgcli, mongosh, clickhouse-client)
- **SRV record support**: MongoDB SRV connection string support
- **Default values**: Smart defaults for ports and hosts
- **Simple management**: Easy-to-use commands for adding, deleting, and viewing configurations

## Prerequisites

### Required Dependencies

- **Bash** 4.0 or higher
- **[yq](https://github.com/mikefarah/yq)** - YAML processor for configuration management
- **[fzy](https://github.com/jhawthorn/fzy)** - Fuzzy text selector for interactive prompts
- **[sops](https://github.com/getsops/sops)** - Encrypts configuration values at rest
- **[gpg](https://www.gnupg.org/)** - Provides the key that sops uses; a running `gpg-agent` caches the passphrase so you aren't prompted on every invocation

### GPG Key

You need an existing GPG keypair. If you don't have one:

```bash
gpg --full-generate-key
# then grab the fingerprint
gpg --list-secret-keys --keyid-format=long
```

Export the fingerprint (no spaces) before the first `init`, or you'll be prompted for it:

```bash
export SOPS_PGP_FP=ABCDEF0123456789ABCDEF0123456789ABCDEF01
```

### Database-Specific Clients

Install the client(s) for the database type(s) you plan to use:

| Database | Client | Installation Guide |
|----------|--------|-------------------|
| MySQL | [mycli](https://www.mycli.net) | https://www.mycli.net/install |
| PostgreSQL | [pgcli](https://www.pgcli.com) | https://www.pgcli.com/install |
| ClickHouse | clickhouse-client | https://clickhouse.com/docs/install |
| MongoDB | [mongosh](https://www.mongodb.com/docs/mongodb-shell/) | https://www.mongodb.com/docs/mongodb-shell/install/ |

## Installation

### 1. Download the script

```bash
# Clone the repository
git clone https://github.com/gishyanart/helper-scripts.git
cd helper-scripts/db_connect

# Or download directly
curl -o db_connect https://raw.githubusercontent.com/gishyanart/helper-scripts/main/db_connect/db_connect
chmod +x db_connect
```

### 2. Make it accessible globally

```bash
# Option 1: Copy to a directory in your PATH
sudo cp db_connect /usr/local/bin/

# Option 2: Create a symlink
sudo ln -s $(pwd)/db_connect /usr/local/bin/db_connect

# Option 3: Add to your PATH in ~/.bashrc or ~/.zshrc
echo 'export PATH="$PATH:/path/to/helper-scripts/db_connect"' >> ~/.bashrc
source ~/.bashrc
```

### 3. Initialize the configuration

```bash
# Optional: provide the fingerprint up front, otherwise init will prompt
export SOPS_PGP_FP=ABCDEF0123456789ABCDEF0123456789ABCDEF01

db_connect init
```

On first run this creates `~/.config/.sops.yaml` (with a `path_regex` scoped to `db_connect.*\.yaml$` so it won't affect other sops users on your machine) and an empty sops-encrypted `~/.config/db_connect.yaml`.

### 4. Enable bash completion (optional but recommended)

```bash
# Add to your ~/.bashrc
eval "$(db_connect completion)"

# Then reload your shell configuration
source ~/.bashrc
```

## Usage

### Adding a Database Configuration

```bash
db_connect add
```

This will interactively prompt you for:

1. **Config name**: A friendly name to identify this connection (e.g., "prod-mysql", "dev-postgres")
2. **Database type**: Select from clickhouse, mysql, postgres, or mongodb (fuzzy search)
3. **Host**: Database server address (default: 127.0.0.1)
4. **SRV record**: Whether the host is an SRV record (MongoDB only, default: no)
5. **Port**: Database port (defaults to standard port for the selected type)
6. **User**: Database username
7. **Database name**: Target database (optional, press Enter to use default)
8. **Password**: Database password (hidden input)

#### Example Session

```bash
$ db_connect add
Config name: prod-mysql
Select database type: mysql
Database host [default: 127.0.0.1]: db.example.com
Is the host record SRV (default: no) [yes/no]: no
Database port [default: 3306]: 3306
Database user: admin
Connect database name: myapp_production
Database password: ********
```

### Connecting to a Database

```bash
db_connect connect <config-name> [additional-args]
```

**Examples:**

```bash
# Connect to a saved configuration
db_connect connect prod-mysql

# Connect and pass additional arguments to the underlying client
db_connect connect dev-postgres --no-password-prompt

# Connect to MongoDB with specific options
db_connect connect prod-mongo --quiet
```

### Viewing Configurations

```bash
# Show all configurations
db_connect show

# Show a specific configuration
db_connect show prod-mysql
```

**Example output:**

```yaml
prod-mysql:
  type: mysql
  host: db.example.com
  port: 3306
  user: admin
  pass: mypassword
  is_srv: no
  db: myapp_production
```

### Deleting a Configuration

```bash
db_connect delete <config-name>
```

**Example:**

```bash
db_connect delete old-dev-db
```

### Available Commands

| Command | Description |
|---------|-------------|
| `init` | Initialize the configuration file and check dependencies |
| `add` | Add a new database configuration |
| `show [name]` | Display all configurations or a specific one |
| `delete <name>` | Delete a configuration |
| `connect <name> [args]` | Connect to a database using a saved configuration |
| `completion` | Output bash completion script |

## Configuration File

Configurations are stored sops-encrypted in `~/.config/db_connect.yaml`. Decrypted (what you edit logically) it looks like:

```yaml
config-name:
  type: mysql
  host: localhost
  port: 3306
  user: username
  pass: password
  is_srv: no
  db: database_name
```

On disk, values are replaced with AES-256-GCM ciphertext and a `sops:` metadata block is appended. The top-level **keys remain cleartext** — that is how bash completion still lists config names without decrypting. Writes go through a `decrypt → yq → re-encrypt` tmp-file cycle; the tmp file is created alongside the config (not in `/tmp`) and is matched by the same `path_regex` in `.sops.yaml`.

### How encryption works

- sops encrypts each value with a random per-file AES-256-GCM data key.
- That data key is wrapped with your GPG public key (via `gpg --encrypt`) and stored inline.
- Decryption calls `gpg --decrypt` through `gpg-agent`; if your agent has cached the passphrase, no prompt appears.

### Inspecting / editing by hand

```bash
sops decrypt ~/.config/db_connect.yaml          # print plaintext
sops ~/.config/db_connect.yaml                  # open $EDITOR on decrypted view; re-encrypts on save
sops updatekeys ~/.config/db_connect.yaml       # re-wrap data key for new recipients after editing .sops.yaml
```

### Default Ports

| Database | Default Port |
|----------|--------------|
| ClickHouse | 9000 |
| MySQL | 3306 |
| PostgreSQL | 5432 |
| MongoDB | 27017 |

## Database-Specific Connection Details

### MySQL (mycli)

```bash
# Connects using:
mycli --ssl -h <host> -u <user> -P <port> -p <password> [-D <database>]
```

SSL is enabled by default for secure connections.

### PostgreSQL (pgcli)

```bash
# Connects using:
PGPASSWORD=<password> pgcli -h <host> -u <user> -p <port> [-d <database>]
```

Password is passed via environment variable to avoid command-line exposure.

### ClickHouse (clickhouse-client)

```bash
# Connects using:
clickhouse-client --host <host> --user <user> --port <port> --password <password> [--database <database>]
```

### MongoDB (mongosh)

```bash
# Standard connection:
mongosh mongodb://<user>:<password>@<host>:<port>[/<database>]

# SRV connection:
mongosh mongodb+srv://<user>:<password>@<host>:<port>[/<database>]
```

## Security Considerations

- Configuration values (including passwords) are encrypted at rest with AES-256-GCM; the data key is wrapped with your GPG public key. Loss of your GPG private key means loss of the config — back up your key.
- `gpg-agent` controls how long your passphrase stays cached. Tune `default-cache-ttl` / `max-cache-ttl` in `~/.gnupg/gpg-agent.conf` to match your threat model.
- During `connect`, decrypted credentials live in shell variables and the arg list of the database client (visible in `ps` output on a multi-user box). For `postgres` we already pass the password via `PGPASSWORD`; for the others, prefer passwordless clients (SSH tunnels, IAM auth, client cert auth) on shared hosts.
- `add` and `delete` write a sibling tmp file (e.g. `~/.config/db_connect.XXXXXX.yaml`) during re-encryption. It is created with `chmod 600` and removed via `trap` on exit. If the process is killed mid-write you may find a stray file — remove it manually.
- Do not commit `~/.config/db_connect.yaml` to a **public** repo; the sops format is safe to store in a **private** repo since values are encrypted.

### Migrating from a pre-encryption config

If you already have a plaintext `~/.config/db_connect.yaml` from an older version:

```bash
# ensure sops + gpg are installed and SOPS_PGP_FP is exported
db_connect init                              # creates ~/.config/.sops.yaml only (existing config untouched)
sops encrypt -i ~/.config/db_connect.yaml    # encrypts in place
chmod 600 ~/.config/db_connect.yaml
```

## Troubleshooting

### Command not found errors

**Problem**: `db_connect: command not found`

**Solution**: Ensure the script is in your PATH and executable:
```bash
chmod +x /path/to/db_connect
echo 'export PATH="$PATH:/path/to/db_connect"' >> ~/.bashrc
source ~/.bashrc
```

### Missing dependency errors

**Problem**: `` `yq` is not installed, install from https://...``

**Solution**: Install the missing dependency following the provided link. `yq`, `fzy`, `sops`, and `gpg` are all required for the tool to function.

### GPG prompts appear on every command

**Problem**: You're asked for the GPG passphrase each time you run `db_connect`.

**Solution**: Ensure `gpg-agent` is running and raise its cache TTL in `~/.gnupg/gpg-agent.conf`:

```
default-cache-ttl 28800
max-cache-ttl 86400
```

Then `gpg-connect-agent reloadagent /bye`.

### "no matching creation rule" during encrypt

**Problem**: `sops` refuses to encrypt saying no creation rule matches.

**Solution**: `~/.config/.sops.yaml` is missing or its `path_regex` doesn't match. Re-run `db_connect init` or check that the file contains `path_regex: db_connect.*\.yaml$`.

### Database client errors

**Problem**: Connection fails or client-specific errors

**Solution**:
- Verify the database client is installed: `which mycli` / `which pgcli` / etc.
- Test the client independently with direct connection parameters
- Check that credentials and connection details are correct using `db_connect show <name>`

### Bash completion not working

**Problem**: Tab completion doesn't work

**Solution**:
```bash
# Ensure completion is loaded in your shell
eval "$(db_connect completion)"

# Or add to ~/.bashrc permanently
echo 'eval "$(db_connect completion)"' >> ~/.bashrc
source ~/.bashrc
```

## Examples

### Quick Start Workflow

```bash
# 1. Initialize
db_connect init

# 2. Add your first database
db_connect add
# Enter: name=dev-db, type=postgres, host=localhost, user=devuser, etc.

# 3. Connect
db_connect connect dev-db

# 4. Later, view your configs
db_connect show

# 5. Remove old configs
db_connect delete old-db
```

### Managing Multiple Environments

```bash
# Add production database
db_connect add  # name: prod-mysql

# Add staging database
db_connect add  # name: staging-mysql

# Add development database
db_connect add  # name: dev-mysql

# Quick switch between environments
db_connect connect prod-mysql
db_connect connect staging-mysql
db_connect connect dev-mysql
```

## Contributing

Contributions are welcome! This is part of the [helper-scripts](https://github.com/gishyanart/helper-scripts) repository.

## License

MIT License - Copyright (c) 2025 gishyanart

See the [LICENSE](../LICENSE) file for full details.

## Acknowledgments

- [yq](https://github.com/mikefarah/yq) - YAML processor by Mike Farah
- [fzy](https://github.com/jhawthorn/fzy) - Fuzzy text selector by John Hawthorn
- [mycli](https://www.mycli.net) - MySQL CLI with auto-completion
- [pgcli](https://www.pgcli.com) - PostgreSQL CLI with auto-completion
- [mongosh](https://www.mongodb.com/docs/mongodb-shell/) - MongoDB Shell
- [ClickHouse](https://clickhouse.com) - ClickHouse database
