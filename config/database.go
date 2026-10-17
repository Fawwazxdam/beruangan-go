package config

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite" // Driver SQLite murni kamu

	// Import library golang-migrate
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var DB *sql.DB

func ConnectDB() {
	var err error
	// Pakai path db kamu, pastikan lokasinya benar
	DB, err = sql.Open("sqlite", "./keuangan.db")
	if err != nil {
		log.Fatal("Gagal buka file SQLite: ", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Database tidak merespon: ", err)
	}

	// Panggil fungsi migrasi otomatis
	runMigrations()

	log.Println("✅ Database siap digunakan!")
}

func runMigrations() {
	// 1. Koneksikan driver golang-migrate dengan *sql.DB milikmu
	driver, err := sqlite.WithInstance(DB, &sqlite.Config{})
	if err != nil {
		log.Fatal("Gagal inisialisasi driver migrasi: ", err)
	}

	// 2. Tentukan letak folder file .sql kamu (menggunakan file://)
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations", // Path ke folder migrasi yang baru kita buat
		"sqlite",               // Nama driver (harus 'sqlite' untuk modernc)
		driver,
	)
	if err != nil {
		log.Fatal("Gagal setup migrasi: ", err)
	}

	// 3. Jalankan migrasi!
	err = m.Up()

	// 4. Cek hasil. Kalau muncul ErrNoChange, berarti gak ada file sql baru, aman.
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("Gagal menjalankan migrasi: ", err)
	}

	log.Println("🚀 Cek Migrasi Selesai: Database up-to-date!")
}