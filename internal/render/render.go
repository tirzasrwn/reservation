package render

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/tirzasrwn/reservation/internal/config"
	"github.com/tirzasrwn/reservation/internal/models"
)

var (
	pathToTemplates string = "./templates"
	functions              = template.FuncMap{}
	app             *config.AppConfig
)

// NewTemplates sets the config for the template package.
func NewTemplates(a *config.AppConfig) {
	app = a
}

// AddDefaultData adds data for all templates.
func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	td.StringMap["quote"] = "All we have is today. Just live it. We don't know about tomorrow."
	return td
}

// RenderTemplate renders templates using html/template.
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) error {
	var tc map[string]*template.Template

	if app.UseCache {
		// Get the template cache from the app config.
		tc = app.TemplateCache
	} else {
		// This is just used for testing, so that we rebuild
		// the cache on every request.
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		return errors.New("could not get template from cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td)

	err := t.Execute(buf, td)
	if err != nil {
		return err
	}
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
		fmt.Println("Error writing template to browser.", err)
		return err
	}
	return nil
}

// CreateTemplateCache creates a template cache as a map.
func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.html", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
