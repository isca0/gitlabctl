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
	"fmt"
	"gitlabctl/model"
	"os"
	"os/exec"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

// project import model.Projects to this package.
//type project model.Projects

// Clone a repositorie
func Clone(p model.Projects, u, t string) (err error) {
	fmt.Printf("git clone --bare %s\r\n", p.WebURL)
	_, err = git.PlainClone(p.Custom.ClonePath, p.Custom.BareRepo, &git.CloneOptions{
		URL:      p.WebURL,
		Progress: os.Stdout,
		Auth: &http.BasicAuth{
			Username: u,
			Password: t,
		},
	})
	if err != nil {
		return err
	}
	return
}

// RemoteChange remove and add a new remote to a cloned repo.
func RemoteChange(p model.Projects) (err error) {

	path := p.Custom.ClonePath
	r, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	fmt.Printf("git remote rm origin.\r\n")
	err = r.DeleteRemote("origin")
	if err != nil {
		return err
	}

	fmt.Printf("git remote add origin %s\r\n", p.Custom.NewRepo)
	_, err = r.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{p.Custom.NewRepo},
	})
	if err != nil {
		return err
	}

	return

}

// Push the entire repositorie
func Push(p model.Projects, u, t string) (err error) {
	path := p.Custom.ClonePath
	r, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	rs := config.RefSpec("+refs/heads/*:refs/remotes/origin/*")
	tg := config.RefSpec("+refs/tags/*:refs/tags/*")
	err = r.Push(&git.PushOptions{
		RefSpecs: []config.RefSpec{rs, tg},
		Progress: os.Stdout,
		Auth: &http.BasicAuth{
			Username: u,
			Password: t,
		}})
	if err != nil {
		return err
	}

	//cmd := exec.Command("git", "push", "--all")
	//cmd.Dir = path
	//cmd.Run()
	//cmd = exec.Command("git", "push", "--tags")
	//cmd.Dir = path
	//cmd.Run()
	cmd := exec.Command("rm", "-rf", path)
	cmd.Run()

	return

}
