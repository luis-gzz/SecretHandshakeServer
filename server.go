package main

import (
    "encoding/json"
    "fmt"
    //"io"
    "net/http"
    //"log"
    "github.com/rs/cors"
    "github.com/turnage/graw/reddit"
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


    err = bot.PostSelf("/r/SecretHandshakeVault", t.Key, t.Post)
    if err != nil {
        fmt.Println("Failed to fetch /r/golang: ", err)
        return
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

    err = bot.PostSelf("/r/SecretHandshakeVault", i.Key, i.Image[0:10])
    if err != nil {
        fmt.Println("Failed to post image ", err)
        return
    }

	fmt.Println(i.Key)
    fmt.Println(i.Image)
}

var bot, errBot = reddit.NewBotFromAgentFile("redditStuff.agent", 0)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/newText", GetText)
    mux.HandleFunc("/newImage", GetImage)


    c := cors.New(cors.Options{
        AllowedOrigins: []string{"*"},
        AllowedMethods: []string{"GET","POST","OPTIONS"},
        AllowCredentials: true,
    })

    if errBot  != nil {
        fmt.Println("Failed to create bot handle: ", errBot)
        return
    }

    handler := c.Handler(mux)
    http.ListenAndServe(":3000", handler)



}
