package plugins

import "golang.org/x/net/html"

// copies every HTML attribute in a map for easier searching
func htmlAttributeToMap(e []html.Attribute) map[string]string {
	m := make(map[string]string)
	for a := range e {
		m[e[a].Key] = e[a].Val
	}
	return m
}
