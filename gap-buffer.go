// SPDX-FileCopyrightText:  Copyright 2024 Roland Csaszar
// SPDX-License-Identifier: MIT
//
// Project:  go-gap-buffer
// File:     gap-buffer.go
// Date:     07.Feb.2024
//
// =============================================================================

// This library implements a gap buffer, which is a data structure to be used as
// the container of the text for a (simple or not so simple) text editor.
// A gap buffer is not ideal for using multiple cursors, as that would involve
// multiple jumps and copying of data in the gap buffer.
//
// This gap buffer includes line movements (up and down a line from the current
// one) but it splits lines based on the newline character '\n'. So
// Windows-style CR LF (`\r\n`) line endings are not supported.
//
// A gap buffer is an array with a gap at the cursor position, where text is to
// be inserted and deleted.
//
// The string "Hello world!" with the cursor at the end of "Hello" -
// "Hello| world!" - looks like this in a gap buffer array:
//
//	Hello|< gap start, the cursor position            gap end >| world!
//
//	['H', 'e', 'l', 'l', 'o', 0, 0, 0, 0, 0, ' ', 'w', 'o', 'r', 'l', 'd', '!']
//	  0    1    2    3    4  |     gap     |  5    6    7    8    9    10   11
//
// Movement in the gap buffer works by moving the start and end of the gap, same
// with deletion of unicode runes in both directions.
//
// Moving the cursor two runes to the left:
//
//	Hel|< gap start, the cursor position            gap end >|lo world!
//
//	['H', 'e', 'l', 0, 0, 0, 0, 0, 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd', '!']
//	  0    1    2   |     gap    |  3    4    5    6    7    8    9    10   11
//
// Deleting three runes to the left:
//
//	|< gap start, the cursor position            gap end >|lo world!
//
//	['H', 'e', 'l', 0, 0, 0, 0, 0, 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd', '!']
//	  |           gap            |  1    2    3    4    5    6    7    8    9
//
// Insertion happens at the cursor position by appending at the start of the gap
// and moving the start of the gap accordingly.
//
// New|< gap start, the cursor position            gap end >|lo world!
//
//	['N', 'e', 'w', 0, 0, 0, 0, 0, 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd', '!']
//	  0    1    2   |   gap      |  3    4    5    6    7    8    9    10   11
package gapbuffer

import (
	"strings"
	"unicode/utf8"
)

// GapBuffer represents a gap buffer.
type GapBuffer struct {
	// The index in the gap buffer `GapBuffer.data` of the start of the gap.
	// The position of the cursor.
	start int

	// The index in the gap buffer `GapBuffer.data` of the end of the gap.
	// The position of the first unicode scalar point after the cursor.
	end int

	// `wantsCol` is the rune column (not byte column!) the cursor wants to hold
	// when going up or down.
	wantsCol int

	// The lineBuffer that stores the line length information of the gap buffer.
	//
	// See [lineBuffer].
	lines lineBuffer

	// The data of the gap buffer.
	data []byte
}

const (
	defaultCapacity = 1024 // The default size of a gap buffer in bytes.

	// 1/10th of the capacity of the gap buffer, default: 102.
	lineCapFactor = 10

	// Minimum size in int of the line buffer `GapBuffer.lines`. A lineBuffer
	// has at least this size, even if the [lineCapFactor] would yield a smaller
	// one.
	minLineCap = 10

	// The factor by which to grow the gap buffer and line buffer, if needed.
	growFactor = 2
)

// Return the contents of the gap buffer as a string.
func (g *GapBuffer) String() string {
	var b strings.Builder
	b.Grow(len(g.data) - (g.end - g.start))
	b.Write(g.data[:g.start])
	b.Write(g.data[g.end:])

	return b.String()
}

// Return the contents of the gap buffer as two strings. The part to the left of
// the cursor is returned in `left` and the part to the right of the cursor is
// returned in `right`.
func (g *GapBuffer) StringPair() (left string, right string) {
	return string(g.data[:g.start]), string(g.data[g.end:])
}

// Return the length in bytes of the contents of the gap buffer.
func (g *GapBuffer) StringLength() int {
	return len(g.data) - (g.end - g.start)
}

// Construct a new GapBuffer from a capacity. The capacity is the number of
// bytes the gap buffer can hold without a resize.
//
// The default size is 1024 bytes, if you know that you need less or more space,
// you can set the initial size to something more appropriate.
//
// See also [New], [NewStr], [NewStrCap].
func NewCap(size int) *GapBuffer {
	return &GapBuffer{
		start:    0,
		end:      size,
		wantsCol: 0,
		data:     make([]byte, size),
		lines:    *newLineBuf(size),
	}
}

// Construct a new, empty GapBuffer with the default capacity.
//
// See also [NewCap], [NewStr], [NewStrCap].
func New() *GapBuffer {
	return NewCap(defaultCapacity)
}

