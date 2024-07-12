package ui

import (
	"strings"

	"github.com/gbin/goncurses"
)

func TextBox(stdscr *goncurses.Window, y int, x int, text string) string {
	stdscr.Move(y, x)
	stdscr.MovePrint(y, x, text)
	stdscr.Refresh()

	result := ""
	total := len(text) + 0

	_, maxX := stdscr.MaxYX()

	newX := x
	for {
		ch := stdscr.GetChar()
		if ch == '\n' {
			break
		} else if ch == goncurses.KEY_BACKSPACE {
			if len(result) > 0 {
				result = result[:len(result)-1]
			}
		} else if IsAscii(ch) {
			result += string(ch)
		} else {
			continue
		}
		total = len(result) + len(text)
		centeredX := newX - (total / 2)
		stdscr.MovePrint(y, 0, strings.Repeat(" ", maxX))
		stdscr.MovePrint(y, centeredX, text)
		stdscr.MovePrint(y, centeredX+len(text), result)
		stdscr.Refresh()
	}

	stdscr.Refresh()

	return result
}

func IsAscii(ch goncurses.Key) bool {
	return ch >= 32 && ch <= 126
}

func PasswordBox(stdscr *goncurses.Window, y int, x int, text string) string {
	stdscr.Move(y, x)
	stdscr.MovePrint(y, x, text)
	stdscr.Refresh()

	_, maxX := stdscr.MaxYX()

	result := ""
	total := len(text) + 0
	newX := x
	for {
		ch := stdscr.GetChar()
		if ch == '\n' {
			break
		} else if ch == goncurses.KEY_BACKSPACE {
			if len(result) > 0 {
				result = result[:len(result)-1]
			}
		} else if IsAscii(ch) {
			result += string(ch)
		} else {
			continue
		}
		total = len(result) + len(text)
		centeredX := newX - (total / 2)
		stdscr.MovePrint(y, 0, strings.Repeat(" ", maxX))
		stdscr.MovePrint(y, centeredX, text)
		stdscr.MovePrint(y, centeredX+len(text), strings.Repeat("*", len(result)))
		stdscr.Refresh()
	}

	stdscr.Refresh()

	return result
}
