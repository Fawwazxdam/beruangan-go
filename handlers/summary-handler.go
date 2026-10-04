package handlers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	"project-keuangan/config"
	"project-keuangan/models"
)

func GetDashboardSummary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 1. Ambil format bulan ini.
	// Karena sekarang bulan Oktober 2026, fungsi ini otomatis ngasilin "2026-10"
	bulanIni := time.Now().Format("2006-01")
	filterBulanIni := bulanIni + "-%" // Jadinya "2026-10-%" (buat nyari di SQLite)

	var summary models.DashboardSummary

	// 2. JURUS SENIOR: 1 Query untuk 4 Metrik Sekaligus!
	// Kita manfaatkan CASE WHEN di SQL.
	query := `
		SELECT 
			-- Total semua hutang
			COALESCE(SUM(amount), 0) as total_hutang,
			
			-- Total hutang lunas
			COALESCE(SUM(CASE WHEN status = 'PAID' THEN amount ELSE 0 END), 0) as total_terbayar,
			
			-- Total tagihan khusus bulan ini (berdasarkan due_date)
			COALESCE(SUM(CASE WHEN due_date LIKE ? THEN amount ELSE 0 END), 0) as total_tagihan_bulan_ini,
			
			-- Total tagihan bulan ini yang sudah lunas
			COALESCE(SUM(CASE WHEN due_date LIKE ? AND status = 'PAID' THEN amount ELSE 0 END), 0) as tagihan_terbayar_bulan_ini
			
		FROM transactions
		WHERE type = 'HUTANG' -- Opsional: Pastikan cuma ngitung yang tipenya hutang
	`

	// Eksekusi query dengan mengirim filterBulanIni 2 kali (karena ada 2 tanda tanya di query)
	err := config.DB.QueryRow(query, filterBulanIni, filterBulanIni).Scan(
		&summary.TotalHutang,
		&summary.TotalHutangTerbayar,
		&summary.TotalTagihanBulanIni,
		&summary.TagihanTerbayarBulanIni,
	)

	if err != nil {
		http.Error(w, "Gagal menghitung summary", http.StatusInternalServerError)
		return
	}

	// 3. Hitung Persentase di Golang (Bukan di SQL)
	// Kenapa di Golang? Biar kita aman dari error "Divide by Zero" (dibagi nol).
	if summary.TotalHutang > 0 {
		persen := (float64(summary.TotalHutangTerbayar) / float64(summary.TotalHutang)) * 100

		// Bulatkan jadi 2 angka di belakang koma (contoh: 33.33)
		summary.PersentaseTerbayar = math.Round(persen*100) / 100
	} else {
		summary.PersentaseTerbayar = 0 // Kalau belum ada hutang sama sekali, persennya 0
	}

	// ---------------------------------------------------------
	// FITUR TAMBAHAN: REMINDER H-3 (Dalam 3 Hari ke Depan)
	// ---------------------------------------------------------

	// 1. Tentukan rentang waktu pencarian
	// Hari ini (misal: 2026-10-04)
	hariIni := time.Now().Format("2006-01-02")
	// Batas H-3 (misal: 2026-10-07) -> AddDate(Tahun, Bulan, Hari)
	batasH3 := time.Now().AddDate(0, 0, 3).Format("2006-01-02")

	// 2. Query untuk ngambil DAFTAR transaksinya
	// Pakai BETWEEN biar dapet semua tanggal dari hari ini sampai 3 hari ke depan
	queryReminder := `
		SELECT id, title, amount, type, due_date, group_id, status 
		FROM transactions 
		WHERE status = 'PENDING' 
		AND type = 'HUTANG'
		AND due_date BETWEEN ? AND ?
		ORDER BY due_date ASC
	`

	rows, err := config.DB.Query(queryReminder, hariIni, batasH3)
	if err != nil {
		// Kalau error log aja, jangan sampai ngerusak dashboard keseluruhan
		fmt.Println("Gagal ambil data H-3:", err)
	} else {
		defer rows.Close()

		// Inisiasi array kosong biar di JSON jadinya [] (bukan null)
		summary.ReminderH3 = []models.Transaction{}
		summary.TotalTagihanH3 = 0

		// 3. Looping datanya
		for rows.Next() {
			var t models.Transaction
			rows.Scan(&t.ID, &t.Title, &t.Amount, &t.Type, &t.DueDate, &t.GroupID, &t.Status)

			// Masukin data ke dalam array list
			summary.ReminderH3 = append(summary.ReminderH3, t)

			// JURUS SENIOR: Daripada bikin query SUM() terpisah ke database,
			// Mending kita tambahin totalnya manual di dalam looping ini (Lebih hemat performa!)
			summary.TotalTagihanH3 += t.Amount
		}
	}

	// Terakhir, pastikan baris ini tetap ada di paling bawah:
	// json.NewEncoder(w).Encode(summary)

	// Kirim hasil ke Frontend
	json.NewEncoder(w).Encode(summary)
}
