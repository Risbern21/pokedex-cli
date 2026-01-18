package config

import (
	"log"

	"github.com/risbern21/pokedex/internal/database"
	"github.com/risbern21/pokedex/models"
)

func AutoMigrate() {
	if err := database.Client().AutoMigrate(&models.Pokemon{}, &models.PokemonStats{}, &models.PokemonTypes{}); err != nil {
		log.Fatalf("unablet to perform db migration")
	}
}
