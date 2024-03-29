// SPDX-FileCopyrightText:  Copyright 2024 Roland Csaszar
// SPDX-License-Identifier: MIT
//
// Project:  go-gap-buffer
// File:     gap-buffer_whitebox_test.go
// Date:     07.Feb.2024
//
// =============================================================================

// White-box testing of the gap buffer library, using the internal
// representation of both the text gap buffer and the line lengths gap buffer.
package gapbuffer //nolint:testpackage // I want to white-box test this

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCurLineStartEndAscending(t *testing.T) {
	t.Parallel()

	lines := lineBuffer{
		start:   8,
		end:     10,
		lengths: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
	}
	s := lines.curLineStart()
	e := lines.curLineEnd()

	assert.Equal(t, 1+2+3+4+5+6+7+8, s, "Start")
	assert.Equal(t, 1+2+3+4+5+6+7+8+8, e, "End")
}

func TestCurLineStartEndAll2s(t *testing.T) {
	t.Parallel()

	lines := lineBuffer{
		start:   8,
		end:     10,
		lengths: []int{2, 2, 2, 2, 2, 2, 2, 2, 2, 0},
	}
	s := lines.curLineStart()
	e := lines.curLineEnd()

	assert.Equal(t, 8*2, s, "Start")
	assert.Equal(t, 8*2+1, e, "End")
}

func TestCurLineStartEndAll2s0(t *testing.T) {
	t.Parallel()

	lines := lineBuffer{
		start:   9,
		end:     10,
		lengths: []int{2, 2, 2, 2, 2, 2, 2, 2, 2, 0},
	}
	s := lines.curLineStart()
	e := lines.curLineEnd()

	assert.Equal(t, 9*2, s, "Start")
	assert.Equal(t, 9*2, e, "End")
}

func TestLineNew3332(t *testing.T) {
	t.Parallel()

	exp := lineBuffer{
		start:   3,
		end:     10,
		lengths: []int{3, 3, 3, 2, 0, 0, 0, 0, 0, 0},
	}
	lb := newLineBufStr("12\n12\n12\n12", 10)
	assert.Equal(t, exp, *lb)
}

func TestLineInsert33310(t *testing.T) {
	t.Parallel()

	exp := lineBuffer{
		start:   8,
		end:     10,
		lengths: []int{3, 3, 3, 3, 3, 3, 3, 3, 10, 0},
	}
	lb := newLineBufStr("12\n12\n12\n12\n12\n12\n12\n12\n12", 20)
	lb.insert("34567890", 25)
	assert.Equal(t, exp, *lb)
}

func TestLineInsert33532(t *testing.T) {
	t.Parallel()

	exp := lineBuffer{
		start:   7,
		end:     10,
		lengths: []int{3, 3, 3, 5, 3, 3, 3, 2, 0, 0},
	}
	lb := newLineBufStr("12\n12\n12\n12", 20)
	lb.insert("12\n12\n12\n12\n12", 11)
	assert.Equal(t, exp, *lb)
}

