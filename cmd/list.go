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

	opts := token + "&per_page=5"
	totalpages := handlers.ScanTotalPages(client, url+opts)
	opts = opts + "&page="

	get := handlers.Requester{
		Client: client,
		Url:    url + token}

	for page := 1; page <= totalpages; page++ {
		url = url + opts + strconv.Itoa(page)
		_, group := get.Req()
		err = json.Unmarshal(group, &pages)
		if err != nil {
			return err
		}
		for _, p := range pages {
			fmt.Println(p.FullPath)
		}

	}

	return nil
}
