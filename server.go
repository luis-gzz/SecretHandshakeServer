package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/rs/cors"
    "github.com/turnage/graw/reddit"
    "math"
    "time"
    //"io"
)

type recievedText struct {
	Key string
    Post string
}

type recievedImage struct {
	Key string
    Image string
}

type retrieval struct {
	Key string
}

func SetText(w http.ResponseWriter, req *http.Request) {
    decoder := json.NewDecoder(req.Body)

    var t recievedText
    err := decoder.Decode(&t)
    if err != nil {
        panic(err)
    }

    ekey := GetKey("E.txt")
    postLen := len([]rune(t.Post))

    if (len([]rune(t.Post)) > 5000) {

        //Post the first 5000 characters in the self text
        t.Key = string(Encode([]uint8(t.Key), ekey ))
        err = bot.PostSelf("/r/SecretHandshakeVault", t.Key, string(Encode([]uint8(t.Post[0:5000]), ekey )))
        postLen = postLen - 5000;

        // Sleep to allow the parent post to be avaliable
        //Harvest the available posts and use the parent
        time.Sleep(7000*time.Millisecond)
        harvest, err := bot.Listing("/r/SecretHandshakeVault", "")
        if err != nil {
            fmt.Println("Failed to fetch /r/golang: ", err)
            return
        }

        parentPost := harvest.Posts[0]
        //Post the remainder of the image in the comments
        for k := 1; k < int(math.Ceil(float64(len([]rune(t.Post))) / 5000)); k++ {

            if (postLen > 5000){
                err = bot.Reply(parentPost.Name, string(Encode([]uint8(t.Post[k*5000:k*2*5000+1]),ekey)))
                postLen = postLen - 5000;
            } else {
                err = bot.Reply(parentPost.Name, string(Encode([]uint8(t.Post[k*5000:]),ekey)))

            }
            if err != nil {
                fmt.Println("Failed to post image ", err)
            }

        }

    } else {
        t.Key = string(Encode([]uint8(t.Key), ekey ))
        err = bot.PostSelf("/r/SecretHandshakeVault", t.Key, string(Encode([]uint8(t.Post), ekey )))

    }

	fmt.Println(t.Key)
    fmt.Println(len(t.Post))
}

func SetImage(w http.ResponseWriter, req *http.Request) {
    decoder := json.NewDecoder(req.Body)

    var i recievedImage
    err := decoder.Decode(&i)
    if err != nil {
        panic(err)
    }

    ekey := GetKey("E.txt")
    postLen := len([]rune(i.Image))

    if (len([]rune(i.Image)) > 5000) {

        //Post the first 5000 characters in the self text
        i.Key = string(Encode([]uint8(i.Key), ekey ))
        err = bot.PostSelf("/r/SecretHandshakeVault", i.Key, string(Encode([]uint8(i.Image[0:5000]), ekey )))
        postLen = postLen - 5000;

        // Sleep to allow the parent post to be avaliable
        //Harvest the available posts and use the parent
        time.Sleep(7000*time.Millisecond)
        harvest, err := bot.Listing("/r/SecretHandshakeVault", "")
        if err != nil {
            fmt.Println("Failed to fetch /r/SecretHandshakeVault: ", err)
            return
        }

        parentPost := harvest.Posts[0]
        //Post the remainder of the image in the comments
        for k := 1; k < int(math.Ceil(float64(len([]rune(i.Image))) / 5000)); k++ {

            if (postLen > 5000){
                err = bot.Reply(parentPost.Name, string(Encode([]uint8(i.Image[k*5000:k*2*5000+1]),ekey)))
                postLen = postLen - 5000;
            } else {
                err = bot.Reply(parentPost.Name, string(Encode([]uint8(i.Image[k*5000:]),ekey)))

            }
            if err != nil {
                fmt.Println("Failed to post image ", err)
            }

        }

    } else {
        i.Key = string(Encode([]uint8(i.Key), ekey ))
        err = bot.PostSelf("/r/SecretHandshakeVault", i.Key, string(Encode([]uint8(i.Image), ekey )))
    }

	fmt.Println(i.Key)
    fmt.Println(len(i.Image))
}

func Retrieve(w http.ResponseWriter, req *http.Request) {
    decoder := json.NewDecoder(req.Body)

    var i retrieval
	err := decoder.Decode(&i)
	if err != nil {
		panic(err)
	}

    ekey := GetKey("E.txt")
    retrieveKey := string(Encode([]uint8(i.Key), ekey))

    harvest, err := bot.Listing("/r/SecretHandshakeVault", "")
    if err != nil {
        fmt.Println("Failed to fetch /r/SecretHandshakeVault ", err)
        return
    }

    var postText string
    var postNum int
    for index, items := range harvest.Posts {
       if retrieveKey == items.Title {
           postText = string(Decode([]uint8(items.SelfText),ekey))
           postNum = index
           //fmt.Println(items.SelfText)
       }
    }
    if int(harvest.Posts[postNum].NumComments) > 0 {
        harvest1, err := bot.Listing("/r/SecretHandshakeVault/comments/" + harvest.Posts[postNum].ID, "")
        if err != nil {
            fmt.Println("Failed to fetch the post: ", err)
            return
        }
        for _, comment := range harvest1.Posts[0].Replies {
            postText = postText + string(Decode([]uint8(comment.Body),ekey))

        }

        //postText = postText;
        //postText = string(Decode([]uint8(postText),ekey))
        fmt.Println(postText)
    }

    \

    //postText = string(Decode([]uint8(postText),ekey))
    i.Key = postText;
    postJson, err := json.Marshal(i)
    if err != nil {
        panic(err)
    }

    // Setup and write the response back to the page
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(postJson)


    //fmt.Println(w)

}



var bot, errBot = reddit.NewBotFromAgentFile("redditStuff.agent", 0)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/newText", SetText)
    mux.HandleFunc("/newImage", SetImage)
    mux.HandleFunc("/retrieve", Retrieve)


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
