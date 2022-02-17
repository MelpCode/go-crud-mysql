package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

func conexionDB() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "root"
	Contrasenia := ""
	Nombre := "gomysql"

	conexion, err := sql.Open(Driver, Usuario+":"+Contrasenia+"@tcp(127.0.0.1)/"+Nombre)
	if err != nil {
		panic(err.Error())
	}
	return conexion
}

//Structure
type Menu struct {
	ID    int     `json:"ID"`
	Name  string  `json:"Name"`
	Price float64 `json:"Price"`
}

//Controllers
func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my API with MYSQL")
}

func getMenus(w http.ResponseWriter, r *http.Request) {
	stablishedConexion := conexionDB()
	registers, err := stablishedConexion.Query("SELECT * FROM menus")
	if err != nil {
		panic(err.Error())
	}

	menu := Menu{}
	arrayMenu := []Menu{}

	for registers.Next() {
		var id int
		var nombre string
		var price float64

		err := registers.Scan(&id, &nombre, &price)
		if err != nil {
			panic(err.Error())
		}
		menu.ID = id
		menu.Name = nombre
		menu.Price = price

		arrayMenu = append(arrayMenu, menu)
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(arrayMenu)
}

func getMenu(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	menuID, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err.Error())
	}

	stablishedConexion := conexionDB()
	register, err := stablishedConexion.Query("SELECT * FROM menus WHERE id=?", menuID)
	if err != nil {
		panic(err.Error())
	}

	menu := Menu{}

	for register.Next() {
		var id int
		var nombre string
		var price float64
		err := register.Scan(&id, &nombre, &price)
		if err != nil {
			panic(err.Error())
		}
		menu.ID = id
		menu.Name = nombre
		menu.Price = price
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(menu)
}

func createMenu(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		var newMenu Menu
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err.Error())
		}
		json.Unmarshal(reqBody, &newMenu)

		stablishedConexion := conexionDB()
		insertRegister, err := stablishedConexion.Prepare("INSERT INTO menus (name, price) VALUES (?,?)")
		if err != nil {
			panic(err.Error())
		}
		insertRegister.Exec(newMenu.Name, newMenu.Price)

		fmt.Fprintf(w, "New Menu added successfully")
	}
}

func deleteMenu(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	menuID, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err.Error())
	}

	stablishedConexion := conexionDB()
	registro, err := stablishedConexion.Prepare("DELETE FROM menus WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	registro.Exec(menuID)
	fmt.Fprintf(w, "The task with ID %v has been deleted successfully", menuID)
}

func updateMenu(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	menuID, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err.Error())
	}

	var updatedMenu Menu
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	json.Unmarshal(reqBody, &updatedMenu)

	stablishedConexion := conexionDB()
	registro, err := stablishedConexion.Prepare("UPDATE menus SET name=?,price=? WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	registro.Exec(updatedMenu.Name, updatedMenu.Price, menuID)
	fmt.Fprintf(w, "The menu with ID %v has been updated", menuID)
}

func main() {
	//Instantiate router
	router := mux.NewRouter().StrictSlash(true)

	//Routes
	router.HandleFunc("/", indexRoute).Methods("GET")
	router.HandleFunc("/api/menus", getMenus).Methods("GET")
	router.HandleFunc("/api/menus/{id}", getMenu).Methods("GET")
	router.HandleFunc("/api/menus", createMenu).Methods("POST")
	router.HandleFunc("/api/menus/{id}", deleteMenu).Methods("DELETE")
	router.HandleFunc("/api/menus/{id}", updateMenu).Methods("PUT")

	//Start the server
	fmt.Println("Servidor corriendo...")
	log.Fatal(http.ListenAndServe(":3000", router))
}
