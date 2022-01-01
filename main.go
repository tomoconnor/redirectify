package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const NEW_URL = "https://blog.devopstom.com/"
const RSS_URL = NEW_URL + "rss.xml"

var (
	validSlugRedirects = promauto.NewCounter(prometheus.CounterOpts{
		Name: "valid_slug_redirects_total",
		Help: "The total number of redirects.",
	})
	rssRedirects = promauto.NewCounter(prometheus.CounterOpts{
		Name: "rss_redirects_total",
		Help: "The total number of RSS redirects.",
	})
	blogishRedirects = promauto.NewCounter(prometheus.CounterOpts{
		Name: "blogish_redirects_total",
		Help: "The total number of blogish redirects.",
	})
	unknownRedirects = promauto.NewCounter(prometheus.CounterOpts{
		Name: "unknown_redirects_total",
		Help: "The total number of unknown redirects.",
	})
	customRedirects = promauto.NewCounter(prometheus.CounterOpts{
		Name: "custom_redirects_total",
		Help: "The total number of custom redirects.",
	})
	rootRedirects = promauto.NewCounter(prometheus.CounterOpts{
		Name: "root_redirects_total",
		Help: "The total number of root redirects.",
	})
)

func redirect(w http.ResponseWriter, r *http.Request) {
	request_url := strings.Split(r.URL.Path, "/")
	log.Println("Request length:", len(request_url))
	if len(request_url) > 1 {
		if r.URL.Path != "/" {
			// log.Println(len(request_url))
			folder := request_url[1]
			log.Println("Folder: ", folder)
			log.Println("Request: ", r.URL.Path)
			if IsValidSlug(folder) { // Check if the folder is a valid slug (mistaken links so /bob/ will redirect to newurl/bob/)
				log.Println("Valid slug")
				log.Println("Redirecting to: ", NEW_URL+folder)
				validSlugRedirects.Inc()
				http.Redirect(w, r, NEW_URL+request_url[1], http.StatusMovedPermanently)
			} else if r.URL.Path == "/rss" {
				log.Println("Redirecting to: ", RSS_URL)
				rssRedirects.Inc()
				http.Redirect(w, r, RSS_URL, http.StatusMovedPermanently)
			} else if r.URL.Path == "/blogish" {
				blogishRedirects.Inc()
				log.Println("Redirecting to: ", NEW_URL)
				http.Redirect(w, r, NEW_URL, http.StatusMovedPermanently)
			} else if r.URL.Path == "/blogish/rss" {
				log.Println("Redirecting to: ", RSS_URL)
				rssRedirects.Inc()
				http.Redirect(w, r, RSS_URL, http.StatusMovedPermanently)
			} else { // invalid slug
				if len(request_url) > 2 { // then there/s a second part to the path like /x/y or /x/y/z
					slug := request_url[2]

					if slug == "epic-fail" { // special case for this as it got renamed foolishly
						log.Println("Redirected to https://blog.devopstom.com/raid-is-not-backup/")
						customRedirects.Inc()
						http.Redirect(w, r, "https://blog.devopstom.com/raid-is-not-backup/", http.StatusMovedPermanently)
					} else if folder == "blogish" {
						redirect_url := NEW_URL + slug
						log.Println("blogish: Redirected to ", redirect_url)
						blogishRedirects.Inc()
						http.Redirect(w, r, redirect_url, http.StatusMovedPermanently)
					} else if folder == "rss" {
						log.Println("RSS: Redirected to ", RSS_URL)
						rssRedirects.Inc()
						http.Redirect(w, r, RSS_URL, http.StatusMovedPermanently)
					} else {
						log.Println("Last Chance: Redirected to ", NEW_URL)
						unknownRedirects.Inc()
						http.Redirect(w, r, NEW_URL, http.StatusMovedPermanently)
					}
				} else { // no second part to the path
					log.Println("Unknown Folder: "+folder+" Redirecting to: ", NEW_URL)
					unknownRedirects.Inc()
					http.Redirect(w, r, NEW_URL, http.StatusMovedPermanently)
				}
			}
		} else { // Path is / so redirect to https://blog.devopstom.com/
			log.Println("Redirected to ", NEW_URL)
			rootRedirects.Inc()
			http.Redirect(w, r, NEW_URL, http.StatusMovedPermanently)
		}
	} else { // Everything else redirect to https://blog.devopstom.com/
		log.Println("Dubious request: ", r.URL.Path)
		log.Println("Redirecting to ", NEW_URL)
		unknownRedirects.Inc()
		http.Redirect(w, r, NEW_URL, http.StatusMovedPermanently)
	}
}

func main() {
	http.HandleFunc("/", redirect)
	http.Handle("/metrics", promhttp.Handler())
	log.Println("Server started on port 8080")
	http.ListenAndServe(":8080", nil)
}
