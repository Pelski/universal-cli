Jasne! Oto zaktualizowana wersja Twojego `README.md`:

---

# Universal CLI - HTTP client for creating tools

This is a command-line tool written in Go that allows you to interact with HTTP APIs using standard HTTP methods (`GET`, `POST`, `PUT`, `DELETE`). It supports configuration via a file, authentication methods, custom headers, and dynamic parameters.

## Features

- **HTTP Operations**: Supports `get`, `create`, `update`, `delete` and their aliases (`list`, `show`, `set`, `drop`).
- **Authentication**:
  - **Bearer Token**: Use a token file for Bearer authentication.
  - **Basic Auth**: Use username and password for Basic Authentication.
- **Configuration**: Load settings from a YAML configuration file.
- **Custom Headers**: Add additional headers from the configuration.
- **Dynamic Flags**: Pass parameters and data via command-line flags.
- **JSON Input Support**: Dynamically parse JSON input as parameters.
- **Auto Parsing**: Handles integer, boolean, float, and JSON types for flags.
- **Debug Mode**: Enable verbose output for debugging purposes.

## Installation

Ensure you have Go installed (version 1.16 or higher is recommended).

```bash
# Clone the repository
git clone [repository_url]
cd [repository_directory]

# Build the executable
go build -o ucli
# or
just build
```

## Usage

```bash
./ucli [--config CONFIG_FILE] [--debug] OPERATION RESOURCE [RESOURCE ...] [--FLAG VALUE ...]
```

### Operations

- `get`, `list`, `show`: Perform a **GET** request.
- `create`, `search`, `find`: Perform a **POST** request.
- `update`, `set`: Perform a **PUT** request.
- `delete`, `drop`: Perform a **DELETE** request.

### Options

- `--config`: Specify a custom configuration file. Defaults to `configuration.yaml` in the current directory.
- `--debug`: Enable debug mode for verbose output.

### Flags

- Pass additional parameters or data using flags prefixed with `--`.
- Flags can be in the form `--key value` or `--key=value`.
- Supports automatic parsing of **integers**, **booleans**, **floats**, and **JSON** data.

### Examples

#### GET Request

```bash
./ucli get users --page 2 --limit 10
```

Performs a **GET** request to `/users` with query parameters `page=2` and `limit=10`.

#### POST Request

```bash
./ucli create users --name "John Doe" --email "john@example.com"
```

Performs a **POST** request to `/users` with JSON body `{"name": "John Doe", "email": "john@example.com"}`.

#### PUT Request

```bash
./ucli update users 123 --email "newemail@example.com"
```

Performs a **PUT** request to `/users/123` with JSON body `{"email": "newemail@example.com"}`.

#### DELETE Request

```bash
./ucli delete users 123
```

Performs a **DELETE** request to `/users/123`.

## Configuration

Create a `configuration.yaml` file in the current directory or specify a custom file with the `--config` option.

### Configuration Options

- `url`: Base URL of the API (e.g., `https://api.example.com`).
- `token`: Path to a file containing the Bearer token.
- `username`: Username for Basic Authentication.
- `password`: Password for Basic Authentication.
- `headers`: A map of additional headers to include in requests.

### Example Configuration

```yaml
url: "https://api.example.com"
# Token
token: "/path/to/token_file"
# Basic Auth
username: "user"
password: "password"
# Custom Headers
headers:
  X-Custom-Header: "CustomValue"
```

## Authentication

### Bearer Token

- Set the `token` field in the configuration.
- The tool reads the token from the specified file and sets the `Authorization` header.

### Basic Authentication

- Set `username` and `password` in the configuration.
- The tool uses these credentials for Basic Authentication.

## Custom Headers

- Specify additional headers under the `headers` section in the configuration file.
- These headers are included in every request.

## Debug Mode

- Enable by adding the `--debug` flag.
- Outputs detailed information about requests, responses, and query parameters.

```bash
./ucli --debug get users
```

## Dynamic Flags

- Parameters passed via flags are automatically parsed as:
  - **integers** (e.g., `--limit=10`),
  - **booleans** (e.g., `--debug=true`),
  - **floats** (e.g., `--price=99.99`),
  - **JSON** objects or arrays (e.g., `--filters='{"name": "John"}'`).

## Dependencies

- **Viper**: For configuration management.

Install dependencies using:

```bash
go get github.com/spf13/viper
```

## Building from Source

```bash
go build -o ucli
# or
just build
```

## License

MIT License

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## Contact

For questions or suggestions, please open an issue.
