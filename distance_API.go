package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
)

/*
	Structures necessary to read from API and 
	unmarshal(split JSON into separate variables)
*/

type distanceDuration struct {
	Distance float64
	Duration float64
}

type Response struct {
	Statuscode string    `json:"code"`
	Details    []Details `json:"routes"`
}

type Details struct {
	Distance float64 `json:"distance"`
	Duration float64 `json:"duration"`
}

/* 
	Structures necessary to convert variables
	into JSON notation
*/

type RoutDetails struct {
	Destination string `json:"destination"`
	Duration    float64 `json:"duration"`
	Distance    float64 `json:"distance"`
}

type Route struct{
	SourceElm string `json:"source"`
	RouteElem []RoutDetails `json:"Routes"`
}

/*
	Function that read data from external API
	based on arguments that user provide
*/

func readFromAPI(source string, destination string) []distanceDuration {
	response, err := http.Get("http://router.project-osrm.org/route/v1/driving/" + source + ";" + destination + "?overview=false")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject Response
	json.Unmarshal(responseData, &responseObject)


	distanceDurations := []distanceDuration{
		{
		Distance: responseObject.Details[0].Distance, 
		Duration: responseObject.Details[0].Duration},
	}

	return distanceDurations
}

/*
	Function that read data from external API
	based on arguments that user provide
*/
func queryParamDisplayHandler(res http.ResponseWriter, r *http.Request) {

	res.Header().Set("Content-Type", "application/json")

	src := r.URL.Query()
	sourceRoute := src["src"] // checking how many src parameters are there,


	if  len(sourceRoute) <= 0 {
		io.WriteString(res,"params not present")
		return
	}
	if len(sourceRoute) > 1 {
		io.WriteString(res, "there can be only one parameter as source")
		return
	}


	// if only one expected
	dst := r.URL.Query().Get("dst")
	if dst != "" {

	}

	query := r.URL.Query()
	destinationRecord := query["dst"] //checking how many destination parameters are there


	if  len(destinationRecord) <= 0 {
		io.WriteString(res,"missed destination parameters\n add for instance : '&dst=13.428555,52.523219' ")
		return
	}


	var routeJSON Route
	routeJSON.SourceElm = sourceRoute[0]
	

	for i := 0; i < len(destinationRecord); i++ {
		readFromAPIResult := readFromAPI(sourceRoute[0], destinationRecord[i])

	routeJSON.RouteElem =append(routeJSON.RouteElem ,RoutDetails{
		Destination: destinationRecord[i],
		Duration: readFromAPIResult[0].Duration,
		Distance: readFromAPIResult[0].Distance,
	})
}

	
/*
 Sorting function taken from Golang documentation
*/
	sort.Slice(routeJSON.RouteElem, func(i, j int) bool {
		return routeJSON.RouteElem[i].Duration < routeJSON.RouteElem[j].Duration 
	})



	json.NewEncoder(res).Encode(routeJSON)
}

func main() {
	 

	http.HandleFunc("/routes", func(res http.ResponseWriter, req *http.Request) {
		queryParamDisplayHandler(res, req)
	})
	http.ListenAndServe(":8080", nil)

}
