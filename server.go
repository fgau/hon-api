package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

type Person struct {
	ID       string `json:"id,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Gender   string `json:"gender,omitempty"`
	PixUrl   string `json:"pixurl,omitempty"`
}

type JsonError struct {
	Error string `json:"error,omitempty"`
}

type JsonResult struct {
	Result string `json:"result,omitempty"`
}

type JsonVote struct {
	ID       string `json:"id,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Vote     string `json:"vote,omitempty"`
}

func logHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		w.Header().Set("Content-Type", "application/json")
		fn.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("[%s] %s %s %q %v\n",
			r.Method, r.Proto, r.RemoteAddr, r.URL.String(), t2.Sub(t1))
	}
}

func parseHon(strhon string, src string) string {
	r, _ := regexp.Compile(strhon)
	result := r.FindAllString(src, -1)
	return result[0]
}

func GetPersonEndpoint(w http.ResponseWriter, req *http.Request) {
	var hon_url string
	params := mux.Vars(req)
	if params["gender"] == "female" {
		hon_url = "http://www.hotornot.de/index.php/?changegender=w"
	} else if params["gender"] == "male" {
		hon_url = "http://www.hotornot.de/index.php/?changegender=m"
	} else {
		var hon_err = new(JsonError)
		hon_err.Error = "gender parameter must be 'female' or 'male'"
		json.NewEncoder(w).Encode(hon_err)
		return
	}
	req, err := http.NewRequest("GET", hon_url, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	req.Header.Set("User-Agent", "Golang_Spider_Bot/3.0")
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer resp.Body.Close()

	htmlData, err := ioutil.ReadAll(resp.Body)

	src := string(htmlData)

	pix_result := parseHon("\\<div style=\"background:transparent url.*?no-repeat;\"\\>", src)
	nick_result := parseHon("\\<title\\>.*?\\</title\\>", src)
	id_result := parseHon("\\<a href=\"/index.php.*?\" class=\"ButtonLink", src)

	var people = new(Person)
	people.ID = id_result[24 : len(id_result)-19]
	people.Nickname = nick_result[68 : len(nick_result)-9]
	people.Gender = params["gender"]
	people.PixUrl = pix_result[40 : len(pix_result)-15]

	json.NewEncoder(w).Encode(people)

	if resp.StatusCode == http.StatusOK {
		log.Println(resp.Header)
	}

	return
}

func PostHonResult(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	var v JsonVote
	err = json.Unmarshal(body, &v)
	if err != nil {
		panic(err)
	}
	log.Println(v.ID, v.Nickname, v.Vote)
	defer req.Body.Close()

	var hon_res = new(JsonResult)
	hon_res.Result = "received data succesfully"
	json.NewEncoder(w).Encode(hon_res)
	return
}

func RestEndpoint(w http.ResponseWriter, req *http.Request) {
	var hon_err = new(JsonError)
	hon_err.Error = "not a valid api parameter"
	json.NewEncoder(w).Encode(hon_err)
	return
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/getperson/{gender}", logHandler(GetPersonEndpoint)).Methods("GET")
	router.HandleFunc("/voteperson", logHandler(PostHonResult)).Methods("POST")
	router.PathPrefix("/").Handler(logHandler(RestEndpoint))
	log.Fatal(http.ListenAndServe(":8090", router))
}
