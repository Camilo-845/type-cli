# tpg — Typing Game for the Terminal

⌨️ A Monkeytype-inspired typing speed test for the terminal.

![demo](demo.gif)

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/Camilo-845/typingame/main/install.sh | bash
```

Requires Go 1.25+. Installs `tpg` to `~/.local/bin`. Ensure it's in your `PATH`.

## Usage

```
$ tpg
```

| Key | Action |
|---|---|
| `space` | Start test / submit word |
| `h` / `l` or `←` / `→` | Cycle settings left/right |
| `↑` / `↓` | Navigate menu fields |
| `backspace` | Delete last character (on empty word: jump back to previous) |
| `tab` | History |
| `esc` | Back to menu |
| `q` / `ctrl+c` | Quit |

## Game Modes

| Mode | Options |
|---|---|
| **Timed** | 15s, 30s, 60s, 120s |
| **Word count** | 10, 25, 50, 100 words |

**Word lists:** English 200, English 1k

The test auto-completes on the last correct character of the last word — no trailing space needed.
