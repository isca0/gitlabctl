package cmd

import (
	"encoding/json"
	"fmt"
	"gitmigrate/handlers"
	"gitmigrate/model"
	"net/http"
	"strconv"
)

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
