package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	"gopkg.in/yaml.v2"
)

type SwapiData struct {
	Type        string   `yaml:"type"`
	ID          int64    `yaml:"id"`
	InfoRequest []string `yaml:"infoRequest"`
}

type Config struct {
	Input []SwapiData
}

func readYml(file string) Config {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("ERROR reading from file: %v\n", err)
	}
	swapi := Config{}
	yaml.Unmarshal([]byte(data), &swapi)
	return swapi

}

func getSwapiData(key SwapiData, url string) map[string]interface{} {
	var jsonMap map[string]interface{}
	url = fmt.Sprintf("%v%v/%v", url, key.Type, key.ID)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	if err := json.Unmarshal(response, &jsonMap); err != nil {
		log.Fatal(err)
	}
	return jsonMap
}

func logResult(result string, file *os.File) error {
	_, err := file.WriteString(fmt.Sprintf("%v\n", result))
	if err != nil {
		return err
	}
	file.Sync()
	return nil
}

func processData(wantedFields []string, dataMap map[string]interface{}) map[string]interface{} {
	structuredData := map[string]interface{}{}
	sort.Strings(wantedFields)
	for _, item := range wantedFields {
		structuredData[fmt.Sprint(item)] = dataMap[item]
	}
	return structuredData
}

func main() {
	log.Println("Starting script....")
	swapiAPI := "https://swapi.dev/api/"
	configData := readYml("input.yaml")
	s := []interface{}{}
	file, err := os.Create("./swapi-output.json")
	if err != nil {
		log.Fatal(err)
		file.Close()
	}
	defer file.Close()
	for _, value := range configData.Input {
		rawMap := getSwapiData(value, swapiAPI)
		data := processData(value.InfoRequest, rawMap)
		s = append(s, data)
		if err != nil {
			log.Fatal(err)
		}
	}
	json, err := json.MarshalIndent(s, "", " ")
	err = logResult(string(json), file)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Script finished successfully...")
	// keeps the pod alive for 30 seconds after script finishes so you exec into pod and cat out swapi-output.json :)
	time.Sleep(time.Second * 30)
}
