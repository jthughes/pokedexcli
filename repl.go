package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jthughes/pokedexcli/internal/pokeapi"
)

var commands map[string]cliCommand

func repl() {
	commands = registerCommands()
	config := Config{}
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
		err := command.callback(&config)
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
	callback    func(*Config) error
}

type Config struct {
	Next     *string
	Previous *string
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
	commands["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	}
	return commands
}

func commandHelp(config *Config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, command := range commands {
		fmt.Println(command.name + ": " + command.description)
	}
	return nil
}

func commandMap(url *string, config *Config) error {
	locations, err := pokeapi.GetResourceList(url)
	if err != nil {
		return err
	}
	config.Next = locations.Next
	config.Previous = locations.Previous
	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func commandMapForward(config *Config) error {
	return commandMap(config.Next, config)
}

func commandMapBack(config *Config) error {
	if config.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}
	return commandMap(config.Previous, config)
}

func commandExit(config *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
