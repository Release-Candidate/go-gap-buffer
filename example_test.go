// SPDX-FileCopyrightText:  Copyright 2024 Roland Csaszar
// SPDX-License-Identifier: MIT
//
// Project:  go-gap-buffer
// File:     example_test.go
// Date:     07.Feb.2024
//
// =============================================================================

package gapbuffer_test

import (
	"fmt"

	gap "github.com/Release-Candidate/go-gap-buffer"
)

func ExampleNew() {
	// Create a new, empty gap buffer.
	gapBuffer := gap.New()

	// Print the content of the gap buffer as a single string.
	fmt.Println(gapBuffer.String())
	// Output:
}

func ExampleNewCap() {
	// Create a new, empty gap buffer with a capacity of 10 bytes.
	gapBuffer := gap.NewCap(10)

	// Print the content size of the gap buffer in bytes.
	fmt.Println(gapBuffer.Size())
	// Output: 10
}

func ExampleNewStr() {
	// Create a new gap buffer containing "Hello, World!".
	gapBuffer := gap.NewStr("Hello, World!")

	// Print the content of the gap buffer as a single string.
	fmt.Println(gapBuffer.String())
	// Output: Hello, World!
}

func ExampleNewStrCap() {
	// Create a new gap buffer containing "Hello, World!" with a capacity of 10
	// bytes.
	// But "Hello, World!" is 13 bytes long, so the gap buffer's size will be
	// grown.
	gapBuffer := gap.NewStrCap("Hello, World!", 10)

	// Print the content size of the gap buffer in bytes.
	fmt.Println(gapBuffer.Size())
	// Output: 26
}

func ExampleGapBuffer_Col() {
	// Create a new gap buffer containing the two lines "Hello, World!" and
	// "My name is 阿保昭則."
	gapBuffer := gap.NewStr("Hello, World!\nMy name is 阿保昭則.")

	// Get the current column position of the cursor in bytes.
	// Numbering starts at column 1!
	col := gapBuffer.Col()

	// Print the current column position of the cursor in bytes, which is 24.
	// The string "My name is 阿保昭則." is 24 bytes long, but contains 16 unicode
	// runes.
	fmt.Println(col)
	// Output: 24
}

func ExampleGapBuffer_RuneCol() {
	// Create a new gap buffer containing the two lines "Hello, World!" and
	// "My name is 阿保昭則."
	gapBuffer := gap.NewStr("Hello, World!\nMy name is 阿保昭則.")

	// Get the current column position of the cursor in runes.
	// Numbering starts at column 1!
	runeCol := gapBuffer.RuneCol()

	// Print the current column position of the cursor in runes, which is 16.
	// The string "My name is 阿保昭則." is 24 bytes long, but contains 16 unicode
	// runes.
	fmt.Println(runeCol)
	// Output: 16
}

func ExampleGapBuffer_UpMv() {
	// Create a new gap buffer containing the two lines "Hello, World!" and
	// "My name is John."
	gapBuffer := gap.NewStr("Hello, World!\nMy name is John.")

	// Move the cursor up one line.
	gapBuffer.UpMv()

	// Get the part of the buffer before - to the left - and after - to the
	// right - of the current cursor position.
	l, r := gapBuffer.StringPair()

	// l should be "Hello, World!" and r should be "\nMy name is John."
	fmt.Printf("left: '%s' |cursor| right: '%s'\n", l, r)
	// Output: left: 'Hello, World!' |cursor| right: '
	// My name is John.'
}

func ExampleGapBuffer_DownMv() {
	// Create a new gap buffer containing the two lines "Hello, World!" and
	// "My name is John."
	gapBuffer := gap.NewStr("Hello, World!\nMy name is John.")

	// Move the cursor up one line.
	gapBuffer.UpMv()

	// And move down again, we are were we started.
	gapBuffer.DownMv()

	// Get the part of the buffer before - to the left - and after - to the
	// right - of the current cursor position.
	l, r := gapBuffer.StringPair()

	// l should be "Hello, World!\nMy name is John." and r should be empty.
	fmt.Printf("left: '%s' |cursor| right: '%s'\n", l, r)
	// Output: left: 'Hello, World!
	// My name is John.' |cursor| right: ''
}

