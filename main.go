package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"project-keuangan/config"
	"project-keuangan/handlers"
	"strings"
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

	// 4. SETUP ROUTE FRONTEND (JURUS SPA)
	dist, err := fs.Sub(frontendFiles, "fe-vue/dist")
	if err != nil {
		log.Fatal("Gagal load folder frontend:", err)
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Kalau user salah ketik rute /api/, kasih error 404 murni (jangan kasih HTML)
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		// Rapikan path URL
		path := r.URL.Path
		if path == "/" {
			path = "index.html"
		}
		path = strings.TrimPrefix(path, "/")

		// Cek apakah file fisik (kayak .css, .js, logo.png) benar-benar ada di folder dist
		_, err := fs.Stat(dist, path)
		if err != nil {
			// JURUS SPA: Kalau nggak ada file-nya (misal user akses /dashboard),
			// paksakan tampilkan index.html biar Vue Router yang ambil alih!
			html, _ := fs.ReadFile(dist, "index.html")
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(html)
			return
		}

		// Kalau filenya beneran ada, jalankan FileServer bawaan
		http.FileServer(http.FS(dist)).ServeHTTP(w, r)
	})

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
