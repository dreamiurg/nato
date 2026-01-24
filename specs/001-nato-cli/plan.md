# Implementation Plan: NATO Phonetic Alphabet CLI

**Branch**: `001-nato-cli` | **Date**: 2026-01-24 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-nato-cli/spec.md`

## Summary

Build a Go CLI tool that converts text to NATO phonetic alphabet with colored output, multiple alphabet support, and Homebrew distribution. The tool prioritizes simplicity, good UX (colors, clear output), and easy installation.

## Technical Context

**Language/Version**: Go 1.21+
**Primary Dependencies**: spf13/cobra (CLI), fatih/color (terminal colors)
**Storage**: N/A (stateless)
**Testing**: Go testing with table-driven tests and golden files
**Target Platform**: macOS (arm64, amd64), Linux (amd64, arm64), Windows (amd64)
**Project Type**: Single binary CLI tool
**Performance Goals**: <100ms for any input (spec SC-001)
**Constraints**: Zero runtime dependencies, single binary distribution
**Scale/Scope**: Personal/small team utility

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

The project constitution is a template (not configured). Proceeding with standard Go best practices:
- [x] Simple, focused functionality (single-purpose tool)
- [x] Comprehensive test coverage (80%+ target per spec)
- [x] Clean code structure with internal packages
- [x] Automated releases via release-please + GoReleaser

## Project Structure

### Documentation (this feature)

```text
specs/001-nato-cli/
├── spec.md              # Feature specification
├── plan.md              # This file
├── research.md          # Phase 0 research output
├── data-model.md        # Phase 1 data model
├── quickstart.md        # Phase 1 quickstart guide
├── checklists/
│   └── requirements.md  # Spec quality checklist
└── tasks.md             # Phase 2 output (created by /speckit.tasks)
```

### Source Code (repository root)

```text
nato/
├── main.go                  # Entry point, minimal - delegates to cmd
├── go.mod
├── go.sum
├── .goreleaser.yaml         # Binary build configuration
├── release-please-config.json  # release-please configuration
├── .release-please-manifest.json  # Version tracking
├── .github/
│   └── workflows/
│       ├── release-please.yml  # Creates release PRs on push to main
│       └── release.yml         # Triggers GoReleaser on release publish
├── internal/
│   ├── alphabet/
│   │   ├── alphabet.go      # Alphabet data structures and lookup
│   │   └── alphabet_test.go
│   ├── converter/
│   │   ├── converter.go     # Text-to-phonetic conversion logic
│   │   └── converter_test.go
│   └── output/
│       ├── output.go        # Formatting and color handling
│       └── output_test.go
├── cmd/
│   └── root.go              # Cobra root command and flags
├── testdata/
│   └── *.golden             # Expected output files for tests
└── README.md
```

**Structure Decision**: Single-project structure with `internal/` for private packages. The `cmd/` directory follows Cobra conventions. Three internal packages separate concerns: alphabet data, conversion logic, and output formatting.

## Complexity Tracking

No constitution violations. Design follows Go idioms for simple CLI tools.

## Design Artifacts

### Key Components

1. **Alphabet Registry** (`internal/alphabet/`)
   - Stores phonetic mappings for NATO, LAPD, Western Union, German
   - Provides lookup by letter/digit
   - Returns "unknown" marker for unsupported characters

2. **Converter** (`internal/converter/`)
   - Accepts input text and alphabet name
   - Iterates through characters, looks up phonetic words
   - Preserves word boundaries (input spaces become output line breaks or separators)
   - Case-insensitive processing

3. **Output Formatter** (`internal/output/`)
   - Detects TTY for automatic color enable/disable
   - Respects `--no-color` flag and `NO_COLOR` environment variable
   - Applies alternating or grouped colors for readability
   - Handles piped output gracefully

4. **CLI Interface** (`cmd/root.go`)
   - Root command: `nato <text>` or `echo text | nato`
   - Flags: `--alphabet/-a`, `--no-color`, `--list-alphabets`, `--version`, `--help`
   - Stdin detection for piped input

### Data Flow

```
Input (args or stdin)
    ↓
Normalize (lowercase, handle spaces)
    ↓
For each character:
    ↓
  Alphabet Lookup → Phonetic word or skip
    ↓
Collect results with word boundaries
    ↓
Format with colors (if TTY)
    ↓
Output to stdout
```

### Color Scheme

Default color scheme for visual distinction:
- First letter of each phonetic word: Bold + color
- Alternating colors per word for readability
- Word boundaries: Clear visual separation (newline or spacing)

## Implementation Phases

### Phase 1: Core Conversion (P1 User Story)
- Alphabet data structures with NATO alphabet
- Basic converter logic
- Plain-text output (no colors)
- Command-line argument input
- Tests for all A-Z and 0-9 conversions

### Phase 2: Enhanced Output (P2 User Stories)
- Color support with fatih/color
- TTY detection and auto-disable
- `--no-color` flag
- Stdin/pipe support
- `--version` and `--help` flags

### Phase 3: Multiple Alphabets (P3 User Story)
- Add LAPD, Western Union, German alphabets
- `--alphabet` flag for selection
- `--list-alphabets` flag
- Error handling for invalid alphabet names

### Phase 4: Distribution (P2 User Story)
- release-please configuration (version bumping, changelog generation)
- GoReleaser configuration (binary builds, Homebrew formula)
- GitHub Actions workflows (release-please PR + GoReleaser on publish)
- Homebrew tap repository setup
- README documentation

## Test Strategy

### Unit Tests
- Alphabet lookup: All characters for all alphabets
- Converter: Various input patterns (words, numbers, mixed, special chars)
- Output: Color enable/disable logic

### Integration Tests
- CLI invocation with arguments
- Piped input handling
- Flag combinations
- Golden file comparisons for output format

### Coverage Target
80%+ code coverage (per spec SC-004)

## Release Strategy

1. **Versioning**: Semantic versioning via release-please (Conventional Commits)
2. **Release Flow**:
   - Push to main with Conventional Commits (feat:, fix:, chore:, etc.)
   - release-please creates/updates a Release PR with changelog and version bump
   - Merge Release PR → GitHub Release created automatically
   - GoReleaser triggers on release publish → builds binaries + updates Homebrew formula
3. **Platforms**: darwin/amd64, darwin/arm64, linux/amd64, linux/arm64, windows/amd64
4. **Distribution**:
   - GitHub Releases (binaries via GoReleaser)
   - Homebrew tap (`dreamiurg/tap/nato`)

## Dependencies

| Package | Purpose | Justification |
|---------|---------|---------------|
| github.com/spf13/cobra | CLI framework | Industry standard, auto-generates help |
| github.com/fatih/color | Terminal colors | TTY detection, NO_COLOR support |

No other runtime dependencies. Build/release tools (release-please, GoReleaser) are CI-only.
