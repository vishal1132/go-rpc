package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strconv"
)

// Coffee is the struct used for request and response
type Coffee struct {
	Name  string
	Price float64
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

func main() {
	cl, err := rpc.DialHTTPPath("tcp", fmt.Sprintf("localhost:%d", getEnvInt("PORT", 8080)), "/rpc")
	if err != nil {
		log.Fatalf("error occurred while creating rpc client %v", err)
	}
	var coffees []Coffee
	var coffee = Coffee{Name: "Adding a coffee", Price: 3.4}
	var response Coffee

	cl.Call("Cafeteria.GetAllCoffees", "", &coffees)
	fmt.Println(coffees)

	cl.Call("Cafeteria.AddCoffee", coffee, &response)
	fmt.Println(response)

	cl.Call("Cafeteria.GetAllCoffees", "", &coffees)
	fmt.Println(coffees)

}
