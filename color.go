package main

import (
    "github.com/mgutz/ansi"
)

var reset_color = ansi.ColorCode("reset")

type ColorTheme struct {
    Verified string
    Bad string
    Plus string
    OK string
    Absent string
    Title string
    Reset string
}

func NewColorTheme(c *ColorTheme) (theme *ColorTheme) {
    color_theme := &ColorTheme{
                Verified: ansi.ColorCode(c.Verified),
                Bad:      ansi.ColorCode(c.Bad),
                Plus:     ansi.ColorCode(c.Plus),
                OK:       ansi.ColorCode(c.OK),
                Absent:   ansi.ColorCode(c.Absent),
                Title:    ansi.ColorCode(c.Title),
                Reset:    ansi.ColorCode("reset"),
            }
    return color_theme
}

