package main

import (
  "html/template"
  ga "github.com/ralfonso-directnic/wgoth"
  "os"
  "github.com/markbates/goth"
  "net/http"
  "log"
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


func main(){

    os.Setenv("SESSION_SECRET", "asdfasfasdf")//Required
    os.Setenv("WGOTH_HOST","127.0.0.1")//Required
    os.Setenv("WGOTH_PORT","8087")//Required

    //each plugin requires its onw pluginname_key,secret etc
    //DIGITALOCEAN_KEY
    //OKTA_KEY


    os.Setenv("GOOGLE_KEY","madeupstringhere.apps.googleusercontent.com")
    os.Setenv("GOOGLE_SECRET","abc123")

    ga.Init("google","","")
    ga.AuthListen(loginTemplate,func(user goth.User,res http.ResponseWriter,req *http.Request){
        
         log.Printf("%+v",user)
         t, _ := template.New("").Parse(userTemplate)
		 t.Execute(res,user)
		 
		 //write out a result, redirect,save a session,etc
        
    })
    
    
}
