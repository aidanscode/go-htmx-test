package http

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Renderer struct {
	templates *template.Template
}

func (r *Renderer) Render(w io.Writer, name string, data interface{}, con echo.Context) error {
	return r.templates.ExecuteTemplate(w, name, data)
}

type User struct {
	Name string
	Email string
}

type IndexData struct {
	Users []User
}

func NewIndexData(users []User) *IndexData {
	return &IndexData{Users: users}
}

func Start() {
	server := echo.New()
	renderer := &Renderer{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
	server.Renderer = renderer

	users := make([]User, 0)
	users = append(users, User{Name: "Aidan", Email: "e@mail.com"})

	server.Static("/static", "static")

	server.GET("/", func(con echo.Context) error {
		return con.Render(http.StatusOK, "index", NewIndexData(users))
	})

	server.POST("/user", func(con echo.Context) error {
		//fmt.Println(con.FormParams())
		newUser := User{Name: con.FormValue("name"), Email: con.FormValue("email")}
		users = append(users, newUser)
		return con.Render(http.StatusCreated, "user", newUser)
	})

	server.Logger.Fatal(server.Start(":8000"))
}
