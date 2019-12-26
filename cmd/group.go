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

// search a group from a received pattern /group/subgroup
// if subgroupd exist will return id only for the subgroup
// if not will return 0
func (g *Groups) search(n, t string, c *http.Client) (id int, grp Groups, err error) {

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
		if n == strings.ToLower(name) {
			id = grp.ID
			return id, grp, nil
		}
	}
	return
}

// list print groups and projects accordanly by the received arguments.
func (g *Groups) list(n []string, t string, c *http.Client) (err error) {

	get := handlers.Requester{
		Meth:   "GET",
		Url:    groupURL + "?private_token=" + t + "&owned=true",
		Client: c,
	}

	_, b, _, err := get.Req()
	handlers.Lerror(err)

	err = json.Unmarshal(b, &g.Pages)
	handlers.Lerror(err)

	p := new(Projects)

	for _, grp := range g.Pages {

		if len(n) != 0 && strings.ToLower(grp.FullPath) != strings.ToLower(n[0]) {
			continue
		}

		fmt.Printf("\nListing ID: %d Group:%s\r\n", grp.ID, grp.FullPath)
		get.Url = groupURL + strconv.Itoa(grp.ID) + "/projects?private_token=" + t + "&owned=true"
		_, b, _, err := get.Req()
		handlers.Lerror(err)

		err = json.Unmarshal(b, &p.Pages)
		handlers.Lerror(err)
		for _, prj := range p.Pages {
			fmt.Println(prj.PathWithNamespace)
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

	fromGid, fromGroup, err := g.search(from[1], ftk, client)
	handlers.Lerror(err)
	if fromGid == 0 {
		log.Fatal("group " + fromGroup.Name + " not found")
		return
	}

	toGid, _, err := g.search(to[1], totk, client)
	handlers.Lerror(err)

	if toGid == 0 {
		g.Name = to[1]
		g.Path = to[1]
		g.Visibility = "private"
		toGid, _, _, err = g.create(totk, client)
		handlers.Lerror(err)
	}

	mID, masterGroup, err := g.search(fromGroup.Name, totk, client)
	if mID == 0 {
		destGroup := Groups{}
		destGroup.Name = fromGroup.Name
		destGroup.Path = fromGroup.Path
		destGroup.Description = fromGroup.Description
		destGroup.Visibility = fromGroup.Visibility
		destGroup.ParentID = toGid
		mID, _, masterGroup, err = destGroup.create(totk, client)
		handlers.Lerror(err)
	}

	get := handlers.Requester{
		Meth:   "GET",
		Client: client,
	}

	get.Url = groupURL + strconv.Itoa(fromGid) + "/subgroups?private_token=" + ftk
	_, b, _, err := get.Req()
	handlers.Lerror(err)
	err = json.Unmarshal(b, &g.Pages)
	handlers.Lerror(err)

	srcGroup := Groups{}
	srcGroup.ID = fromGroup.ID
	srcGroup.Name = masterGroup.Name
	srcGroup.FullPath = masterGroup.FullPath

	g.Pages = append(g.Pages, srcGroup)

	proj := new(Projects)
	for _, grp := range g.Pages {
		fmt.Println(grp.Name)

		group := Groups{}
		group.Name = grp.Name
		group.Path = grp.Path
		group.Description = grp.Description
		group.Visibility = grp.Visibility
		group.ParentID = mID
		if grp.ParentID != fromGid {
			_, pName, _ := handlers.GetSplit(grp.FullPath)
			grp.ParentID, _, _ = g.search(pName, totk, client)
			if grp.ParentID == 0 {
				log.Fatal("parent group " + pName + " doesnt exist")
			}
		}

		gid, _, _ := grp.search(grp.Name, totk, client)
		if gid == 0 {
			gid, _, _, err = group.create(totk, client)
			handlers.Lerror(err)
		}

		get.Url = groupURL + strconv.Itoa(grp.ID) + "/projects?private_token=" + ftk + "&per_page=50"
		_, b, _, _ := get.Req()
		err = json.Unmarshal(b, &proj.Pages)
		handlers.Lerror(err)

		prj := Projects{}
		for _, p := range proj.Pages {
			prj.Name = p.Name
			prj.Description = p.Description
			prj.Visibility = p.Visibility
			newP, err := prj.create(totk, gid, client)
			handlers.Lerror(err)
			copyData := model.Projects{
				HTTPURLToRepo: p.HTTPURLToRepo,
				Custom: model.CustomFlags{
					FromUser:  viper.GetString("FROMUSER"),
					FromToken: ftk,
					ClonePath: "/tmp/gitlabctl/" + p.Name,
					ToUser:    viper.GetString("TOUSER"),
					ToToken:   totk,
				},
			}
			copyData.Custom.ToRepo = newP.HTTPURLToRepo
			handlers.Clone(copyData)
			handlers.RemoteChange(copyData)
			handlers.Push(copyData)
			log.Println("clopying project " + prj.Name)
		}
	}
	return

}

// create a group or subgroup on the destination token.
func (g *Groups) create(token string, client *http.Client) (id, pid int, grp Groups, err error) {

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
	handlers.Lerror(err)

	data := strings.NewReader(string(gJSON))
	post.Url = groupURL + "?private_token=" + token
	if g.ParentID != 0 {
		post.Url = groupURL + "?private_token=" + token + "&parent_id=" + strconv.Itoa(g.ParentID)
	}
	post.Io = data

	_, b, _, err := post.Req()
	if err != nil {
		return 0, 0, grp, err
	}
	err = json.Unmarshal(b, &grp)
	if err != nil {
		log.Fatal(err)
		return grp.ID, grp.ParentID, grp, err
	}
	return grp.ID, grp.ParentID, grp, nil
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
		pid, _, _, err = g.create(token, client)
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

	gid, grp, err := g.search(from[1], token, client)
	handlers.Lerror(err)
	if gid == 0 {
		log.Fatal("group " + grp.Name + " not found")
		return
	}

	del := &handlers.Requester{
		Meth:   "DELETE",
		Client: client,
		Url:    groupURL + strconv.Itoa(gid) + "?private_token=" + token,
	}

	_, _, _, err = del.Req()
	handlers.Lerror(err)
	log.Println("group " + grp.Name + " with id: " + strconv.Itoa(gid) + " marked as deleted")
	return
}
