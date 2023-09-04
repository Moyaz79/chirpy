package main

import (
	"log"
	"net/http"

	"github.com/Moyaz79/chirpy/internal/database"
	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	fileserverHits int
	DB 	*database.DB
}

func main() {
	const filepathRoot = "."
	const  port = "8080"

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiConfig {
		fileserverHits: 0,
		DB: db,
	}

	
	r := chi.NewRouter()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	r.Handle("/app", fsHandler)
	r.Handle("/app/*", fsHandler)
	
	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", apiCfg.metricsHandler)
	r.Mount("/admin", adminRouter)

	apiRouters := chi.NewRouter()
	apiRouters.Get("/healthz", readinessHandler)
	apiRouters.Get("/chirps", apiCfg.chirpsHandlerRetrieve)
	apiRouters.Post("/chirps", apiCfg.chirpsHandlerCreate)
	r.Mount("/api", apiRouters)

	corsMux := middlewareCors(r)

	server := &http.Server {
		Addr: ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n",filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}


