package models

import (
	"github.com/risbern21/pokedex/internal/database"
)

type PokemonStats struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	PokemonID uint   `json:"pokemon_id" gorm:"not null;index"`
	BaseStat  int    `json:"base_stat"`
	Effort    int    `json:"effort"`
	Name      string `json:"name"`
	URL       string `json:"url"`
}

type PokemonTypes struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	PokemonID uint   `json:"pokemon_id" gorm:"not null;index"`
	Slot      int    `json:"slot"`
	Name      string `json:"name"`
	URL       string `json:"url"`
}

type Pokemon struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Name   string `json:"name" gorm:"name"`
	Height int    `json:"height" gorm:"height"`
	Weight int    `json:"weight" gorm:"weight"`

	Stats []PokemonStats `json:"stats" gorm:"foreignKey:PokemonID;constraint:OnDelete:CASCADE"`
	Types []PokemonTypes `json:"types" gorm:"foreignKey:PokemonID;constraint:OnDelete:CASCADE"`
}

type Pokemons []Pokemon

func New() *Pokemon {
	return &Pokemon{}
}

func (p *Pokemon) Add() error {
	if err := database.Client().Create(&p); err != nil {
		return err.Error
	}
	return nil
}

func (p *Pokemon) Get() error {
	if err := database.Client().Preload("Stats").Preload("Types").First(&p, p.ID); err != nil {
		return err.Error
	}
	return nil
}

func (p *Pokemon) GetByName() error {
	if err := database.Client().Preload("Stats").Preload("Types").Where("name=?", p.Name).First(&p); err != nil {
		return err.Error
	}
	return nil
}

func (p *Pokemons) GetAll() error {
	if err := database.Client().Select("id", "name").Find(&p); err != nil {
		return err.Error
	}
	return nil
}

func (p *Pokemon) Release() error {
	if err := database.Client().Delete(&p, p.ID); err != nil {
		return err.Error
	}

	return nil
}
