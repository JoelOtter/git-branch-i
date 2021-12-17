package ui

import (
	"fmt"
	"github.com/JoelOtter/git-branch-i/internal/git"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"strings"
)

func drawStr(screen tcell.Screen, x int, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		screen.SetContent(x, y, c, comb, style)
		x += w
	}
}

func draw(screen tcell.Screen, branches []git.Branch, pointer int, deleteBranch string) {
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
			style = style.Reverse(true)
		}
		drawStr(screen, 3, i+1, style, branch.Name)
	}
	if deleteBranch != "" {
		w, h := screen.Size()
		for i := 1; i < w-1; i++ {
			screen.SetCell(i, h-2, tcell.StyleDefault.Background(tcell.ColorRed))
		}
		drawStr(
			screen,
			2,
			h-2,
			tcell.StyleDefault.Background(tcell.ColorRed).Foreground(tcell.ColorBlack),
			fmt.Sprintf("Delete branch %s (y/n)? ", deleteBranch),
		)
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
	draw(screen, branches, pointer, "")
	return pointer
}

func keyUp(screen tcell.Screen, branches []git.Branch, pointer int) int {
	pointer = pointer - 1
	if pointer < 0 {
		pointer = len(branches) - 1
	}
	draw(screen, branches, pointer, "")
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
	uiOut := &strings.Builder{}
	defer func() {
		if uiOut.Len() > 0 {
			fmt.Print(uiOut.String())
		}
	}()

	pointer := getInitialPointer(branches)
	deleteBranch := ""

	draw(screen, branches, pointer, deleteBranch)

	quit := make(chan struct{})

	defer screen.Fini()
	go func() {
		defer close(quit)
		for {
			ev := screen.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape, tcell.KeyCtrlC:
					return
				case tcell.KeyEnter:
					uiErr = git.ChangeBranch(branches[pointer].Name, uiOut)
					return
				case tcell.KeyUp, tcell.KeyPgUp, tcell.KeyCtrlP:
					pointer = keyUp(screen, branches, pointer)
				case tcell.KeyDown, tcell.KeyPgDn, tcell.KeyCtrlN:
					pointer = keyDown(screen, branches, pointer)
				case tcell.KeyDelete, tcell.KeyBackspace, tcell.KeyDEL:
					deleteBranch = branches[pointer].Name
					draw(screen, branches, pointer, deleteBranch)
				case tcell.KeyRune:
					switch ev.Rune() {
					case 'j':
						pointer = keyDown(screen, branches, pointer)
					case 'k':
						pointer = keyUp(screen, branches, pointer)
					case 'y':
						if deleteBranch != "" {
							branches, uiErr = git.DeleteBranch(deleteBranch, uiOut)
							if uiErr != nil {
								return
							}
							deleteBranch = ""
							pointer = pointer - 1
							if pointer < 0 {
								pointer = 0
							}
							draw(screen, branches, pointer, deleteBranch)
						}
					case 'n':
						if deleteBranch != "" {
							deleteBranch = ""
						}
						draw(screen, branches, pointer, deleteBranch)
					}
				}
			case *tcell.EventResize:
				screen.Sync()
				draw(screen, branches, pointer, deleteBranch)
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
