package Int8

import (
	"fmt"
	"math/rand/v2"
	mathrand "math/rand/v2"

	"github.com/Izumra/RefInatore/app/funcgen/swift/helpers"
)

type Perem struct {
	Title, TypeP, Value string
	Actions             []helpers.Action
}

func randomValue() (string, int8) {
	randValue := mathrand.IntN(127)

	typeSign := mathrand.IntN(2)
	if typeSign == 0 && randValue != 0 {
		randValue = 0 - randValue
	}

	return fmt.Sprintf("%d", randValue), int8(randValue)
}

func New() *Perem {
	p := &Perem{
		Title: helpers.GenRandomPeremTitle(),
		TypeP: "Int8",
	}

	p.Value, _ = randomValue()

	p.Actions = []helpers.Action{
		p.ariphmetic,
		p.assignment,
		p.condition,
		p.bitwise,
		p.cycle,
	}

	return p
}

func (p *Perem) ariphmetic() string {
	countOps := rand.IntN(8)
	if countOps == 0 {
		countOps = 1
	}

	attrs := make([]string, countOps)
	for i := 0; i < len(attrs); i++ {

		if helpers.YesOrNo() {

			newValue, _ := randomValue()
			if newValue[0] == '-' {
				newValue = "(" + newValue + ")"
			}

			attrs[i] = newValue
		} else {
			attrs[i] = p.Title
		}
	}

	ariphmeticStr := p.Title + " = " + helpers.GenRandomAriphmeticStr(attrs)
	if helpers.YesOrNo() {
		ariphmeticStr += fmt.Sprintf("\nprint(\"Ariphmetic operation is equal - \\(%v)\")\n", p.Title)
	}

	return ariphmeticStr
}

func (p *Perem) assignment() string {

	randValue, _ := randomValue()

	action := fmt.Sprintf("%s = %s", p.Title, randValue)

	if helpers.YesOrNo() {
		action += ";"
	}

	return action
}

func (p *Perem) condition() string {
	attrs := make([]string, 10)
	attrs[0] = p.Title
	for i := 1; i < len(attrs); i++ {
		attrs[i], _ = randomValue()
	}

	if helpers.YesOrNo() {
		pattern, countConditions := helpers.PatternForIfElseConditions(5)
		conditions := make([]any, countConditions)

		for i := 0; i < countConditions; i++ {
			conditions[i] = helpers.MakeCondition(attrs)
		}

		return fmt.Sprintf(pattern, conditions...)
	}

	pattern, countValues := helpers.PatternForSwitchConditions(p.Title, 5)
	values := make([]any, countValues)
	for i := 0; i < countValues; i++ {
		values[i], _ = randomValue()
	}

	pattern = fmt.Sprintf(pattern, values...)

	return pattern
}

func (p *Perem) bitwise() string {
	countOps := rand.IntN(4)
	if countOps == 0 {
		countOps = 1
	}

	attrs := make([]string, countOps)
	for i := 0; i < len(attrs); i++ {

		if helpers.YesOrNo() {

			newValue, _ := randomValue()
			if newValue[0] == '-' {
				newValue = "(" + newValue + ")"
			}

			attrs[i] = newValue
		} else {
			attrs[i] = p.Title
		}
	}

	bitwiseStr := p.Title + " = " + helpers.GenRandomBitwiseStr(attrs)
	if helpers.YesOrNo() {
		bitwiseStr += fmt.Sprintf("\nprint(\"Bitwise operation is equal - \\(%v)\")\n", p.Title)
	}

	return bitwiseStr
}

func (p *Perem) cycle() string {
	typeCycle := rand.IntN(2)

	returnedStr := ""

	firstV, firstN := randomValue()
	secondV, secondN := randomValue()
	if secondN < firstN {
		temp := secondV
		secondV = firstV
		firstV = temp
	}

	switch typeCycle {
	case 0:
		returnedStr = fmt.Sprintf("for %v: %v in %v...%v {\n\tINSERT\n}",
			helpers.RandomLatinLetter(),
			"Int8",
			firstV,
			secondV,
		)
	case 1:
		returnedStr = fmt.Sprintf("while %v %v %v {\n\tINSERT\n}",
			p.Title,
			helpers.SelectConditionSign(),
			secondV,
		)
	case 2:
		returnedStr = fmt.Sprintf("repeat {\n\tINSERT\n} while %v %v %v",
			p.Title,
			helpers.SelectConditionSign(),
			secondV,
		)
	}

	return returnedStr
}
