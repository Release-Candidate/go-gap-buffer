// SPDX-FileCopyrightText:  Copyright 2024 Roland Csaszar
// SPDX-License-Identifier: MIT
//
// Project:  go-gap-buffer
// File:     line-buffer.go
// Date:     07.Feb.2024
//
// =============================================================================

package gapbuffer

import "strings"

// This is a gap buffer which holds the line lengths of the lines in `GapBuffer`.
//
// Warning: this data structure does not check it's function arguments, it is
// only meant to be called by [GapBuffer] if needed and the arguments are valid.
type lineBuffer struct {
	// The start of the gap, the index of the current line in `lengths`.
	//
	// Note: the string after `GapBuffer.end` may be part of the current line.
	//
	// See [lineBuffer.end]
	start int

	// The end of the gap, the index of the "next" line after the current one,
	// if such a line exists.
	//
	// Note:
	//
	// See [lineBuffer.start]
	end int

	// The array holding the lengths of the lines in the gap buffer. This
	// includes the new line character at the end of each line. Only the last
	// line does not have a new line character at the end and may have a length
	// of zero, every other line but the last has a length of at least one.
	lengths []int
}

// newLineBuf returns a new line buffer from the given capacity of the parent
// [GapBuffer]. The size is the maximum of the given capacity divided by
// [lineCapFactor] and [minLineCap].
func newLineBuf(c int) *lineBuffer {
	cR := max(c/lineCapFactor, minLineCap)
	lb := &lineBuffer{start: 0, end: cR, lengths: make([]int, cR)}

	return lb
}

// newLineBuf returns a new line buffer with the line lengths of the given
// string and the the given capacity of the parent [GapBuffer]. The size is the
// maximum of the given capacity divided by [lineCapFactor] and [minLineCap].
func newLineBufStr(s string, c int) *lineBuffer {
	l := newLineBuf(c)
	l.insert(s, 0)

	return l
}

// insert inserts the given string at the current line. The absolute position in
// the gap buffer is given by the `pos` parameter and necessary to calculate the
// length of the string part after the inserted string, if such a substring
// exists.
//
//	\nfoo|< start   end >|bar\n
//
// current line length is 7 ("foobar\n")
//
//	\nfoo insert\n newline|< start  end >|bar\n
//
// current line length is 12 = 3 + 9 ("foo insert\n"), next line length is
// 13 = 4 + 9 (" newlinebar\n")
func (l *lineBuffer) insert(str string, pos int) {
	strLen := len(str)

	if strLen == 0 {
		return
	}

	lens := lineLengths(str)
	if l.end-l.start < len(lens)+1 {
		l.grow()
	}

	lens[0] += pos - l.curLineStart()
	lens[len(lens)-1] += l.curLineStart() + l.curLineLength() - pos

	for idx := range lens {
		l.lengths[l.start+idx] = lens[idx]
	}

	l.start += len(lens) - 1
}

// size returns the size of the lineBuffer in ints, including the "empty" space
// of the gap.
func (l *lineBuffer) size() int {
	return len(l.lengths)
}

// lastIdx returns the index of the element of the lineBuffer.
//
// Note: [lineBuffer.end] may have an index one larger than this, if there is
// nothing (no line) right of the gap.
func (l *lineBuffer) lastIdx() int {
	return len(l.lengths) - 1
}

// up reacts to a movement of the curser up one line.
//
// The line length of the current line is set as the length of [lineBuffer.end],
// and the gap is moved one step to the left.
//
// Warning: this function does not check if the cursor is in the first line, if
// it is, this panics!
func (l *lineBuffer) up() {
	l.end--
	l.lengths[l.end] = l.lengths[l.start]
	l.start--
}

// upDel reacts to the deletion of the newline before the cursor.
//
// The gap is widened one step to the left.
//
// Warning: this function does not check if the cursor is in the first line, if
// it is, this panics!
func (l *lineBuffer) upDel() {
	l.start--
}

