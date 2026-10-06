package config

import (
	"database/sql"
	"log"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

func ConnectDB() {
	var err error
	DB, err = sql.Open("sqlite", "./keuangan.db")
	if err != nil {
		log.Fatal("Gagal buka file SQLite: ", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Database tidak merespon: ", err)
	}

	buatTabel()
	log.Println("✅ Database siap digunakan!")
}

func buatTabel() {
	// 1. Bikin tabel (Berlaku kalau file keuangan.db benar-benar baru)
	query := `
    CREATE TABLE IF NOT EXISTS transactions (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id TEXT NOT NULL,
        title TEXT NOT NULL,
        amount INTEGER NOT NULL,
        type TEXT NOT NULL,
        due_date TEXT,
        group_id TEXT,
        status TEXT DEFAULT 'PENDING'
    );
    `
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Gagal bikin tabel: ", err)
	}

	// 2. JURUS MIGRASI AMAN (Untuk data lama yang belum punya user_id)
	// Kita paksa tambah kolom user_id. Kalau error (karena kolomnya udah ada), cuekin aja.
	// Wajib pakai DEFAULT supaya data lama otomatis terisi teks ini.
	_, errAlter := DB.Exec("ALTER TABLE transactions ADD COLUMN user_id TEXT NOT NULL DEFAULT 'dam'")
	if errAlter == nil {
		log.Println("🔧 Kolom user_id berhasil ditambahkan ke tabel lama!")
	}

	// 3. Pastikan semua data lama yang mungkin kosong di-set ke kode rahasiamu
	_, errUpdate := DB.Exec("UPDATE transactions SET user_id = 'dam' WHERE user_id IS NULL OR user_id = ''")
	if errUpdate == nil {
		log.Println("♻️ Data lama berhasil disuntik user_id!")
	}
}
