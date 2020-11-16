package wgoth

import (
	"fmt"
	"github.com/gorilla/pat"
	//"github.com/quasoft/memstore"
	"github.com/gorilla/sessions"
	 gsm "github.com/bradleypeabody/gorilla-sessions-memcache"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/markbates/goth"
	"github.com/ralfonso-directnic/wgoth/gothic"
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
var host string
var port string
var Store sessions.Store

func Session(authkey string,enckey string){
	
       memcacheClient := gsm.NewGoMemcacher(memcache.New("localhost:11211"))
       Store = gsm.NewMemcacherStore(memcacheClient, "auth_session_", []byte(authkey))
	/*
       store := memstore.NewMemStore(
		[]byte(authkey),
		[]byte(enckey),
	)*/

        gothic.Store = Store	
	

}

func Init(provider_nm string, host_str string, port_str string, sslkey_str string, sslcrt_str string, args ...string) {

	provider_name = provider_nm
	sslkey = sslkey_str
	sslcrt = sslcrt_str
	host = host_str
	port = port_str

	if len(sslkey) > 0 && len(sslcrt) > 0 {
		protocol = "https"
	} else {
		protocol = "http"
	}

	if len(host) == 0 {
		host, _ = os.Hostname()
	}

	if len(port) == 0 {

		port = "8080"
	}

	var provider goth.Provider

	switch provider_name {
	case "auth0": //key, secret,domain
		provider = auth0.New(getvar(0, args), getvar(1, args), authurl(provider_name), getvar(2, args))
	case "google": //key,secret
		provider = google.New(getvar(0, args), getvar(1, args), authurl(provider_name))
	case "digitalocean": //key,secret
		provider = digitalocean.New(getvar(0, args), getvar(1, args), authurl(provider_name), "read")
	case "okta": //id,secret,org_url
		provider = okta.New(getvar(0, args), getvar(1, args), getvar(2, args), authurl(provider_name), "openid", "profile", "email")
	case "github": //key,secret
		provider = github.New(getvar(0, args), getvar(1, args), authurl(provider_name))
	default:
		provider = google.New(getvar(0, args), getvar(1, args), authurl(provider_name))

	}

	

	goth.UseProviders(provider)	
	

}

func authurl(name string) string {

	return fmt.Sprintf("%s://%s:%s/auth/%s/callback", protocol, host, port, name)

}

func getvar(key int, val []string) string {

	if key >= len(val) {
		return ""
	}

	return val[key]

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
	
	p.Get("/status/{provider}", func(res http.ResponseWriter, req *http.Request) {

	        gothUser,_ := gothic.CompleteUserAuth(res, req)
		
		fn(gothUser, res, req)

	
	})

	p.Get("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {

		if gothUser, err := gothic.CompleteUserAuth(res, req); err == nil {

			fn(gothUser, res, req)

		} else {
			log.Println(err)
			gothic.BeginAuthHandler(res, req)
		}
	})

	p.Get("/", func(res http.ResponseWriter, req *http.Request) {

		t, _ := template.New("tmpl").Parse(loginTemplate)
		t.Execute(res, map[string]string{"Provider": provider_name})
	})

	log.Println("Listening On:", port)

	if len(sslkey) > 0 && len(sslcrt) > 0 {
		log.Println("OAuth https on")
		log.Fatal(http.ListenAndServeTLS(host+":"+port, sslcrt, sslkey, p))

	} else {
		log.Println("OAuth https off")
		log.Fatal(http.ListenAndServe(host+":"+port, p))
	}

}
