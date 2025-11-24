# newline CLI tool

## Problem

AI agents often do not ensure a single newline at EOF, even when asked.

Why this matters: a single newline at EOF corresponds to Git's default core.whitespace blank-at-eof option.

## User story

As a developer who works mostly on the command line, I want a tool that AI agents can call to ensure a single newline at EOF of any files I wish.

I want to control which files this is applied to. This can be handled by rules in AGENTS.md, and/or be part of specific prompts.

While there are many tools that have built-in support to ensure this (see https://github.com/editorconfig/editorconfig/wiki/Newline-at-End-of-File-Support), having a CLI tool that I can expose to AI agents or use myself as needed seems to be the most flexible solution when working primarily on the command line.

## Tool requirements

- The ability to rewrite a file so it ends with exactly one newline.
- The ability to processes a file or all files recursively in an entire directory.
- When processing all files in dir, I want to skip hidden subdirs unless dir itself is hidden, or unless the option to process all hidden dirs is explicitly opted-into.

## Reference

I created a zsh command `newline` that does exactly what I want.
I also created a zsh command `newline_test` that tests exactly the requirements above, except I did not yet include a way to opt-into recursive hidden dirs.

These two commands are in a `.zshrc` file in this workspace, for reference.

## Rewrite these in Go

zsh/shell commands work well, but I would like something that can be more portable, easier to maintain, and easier to import into agent environments.

I have started a Go port - files in this workspace - but it is not yet functional. Experimenting to see if AI codegen can improve this. If nothing else, the prompting forces me to outline my requirements.

## Testing

go test should be at least as robust as `newline_test` from `.zshrc` in this workspace.

## Makefile

Write a simple Makefile with targets to test and compile Go.

## trailingspace CLI tool

### Problem

AI agents often leave trailing whitespace at the end of lines, even when asked not to.

Why this matters: trailing whitespace at end-of-line corresponds to Git's default core.whitespace blank-at-eol option.

### User story

As a developer who works mostly on the command line, I want a tool that AI agents can call to remove trailing whitespace from any files I wish.

I want to control which files this is applied to. This can be handled by rules in AGENTS.md, and/or be part of specific prompts.

### Tool requirements

- The ability to rewrite files to remove trailing spaces and tabs from the end of each line.
- The ability to process a file or all files recursively in an entire directory.
- When processing all files in dir, skip hidden subdirs unless dir itself is hidden, or unless the option to process all hidden dirs is explicitly opted-into.
- Should preserve file content otherwise (no changes to line endings or other whitespace).

### Reference

I created a zsh command `trailingspace` that does exactly what I want, located in `.zshrc` in this workspace.

The logic:
- Uses `perl -pe 's/[ \t]+$//'` to remove trailing spaces and tabs from each line
- Same directory traversal behavior as newline tool
- Creates temporary files for safe atomic updates

### Implementation

- Part of monorepo with newline CLI tool
- Similar architecture and patterns as newline
- Should support same `--hidden` flag for processing hidden directories
- Both tools released together, installable as a suite

## Process

Running Claude in a container, using this simple project: https://github.com/scottrigby/claude-container

Enabling YOLO mode.

Work until either clarification is needed, or tasks are complete.

A script in your PATH is provided called `notify.sh`.

When clarification is needed: Send notification when clarification needed: `notify.sh "Need clarification: reason"`

For task completion: Send audio notification: `notify.sh "completion message"`
