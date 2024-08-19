package http

import (
	"errors"
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
	CreateUserData *CreateUserData
}

func newIndexData(users []User, createUserData *CreateUserData) *IndexData {
	return &IndexData{Users: users, CreateUserData: createUserData}
}

type CreateUserData struct {
	Name *string
	Email *string
	ErrorMessage *string
}

func newCreateUserData(name, email, errorMessage *string) *CreateUserData {
	return &CreateUserData{Name: name, Email: email, ErrorMessage: errorMessage}
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
		return con.Render(http.StatusOK, "index", newIndexData(users, nil))
	})

	server.POST("/user", func(con echo.Context) error {
		name := con.FormValue("name")
		email := con.FormValue("email")

		_, err := findUserIdWithEmail(email, users)
		if (err == nil) {
			errorMsg := "User already exists with given email"
			return con.Render(http.StatusUnprocessableEntity, "add-user", newCreateUserData(&name, &email, &errorMsg))
		}

		newUser := User{Name: name, Email: email}
		users = append(users, newUser)
		con.Render(http.StatusCreated, "add-user", newCreateUserData(nil, nil, nil))
		return con.Render(http.StatusCreated, "user", newUser)
	})

	server.DELETE("/user", func(con echo.Context) error {
		emailToDelete := con.FormValue("email")

		userIndex, err := findUserIdWithEmail(emailToDelete, users)
		if (err != nil) {
			return con.HTML(http.StatusNotFound, "<span>No user found with given email</span>")
		}

		users = append(users[:userIndex], users[userIndex + 1:]...)
		return con.HTML(http.StatusOK, "")
	})

	server.Logger.Fatal(server.Start(":8000"))
}

func findUserIdWithEmail(email string, users []User) (int, error) {
	for index, user := range users {
		if user.Email == email {
			return index, nil
		}
	}

	return 0, errors.New("User not found")
}
