package main

import (
	"github.com/markbates/goth"
	ga "github.com/ralfonso-directnic/wgoth"
	"html/template"
	"log"
	"net/http"
	"os"
)

var loginTemplate = `<a href="/auth/{{.Provider}}">Sign in</a>`

var userTemplate = `
<p><a href="/logout/{{.Provider}}">logout</a></p>
<p>Name: {{.Name}} [{{.LastName}}, {{.FirstName}}]</p>
<p>Email: {{.Email}}</p>
<p>NickName: {{.NickName}}</p>
<p>Location: {{.Location}}</p>
<p>AvatarURL: {{.AvatarURL}} <img src="{{.AvatarURL}}"></p>
<p>Description: {{.Description}}</p>
<p>UserID: {{.UserID}}</p>
<p>AccessToken: {{.AccessToken}}</p>
<p>ExpiresAt: {{.ExpiresAt}}</p>
<p>RefreshToken: {{.RefreshToken}}</p>`

func main() {

	ga.Session("asdfasfasdf")                                                                            //Required
	ga.Init(
	      "google",
	      "", //empty for autodetect
	      viper.GetString("port"),
	      viper.GetString("ssl_key_path"),
	      viper.GetString("ssl_cert_path"),
	      viper.GetString("google_client"),
	      viper.GetString("google_secret"),
	      )
    
	ga.AuthListen(loginTemplate, func(user goth.User, res http.ResponseWriter, req *http.Request) {

		log.Printf("%+v", user)
		t, _ := template.New("").Parse(userTemplate)
		t.Execute(res, user)

		//write out a result, redirect,save a session,etc

	})

}
