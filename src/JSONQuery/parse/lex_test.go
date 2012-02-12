// Copyright 2012 James Deng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parse

import (
	"reflect"
	"testing"
)

type lexTest struct {
	name  string
	input string
	items []item
}

var (
	tEOF      = item{itemEOF, ""}
	tLeft     = item{itemLeftDelim, "{{"}
	tRight    = item{itemRightDelim, "}}"}
	tRange    = item{itemRange, "range"}
	tPipe     = item{itemPipe, "|"}
	tFor      = item{itemIdentifier, "for"}
	tQuote    = item{itemString, `"abc \n\t\" "`}
	raw       = "`" + `abc\n\t\" ` + "`"
	tRawQuote = item{itemRawString, raw}
)

var lexTests = []lexTest{

	{"variables", "{{$c := printf $ $hello $23 $ $var.Field .Method}}", []item{
		tLeft,
		{itemVariable, "$c"},
		{itemColonEquals, ":="},
		{itemIdentifier, "printf"},
		{itemVariable, "$"},
		{itemVariable, "$hello"},
		{itemVariable, "$23"},
		{itemVariable, "$"},
		{itemVariable, "$var.Field"},
		{itemField, ".Method"},
		tRight,
		tEOF,
	}},
	
	{"pipeline", `intro {{echo hi 1.2 |noargs|args 1 "hi"}} outro`, []item{
		{itemText, "intro "},
		tLeft,
		{itemIdentifier, "echo"},
		{itemIdentifier, "hi"},
		{itemNumber, "1.2"},
		tPipe,
		{itemIdentifier, "noargs"},
		tPipe,
		{itemIdentifier, "args"},
		{itemNumber, "1"},
		{itemString, `"hi"`},
		tRight,
		{itemText, " outro"},
		tEOF,
	}},
	
	
}

// collect gathers the emitted items into a slice.
func collect(t *lexTest, left, right string) (items []item) {
	l := lex(t.name, t.input, left, right)
	for {
		item := l.nextItem()
		items = append(items, item)
		if item.typ == itemEOF || item.typ == itemError {
			break
		}
	}
	return
}

func TestLex(t *testing.T) {
	for _, test := range lexTests {
		items := collect(&test, "", "")
		if !reflect.DeepEqual(items, test.items) {
			t.Errorf("%s: got\n\t%v\nexpected\n\t%v", test.name, items, test.items)
		}
	}
}


