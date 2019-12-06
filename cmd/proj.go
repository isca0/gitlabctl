/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"gitlabctl/handlers"
	"gitlabctl/model"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

// projURL is the main gitlab API endpoint used to manage projects.
var projURL = "https://gitlab.com/api/v4/projects/"

// Projects unifies a slice of projectpages and the model.Projects o this package.
//brings model.Projects to this package
type Projects struct {
	model.Projects
	Pages []Projects
}

// ProjCreation used to create a data to POST new project.
type ProjCreation struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Visibility  string `json:"visibility"`
	ParentID    int    `json:"namespace_id"`
	Avatar      string `json:"avatar"`
}

// search will verify if project exist and returns its id.
func (p *Projects) search(n, t string, c *http.Client) (id int, prj *Projects, err error) {

	name, parent, _ := handlers.GetSplit(n)
	get := handlers.Requester{
		Meth:   "GET",
		Url:    projURL + "?private_token=" + t + "&owned=true&search=" + name,
		Client: c,
	}

	_, body, _, err := get.Req()
	handlers.Lerror(err)

	err = json.Unmarshal(body, &p.Pages)
	handlers.Lerror(err)

	for _, prj := range p.Pages {
		if prj.ID != 0 {
			np, pid, _ := handlers.GetSplit(strings.ToLower(prj.PathWithNamespace))
			if np == name && pid == parent {
				id = prj.ID
				return id, &prj, nil
			}
		}
	}
	id = 0
	return
}

// copy from the received source to the destination.
func (p *Projects) copy(f, t string, client *http.Client) (err error) {

	from := strings.Split(f, ":")
	ftk := viper.GetString(from[0])
	to := strings.Split(t, ":")
	totk := viper.GetString(to[0])

	fromName, _, _ := handlers.GetSplit(from[1])
	toName, _, _ := handlers.GetSplit(to[1])

	// geting information of the source
	fromID, fromProject, err := p.search(from[1], ftk, client)
	handlers.Lerror(err)

	if fromID == 0 {
		log.Fatal("project " + fromName + " not found")
		return
	}

	// geting information about the destination
	toID, _, err := p.search(toName, totk, client)
	handlers.Lerror(err)

	if toID != 0 {
		log.Fatal("the project " + toName + " already exist")
		return
	}

	p.Name = fromProject.Name
	p.Description = fromProject.Description
	p.Visibility = fromProject.Visibility
	log.Printf("copying the project %s", p.Name)

	newPrj := new(Projects)
	g := new(Groups)
	gid, _, err := g.search(to[1], totk, client)
	copyData := model.Projects{
		HTTPURLToRepo: fromProject.HTTPURLToRepo,
		Custom: model.CustomFlags{
			FromUser:  viper.GetString("FROMUSER"),
			FromToken: ftk,
			ClonePath: "/tmp/gitlabctl/" + fromProject.Name,
			ToUser:    viper.GetString("TOUSER"),
			ToToken:   totk,
		},
	}

	switch {
	// if group exist
	case gid != 0:
		newPrj, err = p.create(totk, gid, client)
		handlers.Lerror(err)
		copyData.Custom.ToRepo = newPrj.HTTPURLToRepo
		handlers.Clone(copyData)
		handlers.RemoteChange(copyData)
		handlers.Push(copyData)
	case gid == 0:
		// in case of the destination groups doesnt exist
		// will create all groups and subgroups like a tree creation
		_, _, groupTree := handlers.GetSplit(to[1])
		pid, err := g.treeCreation(groupTree, totk, client)
		handlers.Lerror(err)
		newPrj, err = p.create(totk, pid, client)
		handlers.Lerror(err)
		copyData.Custom.ToRepo = newPrj.HTTPURLToRepo
		handlers.Clone(copyData)
		handlers.RemoteChange(copyData)
		handlers.Push(copyData)
	}
	return

}

// create a project using values from the received values in Projects.
func (p *Projects) create(token string, parentID int, client *http.Client) (proj *Projects, err error) {

	post := &handlers.Requester{
		Meth:   "POST",
		Client: client,
	}

	req := ProjCreation{
		Name:        p.Name,
		Description: p.Description,
		Visibility:  p.Visibility,
		Avatar:      p.AvatarURL,
		ParentID:    parentID,
	}

	pJSON, err := json.Marshal(req)
	handlers.Lerror(err)
	data := strings.NewReader(string(pJSON))
	post.Url = projURL + "?private_token=" + token
	post.Io = data

	_, b, _, err := post.Req()
	handlers.Lerror(err)
	err = json.Unmarshal(b, &proj)
	handlers.Lerror(err)
	return proj, nil
}
