package concordances

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Financial-Times/go-fthealth/v1a"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

// ConcordanceDriver for cypher queries
var ConcordanceDriver Driver
var CacheControlHeader string

// HealthCheck does something
func HealthCheck() v1a.Check {
	return v1a.Check{
		BusinessImpact:   "Unable to respond to Public Concordances api requests",
		Name:             "Check connectivity to Neo4j - neoUrl is a parameter in hieradata for this service",
		PanicGuide:       "https://sites.google.com/a/ft.com/ft-technology-service-transition/home/run-book-library/concordance-read-api",
		Severity:         1,
		TechnicalSummary: "Cannot connect to Neo4j a instance with at least one concordance loaded in it",
		Checker:          Checker,
	}
}

// Checker does more stuff
func Checker() (string, error) {
	err := ConcordanceDriver.CheckConnectivity()
	if err == nil {
		return "Connectivity to neo4j is ok", err
	}
	return "Error connecting to neo4j", err
}

// Ping says pong
func Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

//GoodToGo returns a 503 if the healthcheck fails - suitable for use from varnish to check availability of a node
func GoodToGo(writer http.ResponseWriter, req *http.Request) {
	if _, err := Checker(); err != nil {
		writer.WriteHeader(http.StatusServiceUnavailable)
	}
}

// BuildInfoHandler - This is a stop gap and will be added to when we can define what we should display here
func BuildInfoHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "build-info")
}

// MethodNotAllowedHandler handles 405
func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	return
}

// GetConcordances is the public API
func GetConcordances(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if uuid == "" {
		http.Error(w, "uuid required", http.StatusBadRequest)
		return
	}
	concordance, found, err := ConcordanceDriver.Read(uuid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	if !found {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message":"Concordance not found."}`))
		return
	}
	Jason, _ := json.Marshal(concordance)
	log.Debugf("Concordance(uuid:%s): %s\n", Jason)
	w.Header().Set("Cache-Control", CacheControlHeader)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(concordance)
}
