package renamer

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"math/rand/v2"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"unicode"

	"github.com/brianvoe/gofakeit/v7"
)

var regexpChange = regexp.MustCompile(`[^a-zA-Z0-9]`)

type Info struct {
	Author  string `json:"author"`
	Version int    `json:"version"`
}

type Image struct {
	Filename string `json:"filename"`
	Idiom    string `json:"idiom"`
	Scale    string `json:"scale"`
}

type Contents struct {
	Images []Image `json:"images"`
	Info   Info    `json:"info"`
}

type RenameGroup struct {
	comparisons map[string]string
	images      []string
	content     string
}

func (rg *RenameGroup) UpdateContent(wg *sync.WaitGroup, errors chan<- error, parentFolder string) {
	defer wg.Done()

	suffix := genRandomTitle()

	file, err := os.Open(parentFolder + rg.content)
	if err != nil {
		errors <- fmt.Errorf("Occured and error while opening the 'Contents' file - %v", err)
		return
	}
	defer file.Close()

	var contents Contents
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&contents)
	if err != nil {
		log.Println(parentFolder, "Eba - ", rg.content)
		errors <- fmt.Errorf("Occured and error while parsing the 'Contents' file into struct - %v", err)
		return
	}

	for _, title := range rg.images {
		ext := filepath.Ext(parentFolder + title)
		onlyName := title[:len(title)-len(ext)]

		namePices := strings.Split(onlyName, "@")
		namePices[0] = namePices[0] + "_" + suffix
		newTitle := strings.Join(namePices, "@") + ext

		err := os.Rename(parentFolder+title, parentFolder+newTitle)
		if err != nil {
			errors <- fmt.Errorf("Occured an error while renaming the image - %v", err)
			return
		}

		rg.comparisons[title] = newTitle
	}

	for ind := range contents.Images {
		oldTitle := contents.Images[ind].Filename
		contents.Images[ind].Filename = rg.comparisons[oldTitle]
	}

	data, err := json.MarshalIndent(contents, "", "  ")
	if err != nil {
		errors <- fmt.Errorf("Occured and error while encoding the 'Contens' data - %v", err)
		return
	}

	err = os.WriteFile(parentFolder+rg.content, data, 0644)
	if err != nil {
		errors <- fmt.Errorf("Occured and error while changing content data - %v", err)
		return
	}
}

func NewRenameGroup() *RenameGroup {
	return &RenameGroup{
		images:      make([]string, 0),
		comparisons: make(map[string]string),
	}
}

type Service struct {
	folderPath   string
	renameGroups map[string]*RenameGroup
}

func New(folderPath string) *Service {
	return &Service{
		folderPath:   folderPath,
		renameGroups: make(map[string]*RenameGroup),
	}
}

func (s *Service) Rename() error {
	wg := sync.WaitGroup{}
	chanErr := make(chan error)

	err := filepath.WalkDir(s.folderPath, func(path string, d fs.DirEntry, err error) error {
		name := d.Name()
		parentPath := path[:len(path)-len(name)]
		if strings.HasSuffix(parentPath, ".imageset/") {
			if _, ok := s.renameGroups[parentPath]; !ok {
				s.renameGroups[parentPath] = NewRenameGroup()
			}

			if name == "Contents.json" {
				s.renameGroups[parentPath].content = name
			} else {
				s.renameGroups[parentPath].images = append(
					s.renameGroups[parentPath].images,
					name,
				)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	for title, group := range s.renameGroups {
		wg.Add(1)
		go group.UpdateContent(&wg, chanErr, title)
	}

	go func() {
		wg.Wait()
		close(chanErr)
	}()

	for errFromRenamer := range chanErr {
		return errFromRenamer
	}

	return nil
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
			gofakeit.Word(),
			capitalizeFirstLetter(gofakeit.Word()),
		), "")
	case 2:
		return regexpChange.ReplaceAllString(fmt.Sprintf(
			"%s%s%s",
			gofakeit.AdverbPlace(),
			capitalizeFirstLetter(gofakeit.BeerName()),
			capitalizeFirstLetter(gofakeit.Bird()),
		), "")
	case 3:
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
