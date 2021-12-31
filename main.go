package main

import (
	"log"
	"net/http"
	"strings"
)

const NEW_URL = "https://blog.devopstom.com/"
const RSS_URL = NEW_URL + "rss.xml"

func redirect(w http.ResponseWriter, r *http.Request) {
	request_url := strings.Split(r.URL.Path, "/")
	if len(request_url) > 2 {
		if r.URL.Path != "/" {
			log.Println(len(request_url))
			folder := request_url[1]
			log.Println("Folder: ", folder)
			log.Println("Request: ", r.URL.Path)
			if request_url[2] == "epic-fail" {
				log.Println("Redirected to https://blog.devopstom.com/raid-is-not-backup/")
				http.Redirect(w, r, "https://blog.devopstom.com/raid-is-not-backup/", http.StatusMovedPermanently)
			} else if folder == "blogish" {
				slug := request_url[2]
				redirect_url := NEW_URL + slug
				log.Println("Redirected to ", redirect_url)
				http.Redirect(w, r, redirect_url, http.StatusMovedPermanently)
			} else if folder == "rss" {
				log.Println("Redirected to ", RSS_URL)
				http.Redirect(w, r, RSS_URL, http.StatusMovedPermanently)
			} else {
				http.Redirect(w, r, NEW_URL, http.StatusMovedPermanently)
			}
		} else {
			http.Redirect(w, r, NEW_URL, http.StatusMovedPermanently)
		}
	} else {
		log.Println("Dubious request: ", r.URL.Path)
		http.Redirect(w, r, NEW_URL, http.StatusMovedPermanently)
	}
}

func main() {
	http.HandleFunc("/", redirect)
	log.Println("Server started on port 8080")
	http.ListenAndServe(":8080", nil)
}
