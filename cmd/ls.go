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

	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
	lsCmd.Flags().BoolP("group", "g", true, "List groups.")
	lsCmd.Flags().BoolP("proj", "p", true, "List projects.")
}

//groupPages brings model.Groups to this package
type groupPages model.Groups

//projectPages brings model.Projects to this package
type projectPages model.Projects

//Groups is the appended pagesGroup
type Groups struct {
	Group []groupPages
}

//Projects is the appended pagesGroup
type Projects struct {
	Project []projectPages
}

// list groups on gitlab
func (pg groupPages) list(client *http.Client, url, token string) (box Groups, err error) {

	items := []groupPages{}
	box = Groups{items}

	get := handlers.Requester{
		Client: client,
		Url:    url + token}

	opts := "&per_page=5"
	totalpages := handlers.ScanTotalPages(client, get.Url+opts)
	opts = opts + "&page="

	for page := 1; page <= totalpages; page++ {
		get.Url = url + token + opts + strconv.Itoa(page)
		_, pages := get.Req()
		err = json.Unmarshal(pages, &pg)
		if err != nil {
			return box, err
		}

		for _, g := range pg {
			item := groupPages{g}
			box.Group = append(box.Group, item)
		}

	}

	return box, nil

}

// list projects on gitlab
func (pj projectPages) list(client *http.Client, url, token string) (box Projects, err error) {

	items := []projectPages{}
	box = Projects{items}

	get := handlers.Requester{
		Client: client,
		Url:    url + token}

	opts := "&per_page=5"
	totalpages := handlers.ScanTotalPages(client, get.Url+opts)
	fmt.Println(totalpages)
	opts = opts + "&page="

	//for page := 1; page <= totalpages; page++ {
	//	get.Url = url + token + opts + strconv.Itoa(page)
	//	_, pages := get.Req()
	//	err = json.Unmarshal(pages, &pj)
	//	if err != nil {
	//		return box, err
	//	}
	//	for _, p := range pj {
	//		fmt.Println(p.Name)
	//	}
	//}

	return box, nil

}

// runLists executes the ls command with received flags.
func runList(client *http.Client, fToken, dToken string) {

	fmt.Println("group")
	//switch {
	//case lsType == "group":
	//	g := groupPages{}
	//	groupSearch := 0
	//	groups, _ := g.list(client, getGroups, fToken)
	//	for _, grp := range groups.Group {
	//		if name != "" {
	//			if name == grp[0].Path {
	//				groupSearch = grp[0].ID
	//			}
	//			if grp[0].ID == groupSearch || grp[0].ParentID == groupSearch {
	//				fmt.Println(grp[0].FullPath + "\t\t" + grp[0].Path)
	//			}
	//			continue
	//		}
	//		fmt.Println(grp[0].FullPath + "\t\t" + grp[0].Path)
	//	}
	//case lsType == "proj":
	//	p := projectPages{}
	//	p.list(client, getProj, fToken)
	//	//for _, prj := range projects.Project {
	//	//	fmt.Println(prj[0].Name)
	//	//}
	//case cp == "true":
	//	fmt.Println("thanks for use cp")
	//}

}
