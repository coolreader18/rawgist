package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bluele/gcache"
	"github.com/joho/godotenv"
	"github.com/octokit/go-octokit/octokit"
)

type gistRes struct {
	owner   string
	version string
}

func main() {
	// ignore the error, it may be in the environment already
	godotenv.Load()

	var auth octokit.AuthMethod
	if token, exists := os.LookupEnv("GITHUB_TOKEN"); exists {
		auth = octokit.TokenAuth{
			AccessToken: token,
		}
	}
	client := octokit.NewClient(auth)
	gists := client.Gists()

	cache := gcache.New(50).ARC().Expiration(time.Hour).Build()
	getGist := func(id string) (gistRes, error) {
		cacheVal, _ := cache.Get(id)
		if cacheVal != nil {
			return cacheVal.(gistRes), nil
		}
		var result gistRes
		gist, res := gists.One(nil, octokit.M{
			"gist_id": id,
		})
		if res.HasError() {
			return result, res.Err
		}
		result = gistRes{
			owner:   gist.Owner.Login,
			version: gist.History[0].Version,
		}
		cache.Set(id, result)
		return result, nil
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		params := strings.Split(r.URL.String(), "/")
		if len(params) != 3 {
			http.Error(w, "Bad parameters", http.StatusBadRequest)
			return
		}

		id, file := params[1], params[2]

		gist, err := getGist(id)
		if err != nil {
			http.Error(w, "Gist not found", http.StatusNotFound)
			return
		}

		newURL :=
			"https://gistcdn.githack.com/" +
				gist.owner + "/" + id + "/raw/" + gist.version + "/" + file

		http.Redirect(w, r, newURL, 301)
	})
	port := "3030"
	_, isNow := os.LookupEnv("NOW")
	envPort, portExists := os.LookupEnv("PORT")
	if isNow {
		port = "443"
	} else if portExists {
		port = envPort
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
