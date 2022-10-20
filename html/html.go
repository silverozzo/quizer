package html

import (
	"io"
	"log"

	"golang.org/x/net/html"

	"quizer/model"
)

func Parse(rsp io.Reader, inputs *[]model.Input) {
	doc, err := html.Parse(rsp)
	if err != nil {
		log.Fatalln("завалился парсинг ответа")
	}

	traverse(doc, inputs)
}

func traverse(doc *html.Node, inputs *[]model.Input) {
	if doc.Data == "Test successfully passed" {
		return
	}

	if checkInputRadio(doc, inputs) {
		return
	}

	if checkInputText(doc, inputs) {
		return
	}

	if checkInputSelect(doc, inputs) {
		return
	}

	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		traverse(c, inputs)
	}

	return
}

func checkInputSelect(n *html.Node, inputs *[]model.Input) bool {
	if n.Type != html.ElementNode {
		return false
	}

	if n.Data != "option" {
		return false
	}

	if n.Parent == nil {
		return false
	}

	if n.Parent.Type != html.ElementNode || n.Parent.Data != "select" {
		return false
	}

	nm, ok := getAttribute(n.Parent, "name")
	if !ok {
		return false
	}

	vl, ok := getAttribute(n, "value")
	if !ok {
		return false
	}

	*inputs = append(*inputs, model.Input{
		Type:  "select",
		Name:  nm,
		Value: vl,
	})

	return true

}

func checkInputRadio(n *html.Node, inputs *[]model.Input) bool {
	if n.Type != html.ElementNode {
		return false
	}

	if n.Data != "input" {
		return false
	}

	tp, ok := getAttribute(n, "type")
	if !ok || tp != "radio" {
		return false
	}

	nm, ok := getAttribute(n, "name")
	if !ok {
		return false
	}

	vl, ok := getAttribute(n, "value")
	if !ok {
		return false
	}

	*inputs = append(*inputs, model.Input{
		Type:  "radio",
		Name:  nm,
		Value: vl,
	})

	return true
}

func checkInputText(n *html.Node, inputs *[]model.Input) bool {
	if n.Type != html.ElementNode {
		return false
	}

	if n.Data != "input" {
		return false
	}

	tp, ok := getAttribute(n, "type")
	if !ok || tp != "text" {
		return false
	}

	nm, ok := getAttribute(n, "name")
	if !ok {
		return false
	}

	*inputs = append(*inputs, model.Input{
		Type: "text",
		Name: nm,
	})

	return true
}

func getAttribute(n *html.Node, key string) (string, bool) {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val, true
		}
	}

	return "", false
}
