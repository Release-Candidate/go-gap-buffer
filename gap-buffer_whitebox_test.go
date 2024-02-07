// SPDX-FileCopyrightText:  Copyright 2024 Roland Csaszar
// SPDX-License-Identifier: MIT
//
// Project:  go-gap-buffer
// File:     gap-buffer_whitebox_test.go
// Date:     07.Feb.2024
//
// =============================================================================

package gapbuffer //nolint:testpackage // I want to white-box test this

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var curLineTests = []struct { //nolint:gochecknoglobals // not a global
	name          string
	lines         lineBuffer
	expectedStart int
	expectedEnd   int
}{
	{
		name: "all 2s",
		lines: lineBuffer{
			start:   8,
			end:     10,
			lengths: []int{2, 2, 2, 2, 2, 2, 2, 2, 2, 0},
		},
		expectedStart: 8 * 2,
		expectedEnd:   8*2 + 1,
	},
	{
		name: "ascending",
		lines: lineBuffer{
			start:   8,
			end:     10,
			lengths: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		},
		expectedStart: 1 + 2 + 3 + 4 + 5 + 6 + 7 + 8,
		expectedEnd:   1 + 2 + 3 + 4 + 5 + 6 + 7 + 8 + 8,
	},
}

func TestCurLineStart(t *testing.T) {
	t.Parallel()

	for i := range curLineTests {
		tStrct := &curLineTests[i]
		t.Run(tStrct.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tStrct.expectedStart, tStrct.lines.curLineStart(), "Error in curLineStart!")
			assert.Equal(t, tStrct.expectedEnd, tStrct.lines.curLineEnd(), "Error in curLineEnd!")
		})
	}
}

var lineInsertTests = []struct { //nolint:gochecknoglobals // not a global
	name      string
	startText string
	text      string
	insertPos int
	expected  lineBuffer
}{
	{
		name:      "empty add newlines",
		startText: "",
		text:      "12\n12\n12\n12",
		insertPos: 0,
		expected: lineBuffer{
			start:   3,
			end:     10,
			lengths: []int{3, 3, 3, 2, 0, 0, 0, 0, 0, 0},
		},
	},
	{
		name:      "all 2s",
		startText: "12\n12\n12\n12\n12\n12\n12\n12\n12",
		text:      "34567890",
		insertPos: 26,
		expected: lineBuffer{
			start:   8,
			end:     10,
			lengths: []int{3, 3, 3, 3, 3, 3, 3, 3, 10, 0},
		},
	},
	{
		name:      "add with newlines",
		startText: "12\n12\n12\n12",
		text:      "12\n12\n12\n12\n12",
		insertPos: 11,
		expected: lineBuffer{
			start:   7,
			end:     10,
			lengths: []int{3, 3, 3, 5, 3, 3, 3, 2, 0, 0},
		},
	},
	{
		name:      "resize array",
		startText: "12\n12\n12\n12\n12\n12\n12",
		text:      "12\n12\n12\n12\n12\n12\n12\n12",
		insertPos: 20,
		expected: lineBuffer{
			start:   13,
			end:     20,
			lengths: []int{3, 3, 3, 3, 3, 3, 5, 3, 3, 3, 3, 3, 3, 2, 0, 0, 0, 0, 0, 0},
		},
	},
	{
		name:      "insert newline",
		startText: "12\n12",
		text:      "\n",
		insertPos: 5,
		expected: lineBuffer{
			start:   2,
			end:     10,
			lengths: []int{3, 3, 0, 0, 0, 0, 0, 0, 0, 0},
		},
	},
}

func TestLineInsert(t *testing.T) {
	t.Parallel()

	for i := range lineInsertTests {
		tStrct := &lineInsertTests[i]
		lb := newLineBufStr(tStrct.startText, 10)
		lb.insert(tStrct.text, tStrct.insertPos)
		t.Run(tStrct.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tStrct.expected, *lb, "Error in curLineStart!")
		})
	}
}

