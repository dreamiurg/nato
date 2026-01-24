# Feature Specification: Enhanced Color Modes

**Feature Branch**: `003-color-modes`
**Created**: 2026-01-24
**Status**: Draft
**Input**: User description: "Two color modes - a) only first letter of the word is colored b) the whole word is colored. Colors are sticky per letter with contrast-aware adjacency."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Sticky Letter Colors (Priority: P1)

As a user, I want each letter to consistently receive the same color across the entire output, so that I can visually track repeated letters and learn letter-color associations over time.

**Why this priority**: This is the core feature that changes how colors are assigned. Without sticky colors, the other modes don't provide the expected value.

**Independent Test**: Run `nato food` and verify that both 'O' letters display the same color, and this color is consistent across multiple runs.

**Acceptance Scenarios**:

1. **Given** input contains repeated letters (e.g., "food"), **When** output is displayed, **Then** all instances of the same letter use the same color
2. **Given** input "abc abc", **When** output is displayed, **Then** 'A' has the same color in both words, 'B' has the same color in both words, etc.
3. **Given** input "aa", **When** output is displayed, **Then** both A's are the same color (adjacency exception for identical letters)

---

### User Story 2 - Contrast-Aware Adjacency (Priority: P1)

As a user, I want adjacent letters with different assignments to have visually distinct colors, so that I can easily distinguish between consecutive phonetic words.

**Why this priority**: Without contrast, adjacent letters may be hard to distinguish, defeating the purpose of coloring.

**Independent Test**: Run `nato ab` and verify that A and B have visually distinct colors with good contrast.

**Acceptance Scenarios**:

1. **Given** input "ab" where A and B have adjacent color assignments, **When** output is displayed, **Then** the colors are visually distinct (not same or similar hue)
2. **Given** input "food" (f-o-o-d), **When** output is displayed, **Then** F and O have different colors, O and D have different colors
3. **Given** input "aa", **When** output is displayed, **Then** both A's can be the same color (identical letters exempt from contrast rule)

---

### User Story 3 - Full Word Color Mode (Priority: P2)

As a user, I want to see the entire phonetic word colored (current default behavior), so that each word stands out as a complete unit.

**Why this priority**: This is the existing behavior that should be preserved as an option.

**Independent Test**: Run `nato --color-mode=word hello` and verify the entire word "Hotel" is colored, not just "H".

**Acceptance Scenarios**:

1. **Given** color mode is "word", **When** displaying "Hotel" for H, **Then** the entire word "Hotel" is colored
2. **Given** color mode is "word" with input "hello", **When** output is displayed, **Then** each complete phonetic word has its assigned color

---

### User Story 4 - First Letter Color Mode (Priority: P2)

As a user, I want an option where only the first letter of each phonetic word is colored, so that I can quickly identify the letter while keeping the rest readable in default terminal color.

**Why this priority**: Provides a less visually busy alternative for users who find full-word coloring distracting.

**Independent Test**: Run `nato --color-mode=letter hello` and verify only "H" in "Hotel" is colored, "otel" is default color.

**Acceptance Scenarios**:

1. **Given** color mode is "letter", **When** displaying "Hotel" for H, **Then** only "H" is colored, "otel" is default terminal color
2. **Given** color mode is "letter" with input "hello", **When** output is displayed, **Then** each phonetic word has only its first character colored

---

### Edge Cases

- What happens with single-character input? (Color applies to that character's phonetic word)
- What if all 26 letters appear? (Color palette cycles with contrast maintained between adjacent different letters)
- What happens with digits? (Digits get their own color assignments, same rules apply)
- What if terminal has limited color support? (Fall back gracefully to available colors)
- What about the separator between words? (Spaces/separators are not colored)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST assign a consistent color to each unique letter (A-Z) that persists for the entire output
- **FR-002**: System MUST assign a consistent color to each digit (0-9) that persists for the entire output
- **FR-003**: System MUST ensure adjacent different letters have visually contrasting colors
- **FR-004**: System MUST allow identical adjacent letters to share the same color
- **FR-005**: System MUST support "word" color mode where the entire phonetic word is colored
- **FR-006**: System MUST support "letter" color mode where only the first character of the phonetic word is colored
- **FR-007**: System MUST default to "word" color mode to maintain backward compatibility
- **FR-008**: System MUST provide a `--color-mode` flag to select between "word" and "letter" modes
- **FR-009**: System MUST use a minimum color palette of 6+ distinct colors for adequate contrast options
- **FR-010**: System MUST maintain color assignments deterministically (same input always produces same colors)

### Key Entities

- **Letter-Color Map**: A mapping from each character (A-Z, 0-9) to its assigned color
- **Color Palette**: The set of available colors with contrast relationships defined
- **Color Mode**: The display strategy ("word" = full word colored, "letter" = first character only)

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Same letter always displays same color within a single invocation (100% consistency)
- **SC-002**: Adjacent different letters have colors with sufficient contrast (no two adjacent letters share the same color unless they are identical)
- **SC-003**: Existing `nato hello` command produces visually similar output to before (backward compatible default)
- **SC-004**: Users can switch between modes with a single flag (`--color-mode`)

## Assumptions

- Terminal supports ANSI color codes (standard for modern terminals)
- The fatih/color library handles terminal capability detection
- 6-8 colors provide sufficient variety for 36 characters with contrast constraints
- Users expect deterministic coloring (not random)

## Out of Scope

- Custom color palette configuration by users
- Accessibility modes for color-blind users (could be future enhancement)
- Background colors (only foreground coloring)
- Bold/italic styling variations
