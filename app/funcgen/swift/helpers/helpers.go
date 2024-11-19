package helpers

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"unicode"

	"github.com/brianvoe/gofakeit/v7"
)

var TitlePattern = regexp.MustCompile(`[^a-zA-Z0-9]`)

type (
	Action func() string
	Helper func()
)

// Returns answer on the yes or no question
func YesOrNo() bool {
	return rand.Intn(2) == 0
}

func CapitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}

	runes := []rune(s)
	for i := range runes {
		if i == 0 || !unicode.IsLetter(runes[i-1]) {
			runes[i] = unicode.ToUpper(runes[i])
		}
	}
	return string(runes)
}

func GenRandomPeremTitle() string {
	return TitlePattern.ReplaceAllString(fmt.Sprintf(
		"%s%s",
		gofakeit.Word(),
		CapitalizeFirstLetter(gofakeit.Word()),
	), "")
}

// Returns "level" for inserting
func CountTabsInString(str string, indEnd int) int {
	countTabs := 0

	if indEnd < 0 {
		return -1
	}

	for i := indEnd; i > 0; i-- {
		if str[i] == '\n' {
			return countTabs
		} else if str[i] == '\t' {
			countTabs++
		}
	}

	return countTabs
}

func CheckIfItEmptyCaseOfSwith(str string, indStart, indEnd int) bool {
	baseStr := []rune(str)
	downIsEmpty := false

	placeForSearchAfter := []rune(str[indEnd+1:])
	for i := 0; i < len(placeForSearchAfter); i++ {
		if placeForSearchAfter[i] == '\n' {
			afterString := string(strings.Trim(string(placeForSearchAfter[:i]), "\t"))
			if strings.Contains(afterString, "case ") ||
				strings.HasSuffix(afterString, "default:") ||
				afterString == "}" {
				downIsEmpty = true
				break
			}
		}
	}

	countN := 0
	for j := indStart - 1; j > 0; j-- {
		if baseStr[j] == '\n' {
			countN++
			if countN == 1 {
				continue
			}
			zalupa := make([]rune, indStart-j)
			copy(zalupa, baseStr[j:indStart])
			prevString := string(strings.Trim(string(zalupa), "\t"))
			if (strings.Contains(prevString, "switch ") ||
				strings.Contains(prevString, "default:") ||
				strings.Contains(prevString, "case ")) && downIsEmpty {

				return true
			}
			break
		}
	}
	return false
}

func PatternForSwitchConditions(
	mainCondition string,
	maxConditions int,
) (string, int) {
	var pattern string
	countConditions := rand.Intn(maxConditions)

	if countConditions == 0 {
		countConditions = 1
	}

	pattern = fmt.Sprintf("switch %v {\n", mainCondition)
	for i := 0; i < countConditions-1; i++ {
		pattern += "case %v:\n\tINSERT\n"
	}

	countConditions--
	pattern += "default:\n\tINSERT\n}"

	return pattern, countConditions
}

func PatternForIfElseConditions(maxConditions int) (string, int) {
	var pattern string
	countConditions := rand.Intn(maxConditions)

	if countConditions == 0 {
		countConditions = 1
	}

	pattern = "if (%v) {\n"
	for i := 0; i < countConditions-1; i++ {
		pattern += "\tINSERT\n}\nelse if (%v) {\n"
	}

	pattern += "\tINSERT\n}"

	if YesOrNo() {
		pattern += "\nelse{\n\tINSERT\n}"
	}

	return pattern, countConditions
}

func SelectConditionSign() string {
	selection := rand.Intn(5)
	switch selection {
	case 0:
		return "=="
	case 1:
		return "!="
	case 2:
		return ">"
	case 3:
		return "<"
	case 4:
		return ">="
	case 5:
		return "<="
	default:
		return "=="
	}
}

func SelectUnionSign() string {
	if YesOrNo() {
		return "&&"
	}

	return "||"
}

func MakeCondition(attrs []string) string {
	countComparisons := rand.Intn(4)
	if countComparisons == 0 {
		countComparisons = 1
	}

	condition := ""
	for i := 0; i < countComparisons-1; i++ {
		condition += fmt.Sprintf("(%v %v %v) %v ",
			attrs[rand.Intn(len(attrs))],
			SelectConditionSign(),
			attrs[rand.Intn(len(attrs))],
			SelectUnionSign(),
		)
	}

	condition += fmt.Sprintf(
		"(%v %v %v)",
		attrs[rand.Intn(len(attrs))],
		SelectConditionSign(),
		attrs[rand.Intn(len(attrs))],
	)

	return condition
}

func GenRandomAriphmeticStr(
	attrs []string,
) string {
	var ariphmeticStr string

	for i := 0; i < len(attrs)-1; i++ {
		ariphmeticStr += fmt.Sprintf("%v %v ", attrs[i], randMathSign())
	}

	ariphmeticStr += attrs[len(attrs)-1]

	return ariphmeticStr
}

func randMathSign() string {
	sign := rand.Intn(4)

	switch sign {
	case 0:
		return "-"
	case 1:
		return "+"
	case 2:
		return "*"
	case 3:
		return "/"
	case 4:
		return "%"
	}

	return "+"
}

func GenRandomBitwiseStr(
	attrs []string,
) string {
	var ariphmeticStr string

	isBeforeShift := false
	for i := 0; i < len(attrs)-1; i++ {
		sign := randBitwiseSign()
		if sign == ">>" || sign == "<<" {
			if isBeforeShift {
				ariphmeticStr += fmt.Sprintf("%v) %v ", attrs[i], sign)
			} else {
				ariphmeticStr += fmt.Sprintf("(%v %v ", attrs[i], sign)
			}

			isBeforeShift = !isBeforeShift

		} else {
			ariphmeticStr += fmt.Sprintf("%v %v ", attrs[i], sign)
		}
	}

	ariphmeticStr += attrs[len(attrs)-1]
	if isBeforeShift {
		ariphmeticStr += ")"
	}

	return ariphmeticStr
}

func randBitwiseSign() string {
	sign := rand.Intn(5)

	switch sign {
	case 0:
		return "<<"
	case 1:
		return "&"
	case 2:
		return "^"
	case 3:
		return "|"
	default:
		return ">>"
	}
}

func RandomLatinLetter() string {
	letters := []string{
		"q",
		"w",
		"e",
		"r",
		"t",
		"y",
		"u",
		"i",
		"o",
		"p",
		"a",
		"s",
		"d",
		"f",
		"g",
		"h",
		"j",
		"k",
		"l",
		"z",
		"x",
		"c",
		"v",
		"b",
		"n",
		"m",
	}

	return letters[rand.Intn(len(letters))]
}
