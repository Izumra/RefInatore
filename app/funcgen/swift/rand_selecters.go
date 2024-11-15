package swift

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/Izumra/RefInatore/app/funcgen/swift/helpers"
	"github.com/Izumra/RefInatore/app/funcgen/swift/models"
	"github.com/Izumra/RefInatore/app/funcgen/swift/valueobjects/typesfunction"
	"github.com/brianvoe/gofakeit/v7"
)

// Returns generated pattern and ID of the position for further insertions
func (f *Function) chooseTypeFunction() string {
	typeFunction := rand.Intn(6)
	var funcPattern string

	titleFunc := strings.ReplaceAll(fmt.Sprintf(
		"%s%s%s",
		gofakeit.City(),
		helpers.CapitalizeFirstLetter(gofakeit.Animal()),
		helpers.CapitalizeFirstLetter(gofakeit.VerbAction()),
	), " ", "")

	switch typeFunction {
	case typesfunction.FuncWithOutputParams:
		outputs, returnedParams := randomOutputParams()
		returnsParams := formatReturnsParams(returnedParams)

		funcPattern = fmt.Sprintf("func %s () -> %s {\n\tINSERT\n\treturn %s\n}",
			titleFunc,
			outputs,
			returnsParams,
		)
	case typesfunction.FuncWithAllParams:
		outputs, returnedParams := randomOutputParams()
		returnsParams := formatReturnsParams(returnedParams)

		funcPattern = fmt.Sprintf("func %s (%s) -> %s {\n\tINSERT\n\treturn %s\n}",
			titleFunc,
			f.randomInputParams(false, true),
			outputs,
			returnsParams,
		)

	case typesfunction.FuncWithoutParams:

		funcPattern = fmt.Sprintf("func %s () %s {\n\tINSERT\n}",
			titleFunc,
			randomEmptyOutputParam(),
		)

	case typesfunction.FuncWithInputParams,
		typesfunction.FuncWithDefaultInputParams:

		withDefault := false
		if typeFunction == typesfunction.FuncWithDefaultInputParams {
			withDefault = true
		}

		funcPattern = fmt.Sprintf("func %s (%s)%s {\n\tINSERT\n}",
			titleFunc,
			f.randomInputParams(withDefault, true),
			randomEmptyOutputParam(),
		)

	//TODO доделать рекурсивный вызов функций
	case typesfunction.FuncRecursive:
		outputs, returnedParams := randomOutputParams()
    _ = formatReturnsParams(returnedParams)

		funcPattern = fmt.Sprintf("func %s (%s) -> %s {\n\tINSERT\n\treturn %s(%s)\n}",
			titleFunc,
			f.randomInputParams(false, false),
			outputs,
			titleFunc,
			"TODO: 'Вставить новые параметры того же типа что и входные'",
		)
	}

	//log.Println(typeFunction)

	return funcPattern
}

func formatReturnsParams(returnedParams []*models.Perem) string {

	preparedReturnedParams := make([]string, len(returnedParams))
	for ind, param := range returnedParams {
		if param.Value == "" {
			preparedReturnedParams[ind] = "nil"
			continue
		}

		//TODO: random math expression to generate returned param
		preparedReturnedParams[ind] = param.Value
	}

	formattedReturns := strings.Join(preparedReturnedParams, ", ")

	if len(returnedParams) > 1 {
		formattedReturns = "(" + formattedReturns + ")"
	}

	return formattedReturns
}

func (f *Function)randomInputParams(withDefault, withLoweredSep bool) string {
	countInputParams := rand.Intn(3)
	countInputParams++

	generatedParams := make([]string, countInputParams)

	for i := range countInputParams {
		loweredSepRand := rand.Intn(2)
		withDefaultRand := rand.Intn(2)

		withDefaultValue := false
		if withDefault && withDefaultRand == 0 {
			withDefaultValue = true
		}

		loweredSep := false
		if loweredSepRand == 1 && withLoweredSep {
			loweredSep = true
		}

		generatedParams[i] = f.genRandomInputParam(loweredSep, withDefaultValue)
	}

	return strings.Join(generatedParams, ", ")
}

func randomOutputParams() (string, []*models.Perem) {
	countOutputParams := rand.Intn(2)
	countOutputParams++

	generatedParams := make([]string, countOutputParams)
	generatedParamsStructs := make([]*models.Perem, countOutputParams)

	if countOutputParams > 1 {
		for i := range countOutputParams {
			generatedParams[i], generatedParamsStructs[i] = genRandomOutputParam(true)
		}
	} else {
		generatedParams[0], generatedParamsStructs[0] = genRandomOutputParam(false)
	}

	paramsStr := strings.Join(generatedParams, ", ")

	if len(generatedParams) > 1 {
		paramsStr = "(" + paramsStr + ")"
	}

	return paramsStr, generatedParamsStructs
}

// Generates random output param and return it's struct for further manipulating(return value from the func)
func genRandomOutputParam(specifedTitle bool) (string, *models.Perem) {
	perem := models.NewPerem()
	optionValueRand := rand.Intn(2)

	outputParam := perem.Type
	if optionValueRand == 1 {
		outputParam += "?"
		perem.Value = ""
	}

	if specifedTitle {
		outputParam = perem.Title + ": " + outputParam
	}

	return outputParam, perem
}

func (f *Function) genRandomInputParam(withLowered, withDefault bool) string {
	prefix := ""
	if withLowered {
		prefix = "_ "
	}

	perem := models.NewPerem()

  f.locker.Lock()
  f.Stack = append(f.Stack, perem)
  f.locker.Unlock()

	randomParam := fmt.Sprintf("%s%s: %s",
		prefix,
		perem.Title,
		perem.Type,
	)
	if withDefault {
		randomParam = fmt.Sprintf("%s%s: %s = %v",
			prefix,
			perem.Title,
      perem.Type,
			perem.Value,
		)
	}

	return randomParam
}

func randomEmptyOutputParam() string {
	typeFuncWithoutParams := rand.Intn(3)
	typeOutput := ""
	switch typeFuncWithoutParams {
	case 0:
		typeOutput = ""
	case 1:
		typeOutput = " -> Void"
	case 2:
		typeOutput = " -> ()"
	}

	return typeOutput
}

