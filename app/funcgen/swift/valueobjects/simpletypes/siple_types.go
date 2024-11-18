package simpletypes

import (
	mathrand "math/rand/v2"

	"github.com/Izumra/RefInatore/app/funcgen/swift/helpers"
	"github.com/Izumra/RefInatore/app/funcgen/swift/valueobjects/simpletypes/Int8"
)

const (
	Int8T = iota
	UInt8T
	Int16T
	UInt16T
	Int32T
	UInt32T
	Int64T
	UInt64T
	IntT
	UIntT
	FloatT
	DoubleT
	BoolT
	StringT
)

type SimplePerem struct {
	Title   string
	Value   string
	Type    string
	Actions []helpers.Action
	Helpers []helpers.Helper
}

//Generates new simple perem with one of presented types.
// If passed type param, create perem with passed type
/*
 * Int8
 * UInt8
 * Int16
 * UInt16
 * Int32
 * UInt32
 * Int64
 * UInt64
 * Int
 * UInt
 * Float
 * Double
 * Bool
 * String
 */
func New(t int) *SimplePerem {
	choosenType := mathrand.IntN(14)
	perem := &SimplePerem{}

	switch choosenType {
	default:

		int8P := Int8.New()
		perem.Type = int8P.Type
		perem.Title = int8P.Title
		perem.Value = int8P.Value
		perem.Actions = int8P.Actions
		perem.Helpers = int8P.Helpers

		// case UInt8T:
		//	randValue := mathrand.IntN(254)
		//	perem.Value = fmt.Sprintf("%d", randValue)

		// case Int16T:
		//	randValue := gofakeit.Int16()

		//	typeSign := mathrand.IntN(2)
		//	if typeSign == 0 && randValue != 0 {
		//		randValue = 0 - randValue
		//	}

		//	perem.Value = fmt.Sprintf("%d", randValue)

		// case UInt16T:
		//	randValue := mathrand.IntN(65534)

		//	perem.Value = fmt.Sprintf("%d", randValue)

		// case Int32T:
		//	randValue := gofakeit.Int32()

		//	typeSign := mathrand.IntN(2)
		//	if typeSign == 0 && randValue != 0 {
		//		randValue = 0 - randValue
		//	}

		//	perem.Value = fmt.Sprintf("%d", randValue)

		// case UInt32T:
		//	randValue := mathrand.Uint32()

		//	perem.Value = fmt.Sprintf("%d", randValue)
		// case Int64T:
		//	randValue := mathrand.Int64()

		//	typeSign := mathrand.IntN(2)
		//	if typeSign == 0 && randValue != 0 {
		//		randValue = 0 - randValue
		//	}

		//	perem.Value = fmt.Sprintf("%d", randValue)

		// case UInt64T:
		//	randValue := mathrand.Uint64()

		//	perem.Value = fmt.Sprintf("%d", randValue)

		// case IntT:
		//	randValue := mathrand.Int()

		//	typeSign := mathrand.IntN(2)
		//	if typeSign == 0 && randValue != 0 {
		//		randValue = 0 - randValue
		//	}

		//	perem.Value = fmt.Sprintf("%d", randValue)

		// case UIntT:
		//	randValue := mathrand.Uint64()

		//	perem.Value = fmt.Sprintf("%d", randValue)

		// case FloatT:
		//	randValue := gofakeit.Float32()

		//	typeSign := mathrand.IntN(2)
		//	if typeSign == 0 && randValue != 0 {
		//		randValue = 0 - randValue
		//	}

		//	perem.Value = fmt.Sprintf("%4.f", randValue)
		// case DoubleT:
		//	randValue := gofakeit.Float64()

		//	typeSign := mathrand.IntN(2)
		//	if typeSign == 0 && randValue != 0 {
		//		randValue = 0 - randValue
		//	}

		//	perem.Value = fmt.Sprintf("%4.f", randValue)
		// case BoolT:
		//	randValue := false

		//	typeBool := mathrand.IntN(2)
		//	if typeBool == 0 {
		//		randValue = true
		//	}

		//	perem.Value = fmt.Sprintf("%v", randValue)

		// case StringT:
		// randValue := ""

		//typeString := mathrand.IntN(5)
		//switch typeString {
		//case 0:
		//	randValue = strings.ReplaceAll(fmt.Sprintf(
		//		"%s%s%s",
		//    helpers.CapitalizeFirstLetter(gofakeit.VerbAction()),
		//		gofakeit.City(),
		//		helpers.CapitalizeFirstLetter(gofakeit.Animal()),
		//	), " ", "")
		//case 1:
		//	randValue = strings.ReplaceAll(fmt.Sprintf(
		//		"%s%s%s",
		//		gofakeit.Adjective(),
		//		gofakeit.BookTitle(),
		//		gofakeit.Username(),
		//	), " ", "")
		//case 2:
		//	randValue = strings.ReplaceAll(fmt.Sprintf(
		//		"%s%s%s",
		//		helpers.CapitalizeFirstLetter(gofakeit.Animal()),
		//		gofakeit.Dog(),
		//		gofakeit.Word(),
		//	), " ", "")
		//case 3:
		//	randValue = strings.ReplaceAll(fmt.Sprintf(
		//		"%s%s%s",
		//		gofakeit.Word(),
		//		gofakeit.BeerAlcohol(),
		//    helpers.CapitalizeFirstLetter(gofakeit.Cat()),
		//	), " ", "")
		//case 4:
		//	randValue = strings.ReplaceAll(fmt.Sprintf(
		//		"%s%s%s",
		//		gofakeit.Color(),
		//		gofakeit.VerbAction(),
		//		helpers.CapitalizeFirstLetter(gofakeit.CelebrityActor()),
		//	), " ", "")
		//}
	}

	// log.Printf("Сгенерированная переменная значение: %v, тип: %v", perem.Value, perem.Type)
	return perem
}
