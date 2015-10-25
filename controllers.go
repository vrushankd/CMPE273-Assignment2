package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Location structure
type Location struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zip     string `json:"zip"`
}

//count structure to keep the track of "_id"
type count struct {
	ID  string `bson:"_id"`
	Seq int    `bson:"seq"`
}

//bsonjsonStruct structure to store data in mongodb and json structure for displaying in POSTMAN
type bsonjsonStruct struct {
	ID          int    `bson:"_id" json:"id"`
	Name        string `bson:"name" json:"name"`
	Address     string `bson:"address" json:"address"`
	City        string `bson:"city" json:"city"`
	State       string `bson:"state" json:"state"`
	Zip         string `bson:"zip" json:"zip"`
	Coordinates `json:"coordinate"`
}

//Coordinates struct
type Coordinates struct {
	Latitude  float64 `bson:"lat" json:"lat"`
	Longitude float64 `bson:"long" json:"long"`
}

//UserController structure
type UserController struct{}

//NewUserController function
func NewUserController() *UserController {
	return &UserController{}
}

//UpdateStruct for update operation
type UpdateStruct struct {
	Address string `bson:"address" json:"address"`
	City    string `bson:"city" json:"city"`
	State   string `bson:"state" json:"state"`
	Zip     string `bson:"zip" json:"zip"`
}

//GoogleMapResponse structure for response from Google Maps API
type GoogleMapResponse struct {
	Results []struct {
		AddressComponents []struct {
			LongName  string   `json:"long_name"`
			ShortName string   `json:"short_name"`
			Types     []string `json:"types"`
		} `json:"address_components"`
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
			LocationType string `json:"location_type"`
			Viewport     struct {
				Northeast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"northeast"`
				Southwest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"southwest"`
			} `json:"viewport"`
		} `json:"geometry"`
		PlaceID string   `json:"place_id"`
		Types   []string `json:"types"`
	} `json:"results"`
	Status string `json:"status"`
}

var (
	//CoordinatesResp variable of type GoogleMapResponse
	CoordinatesResp GoogleMapResponse
)

//CreateLocation : POST operation
func (uc UserController) CreateLocation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//connect to mongodb in cloud using mongolab
	uri := "mongodb://vrushankd:Vrushank90@ds045628.mongolab.com:45628/gomongodb"
	if uri == "" {
		fmt.Println("no connection string provided")
		os.Exit(1)
	}

	sess, err := mgo.Dial(uri)
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}
	defer sess.Close()

	sess.SetSafe(&mgo.Safe{})

	collection := sess.DB("gomongodb").C("cmpe273Assgn2")

	// Stub an user to be populated from the body
	u := Location{}

	// Populate the user data
	json.NewDecoder(r.Body).Decode(&u)

	addr := strings.Replace(u.Address, " ", "+", -1)
	city := strings.Replace(u.City, " ", "+", -1)
	response, err := http.Get("https://maps.googleapis.com/maps/api/geocode/json?address=" + addr + "," + city + "," + u.State + "-" + u.Zip + "&sensor=false")
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()

		contents, err := ioutil.ReadAll(response.Body)

		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		json.Unmarshal([]byte(contents), &CoordinatesResp)
	}

	newID := getNextSequence("docid")

	resp := bsonjsonStruct{
		ID:      newID,
		Name:    u.Name,
		Address: u.Address,
		City:    u.City,
		State:   u.State,
		Zip:     u.Zip,
		Coordinates: Coordinates{
			Latitude:  CoordinatesResp.Results[0].Geometry.Location.Lat,
			Longitude: CoordinatesResp.Results[0].Geometry.Location.Lng,
		},
	}

	err = collection.Insert(resp)
	if err != nil {
		fmt.Printf("Can't insert document: %v\n", err)
		os.Exit(1)
	}
	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(resp)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", uj)
}

