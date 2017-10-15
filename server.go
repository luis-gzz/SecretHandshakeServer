package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "log"
    "os"
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

func GetKey(fileName string) []uint8 {
    f, err := os.Open(fileName)
    if err != nil {
        log.Fatal(err)
    }
    b1 := make([]byte, 1)
    key := make([]byte, 300000)
    for i := range key {
        _, err = f.Read(b1)
        if err != nil {
            fmt.Println("ya fucked", err)
        }
        if b1[0] - 48 <= 9 {
            key[i] = b1[0] - 48
        } else {
            for b1[0] - 48 > 9 {
                _, err = f.Read(b1)
                if err != nil {
                    log.Fatal(err)
                }
            }
            key[i] = b1[0] - 48
        }
    }
    //fmt.Println(key)
    return key


}

func Encode(arr []uint8, key []uint8) []uint8 {
    for i := range arr {
        if arr[i] >= 97 && arr[i] <= 122 {
            arr[i] += key[i]
            if arr[i] > 122 {
                arr[i] -= 26
            }
        } else if arr[i] >= 65 && arr[i] <= 90 {
            arr[i] += key[i]
            if arr[i] > 90 {
                arr[i] -= 26
            }
        } else if arr[i] >= 48 && arr[i] <= 57 {
            arr[i] += key[i]
            if arr[i] > 57 {
                arr[i] -= 10
            }
        }

    }
    return arr
}

func Decode(arr []uint8, key []uint8) []uint8 {
    for i := range arr {
        if arr[i] >= 97 && arr[i] <= 122 {
            if arr[i] <= 96 + key[i] {
                arr[i] += 26
            }
            arr[i] -= key[i]
        } else if arr[i] >= 65 && arr[i] <= 90 {
            if arr[i] <= 64 + key[i]{
                arr[i] += 26
            }
            arr[i] -= key[i]
        } else if arr[i] >= 48 && arr[i] <= 57 {
            if arr[i] <= 47 + key[i] {
                arr[i] += 10
            }
            arr[i] -= key[i]
        }
    }
    return arr
}

func GetText(w http.ResponseWriter, req *http.Request) {
    decoder := json.NewDecoder(req.Body)

    var t recievedText
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}

    t.Key = string(Encode([]uint8(t.Key), GetKey("E.txt") ))
    t.Post = string(Encode([]uint8(t.Post), GetKey("E.txt") ))
    err = bot.PostSelf("/r/SecretHandshakeVault", t.Key, t.Post)
    if err != nil {
        fmt.Println("Failed to fetch /r/golang: ", err)
        return
    }

	fmt.Println(t.Key)
    fmt.Println(t.Post)
}


func Retrieve(w http.ResponseWriter, req *http.Request) {

    decoder := json.NewDecoder(req.Body)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    var i retrieval
	err := decoder.Decode(&i)
    theKey := i.Key
	if err != nil {
		panic(err)
	}

    ekey := GetKey("E.txt")
    retrieveKey := string(Encode([]uint8(theKey), ekey))

    harvest, err := bot.Listing("/r/SecretHandshakeVault", "")
    if err != nil {
        fmt.Println("Failed to fetch /r/golang: ", err)
        return
    }

    var postText string
    var postNum int
    for index, items := range harvest.Posts {
       if retrieveKey == items.Title {
           postText = items.SelfText
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
            postText += comment.Body

        }

        postText = postText;
        //postText = string(Decode([]uint8(postText),ekey))
        fmt.Println(postText)
    }
    postText = string(Decode([]uint8(postText),ekey))
    i.Key = postText;
    postJson, err := json.Marshal(i)
    if err != nil {
        panic(err)
    }
    //fmt.Println(i.Key)
    //json.NewEncoder(w).Encode(i)
    w.Write(postJson)

    //fmt.Println(w)

}

func GetImage(w http.ResponseWriter, req *http.Request) {
    decoder := json.NewDecoder(req.Body)

    var i recievedImage
	err := decoder.Decode(&i)
	if err != nil {
		panic(err)
	}

    ekey := GetKey("E.txt")

    i.Key = string(Encode([]uint8(i.Key), ekey ))
    err = bot.PostSelf("/r/SecretHandshakeVault", i.Key, string(Encode([]uint8(i.Image[0:5000]), ekey )))

    time.Sleep(7000*time.Millisecond)
    fmt.Println("hey ho let's go!")
    harvest, err := bot.Listing("/r/SecretHandshakeVault", "")
    if err != nil {
        fmt.Println("Failed to fetch /r/golang: ", err)
        return
    }
    parentPost := harvest.Posts[0]
    for k := 1; k < int(math.Ceil(float64(len(i.Image) / 5000))); k++ {
        err = bot.Reply(parentPost.Name, string(Encode([]uint8(i.Image[k*5000:k*5000+5000]),ekey)))
        if err != nil {
            fmt.Println("Failed to post image ", err)
            return
        }
    }
	fmt.Println(i.Key)
    fmt.Println(len(i.Image))
}

var bot, errBot = reddit.NewBotFromAgentFile("redditStuff.agent", 0)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/newText", GetText)
    mux.HandleFunc("/newImage", GetImage)
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