var insertTests = []struct { //nolint:gochecknoglobals,dupl // not a global, not a duplicate
	name           string
	initialText    string
	insertText     string
	capacity       int
	expectedStruct GapBuffer
}{
	{
		name:        "empty",
		initialText: "",
		insertText:  "",
		capacity:    10,
		expectedStruct: GapBuffer{
			start:    0,
			end:      10,
			wantsCol: 0,
			data:     []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			lines: lineBuffer{
				lengths: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				start:   0,
				end:     10,
			},
		},
	},
	{
		name:        "hello",
		initialText: "hello",
		insertText:  "",
		capacity:    10,
		expectedStruct: GapBuffer{
			start:    5,
			end:      10,
			wantsCol: 5,
			data:     []byte{'h', 'e', 'l', 'l', 'o', 0, 0, 0, 0, 0},
			lines: lineBuffer{
				lengths: []int{5, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				start:   0,
				end:     10,
			},
		},
	},
	{
		name:        "hello world",
		initialText: "hello",
		insertText:  " world!",
		capacity:    10,
		expectedStruct: GapBuffer{
			start:    12,
			end:      20,
			wantsCol: 12,
			data:     []byte{'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd', '!', 0, 0, 0, 0, 0, 0, 0, 0},
			lines: lineBuffer{
				lengths: []int{12, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				start:   0,
				end:     10,
			},
		},
	},
	{
		name:        "newlines",
		initialText: "h\nel\nlo",
		insertText:  "\nwo\nld!",
		capacity:    20,
		expectedStruct: GapBuffer{
			start:    14,
			end:      20,
			wantsCol: 3,
			data:     []byte{'h', '\n', 'e', 'l', '\n', 'l', 'o', '\n', 'w', 'o', '\n', 'l', 'd', '!', 0, 0, 0, 0, 0, 0},
			lines: lineBuffer{
				lengths: []int{2, 3, 3, 3, 3, 0, 0, 0, 0, 0},
				start:   4,
				end:     10,
			},
		},
	},
	{
		name:        "only newlines",
		initialText: "h\nel\nlo",
		insertText:  "\n\n\n\n\n",
		capacity:    20,
		expectedStruct: GapBuffer{
			start:    12,
			end:      20,
			wantsCol: 0,
			data:     []byte{'h', '\n', 'e', 'l', '\n', 'l', 'o', '\n', '\n', '\n', '\n', '\n', 0, 0, 0, 0, 0, 0, 0, 0},
			lines: lineBuffer{
				lengths: []int{2, 3, 3, 1, 1, 1, 1, 0, 0, 0},
				start:   7,
				end:     10,
			},
		},
	},
}

func TestGapBufferInsert(t *testing.T) {
	t.Parallel()

	for i := range insertTests {
		tStrct := &insertTests[i]
		t.Run(tStrct.name, func(t *testing.T) {
			t.Parallel()
			gb := NewStrCap(tStrct.initialText, tStrct.capacity)
			gb.Insert(tStrct.insertText)
			assert.Equal(t, tStrct.expectedStruct, *gb, "Error in GapBuffer initialization!")
		})
	}
}

var insertMvLeftTests = []struct { //nolint:gochecknoglobals,dupl // not a global, not duplicates
	name           string
	initialText    string
	insertText     string
	capacity       int
	expectedStruct GapBuffer
}{
	{
		name:        "empty",
		initialText: "",
		insertText:  "",
		capacity:    10,
		expectedStruct: GapBuffer{
			start:    0,
			end:      10,
			wantsCol: 0,
			data:     []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			lines: lineBuffer{
				lengths: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				start:   0,
				end:     10,
			},
		},
	},
	{
		name:        "hello",
		initialText: "hello",
		insertText:  "",
		capacity:    10,
		expectedStruct: GapBuffer{
			start:    2,
			end:      7,
			wantsCol: 2,
			data:     []byte{'h', 'e', 'l', 'l', 'o', 0, 0, 'l', 'l', 'o'},
			lines: lineBuffer{
				lengths: []int{5, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				start:   0,
				end:     10,
			},
		},
	},
	{
		name:        "hello world",
		initialText: "hello",
		insertText:  " world!",
		capacity:    10,
		expectedStruct: GapBuffer{
			start:    9,
			end:      17,
			wantsCol: 9,
			data:     []byte{'h', 'e', ' ', 'w', 'o', 'r', 'l', 'd', '!', 0, 0, 0, 0, 0, 0, 0, 0, 'l', 'l', 'o'},
			lines: lineBuffer{
				lengths: []int{12, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				start:   0,
				end:     10,
			},
		},
	},
	{
		name:        "newlines",
		initialText: "h\nel\nlo",
		insertText:  "\nwo\nld!",
		capacity:    20,
		expectedStruct: GapBuffer{
			start:    8,
			end:      14,
			wantsCol: 3,
			data:     []byte{'h', '\n', 'w', 'o', '\n', 'l', 'd', '!', 0, 0, 0, 0, 0, 0, '\n', 'e', 'l', '\n', 'l', 'o'},
			lines: lineBuffer{
				lengths: []int{2, 3, 4, 0, 0, 0, 0, 0, 3, 2},
				start:   2,
				end:     8,
			},
		},
	},
	{
		name:        "only newlines",
		initialText: "h\nel\nlo",
		insertText:  "\n\n\n\n\n",
		capacity:    20,
		expectedStruct: GapBuffer{
			start:    6,
			end:      14,
			wantsCol: 0,
			data:     []byte{'h', '\n', '\n', '\n', '\n', '\n', 'o', 0, 0, 0, 0, 0, 0, 0, '\n', 'e', 'l', '\n', 'l', 'o'},
			lines: lineBuffer{
				lengths: []int{2, 1, 1, 1, 1, 1, 0, 0, 3, 2},
				start:   5,
				end:     8,
			},
		},
	},
}

func TestGapBufferMvLeftInsert(t *testing.T) {
	t.Parallel()

	for i := range insertMvLeftTests {
		tStrct := &insertMvLeftTests[i]
		t.Run(tStrct.name, func(t *testing.T) {
			t.Parallel()
			gBuf := NewStrCap(tStrct.initialText, tStrct.capacity)
			gBuf.UpMv()
			gBuf.LeftMv()
			gBuf.LeftMv()
			gBuf.LeftMv()
			gBuf.Insert(tStrct.insertText)
			assert.Equal(t, tStrct.expectedStruct, *gBuf, "Error in GapBuffer initialization!")
		})
	}
}

var insertMvRightTests = []struct { //nolint:gochecknoglobals // not a global
	name           string
	initialText    string
	insertText     string
	capacity       int
	expectedStruct GapBuffer
}{
	{
		name:        "empty",
		initialText: "",
		insertText:  "",
		capacity:    10,
		expectedStruct: GapBuffer{
			start:    0,
			end:      10,
			wantsCol: 0,
			data:     []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			lines: lineBuffer{
				lengths: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				start:   0,
				end:     10,
			},
		},
	},
	{
		name:        "newlines",
		initialText: "h\nel\nlo",
		insertText:  "\nwo\nld!",
		capacity:    20,
		expectedStruct: GapBuffer{
			start:    11,
			end:      17,
			wantsCol: 3,
			data:     []byte{'h', '\n', 'e', 'l', '\n', 'w', 'o', '\n', 'l', 'd', '!', 0, 0, 0, '\n', 'e', 'l', '\n', 'l', 'o'},
			lines: lineBuffer{
				lengths: []int{2, 3, 3, 4, 0, 0, 0, 0, 3, 2},
				start:   3,
				end:     9,
			},
		},
	},
	{
		name:        "only newlines",
		initialText: "h\nel\nlo",
		insertText:  "\n\n\n\n\n",
		capacity:    20,
		expectedStruct: GapBuffer{
			start:    9,
			end:      17,
			wantsCol: 0,
			data:     []byte{'h', '\n', 'e', 'l', '\n', '\n', '\n', '\n', '\n', 0, 0, 0, 0, 0, '\n', 'e', 'l', '\n', 'l', 'o'},
			lines: lineBuffer{
				lengths: []int{2, 3, 1, 1, 1, 1, 1, 0, 3, 2},
				start:   6,
				end:     9,
			},
		},
	},
}

func TestGapBufferMvRightInsert(t *testing.T) {
	t.Parallel()

	for i := range insertMvRightTests {
		tStrct := &insertMvRightTests[i]
		t.Run(tStrct.name, func(t *testing.T) {
			t.Parallel()
			gBuf := NewStrCap(tStrct.initialText, tStrct.capacity)
			gBuf.UpMv()
			gBuf.UpMv()
			gBuf.RightMv()
			gBuf.RightMv()
			gBuf.RightMv()
			gBuf.Insert(tStrct.insertText)
			assert.Equal(t, tStrct.expectedStruct, *gBuf, "Error in GapBuffer initialization!")
		})
	}
}

var upDownEmptyTests = []struct { //nolint:gochecknoglobals // not a global
	name           string
	initialText    string
	capacity       int
	expectedStruct GapBuffer
}{
	{
		name:        "empty",
		initialText: "",
		capacity:    10,
		expectedStruct: GapBuffer{
			start:    0,
			end:      10,
			wantsCol: 0,
			data:     []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			lines: lineBuffer{
				lengths: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				start:   0,
				end:     10,
			},
		},
	},
	{
		name:        "only newlines",
		initialText: "\n",
		capacity:    10,
		expectedStruct: GapBuffer{
			start:    1,
			end:      10,
			wantsCol: 0,
			data:     []byte{'\n', 0, 0, 0, 0, 0, 0, 0, 0, '\n'},
			lines: lineBuffer{
				lengths: []int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				start:   1,
				end:     10,
			},
		},
	},
}

func TestGapBufferUpDownEmpty(t *testing.T) {
	t.Parallel()

	for i := range upDownEmptyTests {
		tStrct := &upDownEmptyTests[i]
		t.Run(tStrct.name, func(t *testing.T) {
			t.Parallel()
			gBuf := NewStrCap(tStrct.initialText, tStrct.capacity)
			gBuf.UpMv()
			gBuf.DownMv()
			assert.Equal(t, tStrct.expectedStruct, *gBuf, "Error in GapBuffer initialization!")
		})
	}
}

var upDownInsertTests = []struct { //nolint:gochecknoglobals // not a global
	name           string
	initialText    string
	insertText     string
	capacity       int
	expectedStruct GapBuffer
}{
	{
		name:        "only newlines",
		initialText: "\n1",
		insertText:  "12",
		capacity:    10,
		expectedStruct: GapBuffer{
			start:    4,
			end:      10,
			wantsCol: 2,
			data:     []byte{'1', '2', '\n', '1', 0, 0, 0, 0, '\n', '1'},
			lines: lineBuffer{
				lengths: []int{3, 1, 0, 0, 0, 0, 0, 0, 0, 1},
				start:   1,
				end:     10,
			},
		},
	},
	{
		name:        "up and down",
		initialText: "12\n",
		insertText:  "",
		capacity:    10,
		expectedStruct: GapBuffer{
			start:    3,
			end:      10,
			wantsCol: 3,
			data:     []byte{'1', '2', '\n', 0, 0, 0, 0, 0, 0, '\n'},
			lines: lineBuffer{
				lengths: []int{2, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				start:   1,
				end:     10,
			},
		},
	},
	{
		name:        "up and down 2",
		initialText: "1\n1",
		insertText:  "\n",
		capacity:    10,
		expectedStruct: GapBuffer{
			start:    3,
			end:      9,
			wantsCol: 0,
			data:     []byte{'1', '\n', '\n', 0, 0, 0, 0, 0, '\n', '1'},
			lines: lineBuffer{
				lengths: []int{2, 1, 1, 0, 0, 0, 0, 0, 0, 1},
				start:   2,
				end:     10,
			},
		},
	},
}

func TestGapBufferUpDownInsert(t *testing.T) {
	t.Parallel()

	for i := range upDownInsertTests {
		tStrct := &upDownInsertTests[i]
		t.Run(tStrct.name, func(t *testing.T) {
			t.Parallel()
			gBuf := NewStrCap(tStrct.initialText, tStrct.capacity)
			// if g.lines.curLine() == 1
			gBuf.UpMv()
			gBuf.Insert(tStrct.insertText)
			// if g.lines.end > g.lines.lastIdx()
			gBuf.DownMv()
			assert.Equal(t, tStrct.expectedStruct, *gBuf, "Error in GapBuffer initialization!")
		})
	}
}

var upDownTests = []struct { //nolint:gochecknoglobals // not a global
	name           string
	initialText    string
	capacity       int
	expectedStruct GapBuffer
}{
	{
		name:        "only newlines",
		initialText: "\n1",
		capacity:    10,
		expectedStruct: GapBuffer{
			start:    2,
			end:      10,
			wantsCol: 1,
			data:     []byte{'\n', '1', 0, 0, 0, 0, 0, 0, '\n', '1'},
			lines: lineBuffer{
				lengths: []int{1, 1, 0, 0, 0, 0, 0, 0, 0, 1},
				start:   1,
				end:     10,
			},
		},
	},
	{
		name:        "up and down",
		initialText: "12\n",
		capacity:    10,
		expectedStruct: GapBuffer{
			start:    3,
			end:      10,
			wantsCol: 3,
			data:     []byte{'1', '2', '\n', 0, 0, 0, 0, '1', '2', '\n'},
			lines: lineBuffer{
				lengths: []int{3, 0, 0, 0, 0, 0, 0, 0, 0, 3},
				start:   1,
				end:     10,
			},
		},
	},
	{
		name:        "up and down 2",
		initialText: "1\n1",
		capacity:    10,
		expectedStruct: GapBuffer{
			start:    3,
			end:      10,
			wantsCol: 1,
			data:     []byte{'1', '\n', '1', 0, 0, 0, 0, 0, '\n', '1'},
			lines: lineBuffer{
				lengths: []int{2, 1, 0, 0, 0, 0, 0, 0, 0, 1},
				start:   1,
				end:     10,
			},
		},
	},
}

func TestGapBufferUpDown(t *testing.T) {
	t.Parallel()

	for i := range upDownTests {
		tStrct := &upDownTests[i]
		t.Run(tStrct.name, func(t *testing.T) {
			t.Parallel()
			gBuf := NewStrCap(tStrct.initialText, tStrct.capacity)
			gBuf.UpMv()
			gBuf.DownMv()
			assert.Equal(t, tStrct.expectedStruct, *gBuf, "Error in GapBuffer initialization!")
		})
	}
}
