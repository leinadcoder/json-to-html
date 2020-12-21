package jsontohtml

import (
	"bytes"
	"reflect"
)

//Parse JSON to start a HTML build
func GetHTML(contents string) string {
	htmlOut := ""
	response, err := toStruct(contents)

	if err != nil {
		htmlOut = "Invalid Json format"
	} else {
		htmlOut = buildComponents(response)
	}

	return htmlOut
}

func buildComponents(components map[string]interface{}) string {
	html := ""

	if reflect.TypeOf(components["element"]).Kind() != reflect.String {
		return "Not valid HTML element!"
	}

	element := components["element"].(string)

	// read attributes
	if reflect.TypeOf(components["attribs"]).Kind() != reflect.Map &&
		reflect.TypeOf(components["attribs"]).Kind() != reflect.String {
		return "Not valid attributes!"
	}

	attribs := make(map[string]string)
	if reflect.TypeOf(components["attribs"]).Kind() == reflect.Map {
		attribs = toMap(components["attribs"])
	}

	// get contents
	contents := ""
	switch reflect.TypeOf(components["contents"]).Kind() {
	case reflect.Map:
		contents = buildComponents(components["contents"].(map[string]interface{}))
	case reflect.Slice:
		if (element == "li") {
			contents = li(attribs, components)
			return contents
		}

		s := reflect.ValueOf(components["contents"])
		for i := 0; i < s.Len(); i++ {
			if reflect.TypeOf(s.Index(i).Interface()).Kind() == reflect.Map {
				contents = contents + buildComponents(s.Index(i).Interface().(map[string]interface{}))
			}

			if reflect.TypeOf(s.Index(i).Interface()).Kind() == reflect.String {
				contents = contents + buildElement("p", make(map[string]string), s.Index(i).Interface().(string))
			}
		}

	case reflect.String :
		contents = components["contents"].(string)
	}

	html = buildElement(element, attribs, contents)

	return html
}

func buildElement(element string, attributes map[string]string, content string) string {
	var htmlOut bytes.Buffer // html writer

	htmlOut.WriteString("<")
	htmlOut.WriteString(element)

	for k, v := range attributes {
		htmlOut.WriteString(" ")
		htmlOut.WriteString(k)
		htmlOut.WriteString("=\"")
		htmlOut.WriteString(v)
		htmlOut.WriteString("\"")
	}

	htmlOut.WriteString(">")

	if element != "img" {
		htmlOut.WriteString(content)

		htmlOut.WriteString("</")
		htmlOut.WriteString(element)
		htmlOut.WriteString(">")
	}

	return htmlOut.String()
}

func li(attributes map[string]string, elements map[string]interface{}) string {
	var htmlElement, htmlOut bytes.Buffer // html writer

	htmlElement.WriteString("<li")

	for k, v := range attributes {
		htmlElement.WriteString(" ")
		htmlElement.WriteString(k)
		htmlElement.WriteString("=\"")
		htmlElement.WriteString(v)
		htmlElement.WriteString("\"")
	}

	htmlElement.WriteString(">")

	s := reflect.ValueOf(elements["contents"])
	for i := 0; i < s.Len(); i++ {
		htmlOut.WriteString(htmlElement.String())

		if reflect.TypeOf(s.Index(i).Interface()).Kind() == reflect.Map {
			htmlOut.WriteString(buildComponents(s.Index(i).Interface().(map[string]interface{})))
		}

		if reflect.TypeOf(s.Index(i).Interface()).Kind() == reflect.String {
			htmlOut.WriteString(s.Index(i).Interface().(string))
		}

		htmlOut.WriteString("</li>")
	}

	return htmlOut.String()
}
