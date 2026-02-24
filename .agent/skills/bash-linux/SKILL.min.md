---
name: bash-linux
description: Bash/Linux terminal patterns. Critical commands, piping, error
  handling, scripting. Use when working on macOS or Linux systems.
allowed-tools: Read Write Edit Glob Grep Bash
---
# Bash Linux Patterns

> Essential patterns for Bash on Linux/macOS.

---

## 1. Operator Syntax

### Chaining Commands

---

## 2. File Operations

### Essential Commands

---

## 3. Process Management

---

## 4. Text Processing

### Core Tools

---

## 5. Environment Variables

---

## 6. Network

---

## 7. Script Template

```bash
#!/bin/bash
set -euo pipefail  # Exit on error, undefined var, pipe fail

# Colors (optional)
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Functions
log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1" >&2; }

# Main
main() {
    log_info "Starting..."
    # Your logic here
    log_info "Done!"
}

main "$@"
```

---

## 8. Common Patterns

### Check if command exists

```bash
if command -v node &> /dev/null; then
    echo "Node is installed"
fi
```

### Default variable value

```bash
NAME=${1:-"default_value"}
```

### Read file line by line

```bash
while IFS= read -r line; do
    echo "$line"
done < file.txt
```

### Loop over files

```bash
for file in *.js; do
    echo "Processing $file"
done
```

---

## 9. Differences from PowerShell

---

## 10. Error Handling

### Set options

```bash
set -e          # Exit on error
set -u          # Exit on undefined variable
set -o pipefail # Exit on pipe failure
set -x          # Debug: print commands
```

### Trap for cleanup

```bash
cleanup() {
    echo "Cleaning up..."
    rm -f /tmp/tempfile
}
trap cleanup EXIT
```

---

> **Remember:** Bash is text-based. Use `&&` for success chains, `set -e` for safety, and quote your variables!