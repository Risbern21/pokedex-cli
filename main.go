package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	res, err := http.Get("https://pokeapi.co/api/v2/pokemon/ditto")
	if err != nil {
		fmt.Println("error while fetching pokedata :", err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(":error reading bytes from resopnse body : ", err)
	}

	fmt.Println(string(body))
}
