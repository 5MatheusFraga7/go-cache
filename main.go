package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	cache := NewCache()

	ctx := context.Background()
	states, err := fetchBrazilianStatesFromCache(ctx, cache.Redis)

	if err != nil {
		panic(err)
	}

	fmt.Println("Estados do Brasil:", states)
}

func fetchBrazilianStatesFromDB() ([]string, error) {
	brazilianStates := []string{"Acre", "Alagoas", "Amapá", "Amazonas", "Bahia", "Ceará", "Espírito Santo", "Goiás", "Maranhão", "Mato Grosso", "Mato Grosso do Sul", "Minas Gerais", "Pará", "Paraíba", "Paraná", "Pernambuco", "Piauí", "Rio de Janeiro", "Rio Grande do Norte", "Rio Grande do Sul", "Rondônia", "Roraima", "Santa Catarina", "São Paulo", "Sergipe", "Tocantins"}
	return brazilianStates, nil
}

func fetchBrazilianStatesFromCache(ctx context.Context, client *redis.Client) ([]string, error) {
	// Verificar se os estados do Brasil estão no cache
	val, err := client.Get(ctx, "brazilian_states").Result()
	if err == redis.Nil {
		states, err := fetchBrazilianStatesFromDB()
		if err != nil {
			fmt.Println("Error! States Not found!")
			return nil, err
		}

		statesJSON, _ := json.Marshal(states)
		err = client.Set(ctx, "brazilian_states", statesJSON, 24*time.Hour).Err()
		if err != nil {
			return nil, err
		}

		return states, nil
	} else if err != nil {
		fmt.Println("Error! Erro ao selecionar valores no redis")
		fmt.Println(err)
		return nil, err
	}

	var states []string
	err = json.Unmarshal([]byte(val), &states)
	if err != nil {
		return nil, err
	}

	return states, nil
}

type Cache struct {
	Redis *redis.Client
}

func NewCache() *Cache {
	return &Cache{
		Redis: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}),
	}
}
