package selma

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// ViewRenderer handles rendering views
type ViewRenderer struct {
	templates   map[string]*template.Template
	Config      *Config
	DebugLogger *log.Logger
}

// NewViewRenderer initializes a new ViewRenderer
func NewViewRenderer(config *Config, debugLogger *log.Logger) *ViewRenderer {
	return &ViewRenderer{
		templates:   make(map[string]*template.Template),
		Config:      config,
		DebugLogger: debugLogger,
	}
}

// LoadTemplates loads templates from the specified directory
func (vr *ViewRenderer) LoadTemplates(dir string) error {
	templates, err := template.ParseGlob(dir + "/*.html")
	if err != nil {
		return err
	}
	for _, tmpl := range templates.Templates() {
		vr.templates[tmpl.Name()] = tmpl

		if vr.Config.Mode == "development" {
			vr.DebugLogger.Printf("Loaded template: %s", tmpl.Name())
		}
	}
	return nil
}

// Render renders the given template with data
func (vr *ViewRenderer) Render(w http.ResponseWriter, r *http.Request, name string, data interface{}) error {
	tmpl, ok := vr.templates[name]
	if !ok {
		log.Printf("Template %s not found", name)
		return fmt.Errorf("template not found")
	}
	err := tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template %s: %v", name, err)
	}
	return err
}
