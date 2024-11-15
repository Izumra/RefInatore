package swift

import (
	"fmt"
	"math/rand/v2"
	"regexp"
	"strings"
	"sync"

	"github.com/Izumra/RefInatore/app/funcgen/swift/helpers"
	"github.com/Izumra/RefInatore/app/funcgen/swift/models"
)

type Function struct {
	counterInsertions sync.WaitGroup
	locker            sync.Mutex
	value             string
	insertPattern     *regexp.Regexp
	maxConditions     int
	maxCycles         int

	Stack []*models.Perem
}

// Create and initialize struct for Swift function
func NewFunction() *Function {
	return &Function{
		insertPattern: regexp.MustCompile(`INSERT`),
		maxConditions: 2,
		maxCycles:     1,
		Stack:         make([]*models.Perem, 0),
	}
}

func (f *Function) GenerateFilling(countOperations uint) string {
	f.value = f.chooseTypeFunction()

	if helpers.YesOrNo() || len(f.Stack) == 0 {
		f.addNewPerems(uint(10))
	}

	for i := 0; i < int(countOperations); i++ {
		f.counterInsertions.Add(1)

		go func() {
			defer f.counterInsertions.Done()

			f.locker.Lock()

			peremId := rand.IntN(len(f.Stack))

			action := f.Stack[peremId].ExecuteRandomActionWithPerem()
			withoutTabs := strings.TrimLeft(action, "\t")
			isCondition := strings.HasPrefix(withoutTabs, "if ") || strings.HasPrefix(withoutTabs, "switch ")
			isCycle := strings.HasPrefix(withoutTabs, "for ") || strings.HasPrefix(withoutTabs, "while ") || strings.HasPrefix(withoutTabs, "repeat {")

			if (isCondition && f.maxConditions <= 0) || (isCycle && f.maxCycles <= 0) {

				for {
					peremId = rand.IntN(len(f.Stack))
					action = f.Stack[peremId].ExecuteRandomActionWithPerem()

					withoutTabs = strings.TrimLeft(action, "\t")
					isCondition = strings.HasPrefix(withoutTabs, "if ") || strings.HasPrefix(withoutTabs, "switch ")
					isCycle = strings.HasPrefix(withoutTabs, "for ") || strings.HasPrefix(withoutTabs, "while ") || strings.HasPrefix(withoutTabs, "repeat {")
					if !isCondition && !isCycle {
						break
					}

				}
			}

			results := f.insertPattern.FindAllStringSubmatchIndex(f.value, -1)
			if len(results) != 0 {

				spotForInsert := results[rand.IntN(len(results))]

				countTabs := helpers.CountTabsInString(f.value, spotForInsert[0])
				unTabbedStr := strings.Split(action, "\n")

				for i := 1; i < len(unTabbedStr); i++ {
					unTabbedStr[i] = strings.Repeat("\t", countTabs) + unTabbedStr[i]
				}

				insertedStr := strings.Join(unTabbedStr, "\n")

				f.value = f.value[:spotForInsert[0]] + insertedStr + f.value[spotForInsert[1]:]

				if isCondition && f.maxConditions > 0 {
					f.maxConditions--
				} else if isCycle && f.maxCycles > 0 {
					f.maxCycles--
				}
			}

			f.locker.Unlock()

		}()
	}

	f.counterInsertions.Wait()

	f.value = strings.ReplaceAll(f.value, "INSERT", "")

	return f.value
}

func (f *Function) addNewPerems(maxPerems uint) {
	countPerems := rand.UintN(maxPerems)
	if countPerems == 0 {
		countPerems = 1
	}

	addedPerems := make([]*models.Perem, countPerems)

	f.locker.Lock()

	for i := 0; i < int(countPerems); i++ {
		addedPerems[i] = models.NewPerem()

		typeInitialize := rand.IntN(5)
		typePerem := rand.IntN(2)

		initializeStr := fmt.Sprintf("var %v", addedPerems[i].Title)
		if typePerem == 1 {
			initializeStr = fmt.Sprintf("let %v", addedPerems[i].Title)
		}

		switch typeInitialize {
		case 0:
			initializeStr += fmt.Sprintf(
				" = %v",
				addedPerems[i].Value,
			)
		case 1:
			if typePerem == 0 {
				initializeStr += fmt.Sprintf(
					": %v\n%v = %v",
					addedPerems[i].Type,
					addedPerems[i].Title,
					addedPerems[i].Value,
				)
			} else {
				initializeStr += fmt.Sprintf(
					": %v = %v",
					addedPerems[i].Type,
					addedPerems[i].Value,
				)
			}
		case 2:
			initializeStr += fmt.Sprintf(
				": %v? = %v",
				addedPerems[i].Type,
				addedPerems[i].Value,
			)
		case 3:
			initializeStr += fmt.Sprintf(
				": () -> %v = {\n\treturn %v\n}",
				addedPerems[i].Type,
				addedPerems[i].Value,
			)
		case 4:
			if typePerem == 0 {
				initializeStr += fmt.Sprintf(
					": %v {\n\tget {\n\t\treturn %v\n\t}\n\tset {\n\t\tprint(%v)\n\t}\n}",
					addedPerems[i].Type,
					addedPerems[i].Value,
					addedPerems[i].Title,
				)
			} else {
				initializeStr += fmt.Sprintf(
					": %v {\n\tget {\n\t\treturn %v\n\t}\n}",
					addedPerems[i].Type,
					addedPerems[i].Value,
				)
			}
		}
		initializeStr += "\nINSERT"

		if i != 0 {
			results := f.insertPattern.FindAllStringSubmatchIndex(f.value, -1)
			if len(results) != 0 {

				spotForInsert := results[rand.IntN(len(results))]

				countTabs := helpers.CountTabsInString(f.value, spotForInsert[0])

				unTabbedStr := strings.Split(initializeStr, "\n")

				for i := 1; i < len(unTabbedStr); i++ {
					unTabbedStr[i] = strings.Repeat("\t", countTabs) + unTabbedStr[i]
				}

				insertedStr := strings.Join(unTabbedStr, "\n")

				f.value = f.value[:spotForInsert[0]] + insertedStr + f.value[spotForInsert[1]:]
			}
		}

	}

	f.Stack = append(f.Stack, addedPerems...)

	f.locker.Unlock()

}
