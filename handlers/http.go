package handlers

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

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

	if resp.StatusCode > 206 {
		err = errors.New(resp.Status + " " + get.Url)
		log.Fatal(err)
		return
	}

	return
}
