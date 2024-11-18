package Int8

import (
	"fmt"
	"math/rand/v2"
	mathrand "math/rand/v2"
	"strings"
	"sync"

	"github.com/Izumra/RefInatore/app/funcgen/swift/helpers"
)

type Perem struct {
	Title, Type, Value string
	Actions            []helpers.Action
	Helpers            []helpers.Helper
	isClosure          bool
	isNilChecker       sync.Once
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
		Type:  "Int8",
	}

	p.Value, _ = randomValue()
	if helpers.YesOrNo() {
		p.Type += "?"
	}

	p.Actions = []helpers.Action{
		func() string { return p.ariphmetic() },
		func() string { return p.assignment() },
		func() string { return p.condition() },
		func() string { return p.bitwise() },
		func() string { return p.cycle() },
	}
	p.Helpers = []helpers.Helper{
		p.unOptionType,
	}

	return p
}

func (p *Perem) unOptionType() {
	p.Type = strings.ReplaceAll(p.Type, "?", "")
}

func (p *Perem) assertion(assertStr string) string {
	var str string
	if p.isClosure {
		title := p.Title[:len(p.Title)-2]
		latinLet := helpers.RandomLatinLetter()

		if helpers.YesOrNo() {
			typeP := strings.ReplaceAll(p.Type, "() ->", "")
			firstV, firstN := randomValue()
			secondV, secondN := randomValue()
			if secondN < firstN {
				temp := secondV
				secondV = firstV
				firstV = temp
			}

			str = fmt.Sprintf(
				"%v = {\n\treturn %v.random(in: %v...%v\n}",
				title,
				typeP,
				firstV,
				secondV,
			)
		} else {
			str = fmt.Sprintf("%v = {\n\t%v = %v\n\treturn %v\n}", title, latinLet, assertStr, latinLet)
		}
	} else {
		str = p.Title + " = " + assertStr
	}

	return str
}

func (p *Perem) genRandomAttr() string {
	newValue, _ := randomValue()
	typeP := strings.ReplaceAll(p.Type, "?", "")

	if helpers.YesOrNo() {

		newValue = "(" + newValue + ")"
		return typeP + newValue

	} else {

		attr := p.Title

		if strings.HasSuffix(p.Type, "?") {
			switch mathrand.IntN(2) {
			case 0:
				attr = p.Title + "!"
			case 1:
				attr = fmt.Sprintf("(%v ?? %v(%v))", p.Title, typeP, newValue)
			}
		}

		return attr
	}
}

func (p *Perem) checkIsClosure() {
	if strings.HasPrefix(p.Type, " () ->") {
		p.Title += "()"
		p.isClosure = true
	}
}

func (p *Perem) checkNil(def string) string {
	var str string

	p.isNilChecker.Do(func() {
		if strings.Contains(p.Type, "?") {
			str += fmt.Sprintf("if (%v != nil) {\n\t", p.Title)

			countTabs := 1
			unTabbedStr := strings.Split(def, "\n")

			for i := 1; i < len(unTabbedStr); i++ {
				unTabbedStr[i] = strings.Repeat("\t", countTabs) + unTabbedStr[i]
			}

			insertedStr := strings.Join(unTabbedStr, "\n")

			str += insertedStr + "\n}"
		}
	})

	if str == "" {
		return def
	}

	return str
}

func (p *Perem) ariphmetic() string {
	p.checkIsClosure()

	var ariphmeticStr string

	countOps := rand.IntN(8)
	if countOps == 0 {
		countOps = 1
	}

	attrs := make([]string, countOps)
	for i := 0; i < len(attrs); i++ {
		attrs[i] = p.genRandomAttr()
	}

	ariphStr := helpers.GenRandomAriphmeticStr(attrs)
	if strings.Trim(ariphStr, " ") == p.Title {
		return p.ariphmetic()
	}

	str := p.assertion(ariphStr)

	if !p.isClosure {
		ariphmeticStr = p.checkNil(str)
	}

	return ariphmeticStr
}

func (p *Perem) assignment() string {
	p.checkIsClosure()

	var action string

	randValue, _ := randomValue()

	str := p.assertion(randValue)

	if !p.isClosure {
		action = p.checkNil(str)
	}

	return action
}

func (p *Perem) condition() string {
	p.checkIsClosure()

	attrs := make([]string, 10)
	for i := 0; i < len(attrs); i++ {
		attrs[i] = p.genRandomAttr()
	}

	if helpers.YesOrNo() {
		pattern, countConditions := helpers.PatternForIfElseConditions(5)
		conditions := make([]any, countConditions)

		for i := 0; i < countConditions; i++ {
			conditions[i] = helpers.MakeCondition(attrs)
		}

		return fmt.Sprintf(pattern, conditions...)
	}

	pattern, countValues := helpers.PatternForSwitchConditions(p.Title, 3)
	values := make([]any, countValues)
	for i := 0; i < countValues; i++ {
		values[i], _ = randomValue()
	}

	pattern = fmt.Sprintf(pattern, values...)

	return pattern
}

func (p *Perem) bitwise() string {
	p.checkIsClosure()

	countOps := rand.IntN(4)
	if countOps == 0 {
		countOps = 1
	}

	attrs := make([]string, countOps)
	for i := 0; i < len(attrs); i++ {
		attrs[i] = p.genRandomAttr()
	}

	var bitwiseStr string
	btwStr := helpers.GenRandomBitwiseStr(attrs)
	if strings.Trim(btwStr, " ") == p.Title {
		return p.bitwise()
	}

	str := p.assertion(btwStr)

	if !p.isClosure {
		bitwiseStr = p.checkNil(str)
	}

	return bitwiseStr
}

func (p *Perem) cycle() string {
	p.checkIsClosure()

	typeCycle := rand.IntN(2)

	returnedStr := ""

	firstV, firstN := randomValue()
	secondV, secondN := randomValue()
	if firstN < 0 {
		firstV = "(" + firstV + ")"
	}
	if secondN < 0 {
		secondV = "(" + secondV + ")"
	}

	if secondN < firstN {
		temp := secondV
		secondV = firstV
		firstV = temp
	}

	if p.isClosure {
		typeCycle = 0
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
			p.genRandomAttr(),
			helpers.SelectConditionSign(),
			secondV,
		)
	case 2:
		returnedStr = fmt.Sprintf("repeat {\n\tINSERT\n} while %v %v %v",
			p.genRandomAttr(),
			helpers.SelectConditionSign(),
			secondV,
		)
	}

	returnedStr = p.checkNil(returnedStr)

	return returnedStr
}
