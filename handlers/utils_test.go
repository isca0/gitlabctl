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
	"testing"
)

func TestGetSplit(t *testing.T) {
	tt := []struct {
		name string
		in   string
		outI string
		outP string
	}{
		{"only group", "tools", "tools", ""},
		{"group with slash", "/tools", "tools", ""},
		{"group with parent", "/parent/group", "group", "parent"},
		{"multiple parents", "/grampa/grama/parent/group", "group", "parent"},
		{"get group with parent", "/grampa/grama/parent/group", "group", "parent"},
	}

	for _, tc := range tt {
		g, p, _ := GetSplit(tc.in)
		if g != tc.outI && p != tc.outP {
			t.Errorf("failed in %s, expected:%s, and: %s received:%s and: %s", tc.name, tc.outI, tc.outP, p, g)
		}
	}
}
