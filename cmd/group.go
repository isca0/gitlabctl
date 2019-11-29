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
	"errors"
	"fmt"
	"gitlabctl/handlers"
	"gitlabctl/model"
	"net/http"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

//groupPages brings model.Groups to this package
type groupPages model.Groups

//groupSolo brings model.Group to this package
type groupSolo model.Group

//Groups is the appended pagesGroup
type Groups struct {
	Group []groupPages
}

// list groups on gitlab
func (gp groupPages) list(client *http.Client, url, token string, args []string) (box Groups, err error) {

	items := []groupPages{}
	box = Groups{items}

	get := handlers.Requester{
		Client: client,
		Url:    url + token}

	opts := "&per_page=5"
	totalpages := handlers.ScanTotalPages(client, get.Url+opts)
	opts = opts + "&page="

	var grp int
	for page := 1; page <= totalpages; page++ {
		get.Url = url + token + opts + strconv.Itoa(page)
		_, pages := get.Req()
		err = json.Unmarshal(pages, &gp)
		if err != nil {
			return box, err
		}

		for _, g := range gp {
			item := groupPages{g}
			if len(args) <= 1 {
				box.Group = append(box.Group, item)
				continue
			}
			if item[0].FullPath == args[1] {
				grp = item[0].ID
			}
			if grp == item[0].ID || grp == item[0].ParentID {
				box.Group = append(box.Group, item)
			}

		}

	}

	return box, nil

}

// groupList execute the method list to loop over groups and print to stdout.
func groupList(args []string, token string, client *http.Client) {
	g := groupPages{}
	groupSearch := 0
	groups, _ := g.list(client, getGroups, token, args)
	for _, grp := range groups.Group {
		if name != "" {
			if name == grp[0].Path {
				groupSearch = grp[0].ID
			}
			if grp[0].ID == groupSearch || grp[0].ParentID == groupSearch {
				fmt.Println(grp[0].FullPath)
			}
			continue
		}
		fmt.Println(grp[0].FullPath)
	}
}

// groupCopy - copy groups from the received source to the destination.
func groupCopy(f, t string, client *http.Client) (err error) {
	from := strings.Split(f, ":")
	ftk := viper.GetString(from[0])
	//to := strings.Split(t, ":")
	//totk := viper.GetString(to[0])

	g := groupPages{}
	//g.createGroup(groupSolo{}, totk, to[1], client)

	get := handlers.Requester{
		Meth:   "GET",
		Client: client,
	}

	gid, _ := g.searchGroup(ftk, from[1], client)
	get.Url = getSubg + strconv.Itoa(gid) + "/subgroups?private_token=" + ftk
	_, b := get.Req()
	err = json.Unmarshal(b, &g)
	if err != nil {
		return err
	}

	//p := model.Projects{}
	for _, grp := range g {
		fmt.Println(grp.FullPath)
		//g.createGroup(grp, totk, to[1], client)

		//get.Url = getSubg + strconv.Itoa(grp.ID) + "/projects?private_token=" + ftk + "&per_page=50"
		//_, b := get.Req()
		//err = json.Unmarshal(b, &p)
		//if err != nil {
		//	return err
		//}
		//for _, prj := range p {
		//	fmt.Println(prj.PathWithNamespace)
		//}
	}

	return
}

func (gp groupPages) searchGroup(t, n string, client *http.Client) (int, error) {

	get := handlers.Requester{
		Client: client,
		Url:    getGroups + t + "&search=" + n,
		Meth:   "GET",
	}

	_, b := get.Req()
	err := json.Unmarshal(b, &gp)
	if err != nil {
		return 0, err
	}
	if len(gp) <= 0 {
		err = errors.New("failed, group inexistent")
		return 0, err
	}
	return gp[0].ID, nil
}

// createGroup - create a group or subgroup on the destination token.
func (gp groupPages) createGroup(g groupSolo, token, to string, client *http.Client) {

	post := &handlers.Requester{
		Meth:   "POST",
		Client: client,
	}

	gid, _ := gp.searchGroup(token, to, client)

	switch {
	case gid == 0 && to != g.Name:
		data := strings.NewReader(`{"description":"Group Created with gitlabctl"}`)
		post.Url = getGroups + token + "&visibility=private&name=" + to
		post.Url = post.Url + "&path=" + to
		post.Io = data
		//post.Req()

	case gid == 0:
		data := strings.NewReader(`{"description":"` + g.Description + `"}`)
		post.Url = getGroups + token + "&visibility=" + g.Visibility + "&name=" + g.Name
		post.Url = post.Url + "&path=" + g.Path + "&parent_id=" + strconv.Itoa(gid)
		post.Io = data
		//post.Req()
	}

}
