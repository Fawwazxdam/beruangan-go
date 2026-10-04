package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"project-keuangan/config"
	"project-keuangan/handlers"
)
// go:embed fe-vue/dist
var frontendFiles embed.FS

func main() {
	// 1. Panggil koneksi DB
	config.ConnectDB()

	// 2. Siapkan Buku Menu (Mux)
	mux := http.NewServeMux()

	// 3. Daftarin Route (Sebutin HTTP Method-nya)
	mux.HandleFunc("GET /api/transactions", handlers.GetTransactions)
	mux.HandleFunc("POST /api/transactions", handlers.CreateTransaction)
	mux.HandleFunc("DELETE /api/transactions/{id}", handlers.DeleteTransaction)
	mux.HandleFunc("PUT /api/transactions/{id}", handlers.UpdateTransaction)
	mux.HandleFunc("DELETE /api/transactions/group/{group_id}", handlers.DeleteByGroup)

	mux.HandleFunc("GET /api/summary", handlers.GetDashboardSummary)

	dist, err := fs.Sub(frontendFiles, "fe-vue/dist")
	if err != nil {
		log.Fatal("Gagal load folder frontend:", err)
	}

	// Kalau user akses URL selain /api, kasih file Vue
	mux.Handle("/", http.FileServer(http.FS(dist)))

	// 4. Nyalakan Server
	log.Println("🚀 Server Backend nyala di http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", enableCORS(mux)))
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Kasih izin ke semua domain (tanda bintang) buat akses API ini
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Kalau browser ngecek jalur dulu (Preflight request pakai OPTIONS), langsung kasih OK
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Kalau aman, lanjutin ke handler aslinya
		next.ServeHTTP(w, r)
	})
}
