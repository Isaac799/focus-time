# focus-time

A cli tool that tracks window focus to help with backlogging time tracking

## Features

- Tracks time spend in focused window over time
- Relational sqlite storage takes minimal overhead

## Limitations

- Currently only supports windows os
- Report generation dumps all records

## Goals

- more analysis tools (e.g. filter / sort)
- ability to remove old records
- os agnostic (low priority due to complexity of various window managers)

## Install

With Go on your machine, install:

```bash
go install github.com/Isaac799/focus-time/cmd/focustime@latest
```

## Usage

The default command is non-verbose watch.

### Watch

Will continuously monitor the currently focused window, and save how long it was focused once focus is lost. Saved in a sqlite db in home dir.

Command

```bash
focustime --action=watch --verbose
```

Outout

```
watcher active
keep window open
press [ctrl + c] to quit
running in verbose mode
0 README.md - focus-time - Visual Studio Code
1 README.md - focus-time - Visual Studio Code
2 README.md - focus-time - Visual Studio Code
0 focustime and 1 more tab - File Explorer
1 focustime and 1 more tab - File Explorer
```

### Print to console

Command

`yyyy-mm-dd`

```bash
focustime --action=print | grep 2025-05-11
```

Outout

```bash
Visual Studio Code  report.go - focus-time      2025-05-11  7m22s
Visual Studio Code  focustime.go - focus-time   2025-05-11  1m37s
Visual Studio Code  launch.json - focus-time    2025-05-11  14s
```

### Export to csv

Command

```bash
focustime --action=csv
```
Output is to current working dir

| Group              | Title                     | When       | Duration |
| ------------------ | ------------------------- | ---------- | -------- |
| Visual Studio Code | report.go - focus-time    | 2025-05-11 | 7m22s    |
| Visual Studio Code | focustime.go - focus-time | 2025-05-11 | 1m37s    |
| Visual Studio Code | launch.json - focus-time  | 2025-05-11 | 14s      |
