package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/risbern21/pokedex/internal/config"
	"github.com/risbern21/pokedex/internal/database"
	"github.com/risbern21/pokedex/internal/model"
)

func main() {
	database.Connect()
	config.AutoMigrate()

	m := model.New()

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatalf("unable to run the tui")
		os.Exit(1)
	}
}
