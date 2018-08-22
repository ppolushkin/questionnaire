// Package classification Questionnaire REST API.
//
//     Schemes: http, https
//     Host: localhost:8080
//     BasePath: /api
//     Version: 1.0.0
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
// swagger:meta
package main

import (
	"database/sql"

	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
	"strconv"
)

const (
	STATIC_DIR = "/static/"
)

//App is program context
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

//Initialize connects to database and initialize routes
func (a *App) Initialize(user, password, host, dbname string) error {
	if err := a.openDatabase(user, password, host, dbname); err != nil {
		return err
	}
	if err := a.initializeDatabase(); err != nil {
		return err
	}
	a.Router = mux.NewRouter()
	a.initializeRoutes()
	return nil
}

//Run runs http listening. Initialize should be run before
func (a *App) Run(addr string) {
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type", "api_key", "Authorization"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})
	//Setup CORS support for Swagger UI as it works on another port
	//https://github.com/swagger-api/swagger-ui/blob/master/docs/usage/cors.md
	log.Fatal(http.ListenAndServe(addr, handlers.CORS(allowedHeaders, allowedOrigins, allowedMethods)(a.Router)))
}

func (a *App) openDatabase(user string, password string, host string, dbname string) (err error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, host, dbname)
	a.DB, err = sql.Open("mysql", connectionString)
	return err
}

func (a *App) initializeDatabase() (err error) {
	const tableCreationQuery = `
    CREATE TABLE IF NOT EXISTS users
		(
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(50) NOT NULL,
			age INT NOT NULL
		)`
	_, err = a.DB.Exec(tableCreationQuery)
	return err
}

func (a *App) initializeRoutes() {

	a.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "." + STATIC_DIR + "index.html")
	})

	a.Router.PathPrefix(STATIC_DIR).
		Handler(http.StripPrefix(STATIC_DIR, http.FileServer(http.Dir("."+STATIC_DIR))))

	a.Router.HandleFunc("/api/users", a.getUsers).Methods("GET")
	a.Router.HandleFunc("/api/users", a.createUser).Methods("POST")
	a.Router.HandleFunc("/api/users/{id:[0-9]+}", a.getUser).Methods("GET")
	a.Router.HandleFunc("/api/users/{id:[0-9]+}", a.updateUser).Methods("PUT")
	a.Router.HandleFunc("/api/users/{id:[0-9]+}", a.deleteUser).Methods("DELETE")
}

func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	u := user{ID: id}
	if err := u.getUser(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "user not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

// swagger:operation GET /users getUsers
//
// Returns all users in the system
//
// ---
// produces:
// - application/json
// parameters:
// - name: start
//   in: query
//   description: start from user #
//   required: false
//   type: integer
//   format: int32
// - name: count
//   in: query
//   description: number of users to return
//   required: false
//   type: integer
//   format: int32
// responses:
//   '200':
//     description: users
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/user"
func (a *App) getUsers(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	users, err := getUsers(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, users)
}

func (a *App) createUser(w http.ResponseWriter, r *http.Request) {
	var u user
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := u.createUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, u)
}

func (a *App) updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var u user
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	u.ID = id

	if err := u.updateUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func (a *App) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	u := user{ID: id}
	if err := u.deleteUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
