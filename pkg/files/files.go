package files

import (
	"io/ioutil"
	"regexp"
	"strings"
)

// Root is the relative path to the files
var Root string

// GetTemplate for page (this is a recursive function to allow for nesting of templates)
func GetTemplate(domain string, templateName string) string {
	data, err := ioutil.ReadFile(Root + domain + "/templates/" + templateName + ".html")
	if err != nil {
		panic(err)
	}
	template := string(data)
	// Look for sub-template tokens
	re, _ := regexp.Compile(`#T_(\w+)#`) // An example token targeting the footer.html template is: #T_FOOTER#
	m := re.FindAllStringSubmatch(template, -1)
	if m != nil {
		for _, sub := range m {
			template = strings.ReplaceAll(template, sub[0], GetTemplate(domain, strings.ToLower(sub[1])))
		}
	}
	return template
}
