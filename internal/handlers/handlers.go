package handlers

import (
	"net/http"

	"github.com/tirzasrwn/reservation/internal/config"
	"github.com/tirzasrwn/reservation/internal/models"
	"github.com/tirzasrwn/reservation/internal/render"
)

// Repo the repository used by the handlers.
var Repo *Repository

// Repository is the repository type.
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository.
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers.
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler.
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	stringMap := make(map[string]string)
	stringMap["testHome"] = "Hello, again. We are in home page."
	render.RenderTemplate(w, "home.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

// About is the about page handler.
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["testAbout"] = "Hello, again. We are in about page"
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP
	render.RenderTemplate(w, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Reservation renders the make a reservation page and displays form.
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	render.RenderTemplate(w, "make-reservation.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Generals renders the room page.
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	render.RenderTemplate(w, "generals.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Majors renders the room page.
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	render.RenderTemplate(w, "majors.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Availability renders the search availability page.
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	render.RenderTemplate(w, "search-availability.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Contact renders the search availability page.
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	render.RenderTemplate(w, "contact.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}
