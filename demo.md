---
title: CLISlides Demo
author: CLISlides
date: 2024
theme: dark
pagination: fraction
---

# Welcome to CLISlides

A terminal-based presentation tool built with Go.

- Render Markdown beautifully in your terminal
- Navigate with arrow keys or vim-style keys
- Execute code blocks interactively
- Share presentations over SSH

---

## Features

### Markdown Support

- **Bold** and *italic* text
- Lists (ordered and unordered)
- Code blocks with syntax highlighting
- Headers at multiple levels
- Tables

---

## Code Example

Here's a simple Go program:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, CLISlides!")
}
```

Press `e` to execute code blocks!

---

## Navigation

| Key           | Action              |
|---------------|---------------------|
| → / l / j     | Next slide          |
| ← / h / k / p | Previous slide     |
| gg            | First slide         |
| G             | Last slide          |
| 5G            | Go to slide 5       |
| 3→            | Forward 3 slides    |
| /pattern      | Regex search        |
| n / N         | Next/prev match     |
| e             | Execute code block  |
| q / Esc       | Quit                |

---

## Preprocessor Demo

The following content is generated at load time:

~~~date
~~~

---

## SSH Sharing

Start a presentation server:

```bash
clislides --serve --port 2222 demo.md
```

Then audience connects with:

```bash
ssh localhost -p 2222
```

---

# Thank You!

> "The best presentations are the ones you can run in a terminal."

Happy presenting!
