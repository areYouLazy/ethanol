package core

import (
	"bytes"
	"text/template"

	"github.com/areYouLazy/ethanol/types"
)

func renderHTMLResultsFromJSON(results []types.SearchResult) string {
	// instantiate a new bytes buffer
	var w bytes.Buffer

	// load template from card.html
	var tmpl *template.Template

	for _, result := range results {
		raw_label := result["raw_label"].(string)
		switch raw_label {
		case "check_mk":
			tmpl, _ = template.New("checkmk").ParseFiles("ui/templates/checkmk.html")
			tmpl.ExecuteTemplate(&w, "checkmk", result)
		case "jira":
			tmpl, _ = template.New("jira").ParseFiles("ui/templates/jira.html")
			tmpl.ExecuteTemplate(&w, "jira", result)
		case "syspass":
			tmpl, _ = template.New("syspass").ParseFiles("ui/templates/syspass.html")
			tmpl.ExecuteTemplate(&w, "syspass", result)
		case "otrs":
			tmpl, _ = template.New("otrs").ParseFiles("ui/templates/otrs.html")
			tmpl.ExecuteTemplate(&w, "otrs", result)
		default:
			tmpl, _ = template.New("card").ParseFiles("ui/templates/card.html")
			tmpl.ExecuteTemplate(&w, "card", result)
		}
	}

	// return rendered template
	return w.String()
}
