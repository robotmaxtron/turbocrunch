# TurboCrunch

TurboCrunch is a high-performance terminal-based calculator (TUI) that leverages the powerful SpeedCrunch evaluation engine. It features a dual-backend system, allowing users to switch between the robust SpeedCrunch C++ engine and a native Go math backend.

![TurboCrunch Demo](demo.gif)

**Note**: This project is a fork of the [SpeedCrunch](https://github.com/speedcrunch/SpeedCrunch) project, adapted for TUI use.

## Features

- **TUI Interface**: A clean, efficient terminal user interface for quick calculations.
- **SpeedCrunch Backend**: Uses the battle-tested SpeedCrunch core for high-precision and complex math operations.
- **`ans` Variable**: Automatically stores and retrieves the last calculation result for use in subsequent expressions.
- **Angle Modes**: Toggle between Radians, Degrees, and Gradians via `Ctrl+A`.
- **User-Defined Functions**: Define and use custom functions like `f(x) = x^2` that persist across sessions.
- **Formula Book, Constants & Units**: Searchable, interactive lists for common formulas (`Ctrl+F`), physical constants (`Ctrl+C`), and units (`Ctrl+U`).
- **Interactive History**: Navigate, select, and re-evaluate previous expressions using a structured table view.
- **Persistent Sessions**: Automatically saves and loads variables and custom functions to `session.json`.
- **Go Math Backend**: A native Go implementation using `cockroachdb/apd` for arbitrary-precision decimal arithmetic (50-digit precision) and `math/cmplx` for complex numbers.
- **Complex Number Support**: Native support for complex arithmetic in both backends.
- **High Accuracy**: Maintains high precision standards with 50-digit decimal accuracy in the Go backend (including transcendental functions) and arbitrary precision in SpeedCrunch.

## TUI Components & Keyboard Shortcuts

TurboCrunch uses several [Bubble Tea](https://github.com/charmbracelet/bubbletea) components to provide a rich interactive experience:

- **`bubbles/list`**: Used for the searchable Formula Book, Constants, and Units lists. Supports fuzzy filtering.
- **`bubbles/table`**: Provides an interactive history view where you can navigate and select previous results.
- **`bubbles/help`**: Dynamic, context-aware help menu (press `?` to toggle).
- **`bubbles/spinner`**: Visual feedback during complex backend evaluations.

### Key Bindings

| Key | Action |
|-----|--------|
| `Enter` | Evaluate expression |
| `Ctrl+T` | Toggle Math Backend (SpeedCrunch / Go) |
| `Ctrl+A` | Toggle Angle Mode (Rad / Deg / Grad) |
| `Ctrl+F` | Open Formula Book |
| `Ctrl+C` | Browse Physical Constants |
| `Ctrl+U` | Browse Units |
| `Ctrl+H` | Toggle History View / Input Focus |
| `Ctrl+Q` | Quit Application |
| `?` | Toggle Help Menu |

## Math Backends

TurboCrunch features a dual-backend system, allowing you to switch between different calculation engines based on your needs. Use `Ctrl+T` in the TUI to toggle between them.

### SpeedCrunch (Default)
The SpeedCrunch backend uses the original C++ core. It is the most robust and mature engine, supporting a wide range of scientific functions, units, and constants with arbitrary precision.

- **Strengths**: Extremely robust, supports advanced scientific functions, handles units and constants.
- **Recommended for**: General scientific work, unit conversions, and when using the full range of SpeedCrunch features.

### Go Backend
The Go backend is a native implementation designed for high accuracy using decimal arithmetic. It uses the `cockroachdb/apd` library for real numbers (providing 50-digit precision) and falls back to `math/cmplx` for complex numbers and trigonometric functions.

- **Strengths**: Native Go implementation, high-precision decimal arithmetic for real numbers (avoiding floating-point errors like `0.1 + 0.2`).
- **Recommended for**: Financial or high-accuracy decimal calculations where floating-point binary representation errors must be avoided.

---

## Prerequisites

### macOS
- **Go**: Version 1.20 or later.
- **Qt 5**: Required for the SpeedCrunch core.
  - Install via Homebrew: `brew install qt@5`
  - Ensure it's linked: `brew link qt@5 --force`
- **Compiler**: `clang` and `clang++` (included in Xcode Command Line Tools).
- **pkg-config**: `brew install pkg-config`

### Ubuntu / Debian
- **Go**: Version 1.20 or later.
- **Qt 5**: Development libraries.
  - Install: `sudo apt-get install qtbase5-dev`
- **Compiler**: `gcc` and `g++`.
- **pkg-config**: `sudo apt-get install pkg-config`

## Building

The project uses a `Makefile` to manage the C++ bridge and Go compilation.

1.  **Clone the repository and initialize submodules**:
    ```bash
    git clone https://github.com/robotmaxtron/turbocrunch.git
    cd turbocrunch
    git submodule update --init --recursive
    ```

2.  **Ensure Qt 5 is in your path** (especially on macOS):
    The Makefile attempts to find Qt 5 using `pkg-config`. On macOS, you might need to export the path if it's not in the default location:
    ```bash
    export PKG_CONFIG_PATH="/usr/local/opt/qt@5/lib/pkgconfig:$(pkg-config --variable=libdir Qt5Core)/pkgconfig"
    ```

3.  **Run Make**:
    ```bash
    make
    ```
    This will:
    - Run `moc` (Qt Meta-Object Compiler) on the necessary headers.
    - Compile the SpeedCrunch core and the C++ bridge.
    - Create a static library `libbridge.a`.
    - Build the `turbocrunch` Go executable.

## Running

Once built, you can run the application directly:

```bash
./turbocrunch
```

## Releases

To package the application for release:

1.  Ensure you have all dependencies installed for your platform.
2.  Run `make clean && make` to ensure a fresh build.
3.  The `turbocrunch` binary is a self-contained executable.
4.  Distribute the `turbocrunch` binary along with any necessary dynamic libraries or instructions to install Qt5.

## Project Structure

TurboCrunch uses a standard Go project layout:

- `cmd/turbocrunch/`: The application entry point (`main.go`) and TUI-related logic.
- `pkg/backend/`: High-level evaluator logic and dual-backend management (`math_wrapper.go`).
- `pkg/bridge/`: The C++ bridge and CGo integration to interface with the SpeedCrunch core.
- `SpeedCrunch/`: A Git submodule containing the original SpeedCrunch source code (see [Git Submodule](#git-submodule)).
- `Makefile`: Automates the multi-stage build process for both C++ and Go.

## Git Submodule

The SpeedCrunch core is included as a Git submodule in the `SpeedCrunch/` directory.

### Justification

1.  **Upstream Compatibility**: Keeping the SpeedCrunch source as a submodule allows us to easily pull updates or security fixes from the original repository.
2.  **Modular Architecture**: By keeping the core separate, we maintain a clear boundary between the terminal UI (TurboCrunch) and the heavy-duty math engine (SpeedCrunch).
3.  **Build Integration**: The `Makefile` is configured to reach into the submodule and compile only the necessary core components (`evaluator`, `functions`, `math` library) without requiring the full SpeedCrunch Qt GUI build.

## Best Practices & Architecture

TurboCrunch follows modern Go and C++ best practices to ensure performance and reliability:

- **Thread-Safe Bridge**: The C++ bridge uses dynamic allocation for result strings, paired with Go's `C.free`, to ensure thread safety and prevent memory leaks.
- **Idiomatic Error Handling**: Methods in the evaluation chain return `(string, error)`, allowing for robust error propagation and user-friendly error messages in the TUI.
- **Dual-Backend System**: The system can seamlessly switch between the SpeedCrunch engine and a native Go backend for flexibility.

## Testing

To verify the functionality across all packages:

```bash
make test
```

Or run individual package tests:

```bash
go test ./pkg/backend
go test ./cmd/turbocrunch
```

## Documentation

The primary documentation for TurboCrunch is available in this `README.md`. To view it or any other markdown file in the terminal with formatting, you can use a tool like `glow`:

```bash
glow README.md
```

Alternatively, you can view the Go-specific documentation for the internal packages using `godoc`:

```bash
go install golang.org/x/tools/cmd/godoc@latest
godoc -http=:6060
```
Then navigate to `http://localhost:6060/pkg/github.com/robotmaxtron/turbocrunch/`.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
