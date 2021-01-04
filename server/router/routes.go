package router

import (
	"github.com/UN0wen/virtualfs/server/api/controllers"
	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func createRoutes(r *chi.Mux) {
	r.Route("/api", func(r chi.Router) {
		r.Get("/{path}", controllers.GetPath)            // Get /
		r.Get("/{path}/exact", controllers.GetExactPath) // Get /
		r.Post("/", controllers.CreatePath)
		// r.Put("/", controllers.UpdateItem) // This doesn't work properly yet, but it's not in the spec
		r.Put("/", controllers.UpdatePath)
		r.Delete("/{path}", controllers.DeletePath)
	})
}

// NewRouter creates a chi Router with all routes and middleware configured
func NewRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.Recoverer)
	router.Use(chimiddleware.URLFormat)
	router.Use(render.SetContentType(render.ContentTypeJSON))
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	createRoutes(router)

	spa := spaHandler{staticPath: "build", indexPath: "index.html"}
	router.Handle("/*", spa)
	return router
}
