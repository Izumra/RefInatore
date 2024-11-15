package models

import (
	"fmt"
	"math/rand"
	"strings"
	"unicode"

	simpletypes "github.com/Izumra/RefInatore/app/domain/valueobjects/sipletypes"
	"github.com/brianvoe/gofakeit/v7"
)

type Perem struct {
	Title string
	Type  string
	Value string
}

func NewPerem() *Perem {
	perem := &Perem{}
	perem.Type, perem.Value, _ = genRandomPerem()

	perem.Title = genRandomPeremTitle()

	return perem
}

// Generate random perem with type and provides interface of actions for generated perem
func genRandomPerem() (string, string, any) {
	typeComplexityPerem := rand.Intn(2)

	switch typeComplexityPerem {
	case 1:
		simpleParam := simpletypes.New(-1)
		return simpleParam.Type, simpleParam.Value, nil
	default:
		simpleParam := simpletypes.New(-1)
		return simpleParam.Type, simpleParam.Value, nil
	}

}

func genRandomPeremTitle() string {
	return strings.ReplaceAll(fmt.Sprintf(
		"%s%s",
		gofakeit.Word(),
		capitalizeFirstLetter(gofakeit.Word()),
	), " ", "")
}

func capitalizeFirstLetter(s string) string {
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
