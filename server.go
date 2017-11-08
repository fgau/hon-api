package main

import (
    "os"
    "fmt"
    "log"
    "time"
    "regexp"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "github.com/gorilla/mux"
)

type Person struct {
    ID        string   `json:"id,omitempty"`
    Nickname  string   `json:"nickname,omitempty"`
    Gender    string   `json:"gender,omitempty"`
    PixUrl    string   `json:"pixurl,omitempty"`
}

var people []Person

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

func GetPersonEndpoint(w http.ResponseWriter, req *http.Request) {
    var hon_url string
    params := mux.Vars(req)
    if params["gender"] == "female" {
        hon_url = "http://www.hotornot.de/index.php/?changegender=w"
    } else if params["gender"] == "male" {
        hon_url = "http://www.hotornot.de/index.php/?changegender=m"
    }
    req, err := http.NewRequest("GET", hon_url, nil)
    if err != nil {
        log.Fatal(err.Error())
    }

    httpClient := &http.Client{}
    resp, err := httpClient.Do(req)
    if err != nil {
        log.Fatal(err.Error())
    }
    defer resp.Body.Close()

    htmlData, err := ioutil.ReadAll(resp.Body)
    fmt.Println(os.Stdout, string(htmlData))

    src := string(htmlData)

    r, _ := regexp.Compile("\\<div style=\"background:transparent url.*?no-repeat;\"\\>")
    result := r.FindAllString(src, -1)
    pix_result := result[0]

    r, _ = regexp.Compile("\\<title\\>.*?\\</title\\>")
    result = r.FindAllString(src, -1)
    nick_result := result[0]

    r, _ = regexp.Compile("\\<a href=\"/index.php.*?\" class=\"ButtonLink")
    result = r.FindAllString(src, -1)
    id_result := result[0]

    var mars = new(Person)
    mars.ID = id_result[24:len(id_result)-19]
    mars.Nickname = nick_result[68:len(nick_result)-9]
    mars.Gender = params["gender"]
    mars.PixUrl = pix_result[40:len(pix_result)-15]

    json.NewEncoder(w).Encode(mars)
    mars = &Person{}

    if resp.StatusCode == http.StatusOK {
        log.Println(resp.Header)
    }

    return
}

func main() {
    router := mux.NewRouter()
    router.HandleFunc("/getperson/{gender}", logHandler(GetPersonEndpoint)).Methods("GET")
    log.Fatal(http.ListenAndServe(":8090", router))
}
