# GoForge 🚀

**GoForge** is a lightweight and efficient build tool designed to simplify the development workflow for **Go (Golang)** projects. It automates repetitive tasks like project initialization, dependency management, optimization, and cross-compilation — so you can focus on writing code, not boilerplate.

---

## ✨ Features

- 🔧 **Project Initialization** – Scaffold Go projects with sensible defaults.
- 📦 **Automatic Module Management** – No need to manually run `go mod tidy` after every change.
- 🚀 **Optimized Builds** – Easy-to-apply build flags for performance.
- 🌍 **Cross Compilation** – Build for multiple platforms concurrently with simple configuration.
- ⚡ **One-liner Build & Run** – Quickly test your binaries with minimal effort.

---

## 🛠️ Configuration

GoForge uses `GoForge.yaml` to define how your project is built.

```yaml
app:
  package: GoForge
  version: 0.11.1
build:
  output: build/goforge.exe
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
1. **Default/Primary OS**: The first OS in the list is treated as the default. Its binary is output directly to the path specified in `build.output` (e.g. `build/goforge.exe`).
2. **Subsequent OS Targets**: The next operating systems are output to platform-specific subdirectories under the output directory (e.g., `build/mac/goforge`, `build/linux/goforge`).
3. **Target Mapping**: Specifying `mac` automatically compiles using Go's `darwin` target but places the binary under the `mac/` directory.

---

## 🚀 Usage

Use `goforge` in your terminal to manage Go project builds and automation:

```bash
goforge <command> [args]
```

| Command                          | Description                                                                |
| -------------------------------- | -------------------------------------------------------------------------- |
| `goforge new <pkg-name>`         | Create a new Go project in the current directory and initialize `go.mod`.  |
| `goforge build`                  | Build the Go project for all configured platforms.                        |
| `goforge run`                    | Run the primary compiled binary (defined as the first OS target).          |
| `goforge build run`              | Build and immediately run the primary binary.                              |
| `goforge clean`                  | Safely removes all build binaries and platform subdirectories.             |
| `goforge install`                | Install the binary to `$GOBIN`. *(Currently experimental)*                 |
| `goforge remove`                 | Remove the installed binary from `$GOBIN`.                                 |
