// SPDX-FileCopyrightText:  Copyright 2024 Roland Csaszar
// SPDX-License-Identifier: MIT
//
// Project:  go-gap-buffer
// File:     gap-buffer.go
// Date:     07.Feb.2024
//
// =============================================================================

package gapbuffer

import (
	"strings"
	"unicode/utf8"
)

// The index in the current line is `start - newLines.start`.
// `wantsCol` is the rune column (not byte column!) the cursor wants to hold when going up or down.
type GapBuffer struct {
	start    int
	end      int
	wantsCol int
	data     []byte
	lines    lineBuffer
}

const (
	defaultCapacity = 1024
	minGap          = 16
	lineCapFactor   = 10 // 1/10th of the capacity of the gap buffer, default: 102
	minLineCap      = 10
	growFactor      = 2
	minLineGap      = 4
)

func (g *GapBuffer) String() string {
	var b strings.Builder
	b.Grow(g.start + len(g.data) - g.end)
	b.Write(g.data[:g.start])
	b.Write(g.data[g.end:])

	return b.String()
}

func (g *GapBuffer) StringPair() (left string, right string) {
	return string(g.data[:g.start]), string(g.data[g.end:])
}

func NewCap(size int) *GapBuffer {
	return &GapBuffer{
		start:    0,
		end:      size,
		wantsCol: 0,
		data:     make([]byte, size),
		lines:    *newLineBuf(size),
	}
}

func New() *GapBuffer {
	return NewCap(defaultCapacity)
}

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

func NewStr(s string) *GapBuffer {
	return NewStrCap(s, defaultCapacity)
}

func (g *GapBuffer) Col() int {
	if g.start < g.lines.curLineStart() {
		return 0
	}

	return g.start - g.lines.curLineStart()
}

func (g *GapBuffer) RuneCol() int {
	if g.start < g.lines.curLineStart() {
		return 0
	}

	return utf8.RuneCount(g.data[g.lines.curLineStart():g.start])
}

func (g *GapBuffer) LineLength() int {
	return g.lines.curLineLength()
}

func (g *GapBuffer) Line() int {
	return g.lines.curLine()
}

func (g *GapBuffer) LineCol() (line int, col int) {
	return g.lines.curLine(), g.Col() + 1
}

func (g *GapBuffer) LineRuneCol() (line int, runeCol int) {
	return g.lines.curLine(), g.RuneCol() + 1
}

func (g *GapBuffer) LeftDel() {
	if g.start < 1 {
		return
	}

	r, d := utf8.DecodeLastRune(g.data[:g.start])
	g.start -= d

	if r == '\n' {
		g.lines.upDel()
	} else {
		g.lines.del()
	}

	g.wantsCol = g.RuneCol()
}

func (g *GapBuffer) RightDel() {
	if g.end > len(g.data)-1 {
		return
	}

	r, d := utf8.DecodeRune(g.data[g.end:])
	g.end += d

	if r == '\n' {
		g.lines.downDel()
	} else {
		g.lines.del()
	}
}

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

func (g *GapBuffer) DownMv() {
	if g.lines.end > g.lines.lastIdx() {
		return
	}

	newLine := g.lines.curLineEnd() + 1 - g.start
	if g.lines.curLineLength() == 0 || g.lines.lastLine() {
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

func (g *GapBuffer) Insert(str string) {
	if g.end-g.start < len(str)+1 {
		g.grow()
	}

	g.lines.insert(str, g.start)
	l := copy(g.data[g.start:], str)
	g.start += l
	g.wantsCol = g.RuneCol()
}
