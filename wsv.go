package wsv

import (
	"errors"
	"strings"
)

const codepoint_LINEFEED = 0x0A
const codepoint_DOUBLEQUOTE = 0x22
const codepoint_HASH = 0x23
const codepoint_SLASH = 0x2F

type basicWsvCharIterator struct {
	chars     []rune
	index     int
	lineIndex int
}

func (x *basicWsvCharIterator) isEnd() bool {
	return x.index >= len(x.chars)
}

func (x *basicWsvCharIterator) is(c rune) bool {
	return rune(x.chars[x.index]) == c
}

func (x *basicWsvCharIterator) isWhitespace() bool {
	return isWhitespace(rune(x.chars[x.index]))
}

func (x *basicWsvCharIterator) next() bool {
	x.index = x.index + 1
	return !x.isEnd()
}

func (x *basicWsvCharIterator) get() rune {
	return rune(x.chars[x.index])
}

func (x *basicWsvCharIterator) getSlice(startIndex int) []rune {
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
	var iterator = basicWsvCharIterator{getCodePoints(lineStrWithoutLinefeed), 0, lineIndex}
	var values []string
	for {
		skipWhitespace(&iterator)
		if iterator.isEnd() {
			break
		}
		if iterator.is(codepoint_HASH) {
			break
		}

		var curValue string
		if iterator.is(codepoint_DOUBLEQUOTE) {

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

func parseValue(iterator *basicWsvCharIterator) (string, error) {
	var startIndex = iterator.index
	for {
		if !iterator.next() {
			break
		}
		if iterator.isWhitespace() || iterator.is(codepoint_HASH) {
			break
		} else if iterator.is(codepoint_DOUBLEQUOTE) {
			return "", errors.New("invalid double quote in value")

		}
	}

	return string(iterator.getSlice(startIndex)), nil

}

func parseDoubleQuotedValue(iterator *basicWsvCharIterator) (string, error) {
	var value = ""
	for {
		if !iterator.next() {
			return "", errors.New("string not closed")
		}
		if iterator.is(codepoint_DOUBLEQUOTE) {
			if !iterator.next() {
				break
			}
			if iterator.is(codepoint_DOUBLEQUOTE) {
				value += "\""
			} else if iterator.is(codepoint_SLASH) {
				if !iterator.next() && iterator.is(codepoint_DOUBLEQUOTE) {
					return "", errors.New("invalid string line break")
				}
				value += "\n"
			} else if iterator.isWhitespace() || iterator.is(codepoint_HASH) {
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

func skipWhitespace(iterator *basicWsvCharIterator) {
	if iterator.isEnd() {
		return
	}

	//Bascially a do-while loop
	for ok := true; ok; ok = iterator.next() {
		if !iterator.isWhitespace() {
			break
		}
	}
}

func needsDoubleQuotes(value string) bool {
	if len(value) == 0 || value == "-" {
		return true
	}

	chars := getCodePoints(value)
	return containsSpecialChar(chars)
}

func containsSpecialChar(chars []rune) bool {
	for _, c := range chars {
		if isWhitespace(c) || c == codepoint_LINEFEED || c == codepoint_DOUBLEQUOTE || c == codepoint_HASH {
			return true
		}
	}
	return false
}

func serializeValue(value string) string {
	if len(value) == 0 {
		return "\"\""
	} else if value == "-" {
		return "\"-\""
	} else {
		chars := getCodePoints(value)
		if containsSpecialChar(chars) {
			var result = "\""
			for _, c := range chars {
				if c == codepoint_LINEFEED {
					result += "\"/\""
				} else if c == codepoint_DOUBLEQUOTE {
					result += "\"\""
				} else {
					result += string(c)
				}
			}
			result += "\""
			return result
		} else {
			return value
		}
	}
}

func SerializeRow(values []string) string {
	isFirst := true
	result := ""
	for _, value := range values {
		if !isFirst {
			result += " "
		} else {
			isFirst = false
		}
		result += serializeValue(value)
	}
	return result
}

func Serialize(values [][]string) string {
	isFirst := true
	result := ""
	for _, line := range values {
		if !isFirst {
			result += "\n"
		} else {
			isFirst = false
		}
		result += SerializeRow(line)
	}
	return result
}
