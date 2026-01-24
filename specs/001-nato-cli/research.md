# Research: NATO Phonetic Alphabet CLI

**Date**: 2026-01-24
**Feature**: 001-nato-cli

## Go CLI Development Decisions

### 1. CLI Framework

**Decision**: Use Cobra for CLI framework

**Rationale**: While the built-in `flag` package would suffice for basic usage, Cobra provides:
- Built-in help and version flag generation
- Clean flag handling with both short and long forms
- Good documentation and wide adoption (Docker, Kubernetes, Hugo)
- Easy extensibility if future features are needed

**Alternatives Considered**:
- `flag` (stdlib): Sufficient for simple tools, but lacks automatic help generation and requires manual flag parsing
- Kong: Clean struct-based approach, but smaller community and less familiar to most Go developers
- urfave/cli: Good middle ground, but Cobra's ecosystem is more mature

### 2. Terminal Color Library

**Decision**: Use fatih/color

**Rationale**:
- Automatic TTY detection via go-isatty package
- Supports `NO_COLOR` environment variable (industry standard)
- Simple API with method chaining
- Wide adoption and active maintenance
- Windows 10+ support

**Alternatives Considered**:
- gookit/color: More feature-rich (256/TrueColor) but overkill for this use case
- mgutz/ansi: Fast but requires manual TTY detection

### 3. Project Structure

**Decision**: Minimal structure with internal packages

```
nato/
├── main.go              # Entry point, wires up CLI
├── go.mod
├── go.sum
├── internal/
│   ├── alphabet/        # Phonetic alphabet data and lookup
│   │   └── alphabet.go
│   └── output/          # Formatting and color output
│       └── output.go
├── testdata/            # Golden files for tests
└── README.md
```

**Rationale**: Go idiom favors shallow hierarchies. The `internal/` directory prevents external imports while allowing clean separation between alphabet logic and output formatting.

**Alternatives Considered**:
- Flat structure: Acceptable for <500 lines, but alphabet data alone will be substantial
- Full project-layout: Excessive for a single-binary tool

### 4. Testing Strategy

**Decision**: Table-driven tests with golden files for complex output

**Rationale**:
- Table-driven tests are Go standard for testing multiple scenarios
- Golden files work well for verifying colored terminal output
- Use `-update` flag pattern to regenerate expected outputs

**Alternatives Considered**:
- Inline expected values: Unwieldy for colored output with ANSI codes
- sebdah/goldie: Purpose-built for golden testing but adds dependency; manual golden files suffice

### 5. Homebrew Distribution

**Decision**: GoReleaser with dedicated homebrew-tap repository

**Setup Required**:
1. Create `dreamiurg/homebrew-tap` repository
2. Configure `.goreleaser.yaml` with brew section
3. Use GitHub Actions with Personal Access Token (not GITHUB_TOKEN)

**Rationale**: GoReleaser automates formula generation and publication on each tagged release.

**Alternatives Considered**:
- Manual formula: Requires manual updates per release
- homebrew-core: Requires notability criteria; better suited for established projects

### 6. Release Management

**Decision**: release-please for versioning + GoReleaser for binary builds

**Rationale**:
- release-please automates version bumping and changelog generation based on Conventional Commits
- Provides human-reviewable Release PRs before publishing
- GoReleaser handles the actual binary builds and Homebrew formula updates
- Clear separation: release-please decides *when* to release, GoReleaser handles *how* to build

**Workflow**:
1. Push commits to main using Conventional Commits (feat:, fix:, chore:, etc.)
2. release-please GitHub Action creates/updates a Release PR
3. Release PR accumulates changes and updates CHANGELOG.md
4. Merge Release PR → GitHub Release created automatically
5. GoReleaser triggers on release publish → builds binaries, updates Homebrew formula

**Configuration**:
- release-please: Go strategy with manifest mode
- GoReleaser: Build for darwin (amd64, arm64), linux (amd64, arm64), windows (amd64)
- Use ldflags to embed version information from release-please
- Archive as tar.gz (zip for Windows)

**Alternatives Considered**:
- GoReleaser only: Works but requires manual version bumping and changelog maintenance
- Semantic Release: JavaScript-focused, less natural fit for Go projects

---

## Phonetic Alphabet Data

### NATO/ICAO Alphabet (Default)

