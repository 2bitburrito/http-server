package main

import (
	"context"
	"database/sql"
	"fmt"
	"http-server/internal/auth"
	"http-server/internal/database"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
	tokenSecret    string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func checkHealth(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) showMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	html := fmt.Sprintf(`
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
	`, cfg.fileserverHits.Load())
	w.Write([]byte(html))
}

func (cfg *apiConfig) resetMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	cfg.fileserverHits.Store(0)
	cfg.dbQueries.DeleteAllUsers(req.Context())
	w.WriteHeader(http.StatusOK)
}

func (cfg *apiConfig) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			returnJsonError(w, "error while getting token: "+err.Error(), 401)
			return
		}
		userID, err := auth.ValidateJWT(token, cfg.tokenSecret)
		if err != nil {
			returnJsonError(w, "error while authorizing token: "+err.Error(), 401)
			return
		}
		ctx := context.WithValue(r.Context(), "UserID", userID.String())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("error loading .env")
	}
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Couldn't open connection to database")
	}

	fmt.Println("Starting...")
	const port = "8080"

	mux := http.NewServeMux()
	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}
	cfg := &apiConfig{
		dbQueries:   database.New(db),
		platform:    os.Getenv("PLATFORM"),
		tokenSecret: os.Getenv("TOKEN_SECRET"),
	}

	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /admin/metrics", cfg.showMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.resetMetrics)

	mux.HandleFunc("GET /api/healthz", checkHealth)

	mux.HandleFunc("POST /api/users", cfg.addUser)
	mux.HandleFunc("PUT /api/users", cfg.updateUser)

	mux.Handle("POST /api/chirps", cfg.authMiddleware(http.HandlerFunc(cfg.postChirp)))
	mux.Handle("GET /api/chirps", cfg.authMiddleware(http.HandlerFunc(cfg.getAllChirps)))
	mux.HandleFunc("GET /api/chirps/{id}", http.HandlerFunc(cfg.getSingleChirp))
	mux.Handle("DELETE /api/chirps/{id}", cfg.authMiddleware(http.HandlerFunc(cfg.deleteChirp)))

	mux.HandleFunc("POST /api/login", cfg.login)
	mux.HandleFunc("POST /api/refresh", cfg.refreshToken)
	mux.HandleFunc("POST /api/revoke", cfg.revokeToken)

	server.ListenAndServe()
}
