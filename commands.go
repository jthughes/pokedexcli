package main

import (
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"time"

	"github.com/jthughes/pokedexcli/internal/pokeapi"
)

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
	for name := range config.Pokedex {
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
