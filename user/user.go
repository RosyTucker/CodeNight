package user

type User struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Description string `json:"description"`
}

// "github.com/google/go-github/github"

// oauthClient := oauthConf.Client(oauth2.NoContext, token)
// client := github.NewClient(oauthClient)
// user, _, err := client.Users.Get("")
// if err != nil {
// 	log.Printf("ERROR: client.Users.Get() failed with '%s'\n", err)
// 	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
// 	return
// }
