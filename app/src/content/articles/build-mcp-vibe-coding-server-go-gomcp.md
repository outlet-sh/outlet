---
title: 'How To Build Your Own MCP Vibe Coding Server in Go Using GoMCP'
slug: 'build-mcp-vibe-coding-server-go-gomcp'
date: '2025-05-27'
author: 'Alma Tuck'
excerpt: 'Learn how to build a blazingly fast MCP server in Go with 20 tools that give Claude everything needed to read, write, execute and debug code.'
tags: ['Go', 'MCP', 'AI', 'Development', 'Tools']
featured: true
readTime: '8 min'
---

I don't know about you but I'm sick and tired of reading "How To" articles that show how to make yet another Todo List or the classic Hello World. I want something of substance. Something that I'll want to use again and again, long after I've finished the article.

So, let's build something useful today.

**Let's build a Vibe Coding MCP server that will write code for you!**

That sounds like it just might have the potential to be something you'll want to use after reading this article.

Let's find out.

## Why Go + MCP?

Back in 2017 I was introduced to Go. It was late Friday afternoon when one of my coworkers suggested we use Go for our next project.

Being a seasoned engineer (couple decades under my belt) I performed my solemn duty of pushing back.

Go wasn't in our existing stack and nobody in the company had ever used it. He pushed back and I finally relented, "Sure I'll try it over the weekend."

I kept my promise and surprised myself. Within 3 hours of trying to learn Go I had created my first API, connected to our database, and was hooked.

I immediately fell in love with Go's syntax and methodology of having "only one way of doing things." As a Go developer you know what I mean.

Already loving Go's speed and syntax, when MCP was announced, I was eager to learn how to write an MCP Server in Go. To my surprise I didn't find much and what I did find felt overly verbose and feature-limited with zero or close to zero testing. Basically not quite ready for enterprise users.

Go is blazingly fast and is known for high-performance APIs… why not use to create powerful MCP servers (and clients). It seemed like a perfect match!

## Introducing GoMCP

Like a good developer, I wanted to understand MCP and what better way to learn than to create a library.

**github.com/localrivet/gomcp**

Of course it has to honor the specs plus all the available features. So, I started with these self-imposed library rules:

1. It has to be intuitive and easy to use
2. It has to adhere to all specs (2024–11–05, 2025–03–26 and the latest draft)
3. It has to be feature complete
4. It has to have real testing behind it
5. It has to have an MIT license or companies can't freely use it in their projects
6. It has be understandable by AI coding tools (really important in 2025!)

After reading the docs myself, and then wising up and having AI read the docs for me, I got GoMCP up and running while keeping true to all my self-imposed rules.

Here's a working example of a server with one tool running in stdio:

```go
package main

import (
    "fmt"
    "log/slog"
    "os"

    "github.com/localrivet/gomcp/server"
)

func main() {
    // Create a logger
    logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }))

    // Create a new server
    s := server.NewServer("example-server",
        server.WithLogger(logger),
    ).AsStdio()

    // Register a tool
    s.Tool("say_hello", "Greet someone", func(ctx *server.Context, args struct {
        Name string `json:"name"`
    }) (string, error) {
        return fmt.Sprintf("Hello, %s!", args.Name), nil
    })

    // Start the server
    s.Run()
}
```

This example creates an entire stdio server in 31 lines of code (27 if you remove the comments)!

I like that.

It could be shorter too since the API is fluent… Here's the same thing in 21 lines of code.

Short is good, …but I like readability more.

```go
package main

import (
    "fmt"
    "log/slog"
    "os"

    "github.com/localrivet/gomcp/server"
)

func main() {
    server.NewServer("example-server",
        server.WithLogger(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
            Level: slog.LevelInfo,
        }))),
    ).AsStdio().Tool("say_hello", "Greet someone", func(ctx *server.Context, args struct {
        Name string `json:"name"`
    }) (string, error) {
        return fmt.Sprintf("Hello, %s!", args.Name), nil
    }).Run()
}
```

Want it to run in SSE (Server Side Events) instead? Just swap the `.AsStdio()` with `.AsSSE()`. Both easy and intuitive.

It's pretty cool, but building a say_hello MPC server isn't all that exciting, so let's create something that is way more fun.

## Building Your Own GoMCP Vibe Coder

If you don't know what Vibe Coding is yet you're a little late to the game.

