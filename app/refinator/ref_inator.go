package refinator

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"math/rand/v2"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

  configparser "github.com/Izumra/RefInatore/utils/config_parser"
	"github.com/brianvoe/gofakeit/v7"
)

var (
	ErrNotInsertions = errors.New("Нет текстов для вставки")
	ErrManyTries     = errors.New("Слишком много попыток выборки текста для ставки, попробуйте перейти к следующему файлу")

	regexpChange = regexp.MustCompile(`[^a-zA-Z0-9]`)
)

type Insertion struct {
	text        string
	fileRepeats int
	maxRepeats  int
}

type Changes struct {
	Classes    map[string]string `yaml:"classes"`
	Funcs      map[string]string `yaml:"funcs"`
	Enums      map[string]string `yaml:"enums"`
	Structs    map[string]string `yaml:"structs"`
	Extensions map[string]string `yaml:"extensions"`
}

type RefInator struct {
	excExts    map[string]bool
	excFiles   map[string]bool
	excFolders []string
	insertions []Insertion
	changes    Changes
	mainFile   string
}

func New(cfg configparser.Config) *RefInator {
	refInator := &RefInator{
		excExts:    make(map[string]bool),
		excFiles:   make(map[string]bool),
		excFolders: cfg.Exclusions.Folders,
	}

	for _, ext := range cfg.Exclusions.Extensions {
		refInator.excExts[ext] = true
	}

	for _, ext := range cfg.Exclusions.Files {
		refInator.excFiles[ext] = true
	}

	changes := Changes{
		Classes:    make(map[string]string),
		Funcs:      make(map[string]string),
		Enums:      make(map[string]string),
		Structs:    make(map[string]string),
		Extensions: make(map[string]string),
	}

	for _, class := range cfg.Changes.Classes {
		changes.Classes[class] = genRandomTitle()
	}

	for _, function := range cfg.Changes.Funcs {
		changes.Funcs[function] = genRandomTitle()
	}

	for _, enum := range cfg.Changes.Enums {
		changes.Enums[enum] = genRandomTitle()
	}

	for _, structure := range cfg.Changes.Structs {
		changes.Structs[structure] = genRandomTitle()
	}

	for _, extension := range cfg.Changes.Extensions {
		changes.Extensions[extension] = genRandomTitle()
	}

	refInator.changes = changes

	insertions := make([]Insertion, len(cfg.Insertions))
	for ind, code := range cfg.Insertions {
		insertion := Insertion{}
		insertion.text = strings.ReplaceAll(code, "\n", "\n\t") + "\n"

		insertions[ind] = insertion
	}

	refInator.insertions = insertions

	return refInator
}