The international standard adopted in 1956 by NATO, ICAO, and ITU.

**Letters**:
| A | Alfa | B | Bravo | C | Charlie | D | Delta |
|---|------|---|-------|---|---------|---|-------|
| E | Echo | F | Foxtrot | G | Golf | H | Hotel |
| I | India | J | Juliett | K | Kilo | L | Lima |
| M | Mike | N | November | O | Oscar | P | Papa |
| Q | Quebec | R | Romeo | S | Sierra | T | Tango |
| U | Uniform | V | Victor | W | Whiskey | X | Xray |
| Y | Yankee | Z | Zulu |

**Spelling Notes**:
- Alfa (not Alpha): Avoids mispronunciation in languages without "ph" digraph
- Juliett (two t's): Prevents French pronunciation as "Juillet"
- Xray (no hyphen): Standardized spelling

**Digits**:
| 0 | Zero | 1 | One | 2 | Two | 3 | Three |
|---|------|---|-----|---|-----|---|-------|
| 4 | Four | 5 | Five | 6 | Six | 7 | Seven |
| 8 | Eight | 9 | Niner |

Note: Some standards use "Tree", "Fower", "Fife", "Niner" for clarity. We'll use standard English numbers with "Niner" for 9 (widely recognized to distinguish from German "nein").

### LAPD/APCO Alphabet

Standard US law enforcement alphabet (1941-1974, still widely used).

**Letters**:
| A | Adam | B | Boy | C | Charles | D | David |
|---|------|---|-----|---|---------|---|-------|
| E | Edward | F | Frank | G | George | H | Henry |
| I | Ida | J | John | K | King | L | Lincoln |
| M | Mary | N | Nora | O | Ocean | P | Paul |
| Q | Queen | R | Robert | S | Sam | T | Tom |
| U | Union | V | Victor | W | William | X | X-ray |
| Y | Young | Z | Zebra |

### Western Union Alphabet

Pre-NATO civilian/telegraphy alphabet using American city names.

**Letters**:
| A | Adams | B | Boston | C | Chicago | D | Denver |
|---|-------|---|--------|---|---------|---|--------|
| E | Easy | F | Frank | G | George | H | Henry |
| I | Ida | J | John | K | King | L | Lincoln |
| M | Mary | N | New York | O | Ocean | P | Peter |
| Q | Queen | R | Roger | S | Sugar | T | Thomas |
| U | Union | V | Victor | W | William | X | X-ray |
| Y | Young | Z | Zero |

### German Alphabet (Traditional - Anton/Berta)

**Letters**:
| A | Anton | B | Berta | C | Cäsar | D | Dora |
|---|-------|---|-------|---|-------|---|------|
| E | Emil | F | Friedrich | G | Gustav | H | Heinrich |
| I | Ida | J | Julius | K | Kaufmann | L | Ludwig |
| M | Martha | N | Nordpol | O | Otto | P | Paula |
| Q | Quelle | R | Richard | S | Samuel | T | Theodor |
| U | Ulrich | V | Viktor | W | Wilhelm | X | Xanthippe |
| Y | Ypsilon | Z | Zacharias |

**Special Characters**:
| Ä | Ärger | Ö | Ökonom | Ü | Übermut |

---

## Dependencies Summary

| Purpose | Package | Version |
|---------|---------|---------|
| CLI framework | github.com/spf13/cobra | latest |
| Color output | github.com/fatih/color | latest |
| TTY detection | github.com/mattn/go-isatty | (transitive via fatih/color) |

## Sources

- [NATO phonetic alphabet - Wikipedia](https://en.wikipedia.org/wiki/NATO_phonetic_alphabet)
- [APCO radiotelephony spelling alphabet - Wikipedia](https://en.wikipedia.org/wiki/APCO_radiotelephony_spelling_alphabet)
- [Western Union Phonetic Alphabet - osric.com](https://osric.com/phonetic/index.php?id=8)
- [DIN 5009 - Wikipedia](https://en.wikipedia.org/wiki/DIN_5009)
- [fatih/color - GitHub](https://github.com/fatih/color)
- [Cobra - GitHub](https://github.com/spf13/cobra)
- [GoReleaser Documentation](https://goreleaser.com/)
- [release-please - GitHub](https://github.com/googleapis/release-please)
