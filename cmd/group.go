/*
Copyright © 2019 Igor Brandao <igorsca at protonmail dot com>

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
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// groupURL is the main gitlab API endpoint used to manage groups and subgroups.
var groupURL = "https://gitlab.com/api/v4/groups/"

//Groups brings model.Group to this package.
type Groups struct {
	model.Group
	Pages []Groups
}

// GroupCreation used to create a data to POST new group.
type GroupCreation struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Description string `json:"description"`
	Visibility  string `json:"visibility"`
}

// list groups on gitlab
//func (gp groupPages) list(client *http.Client, url, token string, args []string) (box Groups, err error) {
//
//	items := []groupPages{}
//	box = Groups{items}
//
//	get := handlers.Requester{
//		Client: client,
//		Url:    url + token}
//
//	opts := "&per_page=5"
//	totalpages := handlers.ScanTotalPages(client, get.Url+opts)
//	opts = opts + "&page="
//
//	var grp int
//	for page := 1; page <= totalpages; page++ {
//		get.Url = url + token + opts + strconv.Itoa(page)
//		_, pages := get.Req()
//		err = json.Unmarshal(pages, &gp)
//		if err != nil {
//			return box, err
//		}
//
//		for _, g := range gp {
//			item := groupPages{g}
//			if len(args) <= 1 {
//				box.Group = append(box.Group, item)
//				continue
//			}
//			if item[0].FullPath == args[1] {
//				grp = item[0].ID
//			}
//			if grp == item[0].ID || grp == item[0].ParentID {
//				box.Group = append(box.Group, item)
//			}
//
//		}
//
//	}
//
//	return box, nil
//
//}

// groupList execute the method list to loop over groups and print to stdout.
//func groupList(args []string, token string, client *http.Client) {
//	g := groupPages{}
//	groupSearch := 0
//	groups, _ := g.list(client, getGroups, token, args)
//	for _, grp := range groups.Group {
//		if name != "" {
//			if name == grp[0].Path {
//				groupSearch = grp[0].ID
//			}
//			if grp[0].ID == groupSearch || grp[0].ParentID == groupSearch {
//				fmt.Println(grp[0].FullPath)
//			}
//			continue
//		}
//		fmt.Println(grp[0].FullPath)
//	}
//}

// search a group from a received pattern /group/subgroup
// if subgroupd exist will return id only for the subgroup
// if not will return 0
func (g *Groups) search(n, t string, c *http.Client) (id int, grp *Groups, err error) {

	name, _, _ := handlers.GetSplit(n)
	get := handlers.Requester{
		Meth:   "GET",
		Url:    groupURL + "?private_token=" + t + "&owned=true&search=" + name,
		Client: c,
	}

	_, b, _, err := get.Req()
	handlers.Lerror(err)

	err = json.Unmarshal(b, &g.Pages)
	handlers.Lerror(err)

	if len(g.Pages) == 0 {
		id = 0
		return
	}

	for _, grp := range g.Pages {
		n, _, _ := handlers.GetSplit(grp.FullPath)
		if n == name {
			id = grp.ID
			return id, &grp, nil
		}
	}
	return
}

// copy  from the received source to the destination.
func (g *Groups) copy(f, t string, client *http.Client) (err error) {

	from := strings.Split(f, ":")
	ftk := viper.GetString(from[0])
	to := strings.Split(t, ":")
	totk := viper.GetString(to[0])

	fromName, _, _ := handlers.GetSplit(from[1])
	toName, _, _ := handlers.GetSplit(to[1])

	fromGid, fromGroup, err := g.search(fromName, ftk, client)
	handlers.Lerror(err)
	if fromGid == 0 {
		log.Fatal("group " + fromName + " not found")
		return
	}

	toGid, toGroup, err := g.search(toName, totk, client)
	handlers.Lerror(err)

	//if toGid != 0 {
	//	newPrj, err = p.create(totk, gid, client)
	//	handlers.Lerror(err)
	//	copyData.Custom.ToRepo = newPrj.HTTPURLToRepo
	//	handlers.Clone(copyData)
	//	handlers.RemoteChange(copyData)
	//	handlers.Push(copyData)
	//	return
	//}
	//_, _, groupTree := handlers.GetSplit(to[1])
	//pid, err := g.treeCreation(groupTree, totk, client)
	//handlers.Lerror(err)
	//newPrj, err = p.create(totk, pid, client)
	//handlers.Lerror(err)
	//copyData.Custom.ToRepo = newPrj.HTTPURLToRepo
	//handlers.Clone(copyData)
	//handlers.RemoteChange(copyData)
	//handlers.Push(copyData)
	//return

	//if toGid == 0 {
	//	g.Name = to[1]
	//	g.Path = to[1]
	//	g.Visibility = "private"
	//	toGid, _, err = g.create(totk, client)
	//	handlers.Lerror(err)
	//}

	//// creating the master destination group to be used as parent
	//// for all subgroups and projects
	//gid, sourceGroup, err := g.search(from[1], ftk, client)
	//handlers.Lerror(err)

	//majorGroup := Groups{}
	//majorGroup.Name = sourceGroup.Name
	//majorGroup.Description = sourceGroup.Description
	//majorGroup.Visibility = sourceGroup.Visibility
	//majorGroup.ParentID = toGid
	//destGroupID, _, err := g.create(totk, client)
	//handlers.Lerror(err)

	//get := handlers.Requester{
	//	Meth:   "GET",
	//	Client: client,
	//}

	//// getting all the subgroups to be created and used to search projects inside
	//get.Url = groupURL + strconv.Itoa(fromGid) + "/subgroups?private_token=" + ftk
	//_, b, _, err := get.Req()
	//handlers.Lerror(err)
	//err = json.Unmarshal(b, &g.Pages)
	//handlers.Lerror(err)

	//masterGroup := Groups{}
	//masterGroup.ID = gid
	//masterGroup.Name = sourceGroup.Path
	//g.Pages = append(g.Pages, masterGroup)

	//proj := new(Projects)
	//for _, grp := range g.Pages {

	//	pgid, _, _ := g.search(grp.Name, totk, client)
	//	if pgid != 0 && grp.ParentID != 0 {
	//		grp.ParentID = destGroupID
	//	}

	//	gid, _, _ := g.create(totk, client)
	//	if gid == 0 {
	//		gid, _, _ = g.search(grp.Name, totk, client)
	//	}

	//	get.Url = getSubg + strconv.Itoa(grp.ID) + "/projects?private_token=" + ftk + "&per_page=50"
	//	_, b, _, _ := get.Req()
	//	err = json.Unmarshal(b, &proj)
	//	if err != nil {
	//		return err
	//	}

	//	for _, p := range proj.Pages {
	//		newP, _ := proj.create(totk, gid, client)
	//		fmt.Println(newP)
	//		//fmt.Println(newP)
	//		//p.Custom.NewRepo = newP.HTTPURLToRepo
	//		//fmt.Println(p.Custom)
	//		//handlers.Clone(p, viper.GetString("FROMUSER"), ftk)
	//		//handlers.RemoteChange(p)
	//		//handlers.Push(p, viper.GetString("TOUSER"), totk)
	//		fmt.Println(p.PathWithNamespace)
	//	}
	//}
	return

}

func (g *Groups) setParentGroup(mid int, name, from, token string, client *http.Client) (id int) {

	id, group, err := g.search(name, token, client)
	handlers.Lerror(err)
	if id != 0 && group.ParentID != 0 {
		id = mid
		return
	}

	name, parent, _ := handlers.GetSplit(group.FullPath)
	id, _, _ = g.search(parent, token, client)
	if parent == from {
		id = mid
	}
	id = mid
	return
}

// create a group or subgroup on the destination token.
func (g *Groups) create(token string, client *http.Client) (id, pid int, err error) {

	post := &handlers.Requester{
		Meth:   "POST",
		Client: client,
	}

	req := GroupCreation{
		Name:        g.Name,
		Path:        g.Path,
		Description: g.Description,
		Visibility:  g.Visibility,
	}
	gJSON, err := json.Marshal(req)
	if err != nil {
		return
	}

	data := strings.NewReader(string(gJSON))
	post.Url = groupURL + "?private_token=" + token
	if g.ParentID != 0 {
		post.Url = groupURL + "?private_token=" + token + "&parent_id=" + strconv.Itoa(g.ParentID)
	}
	post.Io = data

	_, b, _, err := post.Req()
	if err != nil {
		return 0, 0, err
	}
	err = json.Unmarshal(b, &g)
	if err != nil {
		log.Fatal(err)
		return g.ID, g.ParentID, err
	}
	return g.ID, g.ParentID, nil
}

// treeCreation creates a tree of groups wih subgroups when received
// a creation pattern like group/subgroup/subgroup/subgroup
// will create all groups and subgroups and return the last subgroup id
func (g *Groups) treeCreation(s []string, token string, client *http.Client) (pid int, err error) {
	for _, t := range s {
		gid, _, _ := g.search(t, token, client)
		if gid != 0 {
			pid = gid
			continue
		}
		g.Name = t
		g.Path = t
		g.Visibility = "private"
		g.Description = "automatic created by gitlabctl"
		g.ParentID = pid
		pid, _, err = g.create(token, client)
		if err != nil {
			log.Fatal(err)
			return pid, err
		}
		log.Println("group " + t + " created with id: " + strconv.Itoa(pid))
	}
	return pid, nil

}

// rm delete a group
func (g *Groups) rm(f string, client *http.Client) (err error) {

	from := strings.Split(f, ":")
	token := viper.GetString(from[0])

	name, _, _ := handlers.GetSplit(from[1])
	gid, _, err := g.search(name, token, client)
	handlers.Lerror(err)
	if gid == 0 {
		log.Fatal("group " + name + " not found")
		return
	}

	del := &handlers.Requester{
		Meth:   "DELETE",
		Client: client,
		Url:    groupURL + strconv.Itoa(gid) + "?private_token=" + token,
	}

	_, _, _, err = del.Req()
	handlers.Lerror(err)
	log.Println("group " + name + " with id: " + strconv.Itoa(gid) + " marked as deleted")
	return
}
