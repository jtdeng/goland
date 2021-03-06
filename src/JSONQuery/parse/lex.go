// Copyright 2012 James Deng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package parse

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// item represents a token or text string returned from the scanner.
type item struct {
	typ itemType
	val string
}

func (i item) String() string {
	switch {
	case i.typ == itemEOF:
		return "EOF"
	case i.typ == itemError:
		return i.val
	}
	return fmt.Sprintf("%q", i.val)
}

// itemType identifies the type of lex items.
type itemType int



//examples
//"$.store..book[(@.length-1)].title"
//"$.store.book[?(@.price<10)].title"
//"$.store.book[*].title"
//"$.store.book[1,3].title"
//"$.store.book[1:3].title"

const (
	itemError        itemType = iota // error occurred; value is text of error
	itemRoot           //$
	itemChild          // .store
	itemRecursiveChild	// ..book
	itemNumber     // valid number for index 
	itemWildcard        // *
    itemCurrent          //@
    itemString     // quoted string (includes quotes)
    itemIndex		//[2]
	//itemLBracket    // [
    //itemRBracket // ]
    //itemLParentheses // (
    //itemRParentheses // )
    //itemLess        // <
    //itemGreat       // >
    //itemMinus       // -
    //itemPlus        // +
    //itemEqual       //==
    //itemEval        //=
    //itemQuestion    // ?        
	itemEOF
	//itemKeyword    // just a separator
	//itemProperty   //all all property of current object
	//itemLength     //length of current object,
)

// Make the types prettyprint.
var itemName = map[itemType]string{
	itemError:        "error",
	//itemChar:         "char",
	//itemCharConstant: "charconst",
	itemEOF:          "EOF",
	itemRoot:			"$",
	itemChild:        "child",
	itemRecursiveChild:   "..child",
	itemIndex:			"index",
	itemNumber:       "number",
	itemString:       "string",
}


func (i itemType) String() string {
	s := itemName[i]
	if s == "" {
		return fmt.Sprintf("item%d", int(i))
	}
	return s
}

const eof = -1

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*lexer) stateFn


// lexer holds the state of the scanner.
type lexer struct {
	name       string    // the name of the input; used only for error reports.
	input      string    // the string being scanned.
	state      stateFn   // the next lexing function to enter.
	pos        int       // current position in the input.
	start      int       // start position of this item.
	width      int       // width of last rune read from input.
	items      chan item // channel of scanned items.
}


// next returns the next rune in the input.
func (l *lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
}

// emit passes an item back to the client.
func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.start = l.pos
}

// accept consumes the next rune if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

// lineNumber reports which line we're on. Doing it this way
// means we don't have to worry about peek double counting.
func (l *lexer) lineNumber() int {
	return 1 + strings.Count(l.input[:l.pos], "\n")
}

// error returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.run.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{itemError, fmt.Sprintf(format, args...)}
	return nil
}

// nextItem returns the next item from the input.
func (l *lexer) nextItem() item {
	for {
		select {
		case item := <-l.items:
			return item
		default:
			l.state = l.state(l)
		}
	}
	panic("not reached")
}

//trim the leading & trailing white spaces, smartly add the $ to input JSONQuery string 
func preprocess(i string) string {
    o := strings.TrimSpace(i)
    switch c := o[0]; {
		case c == '$':
			return o
		case c == '.' || c == '[': //.store or ['store'] or [2]
            return "$" + o
    }

	return "$." + o
}

// lex creates a new scanner for the input string.
func lex(name, input string) *lexer {
		l := &lexer{
		name:       name,
		input:      preprocess(input),
		state:      lexRoot,
		items:      make(chan item, 2), // Two items of buffering is sufficient for all state functions
	}
	return l
}

// state functions as below

// lex always start from Root $, if $ is not found, add it
func lexRoot(l *lexer) stateFn {

	if r := l.next(); r == '$' {
		//fmt.Println(l.input[l.pos:])
		l.emit(itemRoot)
	} else if r == eof {
		l.emit(itemEOF)
		//fmt.Println("Got EOF...")
		return nil
	} else {
		l.backup()
	}

    if strings.HasPrefix(l.input[l.pos:], "['") {
		return lexQuoteChild
    } else if strings.HasPrefix(l.input[l.pos:], "[") {
        return lexIndex
    } else if strings.HasPrefix(l.input[l.pos:], "..") {
		return lexRecursiveChild
    } else if strings.HasPrefix(l.input[l.pos:], ".") {
        return lexChild
    } else {
		fmt.Println(l.input[l.pos:])
		return l.errorf("Unexpected character following root")
    }

	return lexRoot
}

//lex a direct child field
func lexChild(l *lexer) stateFn {
    if r := l.next(); r != '.' {
		return l.errorf("I should be indexing child, but got %s", r)
    } else {
		l.ignore() // we don't need the dot in our token value
	}

    //absorb the identifier
    for {
		if r := l.next(); ! isAlphaNumeric(r) {
			l.backup()
			break
		}
    }
	l.emit(itemChild)
	return lexRoot

}

func lexQuoteChild(l *lexer) stateFn {
	return nil
}


func lexRecursiveChild(l *lexer) stateFn {
	r := l.input[l.pos:l.pos+2]
	if r != ".." {
		return l.errorf("Expecting .., but got %s", r)
	} else {
		_ = l.next()
		_ = l.next()
		l.ignore()
	}

	for {
		if r := l.next(); ! isAlphaNumeric(r) {
			l.backup()
			break
		}
	}

	l.emit(itemRecursiveChild)

	return lexRoot

}

//Now it's like [numbers]
func lexIndex(l *lexer) stateFn {
	l.next()
	l.ignore() //absort and ignore '['

	for {
		if l.next(); l.input[l.pos] == ']' {
			l.emit(itemIndex)
			l.next()
			l.ignore() //absorb and ignore ']'
			break
		}
	}

	return lexRoot
}

func lexFilter(l *lexer) stateFn {
	return nil
}




// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	switch r {
	case ' ', '\t', '\n', '\r':
		return true
	}
	return false
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}


