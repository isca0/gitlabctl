package cmd

import (
	"flag"
	"fmt"
	"net/http"
)

var (
	ls        string
	getGroups = "https://gitlab.com/api/v4/groups/?private_token="
	getProj   = "https://gitlab.com/api/v4/projects/?private_token="
)

//Flags declares CLI arguments to be used.
func Flags() {

	flag.StringVar(&ls, "ls", "group", "List gitlab resources.")
	defer flag.Parse()

}

func Run(client *http.Client, fToken, dToken string) {

	switch {
	case ls == "group":
		g := groupPages{}
		groups, _ := g.list(client, getGroups, fToken)
		for _, grp := range groups.Group {
			fmt.Println(grp[0].FullPath)
		}
	case ls == "proj":
		p := projectPages{}
		p.list(client, getProj, fToken)
		//for _, prj := range projects.Project {
		//	fmt.Println(prj[0].Name)
		//}
	}

}
