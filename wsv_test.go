package wsv_test

import (
	"os"
	"testing"

	"github.com/gami13/wsv-golang"
)

// a 	U+0061    61            0061        "Latin Small Letter A"
// ~ 	U+007E    7E            007E        Tilde
// ¥ 	U+00A5    C2_A5         00A5        "Yen Sign"
// » 	U+00BB    C2_BB         00BB        "Right-Pointing Double Angle Quotation Mark"
// ½ 	U+00BD    C2_BD         00BD        "Vulgar Fraction One Half"
// ¿ 	U+00BF    C2_BF         00BF        "Inverted Question Mark"
// ß 	U+00DF    C3_9F         00DF        "Latin Small Letter Sharp S"
// ä 	U+00E4    C3_A4         00E4        "Latin Small Letter A with Diaeresis"
// ï 	U+00EF    C3_AF         00EF        "Latin Small Letter I with Diaeresis"
// œ 	U+0153    C5_93         0153        "Latin Small Ligature Oe"
// € 	U+20AC    E2_82_AC      20AC        "Euro Sign"
// 東 	U+6771    E6_9D_B1      6771        "CJK Unified Ideograph-6771"
// 𝄞 	U+1D11E   F0_9D_84_9E   D834_DD1E   "Musical Symbol G Clef"
// 𠀇 	U+20007   F0_A0_80_87   D840_DC07   "CJK Unified Ideograph-20007"
var TEST_FILE_EXPECTED = [][]string{
	{"a", "U+0061", "61", "0061", "Latin Small Letter A"},
	{"~", "U+007E", "7E", "007E", "Tilde"},
	{"¥", "U+00A5", "C2_A5", "00A5", "Yen Sign"},
	{"»", "U+00BB", "C2_BB", "00BB", "Right-Pointing Double Angle Quotation Mark"},
	{"½", "U+00BD", "C2_BD", "00BD", "Vulgar Fraction One Half"},
	{"¿", "U+00BF", "C2_BF", "00BF", "Inverted Question Mark"},
	{"ß", "U+00DF", "C3_9F", "00DF", "Latin Small Letter Sharp S"},
	{"ä", "U+00E4", "C3_A4", "00E4", "Latin Small Letter A with Diaeresis"},
	{"ï", "U+00EF", "C3_AF", "00EF", "Latin Small Letter I with Diaeresis"},
	{"œ", "U+0153", "C5_93", "0153", "Latin Small Ligature Oe"},
	{"€", "U+20AC", "E2_82_AC", "20AC", "Euro Sign"},
	{"東", "U+6771", "E6_9D_B1", "6771", "CJK Unified Ideograph-6771"},
	{"𝄞", "U+1D11E", "F0_9D_84_9E", "D834_DD1E", "Musical Symbol G Clef"},
	{"𠀇", "U+20007", "F0_A0_80_87", "D840_DC07", "CJK Unified Ideograph-20007"},
}
var TEST_SERIALIZATION_ROW = `a U+0061 61 0061 "Latin Small Letter A"`
var TEST_SERIALIZATION_RESULT = `a U+0061 61 0061 "Latin Small Letter A"
~ U+007E 7E 007E Tilde
¥ U+00A5 C2_A5 00A5 "Yen Sign"
» U+00BB C2_BB 00BB "Right-Pointing Double Angle Quotation Mark"
½ U+00BD C2_BD 00BD "Vulgar Fraction One Half"
¿ U+00BF C2_BF 00BF "Inverted Question Mark"
ß U+00DF C3_9F 00DF "Latin Small Letter Sharp S"
ä U+00E4 C3_A4 00E4 "Latin Small Letter A with Diaeresis"
ï U+00EF C3_AF 00EF "Latin Small Letter I with Diaeresis"
œ U+0153 C5_93 0153 "Latin Small Ligature Oe"
€ U+20AC E2_82_AC 20AC "Euro Sign"
東 U+6771 E6_9D_B1 6771 "CJK Unified Ideograph-6771"
𝄞 U+1D11E F0_9D_84_9E D834_DD1E "Musical Symbol G Clef"
𠀇 U+20007 F0_A0_80_87 D840_DC07 "CJK Unified Ideograph-20007"`

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
	result, err := wsv.ParseDocument("a b c\n1 2 3")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	expected := [][]string{{"a", "b", "c"}, {"1", "2", "3"}}
	if !compareTables(expected, result) {
		println("Expected:")
		printTable(expected)
		println("Got:")
		printTable(result)
		t.Errorf("Expected %v, got %v", expected, result)
	}

}

func TestTwo(t *testing.T) {

	result, err := wsv.ParseDocument("a b c\n1 2 3 4\n5 6 7")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	expected := [][]string{{"a", "b", "c"}, {"1", "2", "3", "4"}, {"5", "6", "7"}}
	if !compareTables(expected, result) {
		println("Expected:")
		printTable(expected)
		println("Got:")
		printTable(result)
		t.Errorf("Expected %v, got %v", expected, result)
	}

}

func TestThree(t *testing.T) {

	file, err := os.ReadFile("test_input.wsv")
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	result, err := wsv.ParseDocument(string(file))
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if !compareTables(TEST_FILE_EXPECTED, result) {
		println("Expected:")
		printTable(TEST_FILE_EXPECTED)
		println("Got:")
		printTable(result)
		t.Errorf("Expected %v, got %v", TEST_FILE_EXPECTED, result)
	}

}

func TestFour(t *testing.T) {

	result := wsv.Serialize(TEST_FILE_EXPECTED)
	if TEST_SERIALIZATION_RESULT != result {
		t.Errorf("Expected %v, got %v", TEST_SERIALIZATION_RESULT, result)
	}

}

func TestFive(t *testing.T) {

	result := wsv.SerializeRow(TEST_FILE_EXPECTED[0])
	if TEST_SERIALIZATION_ROW != result {
		t.Errorf("Expected %v, got %v", TEST_SERIALIZATION_ROW, result)
	}

}
