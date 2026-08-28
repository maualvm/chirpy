package main

import (
	"fmt"
	"net/http"
)

func (a *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	totalHits := a.fileserverHits.Load()
	content := fmt.Sprintf(`
<html>
	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>
</html>
	`, totalHits)
	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte(content))
}

func (a *apiConfig) handlerMetricsReset(w http.ResponseWriter, r *http.Request) {
	if a.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset only allowed in dev environments"))
		return
	}
	a.fileserverHits.Store(0)
	err := a.db.DeleteAllUsers(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error resetting users table: " + err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Metrics have been reset and users table truncated."))
}

func (a *apiConfig) middlewareMetricsIncrease(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
