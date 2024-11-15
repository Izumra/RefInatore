package arraycollection

import (
	"fmt"
	"math/rand/v2"
	"strings"

	simpletypes "github.com/Izumra/RefInator/app/domain/valueobjects/sipletypes"
)

type Struct struct {
	Title string
	Value []*simpletypes.Type
}

func New() (*Struct, string) {
	choosenType := rand.IntN(14)
	choosenCountElems := rand.IntN(11)

	arr := make([]*simpletypes.Type, choosenCountElems)
	for i := range arr {
		arr[i] = simpletypes.New(choosenType)
	}

	title := simpletypes.New(13).Value
	t := simpletypes.DetermineType(choosenType)

	code := fmt.Sprintf("var %s: [%s]", title, t)
	randTypeAssertion := rand.IntN(2)
	if randTypeAssertion == 0 {
		code = fmt.Sprintf("var %s: Array<%s>", title, t)
	}

	if choosenCountElems != 0 {
		code += " = ["
		chooseRandCountInitValues := rand.IntN(choosenCountElems)
		initionValues := []string{}
		for _, v := range arr[:chooseRandCountInitValues] {
			initionValues = append(initionValues, v.Value)
		}

		code += strings.Join(initionValues, ", ")
		code += "\n"

		for i, v := range arr[chooseRandCountInitValues:] {
			randAssertion := rand.IntN(2)
			if randAssertion == 0 {
				code += fmt.Sprintf("%s[%d] = %v", title, i, v.UpdateValueByRandomMathOp())
			}
		}

	} else {
		code += " = []"
	}

	return &Struct{
		Value: arr,
	}, code
}

func (s *Struct) InsertElems() {

}
