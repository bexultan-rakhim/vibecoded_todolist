# godo
## What
This is a project fully vibe coded with Ollama and Llama model under 2 hours. Key for this project was to learn how to be productive with coding agent. This project was done under 2 hour timebox (with a bit of a break in between)

## Why
Agentic AI and current coding models are good enough to create really fast mvp. Todolist app is the most repeated project and is the simplest crude application, so there is enough code in training data. Idea here is to explore how to effectively plan, and execute this project.

## How: 
Idea was to first generate a plan with llm to maximize the parallelism and spawn parallel processes to execute each. Then bridge them together in later integration stage. 

Turns out, this is not necessarily as straightforward, and llm made few mistakes here and there, but each file was short enough to execute independently still.

Writing tests and putting rules not to update the test helped to stabilize the code and ensure that later changes do not break anything important. And finally, making some general strategic decisions how to use it.

## Features
A keyboard-centric TUI todo app built with Go and the [Bubble Tea](https://github.com/charmbracelet/bubbletea) ecosystem. Tokyo Night theme, Nerd Font icons, vim-first keybindings.

- Vim-style navigation (`j`/`k`, `g`/`G`)
- Add, edit, delete, and toggle tasks
- Undo/redo for delete and edit (25-step history)
- Atomic JSON persistence — writes are crash-safe
- Responsive layout with pinned header/footer and scrolling list
- Centered modal overlay for add/edit input
- XDG Base Directory compliant config path

## Requirements

- Go 1.22+
- A [Nerd Font](https://www.nerdfonts.com/) in your terminal (for icons)

## Install

```bash
git clone https://github.com/bexultan-rakhim/vibecoded_todolist godo
cd godo
go mod tidy
go build -o godo .
./godo
```

Or run directly without building:

```bash
go run .
```

## Keybindings

### List mode

| Key | Action |
|-----|--------|
| `j` / `↓` | Move down |
| `k` / `↑` | Move up |
| `g` | Jump to top |
| `G` | Jump to bottom |
| `a` | Add task |
| `e` | Edit selected task |
| `x` / `d` | Delete selected task |
| `enter` / `space` | Toggle done/todo |
| `u` | Undo (delete or edit) |
| `r` | Redo |
| `q` / `ctrl+c` | Quit |

### Add / Edit mode

| Key | Action |
|-----|--------|
| `enter` | Confirm |
| `esc` | Cancel |

## Data storage

Tasks are saved to `~/.config/godo/data.json` after every change. If `$XDG_CONFIG_HOME` is set, that is used instead (`$XDG_CONFIG_HOME/godo/data.json`).

Saves use an atomic write-ahead strategy — changes are written to a `.tmp` file first, then renamed into place, so a crash mid-save can never corrupt your task list.

## Structure

The app follows The Elm Architecture (MVU):

```
config/             — XDG config path resolution
internal/
  task/             — Task struct, Status enum (core domain, no dependencies)
  repository/       — Repository interface + JSON and in-memory implementations
  history/          — Bounded undo/redo stack (delete + edit, 25 steps)
  navigation/       — Pure cursor logic with wrapping
  model/            — Single source of truth (Model struct, ViewMode enum)
  update/           — Update function, keybinding handlers, Cmd factories
  ui/               — Stateless rendering (styles, header, list, footer, modal)
app.go              — Bubble Tea interface wiring (avoids import cycles)
main.go             — Entry point
```

## Running tests

```bash
go test ./...
```

