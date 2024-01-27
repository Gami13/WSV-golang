package wsvgolang

import (
	"testing"
)

func compareTables(expected [][]string, result [][]string) bool {
	if len(expected) != len(result) {
		return false
	}
	for i := 0; i < len(expected); i++ {
		if len(expected[i]) != len(result[i]) {
			return false
		}
		for j := 0; j < len(expected[i]); j++ {
			if expected[i][j] != result[i][j] {
				return false
			}
		}
	}
	return true
}
func printTable(table [][]string) {
	for _, row := range table {
		for _, cell := range row {
			print(cell + " ")
		}
		println()
	}
}
func TestOne(t *testing.T) {

	result := (WsvParser).parseDocumentNonPreserving(WsvParser{}, "a b c\n1 2 3")

	expected := [][]string{{"a", "b", "c"}, {"1", "2", "3"}}
	if !compareTables(expected, result) {
		println("Expected:")
		printTable(expected)
		println("Got:")
		printTable(result)
		t.Errorf("Expected %v, got %v", expected, result)
	}

}
