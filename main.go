package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/schema"
	"github.com/xV0lk/htmx-go/internal/api"
	"github.com/xV0lk/htmx-go/internal/db"
	"github.com/xV0lk/htmx-go/internal/localizer"
	"github.com/xV0lk/htmx-go/internal/middleware"
	"github.com/xV0lk/htmx-go/models"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		fmt.Printf("-------------------------\nerr loading config: %s\n", err)
		log.Fatal(err)
		return
	}

	logger := createLogger(&cfg)

	psqlStore, err := db.NewPsql(&cfg.Postgres)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer psqlStore.Close()
	fmt.Println("Connected to postgres")

	if err := db.Migrate(psqlStore, embedMigrations); err != nil {
		log.Fatal(err)
	}

	emailService, err := models.NewDefaultEmailService()
	if err != nil {
		fmt.Printf("-------------------------\nerr creating mail service: %s\n", err)
		log.Fatal(err)
		return
	}
	fmt.Println("Connected to mail Service")

	var (
		// Stores
		decoder      = schema.NewDecoder()
		authStore    = db.NewPsAuthStore(psqlStore)
		taskStore    = db.NewPsTaskStore(psqlStore)
		sessionStore = db.NewPsSessionStore(psqlStore)
		pwResetStore = db.NewPsPwResetStore(psqlStore)
		userStore    = db.NewUserStore(authStore, sessionStore, pwResetStore)
		clientStore  = db.NewPsclientStore(psqlStore)
		studiostore  = db.NewPsstudiosstore(psqlStore)

		// handlers
		tasksHandler  = api.NewTasksHandler(taskStore, decoder)
		authHandler   = api.NewAuthHandler(userStore, decoder, emailService)
		clientHandler = api.Newclienthandler(clientStore, decoder)
		studiohandler = api.Newstudiohandler(studiostore, decoder)
		// Connection
		r = chi.NewRouter()
	)
	r.Use(chiMiddleware.Logger)
	r.Use(middleware.User(userStore))
	r.Use(middleware.ILog(logger))
	r.Use(localizer.I18n)

	// main group of routes that use csrf and sessions
	r.Group(func(r chi.Router) {
		r.Use(middleware.Csrf(&cfg.CSRF))
		r.Get("/", authHandler.HandleRoot)
		r.Route("/login", func(r chi.Router) {
			r.Get("/", api.MakeHandler(authHandler.HandleLogin))
			r.Post("/", api.MakeHandler(authHandler.HandleAuthenticate))
			r.Delete("/", authHandler.HandleLogout)
		})
		r.Get("/home", api.MakeHandler(authHandler.HandleHome))
		r.Get("/forgot-password", authHandler.HandleForgotPassword)
		r.Post("/forgot-password", authHandler.HandleResetPassword)
		r.Get("/reset-password", authHandler.HandleEmailToken)
		r.Post("/reset-password", authHandler.HandleChangePassword)
		r.Route("/tasks", func(r chi.Router) {
			r.Get("/", tasksHandler.HandleGetTasks)
			r.Post("/", tasksHandler.HandlePostTask)
			r.Put("/{id}/toggle", tasksHandler.HandleToggleTask)
			r.Delete("/{id}", tasksHandler.HandleDeleteTask)
			r.Get("/{id}/edit", tasksHandler.HandleEditTask)
			r.Put("/{id}", tasksHandler.HandlePutTask)
		})

	})

	r.Route("/client", func(r chi.Router) {

		r.Get("/{id}", clientHandler.GetClient)
		r.Get("/", clientHandler.GetClients)
		r.Post("/", clientHandler.CreateNewClient)
		r.Delete("/{id}", clientHandler.DeleteClient)
		r.Put("/{id}", clientHandler.UpdateClients)

	})

	r.Route("/studio", func(r chi.Router) {

		r.Get("/{id}", studiohandler.Fetchstudio)
		r.Get("/", studiohandler.Fetchstudios)
		r.Post("/", studiohandler.Createnewstudio)
		r.Delete("/{id}", studiohandler.Deletestudio)
		r.Put("/{id}", studiohandler.Updatestudio)

	})

	// This is a json route and is not intended to be used by the browser
	// TODO: Add middleware to protect this route so it can only be accessed by super admin
	// TODO: Add in the front so it can be used by the browser to work with csrf
	r.Route("/admin-internal", func(r chi.Router) {
		r.Get("/", authHandler.HandleGetUser)
		r.Post("/", api.MakeHandler(authHandler.HandleNewUser))
	})

	// Serve static files
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "assets"))
	api.FileServer(r, "/static/", filesDir)

	port := flag.String("port", cfg.Server.Port, "port to listen on")
	fmt.Printf("Server running on port %s\n", *port)
	err = http.ListenAndServe(*port, r)
	if err != nil {
		log.Fatal(err)
	}
}
