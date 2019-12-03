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
	"net/http"
)

// projURL is the main gitlab API endpoint used to manage projects.
var projURL = "https://gitlab.com/api/v4/projects/"

//projectPages brings model.Projects to this package
type projectPages []model.Projects

type Projects struct {
	Project []projectPages
}

// search will verify if project exist and returns its id.
func (p projectPages) search(n, t string, c *http.Client) (id int, prj model.Projects, err error) {

	name, _ := handlers.GetSplit(n)
	get := handlers.Requester{
		Meth:   "GET",
		Url:    projURL + "?private_token=" + t + "&owned=true&search=" + name,
		Client: c,
	}
	fmt.Println(get.Url)

	_, b, _, err := get.Req()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = json.Unmarshal(b, &p)
	if err != nil {
		return
	}

	for _, prj := range p {
		fmt.Println(prj.Path)
		if prj.ID != 0 {
			id = prj.ID
			return id, prj, nil
		}
		id = 0
		return id, prj, nil
	}
	return

}

// list projects on gitlab
//func (pj projectPages) list(client *http.Client, url, token string) (box Projects, err error) {
//
//	items := []projectPages{}
//	box = Projects{items}
//
//	get := handlers.Requester{
//		Client: client,
//		Url:    url + token}
//
//	opts := "&per_page=40"
//	totalpages := handlers.ScanTotalPages(client, get.Url+opts)
//	opts = opts + "&page="
//
//	for page := 1; page <= totalpages; page++ {
//		get.Url = url + token + opts + strconv.Itoa(page)
//		_, pages := get.Req()
//		err = json.Unmarshal(pages, &pj)
//		if err != nil {
//			return box, err
//		}
//		for _, p := range pj {
//			fmt.Println(p.Namespace.FullPath + "/" + p.Name)
//		}
//	}
//
//	return box, nil
//
//}
//
//func (pj *projectPages) create(p model.Projects, token, pid string, client *http.Client) (proj *model.Projects, err error) {
//
//	post := &handlers.Requester{
//		Meth:   "POST",
//		Client: client,
//	}
//	data := strings.NewReader(`{"description":"` + p.Description + `","visibility":"` + p.Visibility + `","name":"` + p.Name + `","namespace_id":"` + pid + `"}`)
//	post.Url = UrlProj + "?private_token=" + token
//	post.Io = data
//
//	_, b := post.Req()
//	//fmt.Println(string(b))
//	err = json.Unmarshal(b, &proj)
//	if err != nil {
//		return proj, err
//	}
//	return
//}
