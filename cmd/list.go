package cmd

import (
	"encoding/json"
	"fmt"
	"gitmigrate/handlers"
	"gitmigrate/model"
	"net/http"
	"strconv"
)

type groups model.Groups

// list groups on gitlab
func (pages groups) list(client *http.Client, url, token string) (err error) {

	get := handlers.Requester{
		Client: client,
		Url:    url + token}

	opts := "&per_page=5"
	totalpages := handlers.ScanTotalPages(client, get.Url+opts)
	opts = opts + "&page="

	for page := 1; page <= totalpages; page++ {
		get.Url = url + token + opts + strconv.Itoa(page)
		fmt.Println(get.Url)
		_, group := get.Req()
		err = json.Unmarshal(group, &pages)
		if err != nil {
			return err
		}
		syncSearch <- pages
	}

	return

}

func search(group groups) {

	fmt.Println(group)

}
