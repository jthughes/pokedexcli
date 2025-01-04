package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jthughes/pokedexcli/internal/pokecache"
)

var commands map[string]cliCommand

func repl() {
	commands = registerCommands()
	interval, err := time.ParseDuration("5s")
	if err != nil {
		fmt.Println("Unable to set duration:", err)
		os.Exit(1)
	}
	config := Config{
		Cache:   pokecache.NewCache(interval),
		Pokedex: map[string]Pokemon{},
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		words := cleanInput(input)
		command, ok := commands[words[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		err := command.callback(&config, words)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func cleanInput(text string) []string {
	words := strings.Fields(strings.ToLower(text))
	return words
}

type cliCommand struct {
	name        string
	description string
	callback    func(*Config, []string) error
}

type Config struct {
	Next     *string
	Previous *string
	Cache    *pokecache.Cache
	Pokedex  map[string]Pokemon
}

func registerCommands() (commands map[string]cliCommand) {
	commands = map[string]cliCommand{}
	commands["help"] = cliCommand{
		name:        "help",
		description: "Displays list of available commands",
		callback:    commandHelp,
	}
	commands["map"] = cliCommand{
		name:        "map",
		description: "Displays the next 20 map locations",
		callback:    commandMapForward,
	}
	commands["mapb"] = cliCommand{
		name:        "mapb",
		description: "Displays the previous 20 map locations",
		callback:    commandMapBack,
	}
	commands["explore"] = cliCommand{
		name:        "explore",
		description: "Displays the Pokemon found at the provided location",
		callback:    commandExplore,
	}
	commands["pokedex"] = cliCommand{
		name:        "pokedex",
		description: "List all Pokemon in the Pokedex",
		callback:    commandPokedex,
	}
	commands["catch"] = cliCommand{
		name:        "catch",
		description: "Attempt to catch a Pokemon",
		callback:    commandCatch,
	}
	commands["inspect"] = cliCommand{
		name:        "inspect",
		description: "Inspect a Pokemon in the Pokedex",
		callback:    commandInspect,
	}
	commands["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	}
	return commands
}
