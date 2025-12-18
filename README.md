# Pokedex CLI

**A simple Command-Line Pokédex** built in Go that lets you explore Pokémon world areas and catch Pokémon right from your terminal! :contentReference[oaicite:1]{index=1}

## Features

- **Explore locations** — view Pokémon by area
- **Explore map pages** — scroll forward/back through location pages
- **Catch Pokémon** — add Pokémon to your Pokédex
- **Inspect caught Pokémon** — view details for Pokémon you’ve caught
- **View caught list** — show your current Pokédex
- **Help & exit** commands for smoother CLI use

*(Work in progress — more features planned!)* :contentReference[oaicite:2]{index=2}

---

## Installation

To compile and install the CLI locally, make sure you have **Go 1.18+** installed.

```bash
git clone https://github.com/risbern21/pokedex-cli.git
cd pokedex-cli
go build -o pokedex
./pokedex
