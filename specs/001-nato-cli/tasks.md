# Tasks: NATO Phonetic Alphabet CLI

**Input**: Design documents from `/specs/001-nato-cli/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md

**Tests**: Included per FR-016 (comprehensive automated tests) and SC-004 (80%+ coverage)

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and Go module setup

- [ ] T001 Initialize Go module with `go mod init github.com/dreamiurg/nato` in go.mod
- [ ] T002 [P] Create project directory structure: cmd/, internal/alphabet/, internal/converter/, internal/output/, testdata/
- [ ] T003 [P] Add Cobra dependency with `go get github.com/spf13/cobra@latest`
- [ ] T004 [P] Add fatih/color dependency with `go get github.com/fatih/color@latest`
- [ ] T005 Create minimal main.go that calls cmd.Execute()

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core data types and interfaces that ALL user stories depend on

**CRITICAL**: No user story work can begin until this phase is complete

- [ ] T006 Define Alphabet struct (Name, DisplayName, Description, Letters, Digits maps) in internal/alphabet/types.go
- [ ] T007 Define ConversionResult struct (Original, Phonetic, IsSpace, IsUnknown) in internal/converter/types.go
- [ ] T008 Implement alphabet registry with Get(name) and List() functions in internal/alphabet/registry.go
- [ ] T009 Add NATO alphabet data (A-Z letters + 0-9 digits) to internal/alphabet/nato.go
- [ ] T010 Create Cobra root command skeleton with Execute() in cmd/root.go

**Checkpoint**: Foundation ready - user story implementation can now begin

---

## Phase 3: User Story 1 - Basic Text Conversion (Priority: P1)

**Goal**: Convert text input to NATO phonetic alphabet via command-line arguments

**Independent Test**: Run `nato dima` and verify output is "Delta India Mike Alpha"

### Tests for User Story 1

- [ ] T011 [P] [US1] Unit tests for alphabet lookup (all A-Z, 0-9) in internal/alphabet/alphabet_test.go
- [ ] T012 [P] [US1] Unit tests for converter (words, numbers, mixed case, special chars) in internal/converter/converter_test.go

### Implementation for User Story 1

- [ ] T013 [US1] Implement Convert(text, alphabetName) function in internal/converter/converter.go
- [ ] T014 [US1] Add case-insensitive character normalization in converter
- [ ] T015 [US1] Implement word boundary detection (spaces create output groupings) in converter
- [ ] T016 [US1] Handle non-letter characters (skip silently) in converter
- [ ] T017 [US1] Wire converter to Cobra command args in cmd/root.go
- [ ] T018 [US1] Implement plain-text output formatter in internal/output/formatter.go
- [ ] T019 [US1] Add integration test for CLI invocation in cmd/root_test.go

**Checkpoint**: `nato hello` outputs "Hotel Echo Lima Lima Oscar" (plain text)

---

## Phase 4: User Story 2 - Visual Enhancement with Colors (Priority: P2)

**Goal**: Display colored output by default, auto-detect TTY, support --no-color flag

**Independent Test**: Run `nato hello` in terminal and observe colored output; pipe to file and verify no ANSI codes

### Tests for User Story 2

- [ ] T020 [P] [US2] Unit tests for color output enable/disable logic in internal/output/output_test.go
- [ ] T021 [P] [US2] Unit tests for TTY detection behavior in internal/output/output_test.go

### Implementation for User Story 2

- [ ] T022 [US2] Implement ColorFormatter with fatih/color in internal/output/color.go
- [ ] T023 [US2] Add TTY detection (auto-disable colors when not terminal) in internal/output/color.go
- [ ] T024 [US2] Add NO_COLOR environment variable support in internal/output/color.go
- [ ] T025 [US2] Add --no-color flag to Cobra command in cmd/root.go
- [ ] T026 [US2] Implement alternating colors for visual distinction in internal/output/color.go
- [ ] T027 [US2] Add golden file test for colored output in internal/output/output_test.go with testdata/colored.golden

**Checkpoint**: Colors display in terminal, auto-disabled when piped

---

## Phase 5: User Story 3 - Alternative Alphabets (Priority: P3)

**Goal**: Support multiple phonetic alphabets with --alphabet flag and --list-alphabets

**Independent Test**: Run `nato --alphabet lapd hello` and verify LAPD words; run `nato --list-alphabets` to see all options

### Tests for User Story 3

- [ ] T028 [P] [US3] Unit tests for LAPD alphabet lookup in internal/alphabet/alphabet_test.go
- [ ] T029 [P] [US3] Unit tests for Western Union alphabet lookup in internal/alphabet/alphabet_test.go
- [ ] T030 [P] [US3] Unit tests for German alphabet lookup (including umlauts) in internal/alphabet/alphabet_test.go

### Implementation for User Story 3

- [ ] T031 [P] [US3] Add LAPD alphabet data in internal/alphabet/lapd.go
- [ ] T032 [P] [US3] Add Western Union alphabet data in internal/alphabet/western_union.go
- [ ] T033 [P] [US3] Add German alphabet data (including Ä, Ö, Ü) in internal/alphabet/german.go
- [ ] T034 [US3] Add --alphabet/-a flag to Cobra command in cmd/root.go
- [ ] T035 [US3] Add --list-alphabets/-l flag with formatted output in cmd/root.go
- [ ] T036 [US3] Implement helpful error message for invalid alphabet names in cmd/root.go

**Checkpoint**: All four alphabets work; --list-alphabets shows available options

---

## Phase 6: User Story 4 - Piped Input Support (Priority: P3)

**Goal**: Accept input via stdin for shell pipeline integration

**Independent Test**: Run `echo "test" | nato` and verify phonetic output

### Tests for User Story 4

- [ ] T037 [P] [US4] Unit test for stdin reading in cmd/root_test.go
- [ ] T038 [P] [US4] Integration test for piped input in cmd/root_test.go

### Implementation for User Story 4

- [ ] T039 [US4] Detect stdin pipe (check if stdin is terminal) in cmd/root.go
- [ ] T040 [US4] Read input from stdin when no args provided and stdin is pipe in cmd/root.go
- [ ] T041 [US4] Show help when no args AND no piped input in cmd/root.go

**Checkpoint**: `echo "test" | nato` works; `nato` alone shows help

---

## Phase 7: User Story 5 - Easy Installation via Homebrew (Priority: P2)

**Goal**: Distribute via Homebrew with automated releases using release-please + GoReleaser

**Independent Test**: Run `brew install dreamiurg/tap/nato` and verify `nato --version` works

### Implementation for User Story 5

- [ ] T042 [P] [US5] Create release-please-config.json with Go strategy
- [ ] T043 [P] [US5] Create .release-please-manifest.json with initial version "0.0.0"
- [ ] T044 [P] [US5] Create .goreleaser.yaml with multi-platform builds and brew config
- [ ] T045 [P] [US5] Create .github/workflows/release-please.yml for release PR automation
- [ ] T046 [P] [US5] Create .github/workflows/release.yml triggered on release publish for GoReleaser
- [ ] T047 [US5] Add --version flag using Cobra's built-in version support in cmd/root.go
- [ ] T048 [US5] Add ldflags in .goreleaser.yaml to embed version from release-please
- [ ] T049 [US5] Create dreamiurg/homebrew-tap repository (manual step - document in README)

**Checkpoint**: CI creates release PRs; merging triggers binary builds and Homebrew formula update

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Documentation, final tests, and cleanup

- [ ] T050 [P] Write README.md with installation, usage examples, and alphabet list
- [ ] T051 [P] Add CHANGELOG.md placeholder (release-please will maintain it)
- [ ] T052 Run `go test -cover ./...` and verify 80%+ coverage
- [ ] T053 Run `go vet ./...` and fix any issues
- [ ] T054 Run `gofmt -w .` to ensure consistent formatting
- [ ] T055 Validate quickstart.md scenarios work end-to-end
- [ ] T056 Create initial commit with conventional commit message: `feat: initial nato CLI implementation`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-7)**: All depend on Foundational phase completion
- **Polish (Phase 8)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational - No dependencies on other stories
- **User Story 2 (P2)**: Depends on US1 (needs working converter to add colors)
- **User Story 3 (P3)**: Can start after Foundational - Independent of US1/US2
- **User Story 4 (P3)**: Depends on US1 (needs working converter for stdin)
- **User Story 5 (P2)**: Can start after Foundational - Independent (CI/CD setup)

### Recommended Order (Sequential)

1. Phase 1: Setup
2. Phase 2: Foundational
3. Phase 3: User Story 1 (MVP)
4. Phase 4: User Story 2 (colors)
5. Phase 5: User Story 3 (alphabets) - can parallel with US2
6. Phase 6: User Story 4 (stdin)
7. Phase 7: User Story 5 (distribution) - can parallel with US3/US4
8. Phase 8: Polish

### Parallel Opportunities

**Within Setup:**
- T002, T003, T004 can run in parallel

**Within User Story 1:**
- T011, T012 (tests) can run in parallel

**Within User Story 2:**
- T020, T021 (tests) can run in parallel

**Within User Story 3:**
- T028, T029, T030 (tests) can run in parallel
- T031, T032, T033 (alphabet data) can run in parallel

**Within User Story 4:**
- T037, T038 (tests) can run in parallel

**Within User Story 5:**
- T042, T043, T044, T045, T046 (config files) can run in parallel

**Across Stories (if team capacity):**
- US3 and US5 can run in parallel (independent)
- US2 and US3 can run in parallel after US1 MVP

---

## Parallel Example: User Story 3

```bash
# Launch all tests for User Story 3 together:
Task: "Unit tests for LAPD alphabet lookup in internal/alphabet/alphabet_test.go"
Task: "Unit tests for Western Union alphabet lookup in internal/alphabet/alphabet_test.go"
Task: "Unit tests for German alphabet lookup in internal/alphabet/alphabet_test.go"

# Launch all alphabet data files together:
Task: "Add LAPD alphabet data in internal/alphabet/lapd.go"
Task: "Add Western Union alphabet data in internal/alphabet/western_union.go"
Task: "Add German alphabet data in internal/alphabet/german.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: `nato hello` outputs "Hotel Echo Lima Lima Oscar"
5. Tool is usable for basic conversion

### Incremental Delivery

1. Setup + Foundational → Foundation ready
2. Add User Story 1 → Basic conversion works (MVP!)
3. Add User Story 2 → Colors make it pleasant to use
4. Add User Story 3 → Multiple alphabets for specialized users
5. Add User Story 4 → Pipeline integration for power users
6. Add User Story 5 → Easy installation via Homebrew
7. Polish → Documentation and final cleanup

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 → User Story 2
   - Developer B: User Story 3 + User Story 5
3. Developer A adds User Story 4 after US1 complete
4. Team converges on Polish phase

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Commit after each task using Conventional Commits (feat:, test:, chore:)
- Stop at any checkpoint to validate story independently
- Go module path: `github.com/dreamiurg/nato`
- Target 80%+ test coverage per SC-004
