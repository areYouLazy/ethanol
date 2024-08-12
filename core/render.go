package core

import (
	"bytes"
	"text/template"

	"github.com/areYouLazy/ethanol/types"
)

func renderHTMLResultsFromJSON(results []types.SearchResult) string {
	// instantiate a new bytes buffer
	var w bytes.Buffer

	// load templates from files
	tmpl, _ := template.New("card").ParseFiles("ui/templates/*.html")
	// checkmk, _ := template.New("card").ParseFiles("ui/templates/checkmk.html")
	// jira, _ := template.New("jira").ParseFiles("ui/templates/jira.html")
	// syspass, _ := template.New("syspass").ParseFiles("ui/templates/syspass.html")
	// otrs, _ := template.New("otrs").ParseFiles("ui/templates/otrs.html")
	// prtg, _ := template.New("prtg").ParseFiles("ui/templates/prtg.html")
	//

	for _, result := range results {
		raw_label := result["raw_label"].(string)
		tmpl.ExecuteTemplate(&w, raw_label, result)
		// switch raw_label {
		// case "check_mk":
		// 	tmpl.ExecuteTemplate(&w, "checkmk", result)
		// case "jira":
		// 	tmpl.ExecuteTemplate(&w, "jira", result)
		// case "syspass":
		// 	tmpl.ExecuteTemplate(&w, "syspass", result)
		// case "otrs":
		// 	tmpl.ExecuteTemplate(&w, "otrs", result)
		// case "prtg":
		// 	tmpl.ExecuteTemplate(&w, "prtg", result)
		// default:
		// 	tmpl.ExecuteTemplate(&w, "card", result)
		// }
	}

	// return rendered template
	return w.String()
}
