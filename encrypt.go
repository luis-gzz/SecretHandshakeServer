package main

import (
    "fmt"
    "log"
    "os"
)

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
