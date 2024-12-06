# homebox-export

A command-line tool to export attachments from a [Homebox](https://hay-kot.github.io/homebox/) instance.

## Features

- Export all attachments
- Organize downloads into folders by item name

## Output Structure

The exported files will be organized in the following structure:

```
export/
  Item Name_SHORTID/
    attachment1.jpg
    attachment2.pdf
    ...
```

## Installation

### From Source

```bash
git clone https://github.com/kusold/homebox-export
cd homebox-export
go build -o homebox-export ./cmd/homebox-export
```

## Usage

### Basic Usage

```bash
homebox-export export -server http://homebox.local -user admin -pass secret
```

### Using Environment Variables

Create a `.env` file:

```env
HOMEBOX_SERVER=http://homebox.local
HOMEBOX_USER=admin
HOMEBOX_PASS=secret
HOMEBOX_OUTPUT=./my-backup
HOMEBOX_PAGESIZE=50
```

Then run:

```bash
homebox-export export
```

### Command Line Options

```
Usage: homebox-export <command> [options]

Commands:
  export        Download all items and their attachments
  help          Show this help message
  version       Show version information

Export Options:
  -server       Homebox server URL
  -user         Username for authentication
  -pass         Password for authentication
  -output       Output directory (default: ./export)
  -pagesize     Number of items per page (default: 100)

Environment Variables:
  HOMEBOX_SERVER   Server URL
  HOMEBOX_USER     Username
  HOMEBOX_PASS     Password
  HOMEBOX_OUTPUT   Output directory
  HOMEBOX_PAGESIZE Number of items per page
```

## Development

This project uses [Task](https://taskfile.dev) for build automation. To get started:

1. Install Task
2. Create a `.env` file with your Homebox credentials
3. Run development commands:

```bash
# Build and run
task

# Just build
task build

# Just run
task run

# Clean build artifacts
task clean
```

## License

MIT

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -am 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
