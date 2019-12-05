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
	"fmt"
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

// search a group if it exist return the id if not returns id 0.
func (g *Groups) search(n, t string, c *http.Client) (id int, grp *Groups, err error) {

	name, _, _ := handlers.GetSplit(n)
	get := handlers.Requester{
		Meth:   "GET",
		Url:    groupURL + "?private_token=" + t + "&owned=true&search=" + name,
		Client: c,
	}

	_, b, _, err := get.Req()
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &g.Pages)
	if err != nil {
		return
	}
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
	//to := strings.Split(t, ":")
	//totk := viper.GetString(to[0])

	gid, _, err := g.search(from[1], ftk, client)
	if err != nil {
		return
	}
	fmt.Println(gid)

	//p := projectPages{}
	//_, _, err = p.search(from[1], ftk, client)
	//if err != nil {
	//	return
	//}
	return

}

//	g := groupPages{}
//	get := handlers.Requester{
//		Meth:   "GET",
//		Client: client,
//	}
//
//	//master group creation
//	mgid, _ := g.searchGroup(totk, to[1], client)
//	if mgid == 0 {
//		grp := model.Group{
//			Name:       to[1],
//			Path:       to[1],
//			Visibility: "private",
//		}
//		mgid, _, _ = g.createGroup(grp, totk, client)
//		fmt.Printf("creating master group: %s", to[1])
//	}
//	fmt.Printf(" Master partend ID will be: %d\r\n", mgid)
//
//	gid, _ := g.searchGroup(ftk, from[1], client)
//	get.Url = getSubg + strconv.Itoa(gid) + "/subgroups?private_token=" + ftk
//	_, b := get.Req()
//	err = json.Unmarshal(b, &g)
//	if err != nil {
//		return err
//	}
//	mg := model.Group{
//		ID:   gid,
//		Name: to[1],
//	}
//	g = append(g, mg)
//
//	proj := projectPages{}
//	for _, grp := range g {
//
//		pgid := g.setParentGroup(mgid, totk, from[1], grp, client)
//		fmt.Printf("group %s with parent id: %d\r\n", grp.Name, pgid)
//		grp.ParentID = pgid
//		gid, _, _ := g.createGroup(grp, totk, client)
//		if gid == 0 {
//			gid, _ = g.searchGroup(totk, grp.Name, client)
//		}
//
//		get.Url = getSubg + strconv.Itoa(grp.ID) + "/projects?private_token=" + ftk + "&per_page=50"
//		_, b := get.Req()
//		err = json.Unmarshal(b, &proj)
//		if err != nil {
//			return err
//		}
//
//		for _, p := range proj {
//			p.Custom.BareRepo = false
//			p.Custom.ClonePath = "/tmp/gitlabctl/" + p.Name
//			newP, _ := proj.create(p, totk, strconv.Itoa(gid), client)
//			fmt.Println(newP)
//			p.Custom.NewRepo = newP.HTTPURLToRepo
//			fmt.Println(p.Custom)
//			handlers.Clone(p, viper.GetString("FROMUSER"), ftk)
//			handlers.RemoteChange(p)
//			handlers.Push(p, viper.GetString("TOUSER"), totk)
//			fmt.Println(p.PathWithNamespace)
//		}
//	}
//
//	return
//}
//
//func (g groupPages) setParentGroup(mid int, token, from string, grp model.Group, client *http.Client) (id int) {
//
//	id, _ = g.searchGroup(token, grp.Name, client)
//	if id != 0 && grp.ParentID != 0 {
//		id = mid
//		//fmt.Printf("already exist ")
//		return
//	}
//	re := regexp.MustCompile(`/`)
//	if re.MatchString(grp.FullPath) {
//		slash := strings.Split(grp.FullPath, "/")
//		id, _ = g.searchGroup(token, slash[len(slash)-2], client)
//		if slash[len(slash)-2] == from {
//			id = mid
//		}
//		id = mid
//	}
//	return
//}
//
//func (gp groupPages) searchGroup(t, n string, client *http.Client) (int, error) {
//
//	get := handlers.Requester{
//		Client: client,
//		Url:    getGroups + t + "&search=" + n,
//		Meth:   "GET",
//	}
//
//	_, b := get.Req()
//	err := json.Unmarshal(b, &gp)
//	if err != nil {
//		return 0, err
//	}
//	if len(gp) <= 0 {
//		err = errors.New("failed, group inexistent")
//		return 0, err
//	}
//	return gp[0].ID, nil
//}
//
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
