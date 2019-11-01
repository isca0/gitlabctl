/* This is a simple script that aims to migrate a on premisses gitlab to another gitlab
You will need an API token from each gitlab servers and Specify this tokens as variables:
FROMTOKEN="mySourceGitlab"
DESTOKEN="myNewGitlabServer"

With this two variables you can run this simple script which by default will clone everything from groups and
--bare projects and after that push to the new server.

This is a working in progress and when I write this text it's not so portable for any gitlab yet and still use lot
of hardcode and a less reusable code.
*/

package main

import (
	"fmt"
	"gitmigrate/cmd"
	"net/http"
	"os"
	"time"
)

/*func projects(client *http.Client, fToken string) {

	opts := "?private_token=" + fToken + "&per_page=5"
	//totalpages := scanTotalPages(client, fromUrl+opts)
	totalpages := scanTotalPages(client, toUrl+opts)
	opts = opts + "&page="

	for page := 1; page <= totalpages; page++ {
		fromUrl = fromUrl + opts + strconv.Itoa(page)
		_, body := getRequest(client, fromUrl)

		jsonOuts := model.Projects{}
		err := json.Unmarshal(body, &jsonOuts)
		if err != nil {
			fmt.Println("fail to Unmarshal")
			return
		}
		for _, single := range jsonOuts {
			fmt.Println(single.SSHURLToRepo + " --> " + "git@gitlab.com:pgmais/" + single.PathWithNamespace + ".git")
			cmd := exec.Command("git", "clone", "--bare", single.SSHURLToRepo, "ok/"+single.Name)
			cmd.Run()
			cmd = exec.Command("git", "remote", "set-url", "origin", "git@gitlab.com:pgmais/"+single.PathWithNamespace+".git")
			cmd.Dir = "ok/" + single.Name
			cmd.Run()
			cmd = exec.Command("git", "push", "--all")
			cmd.Dir = "ok/" + single.Name
			cmd.Run()
			cmd = exec.Command("git", "push", "--tags")
			cmd.Dir = "ok/" + single.Name
			cmd.Run()
			cmd = exec.Command("rm", "-rf", "ok/"+single.Name)
			cmd.Run()
		}
	}

}
*/

func main() {

	fmt.Println("Starting GitMigrate")
	fToken := os.Getenv("FROMTOKEN")
	dToken := os.Getenv("DESTOKEN")

	client := &http.Client{
		Timeout: time.Second * 30,
	}

	cmd.CreateGroups(client, fToken, dToken)
	//projects(client, fToken)
}
