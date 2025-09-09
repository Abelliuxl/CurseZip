# CurseZip: Your Ultimate CurseForge Plugin Archiver

![CurseZip Logo](assets/cursezip_logo_text.png)

CurseZip is a powerful and flexible command-line tool designed to streamline the packaging of your CurseForge plugin directories into clean, optimized archives. Say goodbye to unwanted `.git` folders, `.DS_Store` files, and other unnecessary clutter in your releases!

## ✨ Features

- **Multi-Directory Archiving**: Package one or more plugin directories into a single, consolidated archive.
- **Flexible Exclusion Rules**: Easily exclude files and directories using a `.cursezip.json` configuration file or directly via command-line arguments.
- **Common Exclusions Built-in**: Comes with sensible default exclusions for common development artifacts (e.g., `.git/`, `.DS_Store`, `*.log`).
- **Multiple Archive Formats**: Supports both `.zip` and `.tar.gz` output formats.
- **Automated Naming**: Archives are automatically named after the primary source directory for convenience.

## 🚀 Installation

To install CurseZip, ensure you have Go (1.16 or newer) installed on your system.

```bash
go install github.com/your-repo/cursezip@latest # Replace with actual repo path
```

Alternatively, you can clone the repository and build it manually:

```bash
git clone https://github.com/your-repo/cursezip.git # Replace with actual repo path
cd cursezip
go build -o cursezip .
```

Place the `cursezip` executable in your system's PATH for easy access.

## 💡 Usage

### Basic Archiving

Archive a single plugin directory:

```bash
cursezip path/to/your/plugin/MyAwesomePlugin
```
This will create `path/to/your/plugin/MyAwesomePlugin.zip` (default format).

### Specifying Archive Format

Use the `-f` or `--format` flag to choose between `zip` (default) or `tar.gz`:

```bash
cursezip -f tar.gz path/to/your/plugin/MyAwesomePlugin
```
This will create `path/to/your/plugin/MyAwesomePlugin.tar.gz`.

### Archiving Multiple Directories

Package several plugin directories into one archive. The output archive will be named after the *first* specified directory.

```bash
cursezip path/to/pluginA path/to/pluginB path/to/pluginC
```
This will create `path/to/pluginA.zip` containing `pluginA`, `pluginB`, and `pluginC` (each in their respective subdirectories within the archive).

### Adding Custom Exclusions

Exclude additional files or patterns using the `-e` or `--exclude` flag. This flag can be used multiple times.

```bash
cursezip -e "*.tmp" -e "**/test/" path/to/your/plugin/MyAwesomePlugin
```

### Using a Configuration File for Exclusions

For more complex or persistent exclusion rules, create a `cursezip.json` file in your project root or specify its path using `-c` or `--config`.

**Example `cursezip.json`:**

```json
{
  "exclude": [
    "node_modules/",
    "*.psd",
    "src/dev-assets/"
  ]
}
```

Then, run CurseZip with your config:

```bash
cursezip -c /path/to/my/custom/cursezip.json path/to/your/plugin/MyAwesomePlugin
```

If no config path is specified, CurseZip will look for `cursezip.json` in the current working directory. Custom exclusions from the command line are *merged* with those from the config file.

## ⚙️ Configuration

CurseZip uses a `cursezip.json` file to manage exclusion patterns. When loading, it first applies a set of sensible default exclusions (e.g., `.git/`, `.DS_Store`, `go.mod`, `main.go`, etc.) and then merges them with any patterns defined in your `cursezip.json` file.

**Default Exclusions (built-in):**
- `.git/`
- `.DS_Store`
- `*.log`
- `go.mod`
- `go.sum`
- `main.go`
- `packer/`
- `archiver/`
- `config/`
- `cursezip.example.json`
- `.*` (hidden files)

You can override or extend these by providing your own `cursezip.json`.

## 🤝 Contributing

Contributions are welcome! Please feel free to open issues or submit pull requests.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
