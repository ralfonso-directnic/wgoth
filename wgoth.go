package wgoth

import (
	"fmt"
	"github.com/gorilla/pat"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/auth0"
	"github.com/markbates/goth/providers/digitalocean"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/okta"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var provider_name string
var protocol string
var sslcrt string
var sslkey string

func Init(provider_nm string,  sslcrt_str string, sslkey_str string) {

	provider_name = provider_nm
	sslkey = sslkey_str
	sslcrt = sslcrt_str

	if len(sslkey) > 0 && len(sslcrt) > 0 {
		protocol = "https"
	} else {
		protocol = "http"
	}


	if os.Getenv("WGOTH_HOST") == "" {
		hst, _ := os.Hostname()
		os.Setenv("WGOTH_HOST", hst)
	}

	if os.Getenv("WGOTH_PORT") == "" {

		os.Setenv("WGOTH_PORT", "8000")
	}

	var provider goth.Provider

	switch provider_name {
	case "auth0":
		provider = auth0.New(os.Getenv("AUTH0_KEY"), os.Getenv("AUTH0_SECRET"),authurl(provider_name), os.Getenv("AUTH0_DOMAIN"))
	case "google":
		provider = google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), authurl(provider_name))
	case "digitalocean":
		provider = digitalocean.New(os.Getenv("DIGITALOCEAN_KEY"), os.Getenv("DIGITALOCEAN_SECRET"), authurl(provider_name), "read")
	case "okta":
		provider = okta.New(os.Getenv("OKTA_ID"), os.Getenv("OKTA_SECRET"), os.Getenv("OKTA_ORG_URL"), authurl(provider_name), "openid", "profile", "email")
	case "github":
		provider = github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), authurl(provider_name))
	default:
		provider = google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), authurl(provider_name))

	}

	goth.UseProviders(provider)

}

func authurl(name string) string {

	return fmt.Sprintf("%s://%s:%s/auth/%s/callback", protocol, os.Getenv("WGOTH_HOST"), os.Getenv("WGOTH_PORT"), name)

}

func AuthListen(loginTemplate string, fn func(user goth.User, res http.ResponseWriter, req *http.Request)) {

	//could be a path or could be a string

	f, err := ioutil.ReadFile(loginTemplate)

	if err == nil {
		loginTemplate = string(f) // we found a file with path loginTemplate so we set this to the new string, otherwise it's just a string template
	}



	p := pat.New()
	p.Get("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {

		user, err := gothic.CompleteUserAuth(res, req)
		if err != nil {
			fmt.Fprintln(res, err)
			return
		}

		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		fn(user, res, req)

	})

	p.Get("/logout/{provider}", func(res http.ResponseWriter, req *http.Request) {
		gothic.Logout(res, req)
		res.Header().Set("Location", "/")
		res.WriteHeader(http.StatusTemporaryRedirect)
	})

	p.Get("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {

		if gothUser, err := gothic.CompleteUserAuth(res, req); err == nil {

			fn(gothUser, res, req)

		} else {
			gothic.BeginAuthHandler(res, req)
		}
	})

	p.Get("/", func(res http.ResponseWriter, req *http.Request) {

		t, _ := template.New("tmpl").Parse(loginTemplate)
		t.Execute(res, map[string]string{"Provider":provider_name})
	})

	log.Println("Listening On:",os.Getenv("WGOTH_PORT"))

	if len(sslkey) > 0 && len(sslcrt) > 0 {
		log.Println("OAuth https on")
		log.Fatal(http.ListenAndServeTLS(":"+os.Getenv("WGOTH_PORT"), sslcrt, sslkey, p))

	} else {
		log.Println("OAuth https off")
		log.Fatal(http.ListenAndServe(":"+os.Getenv("WGOTH_PORT"), p))
	}

}
