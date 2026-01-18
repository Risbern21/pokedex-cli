package services

import (
	"github.com/risbern21/pokedex/internal/dto"
	"github.com/risbern21/pokedex/models"
)

type Pokemon struct {
	ID             uint
	Name           string
	PokemonDetails *dto.PokemonDetails
}

func mapStatsAndTypes(p *dto.PokemonDetails) *models.Pokemon {
	m := models.New()
	m.Name = p.Name
	m.Height = p.Height
	m.Weight = p.Weight

	// map stats
	for _, s := range p.Stats {
		m.Stats = append(m.Stats, models.PokemonStats{
			PokemonID: m.ID,
			Name:      s.Stat.Name,
			BaseStat:  s.BaseStat,
			Effort:    s.Effort,
			URL:       s.Stat.URL,
		})
	}

	// map types
	for _, t := range p.Types {
		m.Types = append(m.Types, models.PokemonTypes{
			PokemonID: m.ID,
			Slot:      t.Slot,
			Name:      t.Type.Name,
			URL:       t.Type.URL,
		})
	}

	return m
}

func New() *Pokemon {
	return &Pokemon{}
}

func (p *Pokemon) Add() error {
	m := mapStatsAndTypes(p.PokemonDetails)
	if err := m.Add(); err != nil {
		return err
	}

	return nil
}
