package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Forge Abacus listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
