package main

import (
	"context"
	"fmt"
	pb "fullstack-go-grpc/protos/user"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

var userClient pb.UserServiceClient

func main() {
	// gRPC client connection
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	userClient = pb.NewUserServiceClient(conn)

	// Echo instance
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Static files
	e.Static("/static", "static")

	// Template renderer
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Renderer = renderer

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/users")
	})
	e.GET("/users", listUsersHandler)
	e.GET("/users/:id", getUserHandler)
	e.GET("/users/new", newUserFormHandler)
	e.GET("/users/edit/:id", editUserFormHandler)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

// --- Handlers ---

func listUsersHandler(c echo.Context) error {
	// fmt.Println("==============================I am in listuserHandler======================")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := userClient.ListUsers(ctx, &pb.ListUsersRequest{})
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error listing users: %v", err))
	}
	// a, err := json.Marshal(resp)
	// if err != nil {
	// 	log.Fatalf("ruk jaa bsdk ")
	// 	return err
	// }
	// fmt.Println("resp ", string(a))

	return c.Render(http.StatusOK, "users.html", map[string]interface{}{
		"Users": resp.Users,
	})
}

func getUserHandler(c echo.Context) error {
	id := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := userClient.GetUser(ctx, &pb.UserGetterRequest{UniqueId: id})
	if err != nil {
		return c.String(http.StatusNotFound, fmt.Sprintf("Error getting user: %v", err))
	}

	return c.Render(http.StatusOK, "user.html", map[string]interface{}{
		"User": resp.User,
	})
}

func newUserFormHandler(c echo.Context) error {
	fmt.Println("======================i am here ========>>>>>>>>>>")
	
	return c.Render(http.StatusOK, "addusers.html", map[string]interface{}{
		"Title":  "Create New User",
		"Action": "/v1/users", // Backend REST endpoint
		"Method": "POST",
		"User":   nil,
	})
}

func editUserFormHandler(c echo.Context) error {
	id := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := userClient.GetUser(ctx, &pb.UserGetterRequest{UniqueId: id})
	if err != nil {
		return c.String(http.StatusNotFound, fmt.Sprintf("Error getting user for edit: %v", err))
	}

	return c.Render(http.StatusOK, "form.html", map[string]interface{}{
		"Title":  "Edit User",
		"Action": fmt.Sprintf("/v1/users/%s", id), // Backend REST endpoint
		"Method": "PUT",
		"User":   resp.User,
	})
}
