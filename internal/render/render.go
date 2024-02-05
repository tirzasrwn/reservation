package render

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/justinas/nosurf"
	"github.com/tirzasrwn/reservation/internal/config"
	"github.com/tirzasrwn/reservation/internal/models"
)

var (
	pathToTemplates string = "./templates"

	functions = template.FuncMap{
		"humanDate":               HumanDate,
		"formatDate":              FormatDate,
		"iterate":                 Iterate,
		"add":                     Add,
		"stringNumberIterate":     StringNumberIterate,
		"stringToInt":             StringToInt,
		"subtract":                Subtract,
		"changeUnderscoreToSlash": ChangeUnderscoreToSlash,
	}

	app *config.AppConfig
)

// NewRenderer sets the config for the template package.
func NewRenderer(a *config.AppConfig) {
	app = a
}

// HumanDate returns time in YYYY-mm-dd.
func HumanDate(t time.Time) string {
	return t.Format(time.DateOnly)
}

// FormatDate returns formatted date.
func FormatDate(t time.Time, f string) string {
	return t.Format(f)
}

// Iterate returns a slice of ints, starting at 1, going to count.
func Iterate(count int) []int {
	var i int
	var items []int
	for i = 0; i < count; i++ {
		items = append(items, i)
	}
	return items
}

// input number in string will return slice of int.
func StringNumberIterate(s string) []int {
	count, _ := strconv.Atoi(s)
	var items []int
	for i := 0; i < count; i++ {
		items = append(items, i+1)
	}
	return items
}

// Add adds beetwen two number.
func Add(a, b int) int {
	return a + b
}

// Subtract beetwen two number.
func Subtract(a, b int) int {
	if a > b {
		return a - b
	}
	return b - a
}

// Convert string to int.
func StringToInt(s string) int {
	number, _ := strconv.Atoi(s)
	return number
}

// Replace underscore to slash.
func ChangeUnderscoreToSlash(s string) string {
	return strings.Replace(s, "_", "/", 1)
}

// AddDefaultData adds data for all templates.
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	if td.StringMap == nil {
		td.StringMap = make(map[string]string)
	}
	td.StringMap["quote"] = "All we have is today. Just live it. We don't know about tomorrow."
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticated = 1
	}
	return td
}

// Template renders templates using html/template.
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
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

	td = AddDefaultData(td, r)

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
