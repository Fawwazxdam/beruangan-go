package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"project-keuangan/config" // Sesuaikan dengan folder DB kamu
)

// Struct untuk membaca pesan JSON dari Telegram
type TelegramUpdate struct {
	Message struct {
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
		Text string `json:"text"`
	} `json:"message"`
}

// GANTI PAKAI TOKEN DARI BOTFATHER TADI
const BotToken = "TARUH_TOKEN_BOT_KAMU_DI_SINI"

func TelegramWebhook(w http.ResponseWriter, r *http.Request) {
	// 1. Tangkap JSON yang dilempar Telegram
	var update TelegramUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "Gagal baca data Telegram", http.StatusBadRequest)
		return
	}

	chatID := update.Message.Chat.ID
	text := update.Message.Text

	// Abaikan kalau pesannya kosong atau cuma perintah /start
	if text == "" || text == "/start" {
		kirimBalasan(chatID, "Halo! Ketik pengeluaranmu dengan format:\n[Nominal] [Nama]\nContoh: 25000 Beli Kopi")
		w.WriteHeader(http.StatusOK)
		return
	}

	// 2. Pecah Pesan (Misal: "25000 Beli Kopi" jadi ["25000", "Beli Kopi"])
	parts := strings.SplitN(text, " ", 2)
	if len(parts) < 2 {
		kirimBalasan(chatID, "Format salah Bang! Harus ada spasi antara nominal dan nama. Contoh: 25000 Nasi Padang")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Konversi nominal (string) jadi angka (int)
	amount, err := strconv.Atoi(parts[0])
	if err != nil {
		kirimBalasan(chatID, "Nominalnya harus angka Bang! Jangan pakai titik atau koma ya. Contoh: 25000")
		w.WriteHeader(http.StatusOK)
		return
	}
	title := parts[1]

	// 3. Simpan ke Database SQLite!
	// JADIKAN CHAT ID SEBAGAI USER ID (KODE RAHASIA)
	userID := fmt.Sprintf("tg-%d", chatID) 
	
	// Format tanggal hari ini
	dueDate := "2026-10-16" // Idealnya pakai: time.Now().Format("2006-01-02")
	groupID := fmt.Sprintf("TRX-%d", chatID) // Asal unik aja untuk single transaction

	query := "INSERT INTO transactions (user_id, title, amount, type, due_date, status, group_id) VALUES (?, ?, ?, ?, ?, ?, ?)"
	_, err = config.DB.Exec(query, userID, title, amount, "PENGELUARAN", dueDate, "PAID", groupID)

	if err != nil {
		kirimBalasan(chatID, "Waduh, gagal nyatet ke database nih!")
		w.WriteHeader(http.StatusOK)
		return
	}

	// 4. Kasih laporan sukses ke Telegram
	kirimBalasan(chatID, fmt.Sprintf("✅ Mantap! Rp %d buat '%s' udah masuk catatan.", amount, title))
	w.WriteHeader(http.StatusOK)
}

// Fungsi pembantu buat ngirim chat balasan dari Bot ke HP Kamu
func kirimBalasan(chatID int64, text string) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", BotToken)
	
	body, _ := json.Marshal(map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
	})
	
	http.Post(url, "application/json", bytes.NewBuffer(body))
}