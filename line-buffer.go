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

type lineBuffer struct {
	start   int
	end     int
	lengths []int
}

func newLineBuf(c int) *lineBuffer {
	cR := max(c/lineCapFactor, defaultLineCap)
	lb := &lineBuffer{start: 0, end: cR, lengths: make([]int, cR)}

	return lb
}

func newLineBufStr(s string, c int) *lineBuffer {
	l := newLineBuf(c)
	l.insert(s, 0)

	return l
}

func (l *lineBuffer) size() int {
	return len(l.lengths)
}

func (l *lineBuffer) lastIdx() int {
	return len(l.lengths) - 1
}

func (l *lineBuffer) up() {
	if l.start < 1 {
		return
	}
	l.end--
	l.lengths[l.end] = l.lengths[l.start]
	l.start--
}

func (l *lineBuffer) upDel() {
	if l.start < 1 {
		return
	}
	l.start--
}

func (l *lineBuffer) down() {
	if l.end > l.lastIdx() {
		return
	}
	l.start++
	l.lengths[l.start] = l.lengths[l.end]
	l.end++
}

func (l *lineBuffer) downDel() {
	if l.end > l.lastIdx() {
		return
	}
	l.end++
}

func (l *lineBuffer) del() {
	if l.lengths[l.start] == 0 {
		return
	}

	l.lengths[l.start]--
}

func newlineSplit(str string) []string {
	line, rest, _ := strings.Cut(str, "\n")
	lines := make([]string, 0)
	lines = append(lines, line)

	for rest != "" {
		line, rest, _ = strings.Cut(rest, "\n")
		lines = append(lines, line)
	}

	if rest == "" && line == "" {
		lines = append(lines, "")
	}

	return lines
}

func lineLengths(str string) []int {
	lines := newlineSplit(str)

	lens := make([]int, 0, len(lines))

	for i := range lines {
		lens = append(lens, len(lines[i])+1)
	}
	if strings.HasSuffix(str, "\n") {
		lens = append(lens, 0)
	}

	return lens
}

// grow resizes the line buffer by `growFactor` times its current size and
// copies the existing data.
func (l *lineBuffer) grow() {
	tmp := make([]int, growFactor*l.size())
	_ = copy(tmp, l.lengths[:l.start])
	nE := len(tmp) - (l.size() - l.end)
	_ = copy(tmp[nE:], l.lengths[l.end:])
	l.end = nE
	l.lengths = tmp
}

func (l *lineBuffer) insert(str string, pos int) {
	strLen := len(str)

	if strLen == 0 {
		return
	}

	lens := lineLengths(str)
	relPos := pos - l.curLineStart()
	lens[0] += relPos

	lens[len(lens)-1] += l.curLineEnd() - pos

	if l.end-l.start < len(lens)+1 {
		l.grow()
	}

	for idx := range lens {
		l.lengths[l.start+idx] = lens[idx]
	}

	l.start += len(lens) - 1
}

func (l *lineBuffer) curLine() int {
	return l.start + 1
}

func (l *lineBuffer) curLineLength() int {
	return l.lengths[l.start]
}

func (l *lineBuffer) curLineEnd() int {
	sum := 0
	for i := range l.lengths[:l.start+1] {
		sum += l.lengths[i]
	}

	return sum - 1 // returns the current last character which is not a newline.
}

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

func (l *lineBuffer) lastLine() bool {
	return l.end == l.size()
}