func (r *RefInator) Refactor(folderPath string) error {
	regexp := regexp.MustCompile(`func .*\(.*\).*{`)
	idInsertion, errChooseInsertion := r.chooseRandomInsertion()
	return filepath.WalkDir(folderPath, func(path string, d fs.DirEntry, err error) error {
		for _, folder := range r.excFolders {
			if strings.HasPrefix(path, folder) {
				return nil
			}
		}

		if _, ok := r.excFiles[path]; ok {
			return nil
		}

		if !d.IsDir() {

			ext := filepath.Ext(path)
			if _, ok := r.excExts[ext]; ok {
				return nil
			}

			newFileName := genRandomTitle()
			newPath := strings.Replace(
				path,
				d.Name(),
				strings.Replace(d.Name(), ext, "", -1)+"_"+newFileName+ext,
				1,
			)
			os.Rename(path, newPath)
			path = newPath

			fileReader, err := os.Open(path)
			if err != nil {
				log.Println(err)
				return nil
			}
			defer fileReader.Close()

			lines := []string{}
			scanner := bufio.NewScanner(fileReader)

			for scanner.Scan() {
				unchanged_line := scanner.Text()
				line := r.changeNamesWorker(unchanged_line)
				if unchanged_line != line {
					log.Printf("\n\nИзменение!!!\nФайл: %s\nИзменилась строка под номером: %d\nИзначальная строка:\n%s\nИзмененная строка:\n%s\n\n", path, len(lines)+1, unchanged_line, line)
				}
				line += "\n"

				lines = append(lines, line)
			}
			if err := scanner.Err(); err != nil {
				log.Printf("where is the error: %s,\n cause: %s", path, err)
				return nil
			}
			fileReader.Close()
			copy(make([]string, len(lines)), lines)

			fileWriter, err := os.OpenFile(path, os.O_WRONLY, 0644)
			if err != nil {
				log.Println(err)
				return nil
			}
			defer fileWriter.Close()

			writer := bufio.NewWriter(fileWriter)

			for i := range lines {
				if errChooseInsertion == nil && i != 0 {
					if regexp.FindString(lines[i-1]) != "" {
						for _, funcForInsert := range r.changes.Funcs {
							if strings.Contains(lines[i-1], funcForInsert) {
								if _, err := writer.WriteString("\t" + r.insertions[idInsertion].text); err == nil {
									r.insertions[idInsertion].fileRepeats++
									r.insertions[idInsertion].maxRepeats++

									log.Printf("\n\nВставка!!!\nФайл: %s\nВставленный текст: \n%s\n", path, r.insertions[idInsertion].text)

									idInsertion, errChooseInsertion = r.chooseRandomInsertion()

								}
								break
							}
						}
					}
				}

				if strings.HasPrefix(lines[i], "@main") {
					r.mainFile = path
				}

				if _, err := writer.WriteString(lines[i]); err != nil {
					log.Println(err)
				}
			}
			writer.Flush()
			fileWriter.Close()

			for idx := range r.insertions {
				r.insertions[idx].fileRepeats = 0
			}

			if errors.Is(errChooseInsertion, ErrManyTries) {
				idInsertion, errChooseInsertion = r.chooseRandomInsertion()
			}
		}
		if r.mainFile != "" {
			err = r.postInsertingInMainFile(
				r.mainFile,
				idInsertion,
				errChooseInsertion,
				regexp,
			)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *RefInator) postInsertingInMainFile(
	path string,
	idInsertion int,
	errChooseInsertion error,
	regexp *regexp.Regexp,
) error {
	fileReader, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fileReader.Close()

	lines := []string{}
	scanner := bufio.NewScanner(fileReader)

	for scanner.Scan() {
		line := scanner.Text() + "\n"
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	fileReader.Close()
	copy(make([]string, len(lines)), lines)

	fileWriter, err := os.OpenFile(path, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer fileWriter.Close()

	writer := bufio.NewWriter(fileWriter)

	for i := range lines {
		if errChooseInsertion == nil && i != 0 {
			if regexp.FindString(lines[i-1]) != "" {
				for _, funcForInsert := range r.changes.Funcs {
					if strings.Contains(lines[i-1], funcForInsert) {
						if _, err := writer.WriteString("\t" + r.insertions[idInsertion].text); err == nil {
							r.insertions[idInsertion].fileRepeats++
							r.insertions[idInsertion].maxRepeats++

							log.Printf("\n\nВставка!!!\nФайл: %s\nВставленный текст: \n%s\n", path, r.insertions[idInsertion].text)

							idInsertion, errChooseInsertion = r.chooseRandomInsertion()

						}
						break
					}
				}
			}
		} else if errors.Is(errChooseInsertion, ErrManyTries) {
			idInsertion, errChooseInsertion = r.chooseRandomInsertion()
		}

		if _, err := writer.WriteString(lines[i]); err != nil {
			log.Println(err)
		}
	}
	writer.Flush()
	fileWriter.Close()

	for idx := range r.insertions {
		r.insertions[idx].fileRepeats = 0
	}

	return nil
}

func (r *RefInator) chooseRandomInsertion() (int, error) {
	tries := 0

	var insertion Insertion
	for len(r.insertions) != 0 {
		randomInsertion := rand.IntN(len(r.insertions))
		insertion = r.insertions[randomInsertion]

		rarestInsert := insertion

		for _, insert := range r.insertions {
			if insert.fileRepeats != 3 {
				if insert.maxRepeats < rarestInsert.maxRepeats {
					rarestInsert = insert
				}
			}
		}

		if insertion.maxRepeats == 7 {
			r.insertions = append(r.insertions[:randomInsertion], r.insertions[randomInsertion+1:]...)
			continue
		} else if insertion.fileRepeats == 3 {
			tries++
			if tries == 10 {
				return -1, ErrManyTries
			}

			continue
		} else if insertion.text != rarestInsert.text {
			continue
		}

		return randomInsertion, nil
	}

	return -1, ErrNotInsertions
}

func (r *RefInator) changeNamesWorker(line string) string {

	for class := range r.changes.Classes {
		if strings.Contains(line, class) {
			line = strings.ReplaceAll(line, class, r.changes.Classes[class])
		}
	}

	for function := range r.changes.Funcs {
		funcRegexp := regexp.MustCompile(function + `(\(| \()`)
		matches := funcRegexp.FindAllString(line, -1)
		if len(matches) != 0 {
			line = funcRegexp.ReplaceAllString(line, r.changes.Funcs[function]+"(")
		}
	}

	for enum := range r.changes.Enums {
		if strings.Contains(line, enum) {
			line = strings.ReplaceAll(line, enum, r.changes.Enums[enum])
		}
	}

	for structure := range r.changes.Structs {
		if strings.Contains(line, structure) {
			line = strings.ReplaceAll(line, structure, r.changes.Structs[structure])
		}
	}

	for extension := range r.changes.Extensions {
		if strings.Contains(line, extension) {
			line = strings.ReplaceAll(line, extension, r.changes.Structs[extension])
		}
	}

	return line
}

func (r *RefInator) MakeFolderCopy(folderPath string) error {
	copyPath := folderPath + "_copy"
	err := os.Mkdir(copyPath, 0755)
	if os.IsExist(err) {
		os.RemoveAll(copyPath)

		err := os.Mkdir(copyPath, 0755)
		if err != nil && !errors.Is(err, os.ErrPermission) {
			return err
		}
	}

	return filepath.WalkDir(folderPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil && !errors.Is(err, os.ErrPermission) {
			return err
		}
		replacedPath := strings.Replace(path, folderPath, folderPath+"_copy", 1)

		if !d.IsDir() {

			fileForReading, err := os.Open(path)
			if err != nil && !errors.Is(err, os.ErrPermission) {
				return err
			}
			defer fileForReading.Close()

			fileForWriting, err := os.Create(replacedPath)
			if err != nil && !errors.Is(err, os.ErrPermission) {
				return err
			}
			defer fileForWriting.Close()

			_, err = io.Copy(fileForWriting, fileForReading)
			if err != nil {
				return err
			}

		} else {
			info, err := d.Info()
			if err != nil {
				return err
			}

			if err = os.MkdirAll(replacedPath, info.Mode()); err != nil && !errors.Is(err, os.ErrPermission) {
				return err
			}
		}

		return nil
	})
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

func genRandomTitle() string {
	randTypeTitle := rand.IntN(5)

	switch randTypeTitle {
	case 0:
		return regexpChange.ReplaceAllString(fmt.Sprintf(
			"%s%s%s",
			gofakeit.Color(),
			capitalizeFirstLetter(gofakeit.Word()),
			capitalizeFirstLetter(gofakeit.VerbAction()),
		), "")
	case 1:
		return regexpChange.ReplaceAllString(fmt.Sprintf(
			"%s%s",
			gofakeit.Adjective(),
			capitalizeFirstLetter(gofakeit.Word()),
		), "")
	case 2:
		return regexpChange.ReplaceAllString(fmt.Sprintf(
			"%s%s",
			gofakeit.Word(),
			capitalizeFirstLetter(gofakeit.Word()),
		), "")
	case 3:
		return regexpChange.ReplaceAllString(fmt.Sprintf(
			"%s%s%s",
			gofakeit.AdverbPlace(),
			capitalizeFirstLetter(gofakeit.BeerName()),
			capitalizeFirstLetter(gofakeit.Bird()),
		), "")
	case 4:
		return regexpChange.ReplaceAllString(fmt.Sprintf(
			"%s%s",
			gofakeit.VerbAction(),
			capitalizeFirstLetter(gofakeit.Adjective()),
		), "")
	default:
		return regexpChange.ReplaceAllString(fmt.Sprintf(
			"%s%s%s",
			gofakeit.Color(),
			capitalizeFirstLetter(gofakeit.BookAuthor()),
			capitalizeFirstLetter(gofakeit.VerbAction()),
		), "")
	}
}
