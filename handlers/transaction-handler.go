package handlers

import (
	"encoding/json"
	"fmt" // Buat ngegabungin teks (string)
	"net/http"
	"time" // Buat ngatur tanggal dan bulan

	"project-keuangan/config"
	"project-keuangan/models"
)

// Handler GET (Menampilkan Data)
func GetTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Header.Get("X-User-ID") // Ambil user_id dari header
	if userID == "" {
		http.Error(w, "Kode pengguna tidak ditemukan", http.StatusUnauthorized)
		return
	}

	rows, err := config.DB.Query("SELECT id, title, amount, type, due_date, status, group_id FROM transactions WHERE user_id = ?", userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		// Perhatikan urutan scan harus sama persis dengan urutan kolom di SELECT
		rows.Scan(&t.ID, &t.Title, &t.Amount, &t.Type, &t.DueDate, &t.Status, &t.GroupID)
		transactions = append(transactions, t)
	}
	fmt.Println("DATA TRANSAKSI:", transactions)

	json.NewEncoder(w).Encode(transactions)
}

// Handler POST (Menyimpan Data + Auto Generate Cicilan)
func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Header.Get("X-User-ID") // Ambil user_id dari header
	if userID == "" {
		http.Error(w, "Kode pengguna tidak ditemukan", http.StatusUnauthorized)
		return
	}

	var req models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Format JSON salah", http.StatusBadRequest)
		return
	}

	// 1. Validasi Tenor (Kalau user nggak ngisi tenor, kita anggap 1 kali bayar)
	if req.Tenor <= 0 {
		req.Tenor = 1
	}

	// 2. Bikin Group ID unik pakai Timestamp (Waktu saat ini dalam nano-detik)
	// Hasilnya bakal kayak gini: "TRX-1698765432100"
	groupID := fmt.Sprintf("TRX-%d", time.Now().UnixNano())

	// 3. Ubah teks tanggal (String) jadi objek Waktu di Golang
	// INGAT: "2006-01-02" ini adalah FORMAT BAKU dari Golang, bukan tanggal transaksi kamu!
	startDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		http.Error(w, "Format tanggal harus YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	// 4. Looping sebanyak jumlah Tenor (Cicilan)
	for i := 0; i < req.Tenor; i++ {

		// Fungsi andalan Golang: AddDate(Tahun, Bulan, Hari)
		// Kalau putaran pertama (i=0), bulan ditambah 0. Putaran kedua (i=1), bulan ditambah 1, dst.
		installmentDate := startDate.AddDate(0, i, 0).Format("2006-01-02")

		// Atur nama Title. Kalau tenor > 1, tambahin embel-embel "(1/3)"
		title := req.Title
		if req.Tenor > 1 {
			title = fmt.Sprintf("%s (%d/%d)", req.Title, i+1, req.Tenor)
		}

		// Masukin ke DB!
		query := "INSERT INTO transactions (user_id, title, amount, type, due_date, group_id) VALUES (?, ?, ?, ?, ?, ?)"
		_, err := config.DB.Exec(query, userID, title, req.Amount, req.Type, installmentDate, groupID)

		if err != nil {
			http.Error(w, "Gagal simpan ke DB", http.StatusInternalServerError)
			return
		}
	}

	// Balikin pesan sukses
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Berhasil membuat %d transaksi!", req.Tenor),
	})
}

// Handler DELETE (Menghapus Data)
func DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Header.Get("X-User-ID") // Ambil user_id dari header
	if userID == "" {
		http.Error(w, "Kode pengguna tidak ditemukan", http.StatusUnauthorized)
		return
	}

	// 1. Tangkap ID dari URL (Contoh: /api/transactions/5)
	// Fitur ini otomatis ada di Go 1.22 ke atas!
	id := r.PathValue("id")

	// Validasi kalau ID kosong
	if id == "" {
		http.Error(w, "ID tidak boleh kosong", http.StatusBadRequest)
		return
	}

	// 2. Eksekusi query hapus ke database
	query := "DELETE FROM transactions WHERE id = ? AND user_id = ?"
	result, err := config.DB.Exec(query, id, userID)

	if err != nil {
		http.Error(w, "Gagal menghapus data di DB", http.StatusInternalServerError)
		return
	}

	// 3. Cek apakah datanya beneran ada yang kehapus?
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Data tidak ditemukan", http.StatusNotFound)
		return
	}

	// 4. Balikin respon sukses
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Transaksi berhasil dihapus!",
	})
}

// ----------------------------------------------------
// FUNGSI UPDATE DATA (Edit Judul, Nominal, atau Tandai Lunas)
// ----------------------------------------------------
func UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// TAMBAHKAN VALIDASI USER ID
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

	var req models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Format JSON salah", http.StatusBadRequest)
		return
	}

	if req.Status == "" {
		req.Status = "PENDING"
	}

	// UBAH QUERY UNTUK MENGECEK user_id
	query := `
        UPDATE transactions 
        SET title = ?, amount = ?, type = ?, due_date = ?, status = ?
        WHERE id = ? AND user_id = ?
    `
	// TAMBAHKAN userID PADA PARAMETER EXEC
	result, err := config.DB.Exec(query, req.Title, req.Amount, req.Type, req.DueDate, req.Status, id, userID)
	if err != nil {
		http.Error(w, "Gagal update data di DB", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Data tidak ditemukan atau Anda tidak berhak mengubahnya", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Data berhasil diupdate!"})
}

// ----------------------------------------------------
// FUNGSI DELETE ALL BY GROUP ID (Hapus 1 set cicilan)
// ----------------------------------------------------
func DeleteByGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// TAMBAHKAN VALIDASI USER ID
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Kode pengguna tidak ditemukan", http.StatusUnauthorized)
		return
	}

	groupID := r.PathValue("group_id")
	if groupID == "" {
		http.Error(w, "Group ID tidak boleh kosong", http.StatusBadRequest)
		return
	}

	// UBAH QUERY UNTUK MENGECEK user_id
	query := "DELETE FROM transactions WHERE group_id = ? AND user_id = ?"
	// TAMBAHKAN userID PADA PARAMETER EXEC
	result, err := config.DB.Exec(query, groupID, userID)
	if err != nil {
		http.Error(w, "Gagal menghapus data grup", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Tidak ada data dengan Group ID tersebut atau Anda tidak berhak menghapusnya", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Berhasil menghapus %d cicilan sekaligus!", rowsAffected),
	})
}
