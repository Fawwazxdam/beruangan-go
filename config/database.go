package config

import (
	"database/sql"
	"io/fs"
	"log"

	_ "modernc.org/sqlite" // Driver SQLite murni kamu

	// Import library golang-migrate
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

var DB *sql.DB

func ConnectDB(migrations fs.FS) {
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
	runMigrations(migrations)

	log.Println("✅ Database siap digunakan!")
}

func runMigrations(migrations fs.FS) {
	// 1. Koneksikan driver golang-migrate dengan *sql.DB milikmu
	driver, err := sqlite.WithInstance(DB, &sqlite.Config{})
	if err != nil {
		log.Fatal("Gagal inisialisasi driver migrasi: ", err)
	}

	// 2. Baca file .sql dari dalam binary (sudah di-embed lewat //go:embed).
	//    Dibanding "file://db/migrations", cara ini tidak tergantung
	//    current working directory, jadi aman dijalankan dari mana pun (pm2, dll).
	sourceDriver, err := iofs.New(migrations, "db/migrations")
	if err != nil {
		log.Fatal("Gagal memuat migrasi dari binary: ", err)
	}

	// 3. Gabungkan source migrasi + driver database yang sudah ada
	m, err := migrate.NewWithInstance("iofs", sourceDriver, "sqlite", driver)
	if err != nil {
		log.Fatal("Gagal setup migrasi: ", err)
	}

	// 4. Jalankan migrasi!
	err = m.Up()

	// 5. Cek hasil. Kalau muncul ErrNoChange, berarti gak ada file sql baru, aman.
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("Gagal menjalankan migrasi: ", err)
	}

	log.Println("🚀 Cek Migrasi Selesai: Database up-to-date!")
}