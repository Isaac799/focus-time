# Focus Time

A cli tool that tracks window focus and provides simple reporting.

Originally designed to help backlog time tracking.

## Features

- Automatically tracks time spent in focused windows
- Reporting
  - Grouped by suffix (typically application name)
  - Reasonable default filters
  - CSV export enables more advanced analysis
  - Console print enables other tools like grep
- Relational sqlite ensures no wasted disk space

## Limitations

- ability to remove old records
- os agnostic (low priority due to complexity of various window managers)

## Install

With Go on your machine, install:

```bash
go install github.com/Isaac799/focus-time/cmd/focustime@latest
```

## Usage

Thanks to go flags this command explains it all.

```bash
focustime --help
```