# focus-time

A cli tool that tracks window focus to help with backlogging time tracking

## Features

- Tracks time spend in focused window over time
- Relational sqlite storage takes minimal overhead

## Limitations

- Currently only supports windows os
- Report generation dumps all records

## Goals

In no particular order

- os agnostic
- improved api
  - restful http
  - cli flags
- html reports using go templating
- more analysis tools (e.g. filter / sort)
- cool tui

## Install

With Go on your machine, install:

```bash
go install github.com/Isaac799/focus-time/cmd/focustime@latest
```

## Usage 

Ran with:

```bash
focustime
```

No cli flags or anything for now. The app is designed to remain open.