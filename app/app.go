package app

import (
	"encoding/json"
	"fmt"
	"github.com/alexrios/challenge-api/requests"
	"github.com/alexrios/challenge-api/starwarsapi"
	"github.com/asaskevich/govalidator"
	"github.com/pytimer/mux-logrus"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"

	"github.com/alexrios/challenge-api/db"
	"github.com/alexrios/challenge-api/models"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	govalidator.SetFieldsRequiredByDefault(true)
}

type App struct {
	DB     db.Transport
	Router *mux.Router
}

func NewApp() (*App, error) {
	mongoURL := os.Getenv("MONGO_URL")
	mongo, err := db.NewMongoTransport(mongoURL)
	if err != nil {
		log.WithFields(log.Fields{
			"cause": err.Error(),
		}).Fatal("Could not start app")
		return nil, err
	}
	app := &App{
		DB: mongo,
	}

	app.MakeRoutes()
	return app, nil
}

func (a *App) Run(addr string) {
	bootPhaseLogger := log.WithFields(log.Fields{
		"phase": "boot",
	})
	bootPhaseLogger.WithFields(log.Fields{
		"addr": addr,
	}).Info("App started and listening")
	a.Router.Use(muxlogrus.NewLogger().Middleware)
	bootPhaseLogger.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) MakeRoutes() {
	routeRegLogger := log.WithFields(log.Fields{
		"phase": "boot",
		"event": "registering route",
	})
	basePath := "/planets"
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc(basePath, RateLimitedHandler(a.insertPlanet)).Methods("POST")
	routeRegLogger.WithFields(log.Fields{
		"http-method":      "POST",
		"function-handler": "app.insertPlanet",
		"path":             basePath,
	}).Info("ADD PLANET")

	router.HandleFunc(basePath, RateLimitedHandler(a.getPlanetById)).Queries("id", "{id}").Methods("GET")
	routeRegLogger.WithFields(log.Fields{
		"http-method":      "GET",
		"function-handler": "app.getPlanetById",
		"path":             basePath,
		"queries":          "?id={id}",
	}).Info("GET PLANET BY ID")

	router.HandleFunc(basePath, RateLimitedHandler(a.getPlanetByName)).Queries("name", "{name}").Methods("GET")
	routeRegLogger.WithFields(log.Fields{
		"http-method":      "GET",
		"function-handler": "app.getPlanetByName",
		"path":             basePath,
		"queries":          "?name={name}",
	}).Info("GET PLANET BY NAME")

	router.HandleFunc(basePath, RateLimitedHandler(a.listPlanets)).Methods("GET")
	routeRegLogger.WithFields(log.Fields{
		"http-method":      "GET",
		"function-handler": "app.listPlanets",
		"path":             basePath,
	}).Info("LIST PLANETS")

	router.HandleFunc(basePath, RateLimitedHandler(a.deletePlanet)).Queries("id", "{id}").Methods("DELETE")
	routeRegLogger.WithFields(log.Fields{
		"http-method":      "DELETE",
		"function-handler": "app.deletePlanet",
		"path":             basePath,
		"queries":          "?id={id}",
	}).Info("GET PLANET BY ID")
	a.Router = router
}

//Recupera os planetas cadastrados, caso ainda nao tenha nenhum planeta retorna uma lista vazia
func (a *App) listPlanets(w http.ResponseWriter, r *http.Request) {
	var planets []*models.Planet

	err := a.DB.FindAll(&planets)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if planets != nil {
		json.NewEncoder(w).Encode(planets)
	} else {
		json.NewEncoder(w).Encode([]models.Planet{})
	}
	return
}

func (a *App) getPlanetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//Valida request
	request := requests.FindByIDRequest{Id: vars["id"]}
	_, err := govalidator.ValidateStruct(request)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Tentamos pegar o planet pelo ID no banco
	var planet *models.Planet
	err = a.DB.FindByID(request.Id, &planet)
	if err != nil {
		if err == mgo.ErrNotFound {
			http.Error(w, "planet not Found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(planet)
	return
}

func (a *App) getPlanetByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//Valida findByNameReq
	findByNameReq := requests.FindByNameRequest{Name: vars["name"]}
	_, err := govalidator.ValidateStruct(findByNameReq)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var planet *models.Planet
	// Tentamos pegar o planet pelo ID no banco
	err = a.DB.FindByName(findByNameReq.Name, &planet)
	if err != nil {
		if err == mgo.ErrNotFound {
			http.Error(w, "planet not Found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(planet)
	return
}

func (a *App) insertPlanet(w http.ResponseWriter, r *http.Request) {
	var newPlanetReq *requests.NewPlanetRequest
	//Le o json de entrada
	err := json.NewDecoder(r.Body).Decode(&newPlanetReq)
	if err != nil {
		log.Print(err)
		http.Error(w, "Invalid body request", http.StatusBadRequest)
		return
	}
	//Valida os campos do json
	_, err = govalidator.ValidateStruct(newPlanetReq)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Tenta recuperar os dados que necessitam da swapi
	p, err := starwarsapi.FetchPlanet(newPlanetReq.Name)
	if err != nil {
		log.Print(err)
	}

	//Cria um novo Planeta com indentificador unico, dados do request e seus aparicoes.
	planet := models.Planet{
		ID:          bson.NewObjectId(),
		Name:        newPlanetReq.Name,
		Terrain:     newPlanetReq.Terrain,
		Climate:     newPlanetReq.Climate,
		Appearances: p.Appearances(),
	}
	err = a.DB.Insert(planet)
	if err != nil {
		if err.(*mgo.LastError).Code == db.DUPLICATED_KEY_CODE {
			http.Error(w, fmt.Sprintf("Planet %v was previously created", planet.Name), http.StatusPreconditionFailed)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&planet)
}

func (a *App) deletePlanet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//Valida findByIDReq
	findByIDReq := requests.FindByIDRequest{Id: vars["id"]}
	_, err := govalidator.ValidateStruct(findByIDReq)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Tentamos pegar o planet pelo ID no banco e remove-lo
	err = a.DB.Delete(findByIDReq.Id)
	if err != nil {
		if err == mgo.ErrNotFound {
			http.Error(w, "planet not Found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