// Construct a new GapBuffer from a string and a capacity. The cursor position
// is set to the end of the string. The capacity is the number of bytes the gap
// buffer can hold without a resize.
//
// The default size is 1024 bytes, if you know that you need less or more space,
// you can set the initial size to something more appropriate.
//
// See also [New], [NewCap], [NewStr], [GapBuffer.Size].
func NewStrCap(s string, c int) *GapBuffer {
	size := max(c, len(s)*growFactor)
	dat := make([]byte, size)
	sIdx := copy(dat, s)
	lines := newLineBufStr(s, size)
	runeCol := 0
	lineStart := lines.curLineStart()

	if lineStart < sIdx {
		runeCol = utf8.RuneCount(dat[lineStart:sIdx])
	}

	return &GapBuffer{
		start:    sIdx,
		end:      size,
		wantsCol: runeCol,
		data:     dat,
		lines:    *lines,
	}
}

// Construct a new GapBuffer from a string. The cursor position is set to the
// end of the string.
//
// See also [New], [NewCap], [NewStrCap].
func NewStr(s string) *GapBuffer {
	return NewStrCap(s, defaultCapacity)
}

// Return the current number of bytes in the buffer, including the "empty" space
// in the gap.
func (g *GapBuffer) Size() int {
	return len(g.data)
}

// Return the byte column of the cursor, the number of bytes from the start of
// the line to the cursor.
//
// Numbering starts from 1.
//
// See also [GapBuffer.RuneCol], [GapBuffer.LineCol], [GapBuffer.LineRuneCol].
func (g *GapBuffer) Col() int {
	if g.start < g.lines.curLineStart() {
		return 0
	}

	return g.start - g.lines.curLineStart()
}

// Return the rune column of the cursor, the number of unicode runes from the
// start of the line to the cursor.
//
// Numbering starts from 1.
//
// See also [GapBuffer.Col], [GapBuffer.LineCol], [GapBuffer.LineRuneCol].
func (g *GapBuffer) RuneCol() int {
	if g.start < g.lines.curLineStart() {
		return 0
	}

	return utf8.RuneCount(g.data[g.lines.curLineStart():g.start])
}

// Return the length of the current line the cursor is in in bytes.
// This returns the "whole" line length, including the part to the right of the
// cursor.
//
// See also [GapBuffer.Col], [GapBuffer.RuneCol], which is the the length to the
// left of the cursor.
func (g *GapBuffer) LineLength() int {
	return g.lines.curLineLength() + 1
}

// Return the line number of the current line the cursor is in.
//
// Numbering starts from 1.
//
// See also [GapBuffer.Col], [GapBuffer.RuneCol], [GapBuffer.LineRuneCol].
func (g *GapBuffer) Line() int {
	return g.lines.curLine()
}

// Return the line and byte column of the cursor. Byte column means the number
// of bytes from the start of the line to the cursor.
//
// Numbering starts from 1 for both the line number and the column number.
//
// See also [GapBuffer.Col], [GapBuffer.RuneCol], [GapBuffer.LineRuneCol],
// [GapBuffer.Line].
func (g *GapBuffer) LineCol() (line int, col int) {
	return g.lines.curLine(), g.Col()
}

// Return the line and rune column of the cursor. Rune column means the number
// of unicode runes from the start of the line to the cursor.
//
// Numbering starts from 1 for both the line number and the column number.
//
// See also [GapBuffer.Line], [GapBuffer.Col], [GapBuffer.RuneCol], [GapBuffer.LineCol].
func (g *GapBuffer) LineRuneCol() (line int, runeCol int) {
	return g.lines.curLine(), g.RuneCol()
}

// Delete the unicode rune to the left of the cursor. Like the "backspace" key.
//
// See also [GapBuffer.RightDel], [GapBuffer.LeftMv], [GapBuffer.RightMv],
// [GapBuffer.UpMv], [GapBuffer.DownMv]
func (g *GapBuffer) LeftDel() {
	if g.start < 1 {
		return
	}

	r, d := utf8.DecodeLastRune(g.data[:g.start])
	g.start -= d

	if r == '\n' {
		g.lines.upDel()
	} else {
		g.lines.del(d)
	}

	g.wantsCol = g.RuneCol()
}

// Delete the unicode rune to the right of the cursor. Like the "delete" key.
//
// See also [GapBuffer.LeftDel], [GapBuffer.RightMv], [GapBuffer.LeftMv],
// [GapBuffer.UpMv], [GapBuffer.DownMv]
func (g *GapBuffer) RightDel() {
	if g.end > len(g.data)-1 {
		return
	}

	r, d := utf8.DecodeRune(g.data[g.end:])
	g.end += d

	if r == '\n' {
		g.lines.downDel()
	} else {
		g.lines.del(d)
	}
}