func TestLineInsert333353332(t *testing.T) {
	t.Parallel()

	exp := lineBuffer{
		start:   10,
		end:     20,
		lengths: []int{3, 3, 3, 5, 3, 3, 3, 3, 3, 3, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	lb := newLineBufStr("12\n12\n12\n12", 20)
	lb.insert("12\n12\n12\n12\n12\n12\n12\n12", 11)
	assert.Equal(t, exp, *lb)
}

func TestLineInsertNewline(t *testing.T) {
	t.Parallel()

	exp := lineBuffer{
		start:   2,
		end:     10,
		lengths: []int{3, 3, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	lb := newLineBufStr("12\n12", 20)
	lb.insert("\n", 5)
	assert.Equal(t, exp, *lb)
}

func TestLineInsertSpecial(t *testing.T) {
	t.Parallel()

	lineBuf := newLineBufStr("Hello ", 20)
	lineBuf.insert("\nfunny\n", 6)

	exp := lineBuffer{
		lengths: []int{7, 6, 0, 0, 0, 0, 0, 0, 0, 0},
		start:   2,
		end:     10,
	}

	assert.Equal(t, exp, *lineBuf)
}

func TestInsertEmpty(t *testing.T) {
	t.Parallel()

	gapBuf := NewStrCap("", 10)
	gapBuf.Insert("")

	exp := GapBuffer{
		start:    0,
		end:      10,
		wantsCol: 0,
		data:     []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		lines: lineBuffer{
			lengths: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			start:   0,
			end:     10,
		},
	}
	assert.Equal(t, exp, *gapBuf)
}

func TestInsertHelloWorld(t *testing.T) {
	t.Parallel()

	gapBuf := NewStrCap("hello ", 20)
	gapBuf.Insert("world!")

	exp := GapBuffer{
		start:    12,
		end:      20,
		wantsCol: 12,
		data:     []byte{'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd', '!', 0, 0, 0, 0, 0, 0, 0, 0},
		lines: lineBuffer{
			lengths: []int{12, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			start:   0,
			end:     10,
		},
	}
	assert.Equal(t, exp, *gapBuf)
}

func TestInsertHelloWorldNLs(t *testing.T) {
	t.Parallel()

	gapBuffer := NewStrCap("h\nel\nlo", 20)
	gapBuffer.Insert("\nwo\nld!")

	exp := GapBuffer{
		start:    14,
		end:      20,
		wantsCol: 3,
		data:     []byte{'h', '\n', 'e', 'l', '\n', 'l', 'o', '\n', 'w', 'o', '\n', 'l', 'd', '!', 0, 0, 0, 0, 0, 0},
		lines: lineBuffer{
			lengths: []int{2, 3, 3, 3, 3, 0, 0, 0, 0, 0},
			start:   4,
			end:     10,
		},
	}
	assert.Equal(t, exp, *gapBuffer)
}

func TestInsertHelloNLs(t *testing.T) {
	t.Parallel()

	gapBuf := NewStrCap("h\nel\nlo", 20)
	gapBuf.Insert("\n\n\n\n\n")

	exp := GapBuffer{
		start:    12,
		end:      20,
		wantsCol: 0,
		data:     []byte{'h', '\n', 'e', 'l', '\n', 'l', 'o', '\n', '\n', '\n', '\n', '\n', 0, 0, 0, 0, 0, 0, 0, 0},
		lines: lineBuffer{
			lengths: []int{2, 3, 3, 1, 1, 1, 1, 0, 0, 0},
			start:   7,
			end:     10,
		},
	}
	assert.Equal(t, exp, *gapBuf)
}

func TestMvLeftInsertEmpty(t *testing.T) {
	t.Parallel()

	gBuf := NewStrCap("", 10)
	gBuf.UpMv()
	gBuf.LeftMv()
	gBuf.LeftMv()
	gBuf.LeftMv()
	gBuf.Insert("")

	exp := GapBuffer{
		start:    0,
		end:      10,
		wantsCol: 0,
		data:     []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		lines: lineBuffer{
			lengths: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			start:   0,
			end:     10,
		},
	}
	assert.Equal(t, exp, *gBuf)
}

func TestMvLeftInsertHello(t *testing.T) {
	t.Parallel()

	gBuf := NewStrCap("hello", 10)
	gBuf.UpMv()
	gBuf.LeftMv()
	gBuf.LeftMv()
	gBuf.LeftMv()
	gBuf.Insert("")

	exp := GapBuffer{
		start:    2,
		end:      7,
		wantsCol: 2,
		data:     []byte{'h', 'e', 'l', 'l', 'o', 0, 0, 'l', 'l', 'o'},
		lines: lineBuffer{
			lengths: []int{5, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			start:   0,
			end:     10,
		},
	}
	assert.Equal(t, exp, *gBuf)
}

func TestMvLeftInsertHelloWorld(t *testing.T) {
	t.Parallel()

	gBuf := NewStrCap("hello", 10)
	gBuf.UpMv()
	gBuf.LeftMv()
	gBuf.LeftMv()
	gBuf.LeftMv()
	gBuf.Insert(" world!")

	exp := GapBuffer{
		start:    9,
		end:      17,
		wantsCol: 9,
		data:     []byte{'h', 'e', ' ', 'w', 'o', 'r', 'l', 'd', '!', 0, 0, 0, 0, 0, 0, 0, 0, 'l', 'l', 'o'},
		lines: lineBuffer{
			lengths: []int{12, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			start:   0,
			end:     10,
		},
	}
	assert.Equal(t, exp, *gBuf)
}

func TestMvLeftInsertHelloWorldNL(t *testing.T) {
	t.Parallel()

	gBuf := NewStrCap("h\nel\nlo", 20)
	gBuf.UpMv()
	gBuf.LeftMv()
	gBuf.LeftMv()
	gBuf.LeftMv()
	gBuf.Insert("\nwo\nld!")

	exp := GapBuffer{
		start:    8,
		end:      14,
		wantsCol: 3,
		data:     []byte{'h', '\n', 'w', 'o', '\n', 'l', 'd', '!', 0, 0, 0, 0, 0, 0, '\n', 'e', 'l', '\n', 'l', 'o'},
		lines: lineBuffer{
			lengths: []int{2, 3, 4, 0, 0, 0, 0, 0, 3, 2},
			start:   2,
			end:     8,
		},
	}
	assert.Equal(t, exp, *gBuf)
}

func TestMvLeftInsertHelloNL(t *testing.T) {
	t.Parallel()

	gBuf := NewStrCap("h\nel\nlo", 20)
	gBuf.UpMv()
	gBuf.LeftMv()
	gBuf.LeftMv()
	gBuf.LeftMv()
	gBuf.Insert("\n\n\n\n\n")

	exp := GapBuffer{
		start:    6,
		end:      14,
		wantsCol: 0,
		data:     []byte{'h', '\n', '\n', '\n', '\n', '\n', 'o', 0, 0, 0, 0, 0, 0, 0, '\n', 'e', 'l', '\n', 'l', 'o'},
		lines: lineBuffer{
			lengths: []int{2, 1, 1, 1, 1, 1, 0, 0, 3, 2},
			start:   5,
			end:     8,
		},
	}
	assert.Equal(t, exp, *gBuf)
}

func TestMvRightEmpty(t *testing.T) {
	t.Parallel()

	gBuf := NewStrCap("", 10)
	gBuf.UpMv()
	gBuf.UpMv()
	gBuf.RightMv()
	gBuf.RightMv()
	gBuf.RightMv()
	gBuf.Insert("")

	exp := GapBuffer{
		start:    0,
		end:      10,
		wantsCol: 0,
		data:     []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		lines: lineBuffer{
			lengths: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			start:   0,
			end:     10,
		},
	}
	assert.Equal(t, exp, *gBuf)
}

func TestMvRightHelloNL(t *testing.T) {
	t.Parallel()

	gBuf := NewStrCap("h\nel\nlo", 20)
	gBuf.UpMv()
	gBuf.UpMv()
	gBuf.RightMv()
	gBuf.RightMv()
	gBuf.RightMv()
	gBuf.Insert("\nwo\nld!")

	exp := GapBuffer{
		start:    11,
		end:      17,
		wantsCol: 3,
		data:     []byte{'h', '\n', 'e', 'l', '\n', 'w', 'o', '\n', 'l', 'd', '!', 0, 0, 0, '\n', 'e', 'l', '\n', 'l', 'o'},
		lines: lineBuffer{
			lengths: []int{2, 3, 3, 4, 0, 0, 0, 0, 3, 2},
			start:   3,
			end:     9,
		},
	}
	assert.Equal(t, exp, *gBuf)
}

func TestMvRightNL(t *testing.T) {
	t.Parallel()

	gBuf := NewStrCap("h\nel\nlo", 20)
	gBuf.UpMv()
	gBuf.UpMv()
	gBuf.RightMv()
	gBuf.RightMv()
	gBuf.RightMv()
	gBuf.Insert("\n\n\n\n\n")

	exp := GapBuffer{
		start:    9,
		end:      17,
		wantsCol: 0,
		data:     []byte{'h', '\n', 'e', 'l', '\n', '\n', '\n', '\n', '\n', 0, 0, 0, 0, 0, '\n', 'e', 'l', '\n', 'l', 'o'},
		lines: lineBuffer{
			lengths: []int{2, 3, 1, 1, 1, 1, 1, 0, 3, 2},
			start:   6,
			end:     9,
		},
	}
	assert.Equal(t, exp, *gBuf)
}

func TestUpDownEmpty(t *testing.T) {
	t.Parallel()

	gBuf := NewStrCap("", 10)
	gBuf.UpMv()
	gBuf.DownMv()

	exp := GapBuffer{
		start:    0,
		end:      10,
		wantsCol: 0,
		data:     []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		lines: lineBuffer{
			lengths: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			start:   0,
			end:     10,
		},
	}

	assert.Equal(t, exp, *gBuf)
}

func TestUpDownNL(t *testing.T) {
	t.Parallel()

	gBuf := NewStrCap("\n", 10)
	gBuf.UpMv()
	gBuf.DownMv()

	exp := GapBuffer{
		start:    1,
		end:      10,
		wantsCol: 0,
		data:     []byte{'\n', 0, 0, 0, 0, 0, 0, 0, 0, '\n'},
		lines: lineBuffer{
			lengths: []int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			start:   1,
			end:     10,
		},
	}

	assert.Equal(t, exp, *gBuf)
}

func TestUpDownInsert12(t *testing.T) {
	t.Parallel()

	gBuf := NewStrCap("\n1", 10)
	gBuf.UpMv()
	gBuf.Insert("12")
	gBuf.DownMv()

	exp := GapBuffer{
		start:    4,
		end:      10,
		wantsCol: 2,
		data:     []byte{'1', '2', '\n', '1', 0, 0, 0, 0, '\n', '1'},
		lines: lineBuffer{
			lengths: []int{3, 1, 0, 0, 0, 0, 0, 0, 0, 1},
			start:   1,
			end:     10,
		},
	}
	assert.Equal(t, exp, *gBuf)
}

func TestUpDownInsert12NL(t *testing.T) {
	t.Parallel()

	gBuf := NewStrCap("12\n", 10)
	gBuf.UpMv()
	gBuf.Insert("")
	gBuf.DownMv()

	exp := GapBuffer{
		start:    3,
		end:      10,
		wantsCol: 0,
		data:     []byte{'1', '2', '\n', 0, 0, 0, 0, '1', '2', '\n'},
		lines: lineBuffer{
			lengths: []int{3, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			start:   1,
			end:     10,
		},
	}
	assert.Equal(t, exp, *gBuf)
}

func TestUpDownInsert11NL(t *testing.T) {
	t.Parallel()

	gBuf := NewStrCap("1\n1", 10)
	gBuf.UpMv()
	gBuf.Insert("\n")
	gBuf.DownMv()

	exp := GapBuffer{
		start:    3,
		end:      9,
		wantsCol: 0,
		data:     []byte{'1', '\n', '\n', 0, 0, 0, 0, 0, '\n', '1'},
		lines: lineBuffer{
			lengths: []int{2, 1, 1, 0, 0, 0, 0, 0, 0, 1},
			start:   2,
			end:     10,
		},
	}
	assert.Equal(t, exp, *gBuf)
}

func TestUpDownNL1(t *testing.T) {
	t.Parallel()

	gBuf := NewStrCap("\n1", 10)
	gBuf.UpMv()
	gBuf.DownMv()

	exp := GapBuffer{
		start:    2,
		end:      10,
		wantsCol: 1,
		data:     []byte{'\n', '1', 0, 0, 0, 0, 0, 0, '\n', '1'},
		lines: lineBuffer{
			lengths: []int{1, 1, 0, 0, 0, 0, 0, 0, 0, 1},
			start:   1,
			end:     10,
		},
	}

	assert.Equal(t, exp, *gBuf)
}

func TestUpDown12NL(t *testing.T) {
	t.Parallel()

	gBuf := NewStrCap("12\n", 10)
	gBuf.UpMv()
	gBuf.DownMv()

	exp := GapBuffer{
		start:    3,
		end:      10,
		wantsCol: 0,
		data:     []byte{'1', '2', '\n', 0, 0, 0, 0, '1', '2', '\n'},
		lines: lineBuffer{
			lengths: []int{3, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			start:   1,
			end:     10,
		},
	}

	assert.Equal(t, exp, *gBuf)
}

func TestUpDown1NL1(t *testing.T) {
	t.Parallel()

	gBuf := NewStrCap("1\n1", 10)
	gBuf.UpMv()
	gBuf.DownMv()

	exp := GapBuffer{
		start:    3,
		end:      10,
		wantsCol: 1,
		data:     []byte{'1', '\n', '1', 0, 0, 0, 0, 0, '\n', '1'},
		lines: lineBuffer{
			lengths: []int{2, 1, 0, 0, 0, 0, 0, 0, 0, 1},
			start:   1,
			end:     10,
		},
	}

	assert.Equal(t, exp, *gBuf)
}
