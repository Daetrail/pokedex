// Package config manages global state
package config

import "github.com/Daetrail/pokedex/internal/pokecache"

type Config struct {
	BaseURL          string
	Commands         map[string]CliCommand
	MapCommandConfig MapCommandConfig
	PokeCache        *pokecache.PokeCache
}
