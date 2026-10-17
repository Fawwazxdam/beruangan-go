package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"project-keuangan/config" // Sesuaikan dengan nama module go.mod kamu
	"project-keuangan/models"
)

// GET: Ambil daftar dompet
func GetWallets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Ambil user_id dari header (sama kayak endpoint transaksi & summary)
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Kode pengguna tidak ditemukan", http.StatusUnauthorized)
		return
	}

	rows, err := config.DB.Query("SELECT id, user_id, name, balance FROM wallets WHERE user_id = ?", userID)
	if err != nil {
		http.Error(w, "Gagal mengambil data dompet", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var wallets []models.Wallet
	for rows.Next() {
		var wallet models.Wallet
		if err := rows.Scan(&wallet.ID, &wallet.UserID, &wallet.Name, &wallet.Balance); err != nil {
			continue
		}
		wallets = append(wallets, wallet)
	}

	// Kalau kosong, kembalikan array kosong, bukan null
	if wallets == nil {
		wallets = []models.Wallet{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wallets)
}

// POST: Bikin dompet baru
func CreateWallet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Ambil user_id dari header (sama kayak endpoint transaksi & summary)
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Kode pengguna tidak ditemukan", http.StatusUnauthorized)
		return
	}

	var payload struct {
		Name    string `json:"name"`
		Balance int    `json:"balance"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Format data salah", http.StatusBadRequest)
		return
	}

	// Bikin ID unik ala kadarnya pakai timestamp biar gampang (contoh: W-1697543...)
	walletID := fmt.Sprintf("W-%d", time.Now().Unix())

	query := "INSERT INTO wallets (id, user_id, name, balance) VALUES (?, ?, ?, ?)"
	_, err := config.DB.Exec(query, walletID, userID, payload.Name, payload.Balance)
	if err != nil {
		http.Error(w, "Gagal menyimpan dompet", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Dompet berhasil dibuat!",
		"id":      walletID,
	})
}

// PUT: Ubah nama / saldo dompet
func UpdateWallet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Ambil user_id dari header (sama kayak endpoint lainnya)
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Kode pengguna tidak ditemukan", http.StatusUnauthorized)
		return
	}

	// Tangkap ID dari URL (Contoh: /api/wallets/W-1697543...)
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "ID tidak boleh kosong", http.StatusBadRequest)
		return
	}

	var payload struct {
		Name    string `json:"name"`
		Balance int    `json:"balance"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Format data salah", http.StatusBadRequest)
		return
	}

	// Validasi nama tidak boleh kosong
	if payload.Name == "" {
		http.Error(w, "Nama dompet tidak boleh kosong", http.StatusBadRequest)
		return
	}

	query := "UPDATE wallets SET name = ?, balance = ? WHERE id = ? AND user_id = ?"
	result, err := config.DB.Exec(query, payload.Name, payload.Balance, id, userID)
	if err != nil {
		http.Error(w, "Gagal update data di DB", http.StatusInternalServerError)
		return
	}

	// Cek apakah datanya beneran ada & milik user yang sama
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Data tidak ditemukan atau Anda tidak berhak mengubahnya", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Dompet berhasil diupdate!"})
}

// DELETE: Hapus dompet
func DeleteWallet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Ambil user_id dari header (sama kayak endpoint lainnya)
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Kode pengguna tidak ditemukan", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "ID tidak boleh kosong", http.StatusBadRequest)
		return
	}

	query := "DELETE FROM wallets WHERE id = ? AND user_id = ?"
	result, err := config.DB.Exec(query, id, userID)
	if err != nil {
		http.Error(w, "Gagal menghapus data di DB", http.StatusInternalServerError)
		return
	}

	// Cek apakah datanya beneran ada yang kehapus
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Data tidak ditemukan atau Anda tidak berhak menghapusnya", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Dompet berhasil dihapus!"})
}