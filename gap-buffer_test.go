// SPDX-FileCopyrightText:  Copyright 2024 Roland Csaszar
// SPDX-License-Identifier: MIT
//
// Project:  go-gap-buffer
// File:     gap-buffer_test.go
// Date:     07.Feb.2024
//
// =============================================================================

// Black-box testing of the gap buffer library.
package gapbuffer_test

import (
	"testing"

	gapbuffer "github.com/Release-Candidate/go-gap-buffer"
	"github.com/stretchr/testify/assert"
)

//==============================================================================
//                       Simple Sanity Checks

func TestEmpty(t *testing.T) {
	t.Parallel()

	gb := gapbuffer.New()
	strLen := gb.StringLength()

	assert.Equal(t, "", gb.String(), "Error, empty gap buffer isn't emtpy!")
	assert.Equal(t, 0, strLen, "Error checking string length!")
}

func TestInitial(t *testing.T) {
	t.Parallel()

	gb := gapbuffer.NewStr("Hello World!")
	l, r := gb.StringPair()
	strLen := gb.StringLength()

	assert.Equal(t, "Hello World!", l, "Error, left part isn't 'Hello World!'!")
	assert.Equal(t, "", r, "Error, right part isn't empty!")
	assert.Equal(t, len(l)+len(r), strLen, "Error checking string length!")
}

func TestMoveLeft(t *testing.T) {
	t.Parallel()

	gb := gapbuffer.NewStr("Hello World!")
	gb.LeftMv()
	gb.LeftMv()
	gb.LeftMv()
	gb.LeftMv()
	gb.LeftMv()
	gb.LeftMv()
	l, r := gb.StringPair()
	strLen := gb.StringLength()

	assert.Equal(t, "Hello ", l, "Error, left part isn't 'Hello '!")
	assert.Equal(t, "World!", r, "Error, right part isn't 'World!'!")
	assert.Equal(t, len(l)+len(r), strLen, "Error checking string length!")
}

func TestDeleteLeft(t *testing.T) {
	t.Parallel()

	gb := gapbuffer.NewStr("Hello World!")
	gb.LeftDel()
	gb.LeftDel()
	gb.LeftDel()
	gb.LeftDel()
	gb.LeftDel()
	gb.LeftDel()
	gb.LeftDel()
	l, r := gb.StringPair()
	strLen := gb.StringLength()

	assert.Equal(t, "Hello", l, "Error, left part isn't 'Hello'!")
	assert.Equal(t, "", r, "Error, right part isn't empty!")
	assert.Equal(t, len(l)+len(r), strLen, "Error checking string length!")
}

func TestDeleteRight(t *testing.T) {
	t.Parallel()

	gb := gapbuffer.NewStr("Hello World!")
	gb.LeftMv()
	gb.LeftMv()
	gb.LeftMv()
	gb.LeftMv()
	gb.LeftMv()
	gb.LeftMv()
	gb.LeftMv()
	gb.RightDel()
	gb.RightDel()
	gb.RightDel()
	gb.RightDel()
	gb.RightDel()
	gb.RightDel()
	gb.RightDel()
	gb.RightDel()
	l, r := gb.StringPair()
	strLen := gb.StringLength()

	assert.Equal(t, "Hello", l, "Error, left part isn't 'Hello '!")
	assert.Equal(t, "", r, "Error, right part isn't empty!")
	assert.Equal(t, len(l)+len(r), strLen, "Error checking string length!")
}

func TestInsertWithNewlines(t *testing.T) {
	t.Parallel()

	gb := gapbuffer.NewStr("Hello World!")
	gb.LeftMv()
	gb.LeftMv()
	gb.LeftMv()
	gb.LeftMv()
	gb.LeftMv()
	gb.LeftMv()
	gb.LeftMv()
	gb.RightDel()
	gb.Insert("\nfunny\n")

	l, r := gb.StringPair()
	strLen := gb.StringLength()

	assert.Equal(t, "Hello\nfunny\n", l, "Error, left part isn't 'Hello\\nfunny\\n'!")
	assert.Equal(t, "World!", r, "Error, right part isn't 'World!'!")
	assert.Equal(t, len(l)+len(r), strLen, "Error checking string length!")
}

func TestMoveUp(t *testing.T) {
	t.Parallel()

	gb := gapbuffer.NewStr("Hello World!")
	gb.LeftMv()
	gb.LeftMv()
	gb.LeftMv()
	gb.LeftMv()
	gb.LeftMv()
	gb.LeftMv()
	gb.Insert("\nfunny\n")
	gb.UpMv()

	l, r := gb.StringPair()
	strLen := gb.StringLength()

	assert.Equal(t, "Hello \n", l, "Error, left part isn't 'Hello \\n'!")
	assert.Equal(t, "funny\nWorld!", r, "Error, right part isn't 'funny\\nWorld!'!")
	assert.Equal(t, len(l)+len(r), strLen, "Error checking string length!")
}
