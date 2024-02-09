// SPDX-FileCopyrightText:  Copyright 2024 Roland Csaszar
// SPDX-License-Identifier: MIT
//
// Project:  go-gap-buffer
// File:     main.go
// Date:     07.Feb.2024
//
// =============================================================================

package main

import (
	"fmt"

	gap "github.com/Release-Candidate/go-gap-buffer"
)

func main() { //nolint:funlen // Yes, it is long.
	// Create a new, empty gap buffer.
	gapBuffer := gap.New()
	// Insert "Hello, World!" at the start of the buffer.
	gapBuffer.Insert("Hello, World!")
	// Print the content of the gap buffer as a single string.
	fmt.Println(gapBuffer.String())
	fmt.Println("================================================================================")

	// This does the same in a single step.
	gapBuffer = gap.NewStr("Hello, World!")
	fmt.Println(gapBuffer.String())
	fmt.Println("================================================================================")

	// Move 6 Unicode runes to the left, before "World!"
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	// We can also get the content of the gap buffer as a pair of strings, one
	// to the left of the "cursor" and one to the right.
	left, right := gapBuffer.StringPair()
	fmt.Printf("left: '%s' |cursor| right: '%s'\n", left, right)
	fmt.Printf("%s<|>%s\n", left, right)
	fmt.Println("================================================================================")

	// From now on, "<|>" marks the current "cursor" position in the output.

	// Insert a Unicode smiley.
	gapBuffer.Insert("🙂")
	left, right = gapBuffer.StringPair()
	fmt.Printf("%s<|>%s\n", left, right)
	fmt.Println("================================================================================")

	// Delete the Unicode smiley with a single `backspace` (delete the Unicode
	// Rune to the left of the cursor).
	gapBuffer.LeftDel()
	left, right = gapBuffer.StringPair()
	fmt.Printf("%s<|>%s\n", left, right)
	fmt.Println("================================================================================")

	// Insert the string "funny" in a line on its own at the current cursor
	// location.
	gapBuffer.Insert("\nfunny\n")
	left, right = gapBuffer.StringPair()
	fmt.Printf("%s<|>%s\n", left, right)
	fmt.Println("================================================================================")

	// Move the cursor up two lines.
	gapBuffer.UpMv()
	gapBuffer.UpMv()
	left, right = gapBuffer.StringPair()
	fmt.Printf("%s<|>%s\n", left, right)
	fmt.Println("================================================================================")

	// Two runes to the right and down two lines again.
	gapBuffer.RightMv()
	gapBuffer.RightMv()
	gapBuffer.DownMv()
	gapBuffer.DownMv()
	left, right = gapBuffer.StringPair()
	fmt.Printf("%s<|>%s\n", left, right)
}