Vibe coding is the methodology of letting AI code for you. Some pretty amazing (and even crazy) things have been created through Vibe Human/AI collaborations.

So, let's get started creating our very own MCP server that will write code for us. We'll then use use Claude Desktop (or any other locally run tool that can execute MCP… like Cursor) to interact with it and start vibing our own cool things.

First things first… Let's talk tools.

MCP allows us to declare tools that LLMs (AI) can use. Tools come complete with a tool name, tool description (that the LLM will understand), the input parameters (and descriptions if needed/wanted) and a return value.

## Creating GoCreate — Your GoLang Vibe Tool

Let's call our vibe coding tool GoCreate.

We're going to need a bunch of tools to make GoCreate effective.

**Configuration Tools:**

- `get_config` - Get the complete server configuration as JSON
- `set_config_value` - Set a specific configuration value by key

**Filesystem Tools:**

- `read_file` - Read the contents of a file with optional start_line and end_line parameters for paging
- `read_multiple_files` - Read the contents of multiple files simultaneously
- `write_file` - Completely replace file contents
- `create_directory` - Create a new directory or ensure a directory exists
- `list_directory` - Get a detailed listing of all files and directories in a specified path
- `move_file` - Move or rename files and directories
- `search_files` - Find files by name using case-insensitive substring matching
- `get_file_info` - Retrieve detailed metadata about a file or directory

**Search Tools:**

- `search_code` - Search for text/code patterns within file contents using ripgrep

**Edit Tools:**

- `edit_block` - Apply surgical text replacements to files
- `precise_edit` - Precisely edit file content based on start and end line numbers

**Terminal Tools:**

- `execute_command` - Execute a terminal command with timeout
- `execute_in_terminal` - Execute a command in the terminal (client-side execution)
- `read_output` - Read new output from a running terminal session
- `force_terminate` - Force terminate a running terminal session
- `list_sessions` - List all active terminal sessions

**Process Tools:**

- `list_processes` - List all running processes
- `kill_process` - Terminate a running process by PID

In total, we're going to create 20 tools for our GoCreate MCP server.

These tools give Claude (or any MCP client) everything it needs to read, write, execute and debug our code — essentially turning any AI assistant into your personal dev partner.

If you want to just jump to the code you can do so by going here:  
**github.com/localrivet/gocreate**

## Your First GoMCP Server

Let's get started by creating a new project and adding the gomcp library.

```bash
mkdir gocreate
cd gocreate
go mod init gocreate
go get github.com/localrivet/gomcp@latest
```

Now, let's write our `main.go` file.

```go
package main

import (
    "log"
    "log/slog"
    "os"

    "github.com/localrivet/gomcp/server"
)

func main() {
    // Create a logger
    logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelError,
    }))

    // Create a new server
    s := server.NewServer("GoCreate",
        server.WithLogger(logger),
    ).AsStdio()

    // Tools will go here

    // Start the server
    logger.Info("Starting GoCreate MCP server...")
    if err := s.Run(); err != nil {
        log.Fatalf("Server exited with error: %v", err)
    }
    logger.Info("Server shutdown complete.")
}
```

