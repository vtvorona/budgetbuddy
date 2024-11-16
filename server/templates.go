package server

import (
	"fmt"
	"html/template"
	"net/http"
)

var templates = make(map[string]*template.Template)

// LoadTemplates загружает все шаблоны
func LoadTemplates() {
	templateNames := []string{"home", "registration", "login", "dashboard", "menu", "success"}
	for index, name := range templateNames {
		t, err := template.ParseFiles(
			"templates/header.html",
			"templates/footer.html",
			"templates/layout.html",
			"templates/menu.html",
			"templates/"+name+".html",
		)
		if err == nil {
			templates[name] = t
			fmt.Println("Loaded template:", index, name)
		} else {
			fmt.Println("Error loading template:", err)
		}
	}
}

// RenderTemplate рендерит заданный шаблон с данными (если переданы)
func RenderTemplate(writer http.ResponseWriter, name string, data ...interface{}) {
	if t, ok := templates[name]; ok {
		var err error
		if len(data) > 0 {
			err = t.ExecuteTemplate(writer, "layout.html", data[0]) // Используем данные, если они переданы
		} else {
			err = t.ExecuteTemplate(writer, "layout.html", nil) // Используем nil, если данные не переданы
		}
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.NotFound(writer, nil)
	}
}
