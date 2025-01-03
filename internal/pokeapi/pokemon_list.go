package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jthughes/pokedexcli/internal/pokecache"
)

type Encounter struct {
	Chance     int        `json:"chance"`
	Conditions []Resource `json:"condition_values"`
	MaxLevel   int        `json:"max_level"`
	Method     Resource   `json:"method"`
	MinLevel   int        `json:"min_level"`
}

type EncounterMethodRates struct {
	Method         Resource `json:"encounter_method"`
	VersionDetails []struct {
		Rate    int      `json:"rate"`
		Version Resource `json:"version"`
	} `json:"version_details"`
}

type VersionEncounterDetail struct {
	Version          Resource    `json:"version"`
	MaxChange        int         `json:"max_chance"`
	EncounterDetails []Encounter `json:"Encounter"`
}

type PokemonEncounter struct {
	Pokemon        Resource                 `json:"pokemon"`
	VersionDetails []VersionEncounterDetail `json:"version_details"`
}

type LocationArea struct {
	EncounterMethodRates []EncounterMethodRates `json:"encounter_method_rates"`
	GameIndex            int                    `json:"game_index"`
	ID                   int                    `json:"id"`
	Location             Resource               `json:"location"`
	Name                 string                 `json:"name"`
	Names                []struct {
		Name     string   `json:"name"`
		Language Resource `json:"language"`
	} `json:"names"`
	Encounters []PokemonEncounter `json:"pokemon_encounters"`
}

func GetPokemonList(locationArea string, cache *pokecache.Cache) ([]PokemonEncounter, error) {
	url := baseURL + "/location-area/" + locationArea

	var data []byte
	data, ok := cache.Get(url)
	if !ok {
		response, err := http.Get(url)
		if err != nil {
			return []PokemonEncounter{}, fmt.Errorf("network error: %w", err)
		}
		defer response.Body.Close()

		if response.StatusCode < 200 || response.StatusCode > 299 {
			return []PokemonEncounter{}, fmt.Errorf("Non-OK HTTP status: %s", response.Status)
		}

		data, err = io.ReadAll(response.Body)
		if err != nil {
			return []PokemonEncounter{}, fmt.Errorf("unable to read response body: %w", err)
		}
		cache.Add(url, data)
	}
	var resources LocationArea
	if err := json.Unmarshal(data, &resources); err != nil {
		return []PokemonEncounter{}, fmt.Errorf("unable to unmarshall data: %w", err)
	}
	return resources.Encounters, nil
}

// func (r ResourceList) print() {
// 	data, err := json.MarshalIndent(r, "", "  ")
// 	if err != nil {
// 		fmt.Println("Unable to marshall ResourceList")
// 		return
// 	}
// 	fmt.Println(string(data))
// }
