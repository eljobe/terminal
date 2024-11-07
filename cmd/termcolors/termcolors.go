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
package main

import (
    "fmt"
    "github.com/eljobe/terminal"
)

func main() {
    cs := terminal.ColorSupport()
    if !cs.Supports16Colors() {
        fmt.Println("No Colors Supported")
    } else {
        fmt.Println("16 Colors Supported")
    }
    print16Colors()

    if cs.Supports256Colors() {
        fmt.Println("256 Colors Supported")
    } else {
        fmt.Println("256 Colors NOT Supported")
    }
    print256Colors()

    if cs.SupportsTrueColor() {
        fmt.Println("Truecolor Supported")
    } else {
        fmt.Println("Truecolor NOT Supported")
    }
    printTruecolors()
}

func printColors(rows, columns int, printNextColor func()) {
    for y :=0; y < rows; y++ {
        for x := 0; x < columns; x++ {
            printNextColor()
        }
        fmt.Print("\x1b[0m\n")
    }
    fmt.Println("\x1b[0m")
}

func print16Colors() {
    printColors(2, 8, printNext4bit())
}

func print256Colors() {
    printColors(8, 32, printNext8bit())
}

func printTruecolors() {
    printColors(32, 48, printNext24bit())
}

func printNext4bit() func() {
    next := gen4bit()
    return func() {
        fg := next()
        fmt.Printf("\x1b[%v;7m    ", fg)
    }
}

func printNext8bit() func() {
    next := gen8bit()
    return func() {
        fmt.Printf("\x1b[48;5;%vm ", next())
    }
}

func printNext24bit() func() {
    next := gen24bit()
    return func() {
        r, g, b := next()
        fmt.Printf("\x1b[48;2;%v;%v;%vm ", r, g, b)
    }
}

func gen4bit() func() uint8 {
    next := uint8(97)
    return func() uint8 {
        if next == 37 {
            next = 90
        } else if next == 97 {
            next = 30
        } else {
            next += 1
        }
        return next
    }
}

func gen8bit() func() uint8 {
    next := uint8(255)
    return func() uint8 {
        if next == 255 {
            next = 0
        } else {
            next += 1
        }
        return next
    }
}

func gen24bit() func() (uint8, uint8, uint8) {
    colors := [3]uint8{255, 0, 1}
    upColor := 1
    downColor := 2
    increasing := false
    return func() (uint8, uint8, uint8) {
        if increasing {
            colors[upColor] += 1
            if colors[upColor] == 255 {
                upColor = (upColor + 1) % 3
                increasing = false
            }
        } else {
            colors[downColor] -= 1
            if colors[downColor] == 0 {
                downColor = (downColor + 1) % 3
                increasing = true
            }
        }
        return colors[0], colors[1], colors[2]
    }
}
