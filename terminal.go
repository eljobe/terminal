// Copyright (c) Pepper Lebeck-Jobe, 2024
// This file is licensed under the MIT License.
// See the LICENSE.md file in the project root.

// Package terminal helps developers determine various attributes of the current
// terminal emulator in which they are running.
package terminal

import (
	"os"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/term"
)

type supportedColors int

const (
	// For when the terminal doesn't support colors.
	noColor supportedColors = 0
	// Sometimes called 16 color support
	color4Bit supportedColors = 4
	// Sometimes called 256 color support
	color8Bit supportedColors = 8
	// Sometimes called 16m (for million) or Truecolor support
	color24Bit supportedColors = 24
)

var (
	pattern256   = regexp.MustCompile(`-256(color)?$`)
	patternBasic = regexp.MustCompile(`^screen|^xterm|^vt100|color|ansi|cygwin|linux`)
)

// Colors has methods for number of colors supported by the current terminal
// emulator.
type Colors struct {
	value supportedColors
}

// Has16 returns true if the terminal supports 4-bit color.
func (t *Colors) Has16() bool {
	return t.value >= color4Bit
}

// Has256 returns true if the terminal supports 8-bit color.
func (t *Colors) Has256() bool {
	return t.value >= color8Bit
}

// HasTrueColor returns true if the terminal supports 24-bit color.
func (t *Colors) HasTrueColor() bool {
	return t.value >= color24Bit
}

// HasMillions is an alias for HasTrueColor.
func (t *Colors) HasMillions() bool {
	return t.HasTrueColor()
}

func NewColors() *Colors {
	return newColors(term.IsTerminal)
}

type isTerminalFunc func(fd int) bool

// NewColors returns a Colors struct for the current terminal emulator.
func newColors(isTerm isTerminalFunc) *Colors {
	// Highest priority is whether or not Stdout is a TTY
	if !isTerm(int(os.Stdout.Fd())) {
		return &Colors{noColor}
	}

	// Then, we look for supported Environment variables
	fromEnvVars := termColorFromEnvVars()
	if fromEnvVars != nil {
		return fromEnvVars
	}

	// Try to guess based on the TERM_PROGRAM variables?
	t, tpSet := os.LookupEnv("TERM_PROGRAM")
	if tpSet {
		switch t {
		case "iTerm.app":
			if termProgramVersion() >= 3 {
				return &Colors{color24Bit}
			} else {
				return &Colors{color8Bit}
			}
		case "ghostty":
			return &Colors{color24Bit}
		case "Hyper":
			return &Colors{color24Bit}
		case "Apple_Terminal":
			return &Colors{color8Bit}
		}
	}

	// Maybe the TERM variable can tell us more?
	term, tSet := os.LookupEnv("TERM")
	if tSet {
		if pattern256.MatchString(term) {
			return &Colors{color8Bit}
		}
		if patternBasic.MatchString(term) {
			return &Colors{color4Bit}
		}
		if term == "dumb" {
			return &Colors{noColor}
		}
	}

	// If none of that worked, let's assume no color support
	return &Colors{noColor}
}

func termProgramVersion() int {
	v, isSet := os.LookupEnv("TERM_PROGRAM_VERSION")
	if !isSet {
		return -1
	}
	version, _ := strconv.Atoi(strings.Split(v, ".")[0])
	return version
}

func termColorFromEnvVars() *Colors {
	termColors := fromEnv(getFirstEnvVar([]string{"COLORTERM", "LC_COLORTERM"}))
	userTermColors := fromEnv(getFirstEnvVar([]string{"USER_COLORTERM", "LC_USER_COLORTERM"}))
	if termColors == nil {
		return userTermColors
	}
	if userTermColors == nil {
		return termColors
	}
	return min(userTermColors, termColors)
}

func getFirstEnvVar(envVars []string) string {
	v := ""
	vSet := false
	for _, envVar := range envVars {
		v, vSet = os.LookupEnv(envVar)
		if vSet {
			break
		}
	}
	return v
}

func fromEnv(tc string) *Colors {
	switch tc {
	case "none":
		return &Colors{noColor}
	case "basic", "4bit":
		return &Colors{color4Bit}
	case "256", "8bit":
		return &Colors{color8Bit}
	case "16m", "Truecolor", "truecolor", "24bit":
		return &Colors{color24Bit}
	}
	// If it was set to something else, then it cannot be determined
	return nil
}

func min(a, b *Colors) *Colors {
	if a.value > b.value {
		return b
	}
	return a
}
