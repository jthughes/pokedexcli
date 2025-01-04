package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"strings"
	"time"

	"github.com/jthughes/pokedexcli/internal/pokeapi"
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

func commandHelp(config *Config, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, command := range commands {
		fmt.Println(command.name + ": " + command.description)
	}
	return nil
}

func commandMap(url *string, config *Config) error {
	locations, err := pokeapi.GetResourceList(url, config.Cache)
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

func commandMapForward(config *Config, args []string) error {
	return commandMap(config.Next, config)
}

func commandMapBack(config *Config, args []string) error {
	if config.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}
	return commandMap(config.Previous, config)
}

func commandExplore(config *Config, args []string) error {
	if len(args) != 2 {
		fmt.Println("Expecting: explore <location-area>")
		return nil
	}
	locationArea := args[1]
	fmt.Println("Exploring " + locationArea + "...")
	pokemonList, err := pokeapi.GetPokemonList(locationArea, config.Cache)
	if err != nil {
		return err
	}
	fmt.Println("Found Pokemon:")
	for _, encounter := range pokemonList {
		fmt.Println(" - " + encounter.Pokemon.Name)
	}
	return nil
}

type Pokemon struct {
	pokeapi.Pokemon
	Species pokeapi.PokemonSpecies
}

func commandPokedex(config *Config, args []string) error {
	if len(args) != 1 {
		fmt.Println("Expecting: pokedex")
		return nil
	}
	if len(config.Pokedex) == 0 {
		fmt.Println("The Pokedex is empty. Catch some Pokemon!")
		return nil
	}
	fmt.Println("Your Pokedex:")
	for name, _ := range config.Pokedex {
		fmt.Println("  -", name)
	}
	return nil
}

func commandCatch(config *Config, args []string) error {
	pokeballs := map[string]float64{
		"Poke Ball":    1.0,
		"Great Ball":   1.5,
		"Ultra Ball":   2.0,
		"Safari Ball":  1.5,
		"Premier Ball": 1.0,
		"Luxury Ball":  1.0,
		"Heal Ball":    1.0,
		"Cherish Ball": 1.0,
	}
	if len(args) != 2 {
		fmt.Println("Expecting: catch <pokemon>")
		return nil
	}
	pokemonName := args[1]
	pokemon, err := pokeapi.GetPokemon(pokemonName, config.Cache)
	if err != nil {
		return err
	}
	pokemonSpecies, err := pokeapi.GetPokemonSpecies(pokemonName, config.Cache)
	if err != nil {
		return err
	}
	ball_type := "Poke Ball"
	fmt.Println("Throwing a " + ball_type + " at " + pokemonName + "...")

	pokeballRate := pokeballs[ball_type]
	catchRate := (1.0 / 3.0) * float64(pokemonSpecies.CaptureRate) * pokeballRate
	shakeRate := int(math.Floor(
		1_048_560 / math.Floor(math.Sqrt(
			math.Floor(math.Sqrt(
				math.Floor(16_711_680/catchRate)))))))
	shakes := [4]int{
		rand.IntN(65_536),
		rand.IntN(65_536),
		rand.IntN(65_536),
		rand.IntN(65_536),
	}
	shakeSuccesses := 0
	for _, shake := range shakes[:3] {
		if shake >= shakeRate {
			break
		}
		shakeSuccesses += 1
		fmt.Println("*Shakes*")
		time.Sleep(1500 * time.Millisecond)
	}
	shakeMessage := map[int]string{
		0: "Oh, no!\nThe Pokemon broke free!",
		1: "Aww!\nIt appeared to be caught!",
		2: "Aargh!\nAlmost had it!",
		3: "Shoot!\nIt was so close, too!",
	}
	if shakeSuccesses == 3 && shakes[3] < shakeRate {
		fmt.Println("Gotcha! " + pokemonName + " was caught!")
		if _, ok := config.Pokedex[pokemonName]; ok {
			return nil
		}
		fmt.Println("Adding " + pokemonName + " to the Pokedex.")
		config.Pokedex[pokemonName] = Pokemon{
			Pokemon: pokemon,
			Species: pokemonSpecies,
		}
	} else {
		fmt.Println(shakeMessage[shakeSuccesses])
	}
	return nil
}

func commandInspect(config *Config, args []string) error {
	if len(args) != 2 {
		fmt.Println("Expecting: inspect <pokemon>")
		return nil
	}
	pokemonName := args[1]
	pokemon, ok := config.Pokedex[pokemonName]
	if !ok {
		fmt.Println(pokemonName + " has not been caught yet.")
		return nil
	}
	fmt.Println("Name:", pokemon.Name)
	fmt.Println("Height:", pokemon.Height)
	fmt.Println("Weight:", pokemon.Weight)
	fmt.Println("Stats:")

	for _, stat := range pokemon.Stats {
		fmt.Println("  -"+stat.Stat.Name+":", stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, pokemonType := range pokemon.Types {
		fmt.Println("  -", pokemonType.Type.Name)
	}

	return nil
}

func commandExit(config *Config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
