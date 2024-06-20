package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
)

type URL struct{
	ID string `json:"id"`
	OriginalURL string `json:"original_url"`
	ShortenedURL string `json:"shortened_url"`
	CreatedDate time.Time `json:"created_date"`
}
var urlDB = make(map[string]URL)

func generated_URL(original_url string) string {
	hasher := md5.New()
	hasher.Write([]byte(original_url))
	fmt.Println(hasher, "hasher")
	data := hasher.Sum(nil)
	hash:= hex.EncodeToString(data)
	fmt.Println(data, "data", hash[:8])
	return hash[:8]
}



func main() {
	fmt.Println("running server ...")
    generated_URL("http://localhost/3000/server")
}