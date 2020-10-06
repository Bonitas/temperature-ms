package app

import (
	"fmt"
	"it.schwarz/city/app/handler"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"it.schwarz/city/config"
)

// App has the mongo database and router instances
type App struct {
	Router *mux.Router
	Config *config.Config
}

// ConfigAndRunApp will create and initialize App structure. App factory function.
func ConfigAndRunApp(config *config.Config) {
	app := new(App)
	app.Initialize(config)
	app.Run(config.AppPort)
}

// Initialize initialize the app with
func (app *App) Initialize(config *config.Config) {
	app.Config = config

	app.Router = mux.NewRouter()
	app.setRouters()
}

// setRouters will register routes in router
func (app *App) setRouters() {

	var api = app.Router.PathPrefix("/api").Subrouter()

	api.HandleFunc("/cities", app.handleRequest(handler.GetCities)).Methods(http.MethodGet)
	api.HandleFunc("/add-city", app.handleRequest(handler.PostNewCity)).Methods(http.MethodPost)
}

// Run will start the http server on host that you pass in. host:<ip:port>
func (app *App) Run(port string) {
	// use signals for shutdown server gracefully.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt, os.Kill)
	headersOk := handlers.AllowedHeaders([]string{"content-type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	go func() {
		log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), handlers.CORS(originsOk, headersOk, methodsOk)(app.Router)))
	}()
	log.Printf("Server is listning on http://0.0.0.0:%s\n", port)
	sig := <-sigs
	log.Println("Signal: ", sig)
}

// UseMiddleware will add global middleware in router
func (app *App) UseMiddleware(middleware mux.MiddlewareFunc) {
	app.Router.Use(middleware)
}

// RequestHandlerFunction is a custom type that help us to pass db arg to all endpoints
type RequestHandlerFunction func(config *config.Config, w http.ResponseWriter, r *http.Request)

// handleRequest is a middleware we create for pass in db connection to endpoints.
func (app *App) handleRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(app.Config, w, r)
	}
}
