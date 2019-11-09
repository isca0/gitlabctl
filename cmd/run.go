package cmd

import (
	"flag"
	"fmt"
	"net/http"
)

var (
	ls        string
	name      string
	cp        string
	getGroups = "https://gitlab.com/api/v4/groups/?private_token="
	getProj   = "https://gitlab.com/api/v4/projects/?private_token="
)

//Flags declares CLI arguments to be used.
func Flags() {

	flag.StringVar(&ls, "ls", "group", "List gitlab resources.")
	flag.StringVar(&name, "name", "", "Set a specific name to be used.")
	flag.StringVar(&cp, "cp", "", "Copy a group or a project")
	defer flag.Parse()

}

// Run the received commands and test the flags.
func Run(client *http.Client, fToken, dToken string) {

	switch {
	case ls == "group":
		g := groupPages{}
		groupSearch := 0
		groups, _ := g.list(client, getGroups, fToken)
		for _, grp := range groups.Group {
			if name != "" {
				if name == grp[0].Path {
					groupSearch = grp[0].ID
				}
				if grp[0].ID == groupSearch || grp[0].ParentID == groupSearch {
					fmt.Println(grp[0].FullPath + "\t\t" + grp[0].Path)
				}
				continue
			}
			fmt.Println(grp[0].FullPath + "\t\t" + grp[0].Path)
		}
	case ls == "proj":
		p := projectPages{}
		p.list(client, getProj, fToken)
		//for _, prj := range projects.Project {
		//	fmt.Println(prj[0].Name)
		//}
	case cp == "true":
		fmt.Println("thanks for use cp")
	}

}
