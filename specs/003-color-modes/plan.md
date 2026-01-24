# Implementation Plan: Enhanced Color Modes

**Feature**: 003-color-modes
**Spec**: specs/003-color-modes/spec.md
**Created**: 2026-01-24

## Phase 0: Documentation & API Reference

### Allowed APIs (fatih/color v1.18.0)

**Color Creation:**
- `color.New(attributes...)` - Creates Color instance
- `color.Attribute` - Type for color constants

**Foreground Colors Available:**
- Standard: `FgBlack`, `FgRed`, `FgGreen`, `FgYellow`, `FgBlue`, `FgMagenta`, `FgCyan`, `FgWhite`
- High-intensity: `FgHiBlack`, `FgHiRed`, `FgHiGreen`, `FgHiYellow`, `FgHiBlue`, `FgHiMagenta`, `FgHiCyan`, `FgHiWhite`

**Color Output:**
- `c.Sprint(text)` - Returns colored string
- `c.SprintFunc()` - Returns function that colors strings

### Existing Codebase Patterns

**Files to Modify:**
- `internal/output/formatter.go` - Main color logic (lines 14-61)
- `cmd/root.go` - CLI flags (lines 16-41, 84)
- `internal/output/output_test.go` - Tests

**Key Data Structure:**
- `converter.ConversionResult.Original` - Contains the original rune (letter) - USE THIS for sticky colors

**Current Color Cycling (to replace):**
```go
// Current: position-based cycling
c := f.colors[i%len(f.colors)]
```

