package swift

import (
	"fmt"
	"math/rand/v2"
	"os"
	"os/exec"
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
	Attrs []*models.Perem
}

// Create and initialize struct for Swift function
func NewFunction() *Function {
	return &Function{
		insertPattern: regexp.MustCompile(`INSERT`),
		maxConditions: 2,
		maxCycles:     1,
		Stack:         []*models.Perem{},
	}
}

func (f *Function) CheckFunction(function string) error {
	temp, err := os.CreateTemp("", "swift_func_*.swift")
	if err != nil {
		return err
	}
	defer os.Remove(temp.Name())

	_, err = temp.WriteString(function)
	if err != nil {
		return err
	}
	temp.Close()

	cmd := exec.Command("swiftc", "-typecheck", temp.Name())
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("\n***Ошибка при проверке синтаксиса функции***\n%v", string(output))
	}

	return nil
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
			isCondition := strings.HasPrefix(withoutTabs, "if ") ||
				strings.HasPrefix(withoutTabs, "switch ")
			isCycle := strings.HasPrefix(withoutTabs, "for ") ||
				strings.HasPrefix(withoutTabs, "while ") ||
				strings.HasPrefix(withoutTabs, "repeat {")

			if (isCondition && f.maxConditions <= 0) || (isCycle && f.maxCycles <= 0) {
				for {
					peremId = rand.IntN(len(f.Stack))
					action = f.Stack[peremId].ExecuteRandomActionWithPerem()

					withoutTabs = strings.TrimLeft(action, "\t")
					isCondition = strings.HasPrefix(withoutTabs, "if ") ||
						strings.HasPrefix(withoutTabs, "switch ")
					isCycle = strings.HasPrefix(withoutTabs, "for ") ||
						strings.HasPrefix(withoutTabs, "while ") ||
						strings.HasPrefix(withoutTabs, "repeat {")
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

	f.Stack = []*models.Perem{}
	f.maxConditions = 2
	f.maxCycles = 1

	return f.value
}

func (f *Function) addNewPerems(maxPerems uint) {
	countPerems := rand.UintN(maxPerems)
	if countPerems == 0 {
		countPerems = 1
	}

	movedToOnlyRead := []string{}
	addedPerems := make([]*models.Perem, int(countPerems))

	f.locker.Lock()

	for i := 0; i < int(countPerems); i++ {
		perem := models.NewPerem()
		addedPerems[i] = perem

		typeInitialize := rand.IntN(5)

		initializeStr := fmt.Sprintf("var %v", perem.Title)

		switch typeInitialize {
		case 0:
			perem.ReplaceIsNilTypeSign()
			perem.Helpers[0]()

			initializeStr += fmt.Sprintf(
				" = %v(%v)",
				perem.Type,
				perem.Value,
			)
		case 1:
			if helpers.YesOrNo() {
				initializeStr += fmt.Sprintf(
					": %v\n%v = %v",
					perem.Type,
					perem.Title,
					perem.Value,
				)
			} else {
				initializeStr += fmt.Sprintf(
					": %v = %v",
					perem.Type,
					perem.Value,
				)
			}
		case 2:
			initializeStr += fmt.Sprintf(
				": %v = %v",
				perem.Type,
				perem.Value,
			)
		case 3:
			perem.ReplaceIsNilTypeSign()
			perem.Helpers[0]()
			defaultType := perem.Type
			perem.Type = " () -> " + perem.Type

			initializeStr += fmt.Sprintf(
				":%v = {\n\treturn %v(%v)\n}",
				perem.Type,
				defaultType,
				perem.Value,
			)

			movedToOnlyRead = append(movedToOnlyRead, fmt.Sprintf("%d", i))
		case 4:
			perem.ReplaceIsNilTypeSign()
			perem.Helpers[0]()

			if helpers.YesOrNo() {
				initializeStr += fmt.Sprintf(
					": %v {\n\tget {\n\t\treturn %v(%v)\n\t}\n\tset {\n\t\tprint(%v)\n\t}\n}",
					perem.Type,
					perem.Type,
					perem.Value,
					perem.Title,
				)
			} else {
				initializeStr += fmt.Sprintf(
					": %v {\n\tget {\n\t\treturn %v(%v)\n\t}\n}",
					perem.Type,
					perem.Type,
					perem.Value,
				)

				movedToOnlyRead = append(movedToOnlyRead, fmt.Sprintf("%d", i))
			}
		}
		initializeStr += "\nINSERT"

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

	perems := []*models.Perem{}
	keys := strings.Join(movedToOnlyRead, ", ")
	for i, v := range addedPerems {
		if strings.Contains(keys, fmt.Sprintf("%d", i)) {
			continue
		}
		perems = append(perems, v)
	}

	if len(perems) == 0 {
		f.locker.Unlock()
		f.addNewPerems(maxPerems)
		return
	}

	f.Stack = append(f.Stack, perems...)

	f.locker.Unlock()
}
