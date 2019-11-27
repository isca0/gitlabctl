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
	"strconv"
)

//projectPages brings model.Projects to this package
type projectPages model.Projects

type Projects struct {
	Project []projectPages
}

// list projects on gitlab
func (pj projectPages) list(client *http.Client, url, token string) (box Projects, err error) {

	items := []projectPages{}
	box = Projects{items}

	get := handlers.Requester{
		Client: client,
		Url:    url + token}

	opts := "&per_page=40"
	totalpages := handlers.ScanTotalPages(client, get.Url+opts)
	opts = opts + "&page="

	for page := 1; page <= totalpages; page++ {
		get.Url = url + token + opts + strconv.Itoa(page)
		_, pages := get.Req()
		err = json.Unmarshal(pages, &pj)
		if err != nil {
			return box, err
		}
		for _, p := range pj {
			fmt.Println(p.Namespace.FullPath + "/" + p.Name)
		}
	}

	return box, nil

}
