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

package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// ErrorMesg binds the received errors from API.
type ErrorMesg struct {
	Message struct {
		Name []string `json:"name"`
	} `json:"message"`
}

// Requester satisfies the Req method.
type Requester struct {
	Client *http.Client
	Url    string
	Data   []byte
	Meth   string
	Io     io.Reader
}

// ScanTotalPages returns the total pages to be paginated.
func ScanTotalPages(client *http.Client, url string) (p int) {

	get := Requester{
		Client: client,
		Url:    url}

	h, _, _, _ := get.Req()
	p, _ = strconv.Atoi(h["X-Total-Pages"][0])

	return

}

// Req a generic http method to create POST or GET Requests.
func (get *Requester) Req() (h http.Header, b []byte, resp *http.Response, err error) {

	req, err := http.NewRequest(get.Meth, get.Url, get.Io)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err = get.Client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	h = resp.Header
	b, _ = ioutil.ReadAll(resp.Body)

	emesg := ErrorMesg{}
	_ = json.Unmarshal(b, &emesg)

	switch {
	case resp.StatusCode == 400:
		for _, msg := range emesg.Message.Name {
			fmt.Println(msg)
		}
	case resp.StatusCode > 302 && resp.StatusCode != 400:
		err = errors.New(resp.Status + "\t" + string(b) + "\t" + get.Url)
		log.Fatal(err)
		return
	}
	return
}
