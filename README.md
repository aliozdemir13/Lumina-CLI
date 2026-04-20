# Lumina: High-Performance CLI Sports Intelligence

[![Go Version](https://img.shields.io/github/go-mod/go-version/aliozdemir13/Lumina-CLI)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go CI](https://github.com/aliozdemir13/Lumina-CLI/actions/workflows/ci.yaml/badge.svg)](https://github.com/aliozdemir13/Lumina-CLI/actions)

**Lumina** is a sophisticated Command Line Interface (CLI) application developed in **Go (Golang)** that provides real-time sports telemetry and scoreboard data. By orchestrating data from the ESPN API, Lumina delivers a high-fidelity terminal experience for sports enthusiasts and engineers alike.

Originally built as a simple interactive tool, **Lumina v2** has been architected using the **Cobra CLI framework** to support a scalable, command-based interface.

---

## Architectural Overview

Lumina focuses on **Consuming API Services**, **Complex Data Normalization**, **Package Utilization and Command-Based UI/UX** and **System Performance**. It handles complex, nested JSON payloads from enterprise sports providers, transforming them into structured, localized, and professionally rendered tables.

### Core Features
- **Command-Based UX:** Powered by the **Cobra Framework** for an idiomatic CLI experience.
- **Enterprise Table Rendering:** Utilizes `go-pretty` for high-performance, structured data presentation.
- **Real-Time Scoreboard:** Instant updates for NBA, NFL, Champions League, and top European Soccer leagues.
- **Racing Telemetry:** Dedicated session tracking for **F1, IndyCar, and NASCAR**, including podium results and session status.
- **Global Localization:** Automated UTC-to-Local time conversion with German-style (02. Jan, 15:04) formatting.
- **Dynamic Highlights:** Visual highlighting of Red/Yellow cards and critical match events using ANSI sequences.

---

## Usage

Lumina supports subcommands and flags for deep data exploration.

### Basketball & NFL
``` bash
./lumina nba              # Default: View scores for Yesterday and Today
./lumina nba -d 5         # View a 5-day range ending today
./lumina nba -s 3 -d 2    # View a 2-day window in the future (3 to 5 days from now)
./lumina nba -s -4 -d 2   # View a 2-day window in the past (4 to 6 days ago)
./lumina nfl
```
./lumina nba -d 2
- start is 0 (Today).
- end = Today.
- start = Today - 2 days.
- Result: You see everything from 2 days ago until right now.

./lumina nba -d 2 -s 4
- start is 4 (4 days in the future).
- end = Today + 4 days.
- start = (Today + 4) - 2 = Today + 2 days.
- Result: You see the schedule for 2 days from now through 4 days from now.

./lumina nba -d 2 -s -4
- start is -4 (4 days ago).
- start = (Today - 4) - 2 = Today - 6 days.
- Result: You see the scores from 6 days ago until 4 days ago.

### Soccer (Football)

Supported leagues: ger, ita, esp, tur, cl (Champions League), eul (Europa League)
``` bash

./lumina football ger
./lumina football tur --weeks 1  # View previous week's results
./lumina football cl
```

### Racing (Motorsports)

Supported: f1, indy, nascar
``` bash

./lumina racing f1
./lumina racing f1 --all       # See full season results
./lumina racing f1 --weeks 1   # View previous week's results
```

# Technical Deep-Dive (Architect's Perspective)

Lumina implements several advanced Go patterns to demonstrate professional engineering standards:
- Command-Pattern Architecture

Utilized the Cobra CLI library to decouple command logic from data services. This allows for modular development and an extensible subcommand structure.
- Data Marshalling & Normalization

Managed enterprise-level API responses by defining precise struct mappings. Used Go's json.NewDecoder for efficient memory usage when parsing large, nested ESPN payloads.
- Structured Data Presentation

Integrated go-pretty to render complex datasets (like racing podiums and soccer highlights) into clear, readable terminal tables with custom padding and styling.
- Time Logic & Normalization

Implemented custom time-handling logic to manage the 2006-01-02T15:04Z layout, ensuring that data is presented in the user's local time zone regardless of the event location.
- Efficient Memory Management

Utilized pointer slices ([]*Score) to ensure efficient data handling and avoid unnecessary copying during high-frequency CLI interactions.

# Installation

## Clone the repository
git clone [https://github.com/aliozdemir13/Lumina-CLI.git](https://github.com/aliozdemir13/Lumina-CLI.git)

## Build the binary and Run
### Linux/MacOS
``` bash
make build
make run
```

### Windows
```bash
.\build.ps1 build
.\build.ps1 run
```

Disclaimer: This project is for educational purposes only. Lumina is not affiliated with, endorsed by, or representative of ESPN. All sports data is sourced from public API endpoints for personal use.

# Contribution
- Every PR requires at least 90% code coverage for entire code base.
- Linter rules can be found inside .golangci.yaml
