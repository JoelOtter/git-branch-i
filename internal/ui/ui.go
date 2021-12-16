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
			screen.SetCell(0, i, tcell.StyleDefault, '*')
		}
		style := tcell.StyleDefault
		if branch.Current {
			style = style.Bold(true)
		}
		if i == pointer {
			style = style.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack)
		}
		screen.SetCell(2, i, style, []rune(branch.Name)...)
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
				case tcell.KeyUp:
					pointer = pointer - 1
					if pointer < 0 {
						pointer = len(branches) - 1
					}
					draw(screen, branches, pointer)
				case tcell.KeyDown:
					pointer = (pointer + 1) % len(branches)
					draw(screen, branches, pointer)
				}
			case *tcell.EventResize:
				screen.Sync()
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
