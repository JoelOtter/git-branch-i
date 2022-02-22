# `git branch-i`

I got cross that there's no `git branch --interactive`, so I made this. It's a 
very (very) simple curses-mode `git branch`/`git checkout` alternative which
lets you switch branches quickly using the arrow and return keys. It also lets
you delete branches without needing to type their name or use (_shudder_) a GUI.

It should capture and pipe git output correctly. The attempt is to make it feel
like as much of a 'real' Git tool as possible.

## Installation

Install it using Go, like so:

```sh
go install github.com/JoelOtter/git-branch-i@latest
```

You can now run `git branch-i` to use the tool, so long as your Go directory is 
on your system path.

## Usage

* Branches can be navigated using the arrow keys, j and k, Pg Up/Down, or
Ctrl+N and Ctrl+P.
* Checkout a branch with the return key.
* Delete a branch with the delete or backspace key, and use y/n to confirm.
* Exit with Escape or Ctrl+C.
