package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
)

// Cafeteria struct
type Cafeteria struct {
	Coffees []Coffee
	Name    string
}

// Coffee is the struct used for request and response
type Coffee struct {
	Name  string
	Price float64
}

var coffees []Coffee

func getEnvString(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}

func getEnvInt(key string, defaultValue int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}
	return intVal
}

// GetAllCoffees returns all the coffees in the cafeteria
func (cafeteria *Cafeteria) GetAllCoffees(_ string, coffees *[]Coffee) error {
	*coffees = cafeteria.Coffees
	return nil
}

// AddCoffee adds a coffee to cafeteria
func (cafeteria *Cafeteria) AddCoffee(coffee Coffee, response *Coffee) error {
	cafeteria.Coffees = append(cafeteria.Coffees, coffee)
	*response = coffee
	return nil
}

func main() {
	coffees = make([]Coffee, 0, 10)
	coffees = []Coffee{{Name: "Cappucino", Price: 1.2}, {Name: "Frappucino", Price: 1.5},
		{Name: "SomeCoffee", Price: 2.2}}

	cafeteria := new(Cafeteria)
	cafeteria.Coffees = coffees

	rpcServer := rpc.NewServer()

	// custom rpc http handler paths
	rpcServer.HandleHTTP("/rpc", "/debug")
	rpcServer.Register(cafeteria)

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", getEnvInt("PORT", 8080)))
	if err != nil {
		log.Fatalf("unable to create a tcp listener %v", err)
	}

	http.Serve(listener, rpcServer)
}
