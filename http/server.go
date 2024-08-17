package http

import (
	"fmt"
	"io"
	"net/http"
	"text/template"

	"github.com/labstack/echo/v4"
)

type Renderer struct {
	templates *template.Template
}

func (r *Renderer) Render(w io.Writer, name string, data interface{}, con echo.Context) error {
	return r.templates.ExecuteTemplate(w, name, data)
}

func Start() {
	server := echo.New()
	renderer := &Renderer{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
	server.Renderer = renderer

	server.GET("/", func(con echo.Context) error {
		return con.String(http.StatusOK, "Hello, Go World!")
	})

	server.GET("/users/:id", func(con echo.Context) error {
		id := con.Param("id")
		return con.String(http.StatusOK, fmt.Sprintf("Viewing profile for user: %s", id))
	})

	server.POST("/users", func(con echo.Context) error {
		email := con.FormValue("email")
		password := con.FormValue("password")
		return con.String(http.StatusOK, fmt.Sprintf("Creating user with email=\"%v\", password=\"%v\"", email, password))
	})

	server.GET("/secret", func(con echo.Context) error {
		return con.Render(http.StatusOK, "secret", "password123")
	})

	server.Logger.Fatal(server.Start(":8000"))
}
