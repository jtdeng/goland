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


var lexTests = []lexTest{
	{"Test 1", "$.store.book[2].title", []item{
		{itemRoot, "$"},
		{itemChild, "store"},
		{itemChild, "book"},
		{itemIndex, "2"},
		{itemChild, "title"},
		{itemEOF, ""},
	}},
}

// collect gathers the emitted items into a slice.
func collect(t *lexTest) (items []item) {
	l := lex(t.name, t.input)
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
		items := collect(&test)
		//t.Errorf("|%v|", items)
		//t.Errorf("|%v|", test.items)
		if !reflect.DeepEqual(items, test.items) {
			t.Errorf("%s: got\n\t%#v\nexpected\n\t%#v", test.name, items, test.items)
		}
	}
}


