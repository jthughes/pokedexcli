package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jthughes/pokedexcli/internal/pokecache"
)

type Pokemon struct {
	ID             int              `json:"id"`
	Name           string           `json:"name"`
	BaseExperience int              `json:"base_experience"`
	Height         int              `json:"height"`
	IsDefault      bool             `json:"is_default"`
	Order          int              `json:"order"`
	Weight         int              `json:"weight"`
	Abilities      []PokemonAbility `json:"abilities"`
	Forms          []Resource       `json:"forms"`
	GameIndicies   []struct {
		GameIndex int      `json:"game_index"`
		Version   Resource `json:"version"`
	} `json:"game_indices"`
	HeldItems []struct {
		Item           Resource `json:"item"`
		VersionDetails []struct {
			Version Resource `json:"version"`
			Rarity  int      `json:"rarity"`
		}
	} `json:"held_items"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		MoveLearnMethod Resource `json:"move_learn_method"`
		VersionGroup    Resource `json:"version_group"`
		LevelLearnedAt  int      `json:"level_learned_at"`
	} `json:"moves"`
	PastTypes []struct {
		Generation Resource      `json:"generation"`
		Types      []PokemonType `json:"types"`
	} `json:"past_types"`
	Sprites PokemonSprites `json:"sprites"`
	Cries   PokemonCries   `json:"cries"`
	Species Resource       `json:"species"`
	Stats   []PokemonStat  `json:"stats"`
	Types   []PokemonType  `json:"types"`
}

type PokemonAbility struct {
	IsHidden bool     `json:"is_hidden"`
	Slot     int      `json:"slot"`
	Ability  Resource `json:"ability"`
}

type PokemonType struct {
	Slot int      `json:"slot"`
	Type Resource `json:"type"`
}

type PokemonSprites struct {
	FrontDefault     string `json:"front_default"`
	FrontShiny       string `json:"front_shiny"`
	FrontFemale      string `json:"front_female"`
	FrontShinyFemale string `json:"front_shiny_female"`
	BackDefault      string `json:"back_default"`
	BackShiny        string `json:"back_shiny"`
	BackFemale       string `json:"back_female"`
	BackShinyFemale  string `json:"back_shiny_female"`
}

type PokemonCries struct {
	Latest string `json:"latest"`
	Legacy string `json:"legacy"`
}

type PokemonStat struct {
	Stat     Resource `json:"stat"`
	Effort   int      `json:"effort"`
	BaseStat int      `json:"base_stat"`
}

func GetPokemon(pokemonName string, cache *pokecache.Cache) (Pokemon, error) {
	url := baseURL + "/pokemon/" + pokemonName

	var data []byte
	data, ok := cache.Get(url)
	if !ok {
		response, err := http.Get(url)
		if err != nil {
			return Pokemon{}, fmt.Errorf("network error: %w", err)
		}
		defer response.Body.Close()

		if response.StatusCode < 200 || response.StatusCode > 299 {
			return Pokemon{}, fmt.Errorf("Non-OK HTTP status: %s", response.Status)
		}

		data, err = io.ReadAll(response.Body)
		if err != nil {
			return Pokemon{}, fmt.Errorf("unable to read response body: %w", err)
		}
		cache.Add(url, data)
	}
	var resource Pokemon
	if err := json.Unmarshal(data, &resource); err != nil {
		return Pokemon{}, fmt.Errorf("unable to unmarshall data: %w", err)
	}
	return resource, nil
}

type PokemonSpecies struct {
	ID                   int        `json:"id"`
	Name                 string     `json:"name"`
	Order                int        `json:"order"`
	GenderRate           int        `json:"gender_rate"`
	CaptureRate          int        `json:"capture_rate"`
	BaseHappiness        int        `json:"base_happiness"`
	IsBaby               bool       `json:"is_baby"`
	IsLegendary          bool       `json:"is_legendary"`
	IsMythical           bool       `json:"is_mythical"`
	HatchCounter         int        `json:"hatch_counter"`
	HasGenderDifferences bool       `json:"has_gender_differences"`
	FormsSwitchable      bool       `json:"forms_switchable"`
	GrowthRate           Resource   `json:"growth_rate"`
	PokedexNumbers       []Resource `json:"pokedex_numbers"`
	EggGroups            []Resource `json:"egg_groups"`
	Color                Resource   `json:"color"`
	Shape                Resource   `json:"shape"`
	EvolvesFromSpecies   Resource   `json:"evolves_from_species"`
	EvolutionChain       struct {
		Url string `json:"url"`
	} `json:"evolution_chain"`
	Habitat    Resource `json:"habitat"`
	Generation Resource `json:"generation"`
	Names      []struct {
		Name     string   `json:"name"`
		Language Resource `json:"language"`
	} `json:"names"`
	PalParkEncounters []struct {
		BaseScore int      `json:"base_score"`
		Rate      int      `json:"rate"`
		Area      Resource `json:"area"`
	} `json:"pal_park_encounters"`
	FlavorTextEntries []struct {
		FlavorText string   `json:"flavor_text"`
		Language   Resource `json:"language"`
		Version    Resource `json:"version"`
	} `json:"flavor_text_entries"`
	FormDescriptions []struct {
		Description string   `json:"description"`
		Language    Resource `json:"language"`
	} `json:"form_descriptions"`
	Genera []struct {
		Genus    string   `json:"genus"`
		Language Resource `json:"language"`
	} `json:"genera"`
	Varieties []struct {
		IsDefault bool     `json:"is_default"`
		Pokemon   Resource `json:"pokemon"`
	} `json:"varieties"`
}

func GetPokemonSpecies(pokemonName string, cache *pokecache.Cache) (PokemonSpecies, error) {
	url := baseURL + "/pokemon-species/" + pokemonName

	var data []byte
	data, ok := cache.Get(url)
	if !ok {
		response, err := http.Get(url)
		if err != nil {
			return PokemonSpecies{}, fmt.Errorf("network error: %w", err)
		}
		defer response.Body.Close()

		if response.StatusCode < 200 || response.StatusCode > 299 {
			return PokemonSpecies{}, fmt.Errorf("Non-OK HTTP status: %s", response.Status)
		}

		data, err = io.ReadAll(response.Body)
		if err != nil {
			return PokemonSpecies{}, fmt.Errorf("unable to read response body: %w", err)
		}
		cache.Add(url, data)
	}
	var resource PokemonSpecies
	if err := json.Unmarshal(data, &resource); err != nil {
		return PokemonSpecies{}, fmt.Errorf("unable to unmarshall data: %w", err)
	}
	return resource, nil
}