//GetLocation : GET operation
func (uc UserController) GetLocation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//connect to mongodb in cloud using mongolab
	uri := "mongodb://vrushankd:Vrushank90@ds045628.mongolab.com:45628/gomongodb"
	if uri == "" {
		fmt.Println("no connection string provided")
		os.Exit(1)
	}

	sess, err := mgo.Dial(uri)
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}
	defer sess.Close()

	sess.SetSafe(&mgo.Safe{})

	id := p.ByName("id")

	u := bsonjsonStruct{}

	intID, _ := strconv.Atoi(id)

	//Fetch data
	finderr := sess.DB("gomongodb").C("cmpe273Assgn2").FindId(intID).One(&u)

	if finderr != nil {
		w.WriteHeader(404)
		return
	}

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}

//UpdateLocation : PUT operation
func (uc UserController) UpdateLocation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	uri := "mongodb://vrushankd:Vrushank90@ds045628.mongolab.com:45628/gomongodb"
	if uri == "" {
		fmt.Println("no connection string provided")
		os.Exit(1)
	}

	sess, err := mgo.Dial(uri)
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}
	defer sess.Close()

	sess.SetSafe(&mgo.Safe{})

	collection := sess.DB("gomongodb").C("cmpe273Assgn2")

	id := p.ByName("id")
	intID, _ := strconv.Atoi(id)

	u := UpdateStruct{}

	// Populate the user data
	json.NewDecoder(r.Body).Decode(&u)

	addr := strings.Replace(u.Address, " ", "+", -1)
	city := strings.Replace(u.City, " ", "+", -1)

	response, err := http.Get("https://maps.googleapis.com/maps/api/geocode/json?address=" + addr + "," + city + "," + u.State + "-" + u.Zip + "&sensor=false")
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()

		contents, err := ioutil.ReadAll(response.Body)

		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		json.Unmarshal([]byte(contents), &CoordinatesResp)
	}

	b := bsonjsonStruct{}
	finderr := collection.FindId(intID).One(&b)
	if finderr != nil {
		w.WriteHeader(404)
		return
	}

	resp := bsonjsonStruct{
		ID:      intID,
		Name:    b.Name,
		Address: u.Address,
		City:    u.City,
		State:   u.State,
		Zip:     u.Zip,
		Coordinates: Coordinates{
			Latitude:  CoordinatesResp.Results[0].Geometry.Location.Lat,
			Longitude: CoordinatesResp.Results[0].Geometry.Location.Lng,
		},
	}

	if err := collection.Update(bson.M{"_id": intID}, resp); err != nil {
		w.WriteHeader(404)
		return
	}

	//Fetch data
	errcode := sess.DB("gomongodb").C("cmpe273Assgn2").FindId(intID).One(&b)

	if errcode != nil {
		w.WriteHeader(404)
		return
	}

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(b)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", uj)
}

//DeleteLocation : DELETE operation
func (uc UserController) DeleteLocation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	uri := "mongodb://vrushankd:Vrushank90@ds045628.mongolab.com:45628/gomongodb"
	if uri == "" {
		fmt.Println("no connection string provided")
		os.Exit(1)
	}

	sess, err := mgo.Dial(uri)
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}
	defer sess.Close()

	sess.SetSafe(&mgo.Safe{})

	collection := sess.DB("gomongodb").C("cmpe273Assgn2")

	id := p.ByName("id")
	intID, _ := strconv.Atoi(id)

	if err := collection.RemoveId(intID); err != nil {
		w.WriteHeader(404)
		return
	}

	// Write status
	w.WriteHeader(200)
	fmt.Fprint(w, "Delete successful")
}

//getNextSequence to track and auto increment the "_id" field each time user performs the POST operation.
func getNextSequence(name string) int {
	var doc count
	sess, err := mgo.Dial("mongodb://vrushankd:Vrushank90@ds045628.mongolab.com:45628/gomongodb")
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}
	defer sess.Close()

	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB("gomongodb").C("counter")

	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"seq": 1}},
		ReturnNew: true,
	}

	_, err1 := collection.Find(bson.M{"_id": "docid"}).Apply(change, &doc)
	if err1 != nil {
		fmt.Println("got an error finding a doc")
		os.Exit(1)
	}
	return doc.Seq
}
