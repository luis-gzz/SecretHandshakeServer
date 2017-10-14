package main

import (
    "encoding/json"
    "fmt"
    //"io"
    "net/http"
    //"log"
    "github.com/rs/cors"
)

type recievedText struct {
	Key string
    Post string
}

type recievedImage struct {
	Key string
    Image string
}


func GetText(w http.ResponseWriter, req *http.Request) {
    decoder := json.NewDecoder(req.Body)

    var t recievedText
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}

	fmt.Println(t.Key)
    fmt.Println(t.Post)
}

func GetImage(w http.ResponseWriter, req *http.Request) {
    decoder := json.NewDecoder(req.Body)

    var i recievedImage
	err := decoder.Decode(&i)
	if err != nil {
		panic(err)
	}

	fmt.Println(i.Key)
    fmt.Println(i.Image)
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/newText", GetText)
    mux.HandleFunc("/newImage", GetImage)

    // cors.Default() setup the middleware with default options being
    // all origins accepted with simple methods (GET, POST). See
    // documentation below for more options.

    c := cors.New(cors.Options{
        AllowedOrigins: []string{"*"},
        AllowedMethods: []string{"GET","POST","OPTIONS"},
        AllowCredentials: true,
    })


    handler := c.Handler(mux)
    http.ListenAndServe(":3000", handler)
}
