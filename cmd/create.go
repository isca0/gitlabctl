package cmd

import "net/http"

var (
	getGroups = "https://gitlab.com/api/v4/groups/?private_token="
)

func CreateGroups(client *http.Client, fToken, dToken string) {

	g := groups{}
	g.list(client, getGroups, fToken)

	//for _, single := range jsonOuts {
	//	fmt.Println("Creating the group ", single.FullPath)
	//	data := model.GroupCreation{}
	//	data.Name = single.Name
	//	data.Path = single.Path
	//	data.Description = single.Description
	//	re := regexp.MustCompile(`/`)
	//	ok := re.MatchString(single.FullPath)

	//	if !ok {
	//		data.ParentID = 5168612
	//		//data.ParentID = 5134907
	//		//data.ParentID = 5134907
	//		jsonMarsh, _ := json.Marshal(data)
	//		//_, body := postRequest(client, toUrl+dToken, jsonMarsh)
	//		fmt.Println("Tools:", jsonMarsh)

	//		fmt.Println(string(body))
	//		continue
	//	}

	//	spl := (strings.Split(single.FullPath, "/"))
	//	searchOpts := "&search=" + spl[0]
	//	searchurl := toUrl + searchOpts
	//	_, body := getRequest(client, searchurl)

	//	if string(body) != "[]" {
	//		parentGroup := model.Groups{}
	//		err := json.Unmarshal(body, &parentGroup)
	//		if err != nil {
	//			fmt.Println("fail to Unmarshal")
	//			return
	//		}
	//		for _, pg := range parentGroup {
	//			fmt.Println("Using the parent group: ", pg.ID)
	//			data.ParentID = pg.ID
	//			//jsonMarsh, _ := json.Marshal(data)
	//			fmt.Println(data)
	//			//_, _ = postRequest(client, toUrl+dToken, jsonMarsh)

	//			continue
	//		}

	//	}

	//	//return

	//	//spl := (strings.Split(single.FullPath, "/"))

	//}
}
