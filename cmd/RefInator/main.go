package main

import (
	"log"
	"math/rand/v2"

	"github.com/Izumra/RefInatore/app/funcgen/swift"
	"github.com/Izumra/RefInatore/app/refinator"
	configparser "github.com/Izumra/RefInatore/utils/config_parser"
)

func main() {
	cfg := configparser.MustLoadByPath("config/settings.yaml")

	if cfg.Insertions == nil {
		cfg.Insertions = make([]string, 0)
	}

	swiftFuncGenerator := swift.NewFunction()
	for i := 0; i < cfg.CountFunctions; i++ {

		countActions := rand.IntN(cfg.MaxActionsPerFunc)
		function := swiftFuncGenerator.GenerateFilling(countActions)
		log.Println("\n***Сгенерирована новая функция***\n\n", function)

		err := swiftFuncGenerator.CheckFunction(function)
		if err != nil {
			panic(err)
		}

		cfg.Insertions = append(cfg.Insertions, function)

	}

	refinator := refinator.New(cfg)
	err := refinator.MakeFolderCopy(cfg.FolderPath)
	if err != nil {
		log.Fatal(err)
	}

	refinator.Refactor(cfg.FolderPath + "_copy")
}
