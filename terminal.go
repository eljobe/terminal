/*
 * Copyright 2017 Pepper Lebeck-Jobe
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package terminal

import (
    "os"
    "regexp"
    "strconv"
    "strings"

    goterm "golang.org/x/crypto/ssh/terminal"
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
    pattern256 = regexp.MustCompile(`-256(color)?$`)
    patternBasic = regexp.MustCompile(`^screen|^xterm|^vt100|color|ansi|cygwin|linux`)
)

type TerminalColor struct {
    value supportedColors
}

func (t *TerminalColor) Supports16Colors() bool {
    return t.value >= color4Bit
}

func (t *TerminalColor) Supports256Colors() bool {
    return t.value >= color8Bit
}

func (t *TerminalColor) SupportsTrueColor() bool {
    return t.value >= color24Bit
}

func Supports() *TerminalColor {
    // Highest priority is whether or Stdout is a TTY
    if !goterm.IsTerminal(int(os.Stdout.Fd())) {
        return &TerminalColor{noColor}
    }

    // Then, we look for supported Environment variables
    termColors := os.Getenv("TERM_COLORS")
    lcTermColors := os.Getenv("LC_TERM_COLORS")
    userTermColors := os.Getenv("USER_TERM_COLORS")
    if termColors != "" {
        if userTermColors != "" {
            return min(fromEnv(userTermColors), fromEnv(termColors))
        }
        return fromEnv(termColors)
    }
    if lcTermColors != "" {
        if userTermColors != "" {
            return min(fromEnv(userTermColors), fromEnv(lcTermColors))
        }
        return fromEnv(lcTermColors)
    }
    if userTermColors != "" {
        return fromEnv(userTermColors)
    }

    // Try to guess based on the TERM_PROGRAM variables?
    t := os.Getenv("TERM_PROGRAM")
    if t != "" {
        v := os.Getenv("TERM_PROGRAM_VERSION")
        version, _ := strconv.Atoi(strings.Split(v, ".")[0])
        switch t {
        case "iTerm.app":
            if version >= 3 {
                return &TerminalColor{color24Bit}
            } else {
                return &TerminalColor{color8Bit}
            }
        case "Hyper":
            return &TerminalColor{color24Bit}
        case "Apple_Terminal":
            return &TerminalColor{color8Bit}
        }
    }

    // Maybe the TERM variable can tell us more?
    term := os.Getenv("TERM")
    if pattern256.MatchString(term) {
        return &TerminalColor{color8Bit}
    }
    if patternBasic.MatchString(term) {
        return &TerminalColor{color4Bit}
    }
    if term == "dumb" {
        return &TerminalColor{noColor}
    }

    // If none of that worked, let's assume no color support
    return &TerminalColor{noColor}
}

func min(a, b *TerminalColor) *TerminalColor {
    if a.value > b.value {
        return b
    }
    return a
}

func fromEnv(tc string) *TerminalColor {
    switch tc {
    case "none":
        return &TerminalColor{noColor}
    case "basic", "4bit":
        return &TerminalColor{color4Bit}
    case "256", "8bit":
        return &TerminalColor{color8Bit}
    case "16m", "Truecolor", "24bit":
        return &TerminalColor{color24Bit}
    }
    // If it was set to something else, consider it "none"
    return &TerminalColor{noColor}
}
