package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jthughes/pokedexcli/internal/pokecache"
)

type ResourceList struct {
	Count    int        `json:"count"`
	Next     *string    `json:"next"`
	Previous *string    `json:"previous"`
	Results  []Resource `json:"results"`
}

type Resource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func GetResourceList(pageURL *string, cache *pokecache.Cache) (ResourceList, error) {
	url := baseURL + "/location-area"
	if pageURL != nil {
		url = *pageURL
	}

	var data []byte
	data, ok := cache.Get(url)
	if !ok {
		response, err := http.Get(url)
		if err != nil {
			return ResourceList{}, fmt.Errorf("network error: %w", err)
		}
		defer response.Body.Close()

		if response.StatusCode < 200 || response.StatusCode > 299 {
			return ResourceList{}, fmt.Errorf("Non-OK HTTP status: %s", response.Status)
		}

		data, err = io.ReadAll(response.Body)
		if err != nil {
			return ResourceList{}, fmt.Errorf("unable to read response body: %w", err)
		}
		cache.Add(url, data)
	}

	var resources ResourceList
	if err := json.Unmarshal(data, &resources); err != nil {
		return ResourceList{}, fmt.Errorf("unable to unmarshall data: %w", err)
	}
	return resources, nil
}

func (r ResourceList) print() {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		fmt.Println("Unable to marshall ResourceList")
		return
	}
	fmt.Println(string(data))
}
