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
    // Sometimes called 16m (for million) color support
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
    if !goterm.IsTerminal(int(os.Stdout.Fd())) {
        return &TerminalColor{noColor}
    }

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

    return &TerminalColor{noColor}
}