func ExampleGapBuffer_Insert() {
	// Create a new gap buffer containing "Hello, World!".
	gapBuffer := gap.NewStr("Hello, World!")

	// Insert " My name is John." at the current position.
	gapBuffer.Insert(" My name is John.")

	// Print the content of the gap buffer as a single string.
	fmt.Println(gapBuffer.String())
	// Output: Hello, World! My name is John.
}

func ExampleGapBuffer_LeftDel() {
	// Create a new gap buffer containing the two lines "Hello, World!" and
	// "My name is John."
	gapBuffer := gap.NewStr("Hello, World!\nMy name is John.")

	// Delete 16 runes to the left of the cursor (like backspace).
	gapBuffer.LeftDel()
	gapBuffer.LeftDel()
	gapBuffer.LeftDel()
	gapBuffer.LeftDel()
	gapBuffer.LeftDel()
	gapBuffer.LeftDel()
	gapBuffer.LeftDel()
	gapBuffer.LeftDel()
	gapBuffer.LeftDel()
	gapBuffer.LeftDel()
	gapBuffer.LeftDel()
	gapBuffer.LeftDel()
	gapBuffer.LeftDel()
	gapBuffer.LeftDel()
	gapBuffer.LeftDel()
	gapBuffer.LeftDel()

	// Print the content of the gap buffer as a single string.
	fmt.Println(gapBuffer.String())
	// Output: Hello, World!
}

func ExampleGapBuffer_LeftMv() {
	// Create a new gap buffer containing the two lines "Hello, World!" and
	// "My name is John."
	gapBuffer := gap.NewStr("Hello, World!\nMy name is John.")

	// Move 16 runes to the left of the cursor.
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()

	// Get the part of the buffer before - to the left - and after - to the
	// right - of the current cursor position.
	l, r := gapBuffer.StringPair()

	// l should be "Hello, World!" and r should be '\nMy name is John.'.
	fmt.Printf("left: '%s' |cursor| right: '%s'\n", l, r)
	// Output: left: 'Hello, World!
	// ' |cursor| right: 'My name is John.'
}

func ExampleGapBuffer_Line() {
	// Create a new gap buffer containing the two lines "Hello, World!" and
	// "My name is John."
	gapBuffer := gap.NewStr("Hello, World!\nMy name is John.")

	// Get the current line position of the cursor.
	// Numbering starts at line 1!
	line := gapBuffer.Line()

	// Print the current line position of the cursor.
	fmt.Println(line)
	// Output: 2
}

func ExampleGapBuffer_LineCol() {
	// Create a new gap buffer containing the two lines "Hello, World!" and
	// "My name is 阿保昭則."
	gapBuffer := gap.NewStr("Hello, World!\nMy name is 阿保昭則.")

	// Get the current line and column position in bytes of the cursor.
	// Numbering starts at line 1 and column 1!
	line, col := gapBuffer.LineCol()

	// Print the current line position of the cursor.
	fmt.Printf("Line: %d column: %d", line, col)
	// Output: Line: 2 column: 24
}

func ExampleGapBuffer_LineRuneCol() {
	// Create a new gap buffer containing the two lines "Hello, World!" and
	// "My name is 阿保昭則."
	gapBuffer := gap.NewStr("Hello, World!\nMy name is 阿保昭則.")

	// Get the current column position of the cursor in runes.
	// Numbering starts at column 1!
	line, col := gapBuffer.LineRuneCol()

	// Print the current line position of the cursor.
	fmt.Printf("Line: %d column: %d", line, col)
	// Output: Line: 2 column: 16
}

