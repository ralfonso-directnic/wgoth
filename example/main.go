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

	os.Setenv("SESSION_SECRET", "asdfasfasdf")                                                                            //Required
	ga.Init("google", "127.0.0.1", "8087", "", "", "madeupstringhere.apps.googleusercontent.com", "abc123_google_secret") //after the first 5 required args, the variable args are used
	ga.AuthListen(loginTemplate, func(user goth.User, res http.ResponseWriter, req *http.Request) {

		log.Printf("%+v", user)
		t, _ := template.New("").Parse(userTemplate)
		t.Execute(res, user)

		//write out a result, redirect,save a session,etc

	})

}
