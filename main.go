package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/BenLubar/memoize"
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

	client := octokit.NewClient(octokit.TokenAuth{
		AccessToken: os.Getenv("GITHUB_TOKEN"),
	})
	gists := client.Gists()

	getGist := memoize.Memoize(func(id string) (gistRes, error) {
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
		return result, nil
	}).(func(id string) (gistRes, error))

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
			"https://cdn.rawgit.com/" +
				gist.owner + "/" + id + "/raw/" + gist.version + "/" + file

		http.Redirect(w, r, newURL, 301)
	})
	_, isNow := os.LookupEnv("NOW")
	if isNow {
		log.Fatal(http.ListenAndServe(":443", nil))
	} else {
		log.Fatal(http.ListenAndServe(":3030", nil))
	}
}