We're going to run our MCP server in stdio mode since we're going to be creating and editing local files. (SSE won't work in our usecase)

Let's create our first tool `read_file` which will (surprise… surprise), read a file.

```go
s.Tool("read_file", "Read the contents of a file. Supports optional start_line and end_line parameters for paging.",
    func (ctx *server.Context, args struct {
        FilePath  string `json:"file_path" description:"The path to the file to read." required:"true"`
        StartLine *int   `json:"start_line,omitempty" description:"Optional starting line number (1-indexed) for paging."`
        EndLine   *int   `json:"end_line,omitempty" description:"Optional ending line number (1-indexed, inclusive) for paging."`
    }) (string, error) {
        ctx.Logger.Info("Handling read_file tool call")

        // Read the file
        content, err := os.ReadFile(args.FilePath)
        if err != nil {
            ctx.Logger.Info("Error reading file", "file_path", args.FilePath, "error", err)
            return "Error reading file", err
        }

        fileContent := string(content)

        // If no line range specified, return the entire file
        if args.StartLine == nil && args.EndLine == nil {
            return fileContent, nil
        }

        // Handle line-based paging
        lines := strings.Split(fileContent, "\n")
        totalLines := len(lines)

        startLine := 1
        if args.StartLine != nil {
            startLine = *args.StartLine
        }

        endLine := totalLines
        if args.EndLine != nil {
            endLine = *args.EndLine
        }

        // Validate line numbers
        if startLine < 1 {
            startLine = 1
        }
        if endLine > totalLines {
            endLine = totalLines
        }
        if startLine > endLine {
            return "Invalid line range: start_line must be <= end_line", nil
        }

        // Extract the requested lines (convert to 0-based indexing)
        selectedLines := lines[startLine-1 : endLine]
        result := strings.Join(selectedLines, "\n")

        // Add line number information
        info := fmt.Sprintf("Lines %d-%d of %d total lines:\n%s", startLine, endLine, totalLines, result)
        return info, nil
    })
```

The Tool function calls for 3 parameters:

1. **The name of the tool**: `read_file`
2. **The tools description** (used to tell AI what it's for): "Read the contents of a file. Supports optional start_line and end_line parameters for paging."
3. **A function** containing the tool's logic and the return type expected: `func(ctx *server.Context, args interface{}) (interface{}, error)`

## Adding Tools

Since we don't want all our tools in a single file let's organize our code a little better.

```
gocreate/
├── main.go                # Server entry point
├── config/                # Configuration management
├── handlers/
│   ├── config/            # Configuration tools
│   ├── edit/              # Text editing tools
│   ├── filesystem/        # File system operations
│   ├── process/           # Process management
│   ├── search/            # Search functionality
│   └── terminal/          # Terminal operations
├── go.mod                 # Go module definition
```

Yay! We're organized. Not perfect but not bad either.

Let's update our `main.go` to use this new structure.

```go
package main

import (
    "log"
    "log/slog"
    "os"

    "gocreate/tools/config"
    "gocreate/tools/edit"
    "gocreate/tools/filesystem"
    "gocreate/tools/process"
    "gocreate/tools/search"
    "gocreate/tools/terminal"

    "github.com/localrivet/gomcp/server"
)

func main() {
    // Create a logger
    logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelError,
    }))

    // Create a new server
    s := server.NewServer("GoCreate",
        server.WithLogger(logger),
    ).AsStdio()

    // Register tools using the API
    // Configuration tools
    s.Tool("get_config", "Get the complete server configuration as JSON.",
        config.HandleGetConfig)

    s.Tool("set_config_value", "Set a specific configuration value by key.",
        config.HandleSetConfigValue)

    // Filesystem tools
    s.Tool("read_file", "Read the contents of a file. Supports optional start_line and end_line parameters for paging.",
        filesystem.HandleReadFile)

    s.Tool("read_multiple_files", "Read the contents of multiple files simultaneously.",
        filesystem.HandleReadMultipleFiles)

    s.Tool("write_file", "Completely replace file contents.",
        filesystem.HandleWriteFile)

    s.Tool("create_directory", "Create a new directory or ensure a directory exists.",
        filesystem.HandleCreateDirectory)

    s.Tool("list_directory", "Get a detailed listing of all files and directories in a specified path.",
        filesystem.HandleListDirectory)

    s.Tool("move_file", "Move or rename files and directories.",
        filesystem.HandleMoveFile)

    s.Tool("search_files", "Finds files by name using a case-insensitive substring matching.",
        filesystem.HandleSearchFiles)

    s.Tool("get_file_info", "Retrieve detailed metadata about a file or directory.",
        filesystem.HandleGetFileInfo)

    s.Tool("search_code", "Search for text/code patterns within file contents using ripgrep.",
        search.HandleSearchCode)

    s.Tool("edit_block", "Apply surgical text replacements to files.",
        edit.HandleEditBlock)

    s.Tool("precise_edit", "Precisely edit file content based on start and end line numbers.",
        edit.HandlePreciseEdit)

    // Terminal tools
    s.Tool("execute_command", "Execute a terminal command with timeout.",
        terminal.HandleExecuteCommand)

    s.Tool("read_output", "Read new output from a running terminal session.",
        terminal.HandleReadOutput)

    s.Tool("force_terminate", "Force terminate a running terminal session.",
        terminal.HandleForceTerminate)

    s.Tool("list_sessions", "List all active terminal sessions.",
        terminal.HandleListSessions)

    s.Tool("execute_in_terminal", "Execute a command in the terminal (client-side execution).",
        terminal.HandleExecuteInTerminal)

    // Process tools
    s.Tool("list_processes", "List all running processes.",
        process.HandleListProcesses)

    s.Tool("kill_process", "Terminate a running process by PID.",
        process.HandleKillProcess)

    // Start the server
    logger.Info("Starting GoCreate MCP server...")
    if err := s.Run(); err != nil {
        log.Fatalf("Server exited with error: %v", err)
    }
    logger.Info("Server shutdown complete.")
}
```

Clean and simple.

## Just Give Me The Code

For the sake of brevity, and the fact that I have included the code in my repo, I'm not going to add all the tools here. Just go to the repo and get 100% of the code. It's got an MIT license too so you can do anything you want to with it.

Clone the repo and build it. We're going to use it next.

```bash
git clone git@github.com:localrivet/gocreate.git
cd gocreate
go build . -o gocreate
```

## Running Our New GoCreate MCP Server

Next we need to get our new GoCreate server running so we can start vibing!

If you haven't yet download Claude Desktop https://claude.ai/download

Once installed go to **Settings > Developer tab**.

Click the **Edit Config** button to get file access to the config file named `claude_desktop_config.json`. Open that file in your favorite editor.

By default it will look like this:

```json
{
	"mcpServers": {}
}
```

Let's add our newly built gocreate compiled mcp server.

```json
{
	"mcpServers": {
		"gocreate": {
			"command": "/path/to/gocreate",
			"args": [],
			"env": {}
		}
	}
}
```

Make sure you add the full path to the compiled binary, not the directory the binary resides. In my case my path looks like this:

```
/Users/keystone/workspaces/rnd/gocreate/gocreate
# the first gocreate is my directory
# the second gocreate is my executable binary
```

Now all we have to do is restart Claude Desktop (if it's already running) for it to pick up the changes.

Once Claude Desktop is back up you can click the **Search and Tools** icon next to the + plus icon and see your new gocreate MCP server!

Click the gocreate menu item and you'll see all the tools!

So far so good!

## Now the fun part, let's use GoCreate and start Vibing!

Let's vibe a website.

For funsies I created a website with this prompt:

> "Using gocreate let's create a website for my epic hiking pictures of the mountains of Utah. Let's use tailwindcss, Sveltekit (in adapter-static mode) and let's call it My Epic Utah Mountain Hikes. Create a slideshow on the home page and ensure the imagery is presented in a way that will be stunning to the visitors.
>
> Write the files in /Users/keystone/epic-utah-hikes"

Make sure your path is writable.

After entering your own prompt (or you can copy mine) you're going to get this popup dialog from Claude Desktop.

I chose **Allow always** so it could keep vibing. For each tool Claude is going to ask permission to use it.

## One-Shot Vibe Coding

Overall, Claude is pretty good at this "vibe coding" thing. It doesn't hurt that Go is super fast and was designed for file and networking operations.

In a single shot this is what Claude created using our new GoCreate MCP server… and it only took about 11 minutes of total time.

The crazy thing is it even included instruction files!

```
➜  epic-utah-hikes tree
.
├── CUSTOMIZATION.md
├── DEPLOYMENT.md
├── IMAGE-SETUP.md
├── package.json
├── postcss.config.js
├── QUICK-START.md
├── README.md
├── src
│   ├── app.css
│   ├── app.html
│   └── routes
│       ├── +layout.svelte
│       ├── +page.svelte
│       ├── about
│       │   └── +page.svelte
│       ├── gallery
│       │   └── +page.svelte
│       └── hikes
│           └── +page.svelte
├── static
├── svelte.config.js
├── tailwind.config.js
├── tsconfig.json
└── vite.config.js

7 directories, 18 files
```

You can download the site I created here to see it for yourself:  
**https://github.com/localrivet/utah-epic-hikes**

One of my favorite things about using Go for MCP servers is it's super fast. Being a compiled language it doesn't have to run through an interpreter like so many popular languages today.

## Wrapping Up

Well, there you have it.

**GoLang + MCP = GoMCP**

GoMCP is my take on building blazingly fast, scaleable, reliable and production ready MCP servers and clients that simply work.

Here's all the code one more time:

**GoMCP library**  
if you like it, Star it — https://github.com/localrivet/gomcp

**GoCreate**  
Our Vibe Coding MCP Server — https://github.com/localrivet/gocreate

**One-Shot Website Example**  
Built by Claude Desktop using GoCreate — https://github.com/localrivet/utah-epic-hikes

Hope you enjoyed this. Happy GoMCP-ing and Vibing.

Let me know if you want more.
