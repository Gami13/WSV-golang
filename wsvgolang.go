package wsvgolang

import (
	"errors"
	"strings"
)

const CODEPOINT_LINEFEED = 0x0A
const CODEPOINT_DOUBLEQUOTE = 0x22
const CODEPOINT_HASH = 0x23
const CODEPOINT_SLASH = 0x2F

type BasicWsvCharIterator struct {
	chars     []rune
	index     int
	lineIndex int
}

func (x *BasicWsvCharIterator) isEnd() bool {
	return x.index >= len(x.chars)
}

func (x *BasicWsvCharIterator) is(c rune) bool {
	return rune(x.chars[x.index]) == c
}

func (x *BasicWsvCharIterator) isWhitespace() bool {
	return (WsvChar).isWhitespace(WsvChar{}, rune(x.chars[x.index]))
}

func (x *BasicWsvCharIterator) next() bool {
	x.index = x.index + 1
	return !x.isEnd()
}

func (x *BasicWsvCharIterator) get() rune {
	return rune(x.chars[x.index])
}

func (x *BasicWsvCharIterator) getSlice(startIndex int) []rune {
	return []rune(x.chars[startIndex:x.index])
}

type WsvChar struct{}

func (WsvChar) isWhitespace(c rune) bool {
	return c == 0x09 ||
		(c >= 0x0B && c <= 0x0D) ||
		c == 0x20 ||
		c == 0x85 ||
		c == 0xA0 ||
		c == 0x1680 ||
		(c >= 0x2000 && c <= 0x200A) ||
		(c >= 0x2028 && c <= 0x2029) ||
		c == 0x202F ||
		c == 0x205F ||
		c == 0x3000
}

func (x WsvChar) getCodePoints(s string) []rune {
	var result []rune
	for _, c := range s {
		result = append(result, rune(c))
	}
	return result
}

type WsvLine struct{}

func (x *WsvLine) parseAsArray(content string) []string {
	return (WsvParser).parseLineAsArray(WsvParser{}, content)
}

type WsvDocument struct {
}

func (x *WsvDocument) parseAsJaggedArray(content string) [][]string {
	return (WsvParser).parseDocumentNonPreserving(WsvParser{}, content)
}

type WsvParser struct {
}

func (x WsvParser) parseDocumentNonPreserving(content string) [][]string {
	var lines = strings.Split(content, "\n")
	var result [][]string
	for i := 0; i < len(lines); i++ {

		var lineStr = lines[i]
		var lineValues = x.parseLine(lineStr, i)
		result = append(result, lineValues)

	}
	return result
}

func (x WsvParser) parseLineAsArray(content string) []string {
	return x.parseDocumentNonPreserving(content)[0]
}

func (x *WsvParser) parseLine(lineStrWithoutLinefeed string, lineIndex int) []string {
	var iterator = BasicWsvCharIterator{(WsvChar).getCodePoints(WsvChar{}, lineStrWithoutLinefeed), 0, lineIndex}
	var values []string
	for {
		x.skipWhitespace(&iterator)
		if iterator.isEnd() {
			break
		}
		if iterator.is(CODEPOINT_HASH) {
			break
		}

		var curValue string
		if iterator.is(CODEPOINT_DOUBLEQUOTE) {

			var err error
			curValue, err = x.parseDoubleQuotedValue(&iterator)
			if err != nil {
				//TODO: error handling
			}
		} else {

			var err error
			curValue, err = x.parseValue(&iterator)
			if err != nil {
				//TODO: error handling
			}
			if curValue == "-" {

				curValue = ""
			}
		}

		values = append(values, curValue)
	}

	return values
}

func (x *WsvParser) parseValue(iterator *BasicWsvCharIterator) (string, error) {
	var startIndex = iterator.index
	for {
		if !iterator.next() {
			break
		}
		if iterator.isWhitespace() || iterator.is(CODEPOINT_HASH) {
			break
		} else if iterator.is(CODEPOINT_DOUBLEQUOTE) {
			return "", errors.New("invalid double quote in value")

		}
	}

	return string(iterator.getSlice(startIndex)), nil

}

func (x *WsvParser) parseDoubleQuotedValue(iterator *BasicWsvCharIterator) (string, error) {
	var value = ""
	for {
		if !iterator.next() {
			return "", errors.New("string not closed")
		}
		if iterator.is(CODEPOINT_DOUBLEQUOTE) {
			if iterator.next() {
				break
			}
			if iterator.is(CODEPOINT_DOUBLEQUOTE) {
				value += "\""
			} else if iterator.is(CODEPOINT_SLASH) {
				if !iterator.next() && iterator.is(CODEPOINT_DOUBLEQUOTE) {
					return "", errors.New("invalid string line break")
				}
				value += "\n"
			} else if iterator.isWhitespace() || iterator.is(CODEPOINT_HASH) {
				break
			} else {
				return "", errors.New("invalid character after string")
			}
		} else {
			value += string(iterator.get())
		}

	}
	return value, nil
}

func (x *WsvParser) skipWhitespace(iterator *BasicWsvCharIterator) {
	if iterator.isEnd() {
		return
	}
	if !iterator.isWhitespace() {
		return
	}
	for {

		if !iterator.isWhitespace() {
			return
		}
		if !iterator.next() {
			break

		}
	}
}
