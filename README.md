# Girocco 🚀

**Girocco** is a lightweight and efficient build tool designed to simplify the development workflow for **Go (Golang)** projects. It automates repetitive tasks like project initialization, dependency management, optimization, and cross-compilation — so you can focus on writing code, not boilerplate.

---

> [!NOTE]
> **Project Rename**
>
> This project was previously known as **GoForge** and has been renamed to **Girocco** to establish a more unique identity and avoid naming conflicts as the project evolves.
>
> The rename only affects the project name, repository, module path, and CLI executable. The functionality, commands, and overall usage remain the same.
>
> If you were using **GoForge**, simply replace references to:
>
> * `GoForge` → `Girocco`
> * `github.com/piyushdugawa/GoForge` → `github.com/piyushdugawa/girocco`
> * `goforge` → `girocco`
>
> Thank you to everyone who has supported and contributed to the project. Your existing knowledge and workflows should transfer seamlessly to Girocco.


---

## ✨ Features

- 🔧 **Project Initialization** – Scaffold Go projects with sensible defaults.
- 📦 **Automatic Module Management** – No need to manually run `go mod tidy` after every change.
- 🚀 **Optimized Builds** – Easy-to-apply build flags for performance.
- 🌍 **Cross Compilation** – Build for multiple platforms concurrently with simple configuration.
- ⚡ **One-liner Build & Run** – Quickly test your binaries with minimal effort.

---

## 🛠️ Configuration

Girocco uses `Girocco.yml` to define how your project is built.

```yaml
app:
  package: github.com/piyushdugawa/girocco
  version: 0.14.0
build:
  output: build/girocco.exe
  optimisation: true

  env:
    GOOS: [windows, mac, linux]  # Can be a YAML list or comma-separated string: "windows,linux,mac"
    GOARCH: amd64

  flags:
    - -ldflags
    - "-s -w"
```

### 🌍 Multi-Platform Compilation Behavior
When `GOOS` contains multiple operating systems:
1. **Default/Primary OS**: The first OS in the list is treated as the default. Its binary is output directly to the path specified in `build.output` (e.g. `build/girocco.exe`).
2. **Subsequent OS Targets**: The next operating systems are output to platform-specific subdirectories under the output directory (e.g., `build/mac/girocco`, `build/linux/girocco`).
3. **Target Mapping**: Specifying `mac` automatically compiles using Go's `darwin` target but places the binary under the `mac/` directory.

---

## 🚀 Usage

Use `girocco` in your terminal to manage Go project builds and automation:

```bash
girocco <command> [args]
```

| Command                          | Description                                                                |
| -------------------------------- | -------------------------------------------------------------------------- |
| `girocco new <pkg-name>`         | Create a new Go project in the current directory and initialize `go.mod`.  |
| `girocco build`                  | Build the Go project for all configured platforms.                        |
| `girocco run`                    | Run the primary compiled binary (defined as the first OS target).          |
| `girocco build run`              | Build and immediately run the primary binary.                              |
| `girocco clean`                  | Safely removes all build binaries and platform subdirectories.             |
| `girocco install`                | Install the binary to `$GOBIN`. *(Currently experimental)*                 |
| `girocco remove`                 | Remove the installed binary from `$GOBIN`.                                 |
