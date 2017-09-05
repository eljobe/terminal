package main

import (
    "fmt"
    "github.com/eljobe/terminal"
)

func main() {
    cs := terminal.Supports()
    if !cs.Supports16Colors() {
        fmt.Println("No Colors Supported")
    } else {
        fmt.Println("16 Colors Supported")
    }
    if cs.Supports256Colors() {
        fmt.Println("256 Colors Supported")
    }
    if cs.SupportsTrueColor() {
        fmt.Println("Truecolor Supported")
    }
}
