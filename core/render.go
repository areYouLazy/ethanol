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
	tmpl, _ := template.New("card").ParseFiles("ui/templates/card.html")

	// generate results based on card template
	for _, v := range results {
		// convert v to a map
		tmpl.ExecuteTemplate(&w, "card", v)
	}

	// return rendered template
	return w.String()
}