// Move the cursor one unicode rune to the left.
//
// See also [GapBuffer.RightMv], [GapBuffer.LeftDel], [GapBuffer.RightDel],
// [GapBuffer.UpMv], [GapBuffer.DownMv]
func (g *GapBuffer) LeftMv() {
	if g.start < 1 {
		return
	}

	rChar, d := utf8.DecodeLastRune(g.data[:g.start])
	g.end -= d

	_ = copy(g.data[g.end:], g.data[g.start-d:g.start])
	g.start -= d

	if rChar == '\n' {
		g.lines.up()
	}

	g.wantsCol = g.RuneCol()
}

// Move the cursor one unicode rune to the right.
//
// See also [GapBuffer.LeftMv], [GapBuffer.LeftDel], [GapBuffer.RightDel],
// [GapBuffer.UpMv], [GapBuffer.DownMv]
func (g *GapBuffer) RightMv() {
	if g.start > len(g.data)-2 {
		return
	}

	if g.end > len(g.data)-1 {
		return
	}

	r, d := utf8.DecodeRune(g.data[g.end:])
	_ = copy(g.data[g.start:], g.data[g.end:g.end+d])
	g.start += d
	g.end += d

	if r == '\n' {
		g.lines.down()
	}

	g.wantsCol = g.RuneCol()
}

// Move the cursor up one line.
//
// The cursor "tries" to hold the current position in the new line, like we are
// used to in text editors.
//
// Before the moves:
//
//	Some text
//	No
//	More |text
//
// After the first move:
//
//	Some text
//	No|
//	More text
//
// After the second move:
//
//	Some |text
//	No
//	More text
//
// See also [GapBuffer.DownMv], [GapBuffer.LeftMv], [GapBuffer.RightMv],
// [GapBuffer.LeftDel], [GapBuffer.RightDel]
func (g *GapBuffer) UpMv() {
	if g.lines.curLine() == 1 {
		return
	}

	g.lines.up()
	lineStart := g.lines.curLineStart()
	newStart := lineStart
	max := g.lines.curLineEnd()
	runeCnt := 0

	for idx := lineStart; idx < max+1; {
		newStart = idx

		if runeCnt == g.wantsCol {
			break
		}

		_, d := utf8.DecodeRune(g.data[idx:])
		idx += d
		runeCnt++
	}

	g.end -= (g.start - newStart)
	_ = copy(g.data[g.end:], g.data[newStart:g.start])
	g.start = newStart
}

// Move the cursor down one line.
//
// The cursor "tries" to hold the current position in the new line, like we are
// used to in text editors.
//
// Before the moves:
//
//	Some |text
//	No
//	More text
//
// After the first move:
//
//	Some text
//	No|
//	More text
//
// After the second move:
//
//	Some text
//	No
//	More |text
//
// See also [GapBuffer.UpMv], [GapBuffer.LeftMv], [GapBuffer.RightMv],
// [GapBuffer.LeftDel], [GapBuffer.RightDel]
func (g *GapBuffer) DownMv() {
	if g.lines.end > g.lines.lastIdx() {
		return
	}

	newLine := g.lines.curLineEnd() + 1 - g.start
	if g.lines.curLineLength() == 0 || g.lines.isLastLine() {
		newLine--
	}

	idx := newLine
	runeCnt := 0

	for g.end+idx < len(g.data) && g.data[g.end+idx] != '\n' {
		if runeCnt == g.wantsCol {
			break
		}

		_, d := utf8.DecodeRune(g.data[g.end+idx:])
		if g.end+idx+d > len(g.data)-1 {
			break
		}

		idx += d
		runeCnt++
	}

	// runtime error: slice bounds out of range [1014:1013]
	_ = copy(g.data[g.start:], g.data[g.end:g.end+idx])
	g.start += idx
	g.end += idx

	if idx == 0 {
		g.start++
		g.end++
	}
	// if we are at the last line, we can move one step further to the right
	if g.end == len(g.data)-1 && runeCnt < g.wantsCol {
		g.data[g.start] = g.data[g.end]
		g.start++
		g.end++
	}

	g.lines.down()
}

// grow resizes the gap buffer by `growFactor` times its current size and copies
// the existing data.
func (g *GapBuffer) grow() {
	tmp := make([]byte, len(g.data)*growFactor)
	_ = copy(tmp, g.data[:g.start])
	nE := len(tmp) - (len(g.data) - g.end)
	_ = copy(tmp[nE:], g.data[g.end:])
	g.end = nE
	g.data = tmp
}

// Insert inserts the given string at the current cursor position.
// The string can be a single unicode scalar point or text of arbitrary size and
// anything in between (like a single unicode rune).
//
// The cursor is moved to the end of the inserted text.
func (g *GapBuffer) Insert(str string) {
	if g.end-g.start < len(str)+1 {
		g.grow()
	}

	g.lines.insert(str, g.start)
	l := copy(g.data[g.start:], str)
	g.start += l
	g.wantsCol = g.RuneCol()
}
