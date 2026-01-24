# Data Model: NATO Phonetic Alphabet CLI

**Date**: 2026-01-24
**Feature**: 001-nato-cli

## Overview

This is a stateless CLI tool with no persistent storage. The data model describes the in-memory structures for phonetic alphabet mappings.

## Entities

### Alphabet

An alphabet represents a complete phonetic mapping system.

| Field | Type | Description |
|-------|------|-------------|
| Name | string | Unique identifier (e.g., "nato", "lapd") |
| DisplayName | string | Human-readable name (e.g., "NATO/ICAO") |
| Description | string | Brief description of the alphabet's origin/use |
| Letters | map[rune]string | Mapping of A-Z to phonetic words |
| Digits | map[rune]string | Mapping of 0-9 to phonetic words |

### Supported Alphabets

| Name | DisplayName | Character Set |
|------|-------------|---------------|
| nato | NATO/ICAO | A-Z, 0-9 |
| lapd | LAPD/APCO | A-Z |
| western-union | Western Union | A-Z |
| german | German (Traditional) | A-Z, Ä, Ö, Ü |

### ConversionResult

Represents the result of converting a single character.

| Field | Type | Description |
|-------|------|-------------|
| Original | rune | The input character |
| Phonetic | string | The phonetic word (empty if not convertible) |
| IsSpace | bool | True if this was a word boundary |
| IsUnknown | bool | True if character has no mapping |

### ConversionOutput

Represents the complete conversion of an input string.

| Field | Type | Description |
|-------|------|-------------|
| Input | string | Original input text |
| AlphabetUsed | string | Name of alphabet used |
| Words | [][]ConversionResult | Results grouped by input words |

## Alphabet Data

### NATO/ICAO (Default)

```
Letters:
A→Alfa, B→Bravo, C→Charlie, D→Delta, E→Echo, F→Foxtrot,
G→Golf, H→Hotel, I→India, J→Juliett, K→Kilo, L→Lima,
M→Mike, N→November, O→Oscar, P→Papa, Q→Quebec, R→Romeo,
S→Sierra, T→Tango, U→Uniform, V→Victor, W→Whiskey, X→Xray,
Y→Yankee, Z→Zulu

Digits:
0→Zero, 1→One, 2→Two, 3→Three, 4→Four, 5→Five,
6→Six, 7→Seven, 8→Eight, 9→Niner
```

### LAPD/APCO

```
Letters:
A→Adam, B→Boy, C→Charles, D→David, E→Edward, F→Frank,
G→George, H→Henry, I→Ida, J→John, K→King, L→Lincoln,
M→Mary, N→Nora, O→Ocean, P→Paul, Q→Queen, R→Robert,
S→Sam, T→Tom, U→Union, V→Victor, W→William, X→X-ray,
Y→Young, Z→Zebra

Digits: (uses NATO digits)
0→Zero, 1→One, 2→Two, 3→Three, 4→Four, 5→Five,
6→Six, 7→Seven, 8→Eight, 9→Niner
```

### Western Union

```
Letters:
A→Adams, B→Boston, C→Chicago, D→Denver, E→Easy, F→Frank,
G→George, H→Henry, I→Ida, J→John, K→King, L→Lincoln,
M→Mary, N→New York, O→Ocean, P→Peter, Q→Queen, R→Roger,
S→Sugar, T→Thomas, U→Union, V→Victor, W→William, X→X-ray,
Y→Young, Z→Zero

Digits: (uses NATO digits)
0→Zero, 1→One, 2→Two, 3→Three, 4→Four, 5→Five,
6→Six, 7→Seven, 8→Eight, 9→Niner
```

### German (Traditional)

```
Letters:
A→Anton, B→Berta, C→Cäsar, D→Dora, E→Emil, F→Friedrich,
G→Gustav, H→Heinrich, I→Ida, J→Julius, K→Kaufmann, L→Ludwig,
M→Martha, N→Nordpol, O→Otto, P→Paula, Q→Quelle, R→Richard,
S→Samuel, T→Theodor, U→Ulrich, V→Viktor, W→Wilhelm, X→Xanthippe,
Y→Ypsilon, Z→Zacharias

Special:
Ä→Ärger, Ö→Ökonom, Ü→Übermut

Digits: (uses NATO digits)
0→Zero, 1→One, 2→Two, 3→Three, 4→Four, 5→Five,
6→Six, 7→Seven, 8→Eight, 9→Niner
```

## Character Handling Rules

| Input Type | Behavior |
|------------|----------|
| Letters A-Z (any case) | Convert to phonetic word |
| Digits 0-9 | Convert to phonetic word |
| Space/Tab/Newline | Word boundary marker |
| Special characters (Ä, Ö, Ü) | Convert if alphabet supports, else skip |
| Other characters | Skip silently |

## State Transitions

N/A - This tool is stateless. Each invocation is independent.

## Validation Rules

1. Alphabet name must match a registered alphabet (case-insensitive)
2. Input text has no length limit (streams are not supported in v1)
3. Empty input shows help message
