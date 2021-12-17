package ui

import (
	"fmt"
	"github.com/JoelOtter/git-branch-i/internal/git"
	"github.com/gdamore/tcell/v2"
)

func draw(screen tcell.Screen, branches []git.Branch, pointer int) {
	screen.Clear()
	for i, branch := range branches {
		if branch.Current {
			screen.SetCell(1, i+1, tcell.StyleDefault, '*')
		}
		style := tcell.StyleDefault
		if branch.Current {
			style = style.Bold(true)
		}
		if i == pointer {
			style = style.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack)
		}
		screen.SetCell(3, i+1, style, []rune(branch.Name)...)
	}
	screen.Show()
}

func getInitialPointer(branches []git.Branch) int {
	for i, branch := range branches {
		if branch.Current {
			return i
		}
	}
	return 0
}

func keyDown(screen tcell.Screen, branches []git.Branch, pointer int) int {
	pointer = (pointer + 1) % len(branches)
	draw(screen, branches, pointer)
	return pointer
}

func keyUp(screen tcell.Screen, branches []git.Branch, pointer int) int {
	pointer = pointer - 1
	if pointer < 0 {
		pointer = len(branches) - 1
	}
	draw(screen, branches, pointer)
	return pointer
}

func ShowUI(branches []git.Branch) error {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	screen, err := tcell.NewScreen()
	if err != nil {
		return fmt.Errorf("failed to get screen: %w", err)
	}
	if err := screen.Init(); err != nil {
		return fmt.Errorf("failed to init screen: %w", err)
	}

	var uiErr error
	var uiOut string
	defer func() {
		if uiOut != "" {
			fmt.Print(uiOut)
		}
	}()

	pointer := getInitialPointer(branches)
	draw(screen, branches, pointer)

	quit := make(chan struct{})

	defer screen.Fini()
	go func() {
		for {
			ev := screen.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape, tcell.KeyCtrlC:
					close(quit)
					return
				case tcell.KeyEnter:
					uiOut, uiErr = git.ChangeBranch(branches[pointer].Name)
					close(quit)
					return
				case tcell.KeyUp, tcell.KeyPgUp, tcell.KeyCtrlP:
					pointer = keyUp(screen, branches, pointer)
				case tcell.KeyDown, tcell.KeyPgDn, tcell.KeyCtrlN:
					pointer = keyDown(screen, branches, pointer)
				case tcell.KeyRune:
					switch ev.Rune() {
					case 'j':
						pointer = keyDown(screen, branches, pointer)
					case 'k':
						pointer = keyUp(screen, branches, pointer)
					}
				}
			case *tcell.EventResize:
				screen.Sync()
				draw(screen, branches, pointer)
			}
		}
	}()

	for {
		select {
		case <-quit:
			return uiErr
		}
	}
}
