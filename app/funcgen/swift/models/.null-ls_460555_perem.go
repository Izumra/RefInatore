package models

import (
	"math/rand"

	"github.com/Izumra/RefInatore/app/funcgen/swift/helpers"
	simpletypes "github.com/Izumra/RefInatore/app/funcgen/swift/valueobjects/simpletypes"
	"github.com/brianvoe/gofakeit/v7"
)

type Perem struct {
	Title   string
	Type    string
	Value   string
	Actions []helpers.Action
}

func NewPerem() *Perem {
	perem := &Perem{}

	typeComplexityPerem := rand.Intn(2)

	switch typeComplexityPerem {
	case 1:
		simplePerem := simpletypes.New(-1)

		perem.Title = simplePerem.Title
		perem.Type = simplePerem.Type
		perem.Value = simplePerem.Value
		perem.Actions = simplePerem.Actions
	default:
		simplePerem := simpletypes.New(-1)

		perem.Title = simplePerem.Title
		perem.Type = simplePerem.Type
		perem.Value = simplePerem.Value
		perem.Actions = simplePerem.Actions
	}

	return perem
}

func (p *Perem) ExecuteRandomActionWithPerem() string {
	actionId := gofakeit.IntN(len(p.Actions))
	action := p.Actions[actionId]

	//log.Println("Индекс действия на котором ступор: ", actionId)

	return action() + "\nINSERT"
}

func (p *Perem) RandomLatinLetter() string {
	letters := []string{"q", "w", "e", "r", "t", "y", "u", "i", "o", "p", "a", "s", "d", "f", "g", "h", "j", "k", "l", "z", "x", "c", "v", "b", "n", "m"}

	return letters[rand.Intn(len(letters))]
}