func ExampleGapBuffer_LineLength() {
	// Create a new gap buffer containing "Hello, Wôrld!"
	gapBuffer := gap.NewStr("Hello, Wôrld!")

	// Get the length of the current line in bytes.
	length := gapBuffer.LineLength()

	// Print the length of the current line in bytes.
	// The string "Hello, Wôrld!" is 14 bytes long, but contains 13 unicode
	// runes.
	fmt.Println(length)
	// Output: 14
}

func ExampleGapBuffer_RightMv() {
	// Create a new gap buffer containing the two lines "Hello, World!" and
	// "My name is John."
	gapBuffer := gap.NewStr("Hello, World!\nMy name is John.")

	// Move 16 runes to the left of the cursor.
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()

	// Move 8 runes to the right.
	gapBuffer.RightMv()
	gapBuffer.RightMv()
	gapBuffer.RightMv()
	gapBuffer.RightMv()
	gapBuffer.RightMv()
	gapBuffer.RightMv()
	gapBuffer.RightMv()
	gapBuffer.RightMv()

	// Get the part of the buffer before - to the left - and after - to the
	// right - of the current cursor position.
	l, r := gapBuffer.StringPair()

	// l should be "Hello, World!\nMy name " and r should be 'is John.'.
	fmt.Printf("left: '%s' |cursor| right: '%s'\n", l, r)
	// Output: left: 'Hello, World!
	// My name ' |cursor| right: 'is John.'
}

func ExampleGapBuffer_RightDel() {
	// Create a new gap buffer containing the two lines "Hello, World!" and
	// "My name is John."
	gapBuffer := gap.NewStr("Hello, World!\nMy name is John.")

	// Move 5 runes to the left of the cursor.
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()

	// Delete 5 runes to the right, like the "delete" key does.
	gapBuffer.RightDel()
	gapBuffer.RightDel()
	gapBuffer.RightDel()
	gapBuffer.RightDel()
	gapBuffer.RightDel()

	// Get the part of the buffer before - to the left - and after - to the
	// right - of the current cursor position.
	l, r := gapBuffer.StringPair()

	// l should be "Hello, World!\nMy name is" and r should be empty.
	fmt.Printf("left: '%s' |cursor| right: '%s'\n", l, r)
	// Output: left: 'Hello, World!
	// My name is ' |cursor| right: ''
}

func ExampleGapBuffer_Size() {
	// Create a new gap buffer containing the two lines "Hello, World!" and
	// "My name is 阿保昭則."
	gapBuffer := gap.NewStr("Hello, World!\nMy name is 阿保昭則.")

	// Get the current size in bytes of the gap buffer.
	// This includes (unused bytes in) the gap
	size := gapBuffer.Size()

	// Print the current line position of the cursor.
	fmt.Println(size)
	// Output: 1024
}

func ExampleGapBuffer_String() {
	// Create a new gap buffer containing "Hello, World!".
	gapBuffer := gap.NewStr("Hello, World!")

	// Print the content of the gap buffer as a single string.
	fmt.Println(gapBuffer.String())
	// Output: Hello, World!
}

func ExampleGapBuffer_StringLength() {
	// Create a new gap buffer containing "Hello, World!".
	gapBuffer := gap.NewStr("Hello, World!")

	// Get the length of the content of the gap buffer in bytes.
	length := gapBuffer.StringLength()

	// Print the length of the content of the gap buffer in bytes.
	fmt.Println(length)
	// Output: 13
}

func ExampleGapBuffer_StringPair() {
	// Create a new gap buffer containing the two lines "Hello, World!" and
	// "My name is John."
	gapBuffer := gap.NewStr("Hello, World!\nMy name is John.")

	// Move 16 runes to the left of the cursor.
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()
	gapBuffer.LeftMv()

	// Get the part of the buffer before - to the left - and after - to the
	// right - of the current cursor position.
	l, r := gapBuffer.StringPair()

	// l should be "Hello, World!" and r should be '\nMy name is John.'.
	fmt.Printf("left: '%s' |cursor| right: '%s'\n", l, r)
	// Output: left: 'Hello, World!
	// ' |cursor| right: 'My name is John.'
}
