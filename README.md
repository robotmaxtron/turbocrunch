# TurboCrunch

TurboCrunch is a high-performance terminal-based calculator (TUI) that leverages the powerful SpeedCrunch evaluation engine. It features a dual-backend system, allowing users to switch between the robust SpeedCrunch C++ engine and a native Go math backend.

**Note**: This project is a fork of the [SpeedCrunch](https://github.com/speedcrunch/SpeedCrunch) project, adapted for TUI use.

## Features

- **TUI Interface**: A clean, efficient terminal user interface for quick calculations.
- **SpeedCrunch Backend**: Uses the battle-tested SpeedCrunch core for high-precision and complex math operations.
- **Go Math Backend**: A lightweight alternative using Go's `math/cmplx` package.
- **Complex Number Support**: Native support for complex arithmetic in both backends.
- **High Precision**: Maintains the high precision standards of SpeedCrunch for critical calculations.

## Prerequisites

### macOS
- **Go**: Version 1.18 or later.
- **Qt 5**: Required for the SpeedCrunch core.
  - Install via Homebrew: `brew install qt@5`
- **Compiler**: `clang` and `clang++` (included in Xcode Command Line Tools).
- **pkg-config**: `brew install pkg-config`

### Ubuntu / Debian
- **Go**: Version 1.18 or later.
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

## Testing

To verify the backend functionality:

```bash
go test -v backend_test.go math_wrapper.go
```

## Project Structure

- `SpeedCrunch/`: Git submodule containing the SpeedCrunch source code.
- `bridge/`: C++ bridge to interface Go with the SpeedCrunch core.
- `main.go`: TUI implementation and entry point.
- `math_wrapper.go`: Go-side wrapper for backend management and the Go math backend.
- `Makefile`: Build automation.
