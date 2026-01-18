package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/risbern21/pokedex/internal/cache"
	"github.com/risbern21/pokedex/internal/dto"
	"github.com/risbern21/pokedex/models"
	"github.com/risbern21/pokedex/services"
	"gorm.io/gorm"
)

type cliCommand struct {
	Title       string
	Description string
	Callback    func(config *Config) (string, error)
}

type Config struct {
	C            *cache.Cache
	Next         string
	Previous     *string
	LocationName *string
	PokemonID    *uint
	PokemonName  *string
	Pokeballs    map[string]dto.PokemonDetails
	Red          *color.Color
	Blue         *color.Color
	Green        *color.Color
}

func Newconfig() *Config {
	return &Config{C: cache.NewCache(5 * time.Second), Pokeballs: make(map[string]dto.PokemonDetails), Next: "https://pokeapi.co/api/v2/location-area/", Red: color.New(color.FgRed), Blue: color.New(color.FgBlue), Green: color.New(color.FgGreen)}
}

var Commands map[string]cliCommand

func init() {
	Commands = map[string]cliCommand{
		"map":      {Title: "map", Description: "Locations available to explore (paginated)", Callback: commandMap},
		"map back": {Title: "map back", Description: "Goes to previous page of areas", Callback: commandMapBack},
		"explore":  {Title: "explore", Description: "Returns Pokemon that can be encountered in an area", Callback: commandExplore},
		"catch":    {Title: "catch", Description: "Tries to catch a pokemon", Callback: commandCatch},
		"inspect":  {Title: "inspect", Description: "Allows you to inspect you caught pokemon", Callback: commandInspect},
		"pokedex":  {Title: "pokedex", Description: "Displays all your caught pokemon", Callback: commandPokedex},
	}
}

type responseBody struct {
	Count    int
	Next     string
	Previous *string
	Results  []location
}

type location struct {
	Name string
	URL  string
}

func commandMap(config *Config) (string, error) {
	response := responseBody{}

	val, ok := config.C.Get(config.Next)
	if ok {
		if err := json.Unmarshal(val, &response); err != nil {
			return "", err
		}
	} else {
		res, err := http.Get(config.Next)
		if err != nil {
			return "", err
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return "", err
		}
		if err := json.Unmarshal(body, &response); err != nil {
			return "", err
		}
		config.C.Add(config.Next, body)
	}
	var s string

	for _, result := range response.Results {
		s += fmt.Sprintf("%s\n", result.Name)
	}

	config.Next = response.Next
	config.Previous = response.Previous
	return s, nil
}

func commandMapBack(config *Config) (string, error) {
	response := responseBody{}
	prev := config.Previous
	if prev == nil {
		return "", nil
	}

	val, ok := config.C.Get(*prev)
	if ok {
		if err := json.Unmarshal(val, &response); err != nil {
			return "", err
		}
	} else {
		res, err := http.Get(*prev)
		if err != nil {
			return "", err
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return "", err
		}
		defer res.Body.Close()

		if err := json.Unmarshal(body, &response); err != nil {
			return "", err
		}
		config.C.Add(*prev, body)
	}
	var s string

	for _, result := range response.Results {
		s += fmt.Sprintf("%s\n", result.Name)
	}

	config.Next = response.Next
	config.Previous = response.Previous
	return s, nil
}

func commandExplore(config *Config) (string, error) {
	location := &dto.LocationDetails{}
	endpoint := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", *config.LocationName)

	l, ok := config.C.Get(endpoint)
	if ok {
		if err := json.Unmarshal(l, location); err != nil {
			return "", err
		}
	} else {
		res, err := http.Get(endpoint)
		if res.StatusCode == 404 {
			return "invalid area name\nplease use the map option to see list of availbel areas to explore", nil
		}

		if err != nil {
			return "", err
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return "", err
		}

		if err := json.Unmarshal(body, location); err != nil {
			return "", err
		}

		config.C.Add(endpoint, body)
	}
	var s string

	for _, p := range location.PokemonEncounters {
		s += fmt.Sprintf("%s\n", p.Pokemon.Name)
	}
	return s, nil
}
func commandCatch(config *Config) (string, error) {
	endpoint := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", *config.PokemonName)
	pokemon := &dto.PokemonDetails{}

	if p, ok := config.C.Get(endpoint); ok {
		if err := json.Unmarshal(p, pokemon); err != nil {
			return "", err
		}
	} else {
		res, err := http.Get(endpoint)
		if res.StatusCode == 404 {
			return "invalid pokemon name", nil
		}
		if err != nil {
			return "", err
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return "", err
		}

		if err := json.Unmarshal(body, pokemon); err != nil {
			return "", err
		}

		config.C.Add(endpoint, body)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	chance := math.Mod(float64(pokemon.BaseExperience), float64(r.Float32()))
	if chance > 0.5 {
		config.Pokeballs[pokemon.Name] = *pokemon
		p := services.New()
		p.PokemonDetails = pokemon
		if err := p.Add(); err != nil {
			return "", err
		}

		return fmt.Sprintf("%s was caught!\n", pokemon.Name), nil
	} else {
		return fmt.Sprintf("%s escaped!\n", pokemon.Name), nil
	}
}

func commandInspect(config *Config) (string, error) {
	pokemon := models.New()
	pokemon.Name = *config.PokemonName
	if err := pokemon.GetByName(); err != nil {
		if err == gorm.ErrRecordNotFound {
			return "Invalid pokemon name or you have'nt caught this pokemon yet...\nCatch a Pokemon first to inspect it", nil
		}
		return "", err
	}

	var s string

	s += fmt.Sprintf("%s\nHeight: %d\nWeight: %d\nStats\n", pokemon.Name, pokemon.Height, pokemon.Weight)
	for _, stat := range pokemon.Stats {
		s += fmt.Sprintf("\t- %s: %v\n", stat.Name, stat.BaseStat)
	}

	s += "Types:\n"
	for _, t := range pokemon.Types {
		s += fmt.Sprintf("\t- %s\n", t.Name)
	}

	return s, nil
}

func commandPokedex(config *Config) (string, error) {
	pokemon := &models.Pokemons{}
	if err := pokemon.GetAll(); err != nil {
		return "", err
	}

	var s string
	s += "Your pokedex:\n"
	for _, pokemon := range *pokemon {
		s += fmt.Sprintf("\t- %s\n", pokemon.Name)
	}

	return s, nil
}
