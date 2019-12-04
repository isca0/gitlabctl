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
	"fmt"
	"gitlabctl/handlers"
	"gitlabctl/model"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// projURL is the main gitlab API endpoint used to manage projects.
var projURL = "https://gitlab.com/api/v4/projects/"

// Projects unifies a slice of projectpages and the model.Projects o this package.
//brings model.Projects to this package
type Projects struct {
	*model.Projects
	Pages []Projects
}

// search will verify if project exist and returns its id.
func (p *Projects) search(n, t string, c *http.Client) (id int, prj *Projects, err error) {

	name, parent := handlers.GetSplit(n)
	get := handlers.Requester{
		Meth:   "GET",
		Url:    projURL + "?private_token=" + t + "&owned=true&search=" + name,
		Client: c,
	}
	//fmt.Println(get.Url)

	_, b, _, err := get.Req()
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &p.Pages)
	if err != nil {
		return
	}

	for _, prj := range p.Pages {
		if prj.ID != 0 {
			np, pid := handlers.GetSplit(strings.ToLower(prj.PathWithNamespace))
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

	fromProject, fromParent := handlers.GetSplit(from[1])
	toProject, toParent := handlers.GetSplit(to[1])

	// verification of the source
	fid, fromP, err := p.search(from[1], ftk, client)
	if err != nil {
		return
	}
	if fid == 0 {
		log.Fatal("project " + fromProject + " not found")
		return
	}

	// verification of the destination
	tid, _, err := p.search(toProject, totk, client)
	if err != nil {
		return
	}
	if tid != 0 {
		log.Fatal("the project " + toProject + " already exist")
		return
	}
	g := new(Groups)
	gid, _, err := g.search(to[1], totk, client)
	if gid != 0 {
		fmt.Println("copying the project")
		return
	}
	gname, _ := handlers.GetSplit(to[1])
	g.Name = gname
	//g.Path = gname
	//g.Visibility = "private"
	//g.Description = "automatic created by gitlabctl"
	//fmt.Println(g.Name)
	gid, _, err = g.create(totk, client)
	if err != nil {
		return
	}
	log.Println("group " + gname + " created with id: " + strconv.Itoa(gid))

	//starting group creation
	log.Println("creating group"+toProject, fromP, toParent, fromParent)
	//fromP.create(totk, toParent, client)
	return

}

// create a project using values from the received projectPages.
func (p *Projects) create(token string, parentID int, client *http.Client) (err error) {

	post := &handlers.Requester{
		Meth:   "POST",
		Client: client,
	}

	req := model.Projects{
		Name:        p.Name,
		Description: p.Description,
		Visibility:  p.Visibility,
		Namespace: model.Namespace{
			ID: parentID,
		},
	}

	pJSON, err := json.Marshal(req)
	if err != nil {
		return err
	}
	data := strings.NewReader(string(pJSON))
	post.Url = projURL + "?private_token=" + token
	post.Io = data
	fmt.Println(string(pJSON))

	//data := strings.NewReader(`{"description":"` + p.Description + `","visibility":"` + p.Visibility + `","name":"` + p.Name + `","namespace_id":"` + parentId + `"}`)
	//post.Url = projURL + "?private_token=" + token
	//post.Io = data

	//_, b, _, err := post.Req()
	//if err != nil {
	//	return
	//}
	////fmt.Println(string(b))
	//err = json.Unmarshal(b, &proj)
	//if err != nil {
	//	return proj, err
	//}
	return
}
