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
	// Tambah kolom group_id TEXT
	// Di dalam fungsi buatTabel()
	query := `
	CREATE TABLE IF NOT EXISTS transactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		amount INTEGER NOT NULL,
		type TEXT NOT NULL,
		due_date TEXT,
		group_id TEXT,
		status TEXT DEFAULT 'PENDING' -- TAMBAHIN BARIS INI
	);
	`
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Gagal bikin tabel: ", err)
	}
}