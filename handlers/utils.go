/*
Copyright © 2019 Isca <igorsca at protonmail dot com>

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
	"regexp"
	"strings"
)

// getSplit receive a string in pattern some/thing
// and a int to return on object of the splited array.
// if 1 will return the last object "thing"
// if 2 the parent object "some"
func getSplit(s string, i int) (o string) {
	re := regexp.MustCompile(`/`)
	if re.MatchString(s) {
		split := strings.Split(s, "/")
		o = split[len(split)-i]
		return o
	}
	o = s
	return
}
