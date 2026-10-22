package config

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Daetrail/pokedex/internal/util"
)

type CliCommand struct {
	Name        string
	Description string
	Callback    func(*Config, ...string) error
}

func GetCommands() map[string]CliCommand {
	return map[string]CliCommand{
		"exit": {
			Name:        "exit",
			Description: "Exit the Pokedex",
			Callback:    commandExit,
		},
		"help": {
			Name:        "help",
			Description: "Displays a help message",
			Callback:    commandHelp,
		},
		"map": {
			Name:        "map",
			Description: "Displays the names of 20 location areas in the Pokemon world",
			Callback:    commandMap,
		},
		"mapb": {
			Name:        "mapb",
			Description: "Displays the names of the previous 20 location areas in the Pokemon world",
			Callback:    commandMapb,
		},
		"explore": {
			Name:        "explore",
			Description: "Displays the Pokemons in a user-specified location",
			Callback:    commandExplore,
		},
		"catch": {
			Name:        "catch",
			Description: "Attempts to catch a user-specified Pokemon",
			Callback:    commandCatch,
		},
	}
}

func commandExit(cfg *Config, _ ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *Config, _ ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, validCommand := range cfg.Commands {
		fmt.Printf("%s: %s\n", validCommand.Name, validCommand.Description)
	}
	return nil
}

func httpGet(cfg *Config, fullURL string) ([]byte, bool, error) {
	cachedData, ok := cfg.PokeCache.Get(fullURL)
	if ok {
		return cachedData, true, nil
	}

	res, err := http.Get(fullURL)
	if err != nil {
		return []byte{}, false, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, false, err
	}

	cfg.PokeCache.Add(fullURL, body)
	success := res.StatusCode > 299
	return body, !success, nil
}

func printErrorResponse(data []byte) error {
	var errResponse ErrorResponse
	if err := json.Unmarshal(data, &errResponse); err != nil {
		return err
	}
	fmt.Println(errResponse.Message)
	return nil
}

func commandMap(cfg *Config, _ ...string) error {
	var fullURL string
	if cfg.MapCommandConfig.Next != "" {
		fullURL = cfg.MapCommandConfig.Next
	} else {
		fullURL = cfg.BaseURL + "/location-area"
	}

	data, _, err := httpGet(cfg, fullURL)
	if err != nil {
		return err
	}

	var jsonData LocationResponse
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return err
	}

	if jsonData.Next != nil {
		cfg.MapCommandConfig.Next = *jsonData.Next
	}

	if jsonData.Previous != nil {
		cfg.MapCommandConfig.Previous = *jsonData.Previous
	}

	for _, result := range jsonData.Results {
		fmt.Println(result.Name)
	}

	return nil
}

func commandMapb(cfg *Config, _ ...string) error {
	var fullURL string
	if cfg.MapCommandConfig.Previous != "" {
		fullURL = cfg.MapCommandConfig.Previous
	} else {
		fullURL = cfg.BaseURL + "/location-area"
	}

	data, _, err := httpGet(cfg, fullURL)
	if err != nil {
		return err
	}

	var jsonData LocationResponse
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return err
	}

	if jsonData.Next != nil {
		cfg.MapCommandConfig.Next = *jsonData.Next
	}

	if jsonData.Previous != nil {
		cfg.MapCommandConfig.Previous = *jsonData.Previous
	}

	for _, result := range jsonData.Results {
		fmt.Println(result.Name)
	}

	return nil
}

func commandExplore(cfg *Config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("must have only 1 argument (name of location)")
	}

	fullURL := fmt.Sprintf("%v/%v/%v", cfg.BaseURL, "location-area", args[0])

	data, success, err := httpGet(cfg, fullURL)
	if err != nil {
		return err
	}
	if !success {
		err := printErrorResponse(data)
		return err
	}

	var pokemons Encounters
	if err := json.Unmarshal(data, &pokemons); err != nil {
		return err
	}

	fmt.Println("Exploring", args[0]+"...")
	fmt.Println("Found Pokemon:")
	for _, encounter := range pokemons.PokemonEncounters {
		pokemonName := encounter.Pokemon.Name
		fmt.Println("-", pokemonName)
	}
	return nil
}

func commandCatch(cfg *Config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("must have only 1 argument (name of pokemon)")
	}

	fullURL := cfg.BaseURL + "/pokemon/" + args[0]

	data, success, err := httpGet(cfg, fullURL)
	if err != nil {
		return err
	}
	if !success {
		err := printErrorResponse(data)
		return err
	}

	var baseExperience struct {
		BaseExperience int `json:"base_experience"`
	}
	if err := json.Unmarshal(data, &baseExperience); err != nil {
		return err
	}
	if baseExperience.BaseExperience == 0 {
		baseExperience.BaseExperience = 1
	}

	fmt.Printf("Throwing a Pokeball at %v...\n", args[0])

	if util.ChanceInverse(float64(baseExperience.BaseExperience), 50.0) {
		fmt.Printf("%v was caught!\n", args[0])
		return nil
	}

	fmt.Printf("%v escaped!\n", args[0])
	return nil
}
