# db_connect

A unified command-line interface for managing and connecting to multiple database types with secure, reusable connection configurations.

## Description

`db_connect` is a Bash-based database connection manager that simplifies the process of connecting to various databases by storing connection credentials in a centralized configuration file. Instead of remembering complex connection strings and credentials for each database, you can store them once and connect with a simple command.

The tool provides an interactive configuration experience with fuzzy search for database selection and supports multiple database types including ClickHouse, MySQL, PostgreSQL, and MongoDB.

## Features

- **Multi-database support**: ClickHouse, MySQL, PostgreSQL, and MongoDB
- **Secure credential storage**: Passwords stored in YAML configuration file at `~/.config/db_connect.yaml`
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
db_connect init
```

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

Configurations are stored in `~/.config/db_connect.yaml` in the following format:

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

- Configuration file contains passwords in plain text. Ensure proper file permissions:
  ```bash
  chmod 600 ~/.config/db_connect.yaml
  ```
- Consider using SSH tunnels for remote database connections
- Some database clients support password managers or credential helpers
- Do not commit the configuration file to version control

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

**Solution**: Install the missing dependency following the provided link. Both `yq` and `fzy` are required for the tool to function.

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
