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
	"gitlabctl/model"
	"os/exec"
	"strings"
)

// Clone a repositorie
func Clone(p model.Projects) (err error) {
	repo := strings.Replace(p.HTTPURLToRepo, "https://", "https://"+p.Custom.FromUser+":"+p.Custom.FromToken+"@", -1)
	cmd := exec.Command("git", "clone", "--bare", repo, p.Custom.ClonePath)
	cmd.Run()
	return
}

// RemoteChange remove and add a new remote to a cloned repo.
func RemoteChange(p model.Projects) (err error) {
	path := p.Custom.ClonePath
	repo := strings.Replace(p.Custom.ToRepo, "https://", "https://"+p.Custom.ToUser+":"+p.Custom.ToToken+"@", -1)
	cmd := exec.Command("git", "remote", "set-url", "origin", repo)
	cmd.Dir = path
	cmd.Run()
	return
}

// Push the entire repositorie
func Push(p model.Projects) (err error) {
	path := p.Custom.ClonePath
	cmd := exec.Command("git", "push", "--all")
	cmd.Dir = path
	cmd.Run()
	cmd = exec.Command("git", "push", "--tags")
	cmd.Dir = path
	cmd.Run()
	cmd = exec.Command("rm", "-rf", path)
	cmd.Run()
	return
}
