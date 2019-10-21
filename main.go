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
	"bytes"
	"encoding/json"
	"fmt"
	"gitmigrate/model"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func scanTotalPages(client *http.Client, url string) (p int) {

	h, _ := getRequest(client, url)
	p, _ = strconv.Atoi(h["X-Total-Pages"][0])
	return p

}

func getRequest(client *http.Client, url string) (h http.Header, b []byte) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	h = resp.Header
	b, _ = ioutil.ReadAll(resp.Body)

	return h, b

}

func postRequest(client *http.Client, url string, data []byte) (h http.Header, b []byte) {

	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		fmt.Println("AHHHH ", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("AHHHH ", err)
		return
	}
	h = resp.Header
	b, _ = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return

}

func projects(client *http.Client, fToken string) {

	url := "https://git.pgd.to/api/v4/projects"
	opts := "?private_token=" + fToken + "&per_page=5"
	totalpages := scanTotalPages(client, url+opts)
	opts = opts + "&page="

	for page := 1; page <= totalpages; page++ {
		url = url + opts + strconv.Itoa(page)
		_, body := getRequest(client, url)

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

func groups(client *http.Client, fToken, dToken string) {

	url := "https://git.pgd.to/api/v4/groups"
	opts := "?private_token=" + fToken + "&per_page=5"
	totalpages := scanTotalPages(client, url+opts)
	opts = opts + "&page="

	exturl := "https://gitlab.com/api/v4/groups/?private_token=" + dToken
	for page := 1; page <= totalpages; page++ {
		url = url + opts + strconv.Itoa(page)
		_, body := getRequest(client, url)
		jsonOuts := model.Groups{}
		err := json.Unmarshal(body, &jsonOuts)
		if err != nil {
			fmt.Println("fail to Unmarshal")
			return
		}
		for _, single := range jsonOuts {
			fmt.Println("Creating the group ", single.FullPath)
			data := model.GroupCreation{}
			data.Name = single.Name
			data.Path = single.Path
			data.Description = single.Description
			re := regexp.MustCompile(`/`)
			ok := re.MatchString(single.FullPath)

			if !ok {
				data.ParentID = 5134907
				//data.ParentID = 5134907
				jsonMarsh, _ := json.Marshal(data)
				_, body := postRequest(client, exturl, jsonMarsh)

				fmt.Println(string(body))
				continue
			}

			spl := (strings.Split(single.FullPath, "/"))
			searchOpts := "&search=" + spl[0]
			searchurl := exturl + searchOpts
			_, body := getRequest(client, searchurl)

			if string(body) != "[]" {
				parentGroup := model.Groups{}
				err := json.Unmarshal(body, &parentGroup)
				if err != nil {
					fmt.Println("fail to Unmarshal")
					return
				}
				for _, pg := range parentGroup {
					fmt.Println("Using the parent group: ", pg.ID)
					data.ParentID = pg.ID
					jsonMarsh, _ := json.Marshal(data)
					_, _ = postRequest(client, exturl, jsonMarsh)
					continue
				}

			}

			//return

			//spl := (strings.Split(single.FullPath, "/"))

		}

	}
}

func main() {

	fmt.Println("Starting GitMigrate")
	fToken := os.Getenv("FROMTOKEN")
	dToken := os.Getenv("DESTOKEN")

	client := &http.Client{
		Timeout: time.Second * 30,
	}

	groups(client, fToken, dToken)
	projects(client, fToken)
}
