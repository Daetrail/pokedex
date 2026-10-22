package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Daetrail/pokedex/internal/config"
	"github.com/Daetrail/pokedex/internal/pokecache"
	"github.com/Daetrail/pokedex/internal/util"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	cfg := &config.Config{
		BaseURL:          "https://pokeapi.co/api/v2",
		Commands:         config.GetCommands(),
		MapCommandConfig: config.MapCommandConfig{},
		PokeCache:        pokecache.NewPokeCache(10),
	}

	for {
		fmt.Print("Pokedex > ")

		scanner.Scan()
		if err := scanner.Err(); err != nil {
			fmt.Println("error reading input:", err)
			continue
		}

		input := scanner.Text()
		cleanedInput := util.CleanInput(input)
		inputCommand := cleanedInput[0]
		command, ok := cfg.Commands[inputCommand]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		err := command.Callback(cfg, cleanedInput[1:]...)
		if err != nil {
			fmt.Printf("error in executing command '%s': %v\n", command.Name, err)
		}
	}
}
