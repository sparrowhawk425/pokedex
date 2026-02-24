package main

import (
	"fmt"
	"math/rand/v2"
	"os"

	"github.com/sparrowhawk425/pokedexcli/internal/pokeapi"
	"github.com/sparrowhawk425/pokedexcli/internal/pokecache"
)

const LOCATION_URL string = "https://pokeapi.co/api/v2/location-area/"
const POKEMON_URL string = "https://pokeapi.co/api/v2/pokemon/"

type cliCommand struct {
	name        string
	description string
	callback    func(*Config, []string) error
}

type Config struct {
	Next     *string
	Previous *string
	cache    pokecache.Cache
	pokedex  map[string]pokeapi.Pokemon
}

func getCommandMap() map[string]cliCommand {
	commandMap := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "List locations in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Cycle back through list of locations in the Pokemon world",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore <location>",
			description: "Lists available Pokemon from a provided <location>",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch <pokemon>",
			description: "Attempts to catch a Pokemon and add it to the Pokedex",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect <pokemon>",
			description: "See the stats for a Pokemon if you have successfully caught it (and added it to the Pokedex)",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List all the currently caught Pokemon",
			callback:    commandPokedex,
		},
	}
	return commandMap
}

// Command Callbacks

func commandExit(config *Config, _ []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config, _ []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for _, cmd := range getCommandMap() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(config *Config, _ []string) error {
	url := LOCATION_URL
	if config != nil && config.Next != nil {
		url = *config.Next
	}

	data, err := makeRequest(config, url)
	if err != nil {
		return fmt.Errorf("Error making request %v", err)
	}
	areaMap, err := pokeapi.ParseMapData(data)
	if err != nil {
		return fmt.Errorf("Error parsing area map data %v", err)
	}
	printMapData(areaMap, config)
	return nil
}

func commandMapBack(config *Config, _ []string) error {

	if config == nil || config.Previous == nil {
		fmt.Println("You're on the first page")
		return nil
	}
	url := *config.Previous

	data, err := makeRequest(config, url)
	if err != nil {
		return fmt.Errorf("Error making request %v", err)
	}
	areaMap, err := pokeapi.ParseMapData(data)
	if err != nil {
		return fmt.Errorf("Error parsing area map %v", err)
	}
	printMapData(areaMap, config)
	return nil
}

func commandExplore(config *Config, params []string) error {

	if len(params) < 1 {
		return fmt.Errorf("Missing required parameter <location>")
	}
	url := LOCATION_URL + params[0]
	data, err := makeRequest(config, url)
	if err != nil {
		return fmt.Errorf("Error making request %v", err)
	}
	location, err := pokeapi.ParseLocationData(data)
	if err != nil {
		return fmt.Errorf("Error parsing location data %v", err)
	}
	printLocation(location)
	return nil
}

func commandCatch(config *Config, params []string) error {

	if len(params) < 1 {
		return fmt.Errorf("Missing required parameter <pokemon>")
	}
	url := POKEMON_URL + params[0]
	data, err := makeRequest(config, url)
	if err != nil {
		return fmt.Errorf("Error making request %v", err)
	}
	pokemon, err := pokeapi.ParsePokemon(data)
	if err != nil {
		return fmt.Errorf("Error parsing pokemon data %v", err)
	}
	printCatchPokemon(pokemon, config)
	return nil
}

func commandInspect(config *Config, params []string) error {

	if len(params) < 1 {
		return fmt.Errorf("Missing required parameter <pokemon>")
	}
	name := params[0]
	pokemon, ok := config.pokedex[name]
	if !ok {
		fmt.Printf("You haven't caught %s!\n", name)
		return nil
	}
	fmt.Printf("Name: %s\n", name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf(" -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, tp := range pokemon.Types {
		fmt.Printf(" - %s\n", tp.Type.Name)
	}
	return nil
}

func commandPokedex(config *Config, params []string) error {
	fmt.Println("Your Pokemon:")
	for _, pokemon := range config.pokedex {
		fmt.Printf(" - %s\n", pokemon.Name)
	}
	return nil
}

func makeRequest(config *Config, url string) ([]byte, error) {

	data, ok := config.cache.Get(url)
	if ok {
		return data, nil
	}
	data, err := pokeapi.MakeHTTPGetRequest(url)
	if err != nil {
		return nil, fmt.Errorf("Error fetching area map %v", err)
	}
	config.cache.Add(url, data)
	return data, nil
}

// Print functions

func printMapData(areaMap pokeapi.PokedexAreaMap, config *Config) {

	for _, result := range areaMap.Results {
		fmt.Printf("%s\n", result.Name)
	}
	config.Next = areaMap.Next
	config.Previous = areaMap.Previous
}

func printLocation(location pokeapi.PokedexArea) {
	fmt.Printf("Exploring %s...\n", location.Name)
	fmt.Println("Found Pokemon:")
	for _, encounter := range location.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}
}

func printCatchPokemon(pokemon pokeapi.Pokemon, config *Config) {

	name := pokemon.Name
	fmt.Printf("Throwing a Pokeball at %s...\n", name)
	chance := rand.IntN(350)
	if chance > pokemon.BaseExperience {
		fmt.Printf("%s was caught!\n", name)
		config.pokedex[name] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", name)
	}
}