### Anti-Patterns to Avoid
- Do NOT use `color.NoColor` global flag (use formatter's noColor field)
- Do NOT use 256-color or RGB modes (stick to 16-color for terminal compatibility)
- Do NOT pass color mode through converter - keep it in formatter only

---

## Phase 1: Expand Color Palette

**Objective**: Increase color palette from 4 to 8+ colors for better contrast options.

### Tasks

1. **Expand color palette in `internal/output/formatter.go`**

   Replace lines 37-42:
   ```go
   colors: []*color.Color{
       color.New(color.FgCyan),
       color.New(color.FgYellow),
       color.New(color.FgGreen),
       color.New(color.FgMagenta),
   },
   ```

   With expanded palette (8 colors with good visual contrast):
   ```go
   colors: []*color.Color{
       color.New(color.FgRed),
       color.New(color.FgGreen),
       color.New(color.FgYellow),
       color.New(color.FgBlue),
       color.New(color.FgMagenta),
       color.New(color.FgCyan),
       color.New(color.FgHiRed),
       color.New(color.FgHiGreen),
   },
   ```

### Verification
- [ ] `go build ./...` succeeds
- [ ] `go test ./...` passes
- [ ] `./nato hello` shows colors from new palette

---

## Phase 2: Add ColorMode Type and Flag

**Objective**: Add `--color-mode` CLI flag with "word" (default) and "letter" modes.

### Tasks

1. **Define ColorMode type in `internal/output/formatter.go`**

   Add after imports (before Formatter struct):
   ```go
   // ColorMode determines how colors are applied to output.
   type ColorMode string

   const (
       ColorModeWord   ColorMode = "word"   // Color entire phonetic word
       ColorModeLetter ColorMode = "letter" // Color only first character
   )
   ```

2. **Update Formatter struct** (lines 14-18)

   Add `mode` field:
   ```go
   type Formatter struct {
       noColor bool
       isTTY   bool
       colors  []*color.Color
       mode    ColorMode
   }
   ```

3. **Update NewFormatter signature** (line 21)

   Change from:
   ```go
   func NewFormatter(noColor bool) *Formatter
   ```

   To:
   ```go
   func NewFormatter(noColor bool, mode ColorMode) *Formatter
   ```

   Add in constructor body:
   ```go
   mode: mode,
   ```

4. **Add CLI flag in `cmd/root.go`**

   Add variable (after line 19):
   ```go
   colorMode string
   ```

   Add flag in `init()` (after line 40):
   ```go
   rootCmd.Flags().StringVar(&colorMode, "color-mode", "word", "color mode: 'word' (full word) or 'letter' (first char only)")
   ```

   Update formatter call (line 84):
   ```go
   formatter := output.NewFormatter(noColor, output.ColorMode(colorMode))
   ```

5. **Update tests in `internal/output/output_test.go`**

   Update all `NewFormatter` calls to include mode:
   ```go
   f := NewFormatter(true, ColorModeWord)
   ```

### Verification
- [ ] `go build ./...` succeeds
- [ ] `go test ./...` passes
- [ ] `./nato --help` shows `--color-mode` flag
- [ ] `./nato --color-mode=word hello` works
- [ ] `./nato --color-mode=letter hello` works (output same as word for now)

---

## Phase 3: Implement Sticky Letter Colors

**Objective**: Build letter→color mapping so same letter always gets same color.

### Tasks

1. **Add letter-to-color-index mapping in `internal/output/formatter.go`**

   Add new method to Formatter:
   ```go
   // letterColorIndex returns a deterministic color index for a given rune.
   // Same letter always maps to same color index.
   func (f *Formatter) letterColorIndex(r rune) int {
       // Normalize to uppercase
       upper := unicode.ToUpper(r)

       // Map A-Z to 0-25, 0-9 to 26-35
       var index int
       if upper >= 'A' && upper <= 'Z' {
           index = int(upper - 'A')
       } else if upper >= '0' && upper <= '9' {
           index = 26 + int(upper-'0')
       } else {
           index = 0
       }

       // Map to color palette using modulo
       return index % len(f.colors)
   }
   ```

   Add import: `"unicode"`

2. **Update Print method to use sticky colors** (lines 47-61)

   Replace the color selection logic:
   ```go
   func (f *Formatter) Print(words []converter.Word) {
       for _, word := range words {
           var phonetics []string
           for _, result := range word.Results {
               if f.noColor {
                   phonetics = append(phonetics, result.Phonetic)
               } else {
                   // Use sticky color based on original letter
                   colorIdx := f.letterColorIndex(result.Original)
                   c := f.colors[colorIdx]
                   phonetics = append(phonetics, c.Sprint(result.Phonetic))
               }
           }
           fmt.Println(strings.Join(phonetics, " "))
       }
   }
   ```

### Verification
- [ ] `go build ./...` succeeds
- [ ] `go test ./...` passes
- [ ] `./nato food` - both O's have the same color
- [ ] `./nato abc abc` - matching letters across words have same color
- [ ] `./nato abcdefghij` - colors cycle through palette

---

## Phase 4: Implement Contrast-Aware Adjacency

**Objective**: Ensure adjacent different letters don't share the same color.

### Tasks

1. **Refactor Print to handle adjacency** in `internal/output/formatter.go`

   Replace Print method with contrast-aware version:
   ```go
   func (f *Formatter) Print(words []converter.Word) {
       for _, word := range words {
           var phonetics []string
           var prevColorIdx = -1
           var prevRune rune

           for _, result := range word.Results {
               if f.noColor {
                   phonetics = append(phonetics, result.Phonetic)
                   continue
               }

               // Get base color index for this letter
               colorIdx := f.letterColorIndex(result.Original)

               // Contrast-aware: if adjacent different letters would have same color, shift
               if prevColorIdx >= 0 && colorIdx == prevColorIdx && result.Original != prevRune {
                   colorIdx = (colorIdx + 1) % len(f.colors)
               }

               c := f.colors[colorIdx]
               phonetics = append(phonetics, c.Sprint(result.Phonetic))

               prevColorIdx = colorIdx
               prevRune = result.Original
           }
           fmt.Println(strings.Join(phonetics, " "))
       }
   }
   ```

### Verification
- [ ] `go build ./...` succeeds
- [ ] `go test ./...` passes
- [ ] `./nato ab` - A and B have different colors
- [ ] `./nato aa` - both A's have the same color (identical letters exempt)
- [ ] `./nato food` - F≠O, O=O, O≠D colors

---

## Phase 5: Implement Letter-Only Color Mode

**Objective**: In "letter" mode, only the first character of each phonetic word is colored.

### Tasks

1. **Update Print to apply mode** in `internal/output/formatter.go`

   Modify the color application section:
   ```go
   func (f *Formatter) Print(words []converter.Word) {
       for _, word := range words {
           var phonetics []string
           var prevColorIdx = -1
           var prevRune rune

           for _, result := range word.Results {
               if f.noColor {
                   phonetics = append(phonetics, result.Phonetic)
                   continue
               }

               // Get base color index for this letter
               colorIdx := f.letterColorIndex(result.Original)

               // Contrast-aware: if adjacent different letters would have same color, shift
               if prevColorIdx >= 0 && colorIdx == prevColorIdx && result.Original != prevRune {
                   colorIdx = (colorIdx + 1) % len(f.colors)
               }

               c := f.colors[colorIdx]

               // Apply color based on mode
               var colored string
               if f.mode == ColorModeLetter && len(result.Phonetic) > 0 {
                   // Color only first character
                   first := string(result.Phonetic[0])
                   rest := result.Phonetic[1:]
                   colored = c.Sprint(first) + rest
               } else {
                   // Color entire word (default)
                   colored = c.Sprint(result.Phonetic)
               }

               phonetics = append(phonetics, colored)

               prevColorIdx = colorIdx
               prevRune = result.Original
           }
           fmt.Println(strings.Join(phonetics, " "))
       }
   }
   ```

### Verification
- [ ] `go build ./...` succeeds
- [ ] `go test ./...` passes
- [ ] `./nato --color-mode=word hello` - entire words colored
- [ ] `./nato --color-mode=letter hello` - only first char of each word colored
- [ ] `./nato hello` - defaults to word mode

---

## Phase 6: Add Tests for New Functionality

**Objective**: Add comprehensive tests for sticky colors, adjacency, and modes.

### Tasks

1. **Add tests in `internal/output/output_test.go`**

   ```go
   func TestColorModeDefaults(t *testing.T) {
       f := NewFormatter(true, ColorModeWord)
       if f.mode != ColorModeWord {
           t.Errorf("expected ColorModeWord, got %v", f.mode)
       }
   }

   func TestLetterColorIndexConsistency(t *testing.T) {
       f := NewFormatter(true, ColorModeWord)
       // Same letter should always return same index
       idx1 := f.letterColorIndex('A')
       idx2 := f.letterColorIndex('A')
       if idx1 != idx2 {
           t.Error("letterColorIndex should be deterministic")
       }
       // Different letters should potentially differ
       idxA := f.letterColorIndex('A')
       idxB := f.letterColorIndex('B')
       if idxA == idxB {
           t.Error("A and B should have different color indices")
       }
   }

   func TestLetterColorIndexCaseInsensitive(t *testing.T) {
       f := NewFormatter(true, ColorModeWord)
       if f.letterColorIndex('a') != f.letterColorIndex('A') {
           t.Error("letterColorIndex should be case-insensitive")
       }
   }
   ```

### Verification
- [ ] `go test ./...` passes
- [ ] Tests cover: ColorMode type, sticky colors, case insensitivity, adjacency

---

## Phase 7: Final Verification

### Checklist

**Build & Test:**
- [ ] `go build ./...` succeeds
- [ ] `go test ./...` passes
- [ ] `go vet ./...` passes

**Functional Requirements (from spec.md):**
- [ ] FR-001: Same letter → same color (`nato food` both O's match)
- [ ] FR-002: Same digit → same color (`nato 123 123` matching digits)
- [ ] FR-003: Adjacent different letters → different colors (`nato ab`)
- [ ] FR-004: Identical adjacent letters → same color allowed (`nato aa`)
- [ ] FR-005: Word mode colors entire word (`--color-mode=word`)
- [ ] FR-006: Letter mode colors first char only (`--color-mode=letter`)
- [ ] FR-007: Default is word mode (backward compatible)
- [ ] FR-008: `--color-mode` flag works
- [ ] FR-009: 6+ colors in palette (8 implemented)
- [ ] FR-010: Deterministic colors (multiple runs produce same output)

**User Story Tests:**
- [ ] `nato food` - O's same color ✓
- [ ] `nato abc abc` - matching letters across words ✓
- [ ] `nato ab` - A and B visually distinct ✓
- [ ] `nato --color-mode=word hello` - full word colored ✓
- [ ] `nato --color-mode=letter hello` - only H colored in "Hotel" ✓

**Backward Compatibility:**
- [ ] `nato hello` (no flags) produces colored output
- [ ] `nato --no-color hello` produces plain output
- [ ] Piped output still disables colors

### Anti-Pattern Check
- [ ] No hardcoded 256-color codes
- [ ] No global color state modifications
- [ ] No changes to converter package
