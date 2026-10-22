package config

type Encounters struct {
	PokemonEncounters []PokemonEncountersStruct `json:"pokemon_encounters"`
}

type PokemonEncountersStruct struct {
	Pokemon PokemonStruct `json:"pokemon"`
}

type PokemonStruct struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
