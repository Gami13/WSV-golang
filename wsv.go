package wsv

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
	return isWhitespace(rune(x.chars[x.index]))
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

func isWhitespace(c rune) bool {
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

func getCodePoints(s string) []rune {
	var result []rune
	for _, c := range s {
		result = append(result, rune(c))
	}
	return result
}

// Parses a WSV document's line as an array of strings.
func ParseAsArray(content string) ([]string, error) {
	return ParseLineAsArray(content)
}

// Parses a WSV document as a jagged array of strings.
func ParseAsJaggedArray(content string) ([][]string, error) {
	return ParseDocument(content)
}

func ParseDocument(content string) ([][]string, error) {

	var lines = strings.Split(content, "\n")

	var result [][]string
	for i := 0; i < len(lines); i++ {

		var lineStr = lines[i]
		lineValues, err := parseLine(lineStr, i)
		if err != nil {
			return nil, err
		}
		result = append(result, lineValues)

	}
	return result, nil
}

func ParseLineAsArray(content string) ([]string, error) {
	result, err := ParseDocument(content)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, errors.New("empty document")
	}
	return result[0], nil
}

func parseLine(lineStrWithoutLinefeed string, lineIndex int) ([]string, error) {
	var iterator = BasicWsvCharIterator{getCodePoints(lineStrWithoutLinefeed), 0, lineIndex}
	var values []string
	for {
		skipWhitespace(&iterator)
		if iterator.isEnd() {
			break
		}
		if iterator.is(CODEPOINT_HASH) {
			break
		}

		var curValue string
		if iterator.is(CODEPOINT_DOUBLEQUOTE) {

			var err error
			curValue, err = parseDoubleQuotedValue(&iterator)
			if err != nil {
				return nil, err
			}
		} else {

			var err error
			curValue, err = parseValue(&iterator)
			if err != nil {
				return nil, err
			}
			if curValue == "-" {

				curValue = ""
			}
		}

		values = append(values, curValue)
	}

	return values, nil
}

func parseValue(iterator *BasicWsvCharIterator) (string, error) {
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

func parseDoubleQuotedValue(iterator *BasicWsvCharIterator) (string, error) {
	var value = ""
	for {
		if !iterator.next() {
			return "", errors.New("string not closed")
		}
		if iterator.is(CODEPOINT_DOUBLEQUOTE) {
			if !iterator.next() {
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

func skipWhitespace(iterator *BasicWsvCharIterator) {
	if iterator.isEnd() {
		return
	}
	// if iterator.next() {
	// 	if !iterator.isWhitespace() {
	// 		return
	// 	}
	// }
	// for {

	// 	if !iterator.isWhitespace() {
	// 		return
	// 	}
	// 	if !iterator.next() {
	// 		break

	// 	}
	// }

	for ok := true; ok; ok = iterator.next() {
		if !iterator.isWhitespace() {
			break
		}
	}
}
