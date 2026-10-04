package models

type DashboardSummary struct {
	TotalHutang             int     `json:"total_hutang"`
	TotalHutangTerbayar     int     `json:"total_hutang_terbayar"`
	TotalTagihanBulanIni    int     `json:"total_tagihan_bulan_ini"`
	TagihanTerbayarBulanIni int     `json:"tagihan_terbayar_bulan_ini"`
	PersentaseTerbayar      float64 `json:"persentase_terbayar"` // Pakai float biar bisa desimal (misal 50.5%)

	TotalTagihanH3          int           `json:"total_tagihan_h3"` // Total rupiahnya
	ReminderH3              []Transaction `json:"reminder_h3"`      // List data buat di Modal
}