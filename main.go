package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/sparrowhawk425/pokedexcli/internal/pokeapi"
	"github.com/sparrowhawk425/pokedexcli/internal/pokecache"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	config := Config{}
	config.cache = pokecache.NewCache(time.Second * 5)
	config.pokedex = map[string]pokeapi.Pokemon{}
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		cleanText := cleanInput(scanner.Text())
		commands := getCommandMap()
		cmd, exists := commands[cleanText[0]]
		if exists {
			if err := cmd.callback(&config, cleanText[1:]); err != nil {
				fmt.Printf("%v\n", err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
