package worker

import (
	"fmt"
	"log"
	"time"

	"reserveflow-v1/dao"
)

// expiredJobs: işçilerin (worker goroutine'lerin) iş alacağı kanal (tepsi).
// Buffered Channel — aynı anda 100 rezervasyon ID'si bekleyebilir, bloke olmaz.
var expiredJobs chan uint

// InitWorkerPool sunucu başlarken bir kez çağrılır.
// numWorkers: kaç adet paralel işçi goroutine başlatılacağı.
func InitWorkerPool(numWorkers int) {
	expiredJobs = make(chan uint, 100)

	log.Printf("[WorkerPool] %d işçi başlatılıyor...\n", numWorkers)

	// Her işçiyi ayrı bir goroutine olarak başlat
	for i := 1; i <= numWorkers; i++ {
		go processJob(i, expiredJobs)
	}

	// Şefi (dispatcher) başlat — periyodik DB taraması yapar
	go dispatcher()

	log.Println("[WorkerPool] Sistem hazır. Süresi dolan rezervasyonlar izleniyor.")
}

// processJob: Her işçinin döngüsü.
// Kanaldan bir ID geldiğinde o rezervasyonu "expired" yapar.
// Kanal kapanana kadar beklemede kalır (for range pattern).
func processJob(workerID int, jobs <-chan uint) {
	for resID := range jobs {
		err := dao.MarkReservationAsExpired(resID)
		if err != nil {
			log.Printf("[Worker-%d] HATA: Rezervasyon #%d güncellenemedi: %v\n", workerID, resID, err)
		} else {
			fmt.Printf("[Worker-%d] Rezervasyon #%d → 'expired' yapıldı.\n", workerID, resID)
		}
	}
}

// dispatcher: Periyodik olarak (her 1 dakika) DB'yi tarar.
// Süresi geçmiş ama hâlâ "held" durumundaki rezervasyonları bulur
// ve ID'lerini expiredJobs kanalına gönderir.
func dispatcher() {
	scanAndDispatch()

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		scanAndDispatch()
	}
}

// scanAndDispatch: tek bir DB tarama turunu gerçekleştirir.
// Pluck sadece ID sütununu çeker — büyük tablolarda RAM dostudur.
func scanAndDispatch() {
	expiredIDs, err := dao.GetExpiredReservationIDs(time.Now())
	if err != nil {
		log.Printf("[Dispatcher] DB tarama hatası: %v\n", err)
		return
	}

	if len(expiredIDs) == 0 {
		return
	}

	log.Printf("[Dispatcher] %d adet süresi dolmuş rezervasyon bulundu, işçilere gönderiliyor...\n", len(expiredIDs))

	for _, id := range expiredIDs {
		expiredJobs <- id
	}
}
