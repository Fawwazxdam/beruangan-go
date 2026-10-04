# Tech Stack

**Backend**

- Go 1.27
- `net/http` (standard library)
- SQLite via `modernc.org/sqlite` (pure Go, tanpa CGO)

**Frontend** (`fe-vue/`)

- Vue 3
- Vite 8
- Vue Router 5
- Pinia 4
- Tailwind CSS 4
- @phosphor-icons/vue
- oxfmt (formatter)

# Prasyarat

- Go >= 1.27
- Node.js ^22.18.0 atau >= 24.12.0
- npm

# Menjalankan

## 1. Backend

```bash
go mod tidy
go run .
```

Server jalan di `http://localhost:8080`. Database SQLite dibuat otomatis saat pertama kali start.

## 2. Frontend

```bash
cd fe-vue
npm install
npm run dev
```

Vite jalan di `http://localhost:5173` (atau sesuai output terminal).

# Perintah Lainnya

| Perintah              | Deskripsi                  |
| --------------------- | -------------------------- |
| `go build -o app .`   | Build binary backend       |
| `npm run build`       | Build produksi frontend    |
| `npm run preview`     | Preview hasil build        |
| `npm run format`      | Format kode frontend       |
