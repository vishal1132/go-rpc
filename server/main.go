package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
)

// CafeManager struct
type CafeManager struct {
	Cafes []Cafeteria
}

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

func (cafe *Cafeteria) Rename(name string, status *string) error {
	cafe.Name = name
	*status = cafe.Name
	return nil
}

func (cafeteria *Cafeteria) GetName(_ string, resp *string) error {
	*resp = cafeteria.Name
	return nil
}

func (cm *CafeManager) OpenCafe(name string, response *string) error {
	for _, v := range cm.Cafes {
		if v.Name == name {
			*response = "fail"
			return errors.New("OOps cafe with this name already exists")
		}
	}
	cm.Cafes = append(cm.Cafes, Cafeteria{Name: name, Coffees: []Coffee{}})
	*response = "success"
	return nil
}

func main() {
	cmanager := new(CafeManager)
	cmanager.Cafes = make([]Cafeteria, 0, 10)
	rpcServer := rpc.NewServer()

	// custom rpc http handler paths
	rpcServer.HandleHTTP("/rpc", "/debug")
	rpcServer.Register(cmanager)

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", getEnvInt("PORT", 8080)))
	if err != nil {
		log.Fatalf("unable to create a tcp listener %v", err)
	}

	http.Serve(listener, rpcServer)

}