// down reacts to a movement of the curser down one line.
//
// The line length of the next line - of [lineBuffer.end is set as the length of
// the current line, and the gap is moved one step to the right.
//
// Warning: this function does not check if the cursor is in the last line, if
// it is, this panics!
func (l *lineBuffer) down() {
	l.start++
	l.lengths[l.start] = l.lengths[l.end]
	l.end++
}

// downDel reacts to the deletion of the newline after the cursor.
//
// The gap is widened one step to the right.
//
// Warning: this function does not check if the cursor is in the last line, if
// it is, this panics!
func (l *lineBuffer) downDel() {
	l.end++
}

// del reacts to the deletion of a rune by shortening the line length by the
// number of bytes given. If the current line length already is zero, nothing
// happens.
func (l *lineBuffer) del(b int) {
	if l.lengths[l.start] == 0 {
		return
	}

	l.lengths[l.start] -= b
}

// newlineSplit splits the given string into lines, by splitting on newline
// characters '\n' and returning the substrings in a slice.
//
// If a substring between two newlines is empty, the empty string is added to
// the slice. The number of elements, (the length) of the returned slice is at
// least the number of newline characters in the given string.
//
// Example:
//
//	newlineSplit("\nfunny\n") == ["", "funny"]
func newlineSplit(str string) []string {
	line, rest, _ := strings.Cut(str, "\n")
	lines := make([]string, 0)
	lines = append(lines, line)

	for rest != "" {
		line, rest, _ = strings.Cut(rest, "\n")
		lines = append(lines, line)
	}

	return lines
}

// lineLengths returns the lengths of the lines in bytes in the given string in
// a slice.
//
// The newline character is included in the length of each line. If the string
// ends in a newline, 0 (zero) is returned as the length of the last line. So,
// every line length but the last is at least 1.
//
// Example:
//
//	lineLengths("\nfunny\n") == [1, 6, 0]
func lineLengths(str string) []int {
	lines := newlineSplit(str)

	lens := make([]int, 0, len(lines))

	for i := range lines {
		lens = append(lens, len(lines[i])+1)
	}
	if strings.HasSuffix(str, "\n") {
		lens = append(lens, 0)
	} else {
		lens[len(lens)-1]--
	}

	return lens
}

// grow resizes the line buffer by `growFactor` times its current size and
// copies the existing data.
func (l *lineBuffer) grow() {
	tmp := make([]int, growFactor*l.size())
	_ = copy(tmp, l.lengths[:l.start+1])
	nE := len(tmp) - (l.size() - l.end)
	_ = copy(tmp[nE:], l.lengths[l.end:])
	l.end = nE
	l.lengths = tmp
}

// curLine returns the number of the current line, starting from 1.
func (l *lineBuffer) curLine() int {
	return l.start + 1
}

// curLineLength returns the length of the current line, including the final
// newline character, if it isn't the last line.
func (l *lineBuffer) curLineLength() int {
	return l.lengths[l.start]
}

// isLastLine returns true if the cursor is in the last line.
func (l *lineBuffer) isLastLine() bool {
	return l.end == l.size()
}

// curLineStart returns the index in the gap buffer of the first character in
// the current line. This is the sum of all line length before the current line.
func (l *lineBuffer) curLineStart() int {
	if l.start == 0 {
		return 0
	}

	sum := 0
	for i := range l.lengths[:l.start] {
		sum += l.lengths[i]
	}

	return sum
}

// curLineEnd returns the index in the gap buffer of the last character in the
// current line, including the newline character. This is the sum of all
// line lengths before the current line and the length of the current line minus
// one to get the index instead of the length.
func (l *lineBuffer) curLineEnd() int {
	sum := 0
	for i := range l.lengths[:l.start+1] {
		sum += l.lengths[i]
	}

	// do not subtract from a zero length line.
	if l.curLineLength() == 0 {
		return sum
	}

	// returns the index of the last character, so subtract one.
	return sum - 1
}
