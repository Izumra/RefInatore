package models

import (
	"math/rand"
	"strings"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/Izumra/RefInatore/app/funcgen/swift/helpers"
	simpletypes "github.com/Izumra/RefInatore/app/funcgen/swift/valueobjects/simpletypes"
)

type Perem struct {
	Title   string
	Type    string
	Value   string
	Actions []helpers.Action
	Helpers []helpers.Helper
}

func NewPerem() *Perem {
	perem := &Perem{}

	typeComplexityPerem := rand.Intn(2)

	switch typeComplexityPerem {
	default:
		simplePerem := simpletypes.New(-1)

		perem.Title = simplePerem.Title
		perem.Type = simplePerem.Type
		perem.Value = simplePerem.Value
		perem.Actions = simplePerem.Actions
		perem.Helpers = simplePerem.Helpers
	}

	return perem
}

func (p *Perem) ExecuteRandomActionWithPerem() string {
	actionId := gofakeit.IntN(len(p.Actions))
	action := p.Actions[actionId]

	// log.Println("Индекс действия на котором ступор: ", actionId)

	return action() + "\nINSERT"
}

func (p *Perem) ReplaceIsNilTypeSign() {
	if strings.HasSuffix(p.Type, "?") {
		p.Type = strings.ReplaceAll(p.Type, "?", "")
	}
}
