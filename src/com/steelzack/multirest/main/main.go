package main

import (
	"com/steelzack/multirest/searcher"
	"path/filepath"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"log"
	"html"
	"github.com/gorilla/mux"
	"encoding/json"
	"fmt"
)


type Value1 struct {
	Value1 string
}

type Value2 struct {
	Value2 string
}

type Value1s []Value1

type Value2s []Value2

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/testservice", Index)
	router.HandleFunc("/value1Id{value2Id}", Value1FromValue2).Methods("GET")
	router.HandleFunc("/value2Id{value1Id}", Value2FromValue1).Methods("GET")

	log.Fatal(http.ListenAndServe(":8085", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is a gorilla REST service: %q", html.EscapeString(r.URL.Path))
}

func Value1FromValue2(w http.ResponseWriter, r *http.Request) {
	config, err := GetConfiguration()
	fmt.Println(config.CASSANDRA.HOST)
	keystorage := new(searcher.KeyStorage)
    keystorage.OpenDatabase(config.CASSANDRA.HOST, config.CASSANDRA.PORT)
	keystorage.Init()
	if err != nil {
		panic(err)
	}
	vars := mux.Vars(r)
	value2Id := vars["value2Id"]
	fmt.Println(value2Id)
	fmt.Println(keystorage.GetValue1FromValue2(value2Id))

	response := Value1s{Value1{Value1: keystorage.GetValue1FromValue2(value2Id)}}
	json.NewEncoder(w).Encode(response)
	keystorage.CloseDatabase()
}
func Value2FromValue1(w http.ResponseWriter, r *http.Request) {
	config, err := GetConfiguration()
	fmt.Println(config.CASSANDRA.HOST)
	keystorage := new(searcher.KeyStorage)
	keystorage.OpenDatabase(config.CASSANDRA.HOST, config.CASSANDRA.PORT)
	keystorage.Init()
	if err != nil {
		panic(err)
	}
	vars := mux.Vars(r)
	value1Id := vars["value1Id"]
	fmt.Println(value1Id)
	fmt.Println(keystorage.GetValue2sFromValue1(value1Id))

	rawvalue2s:=keystorage.GetValue2sFromValue1(value1Id)
	value2s := make([]Value2, rawvalue2s.Len())

	i:=0
	for e := rawvalue2s.Front(); e != nil; e = e.Next() {
		value2s[i] = Value2{e.Value.(string)}
		i++;
	}

	response := value2s
	json.NewEncoder(w).Encode(response)
	keystorage.CloseDatabase();
}

func GetConfiguration() (searcher.Config, error) {
	configuration := searcher.Config{}
	filename, err := filepath.Abs("./properties.yml")
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
	}
	err = yaml.Unmarshal([]byte(string(yamlFile)), &configuration)
	if err != nil {
		log.Println(err)
	}
	return configuration, err
}

