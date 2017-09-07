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
