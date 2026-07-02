package worker

import (
	"fmt"
	"log"
	"time"

	"reserveflow-v1/service"
)

var resService = service.NewReservationService()

// expiredJobs: işçi goroutine'lerin iş aldığı kanal (max 100 ID tamponlu).
var expiredJobs chan uint

// InitWorkerPool sunucu başlarken çağrılır; n adet işçi + dispatcher başlatır.
func InitWorkerPool(numWorkers int) {
	expiredJobs = make(chan uint, 100)

	log.Printf("[WorkerPool] %d işçi başlatılıyor...\n", numWorkers)

	for i := 1; i <= numWorkers; i++ {
		go processJob(i, expiredJobs)
	}

	// Her 1 dakikada bir otomatik tarama yapacak dispatcher'ı başlatıyoruz
	go startDispatcher()

	log.Println("[WorkerPool] Sistem hazır. Süresi dolan rezervasyonlar izleniyor.")
}

// processJob kanaldan gelen rezervasyon ID'lerini "expired" olarak işler.
func processJob(workerID int, jobs <-chan uint) {
	for resID := range jobs {
		err := resService.ExpireReservation(resID)
		if err != nil {
			log.Printf("[Worker-%d] HATA: Rezervasyon #%d güncellenemedi: %v\n", workerID, resID, err)
		} else {
			fmt.Printf("[Worker-%d] Rezervasyon #%d → 'expired' yapıldı.\n", workerID, resID)
		}
	}
}

// startDispatcher her 1 dakikada bir veritabanını tarar.
func startDispatcher() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		_, err := runExpirationCheck()
		if err != nil {
			log.Printf("[Dispatcher] Hata oluştu: %v\n", err)
		}
	}
}

// runExpirationCheck süresi dolmuş rezervasyonları bulup işçilere gönderir.
func runExpirationCheck() (int, error) {
	expiredIDs, err := resService.GetExpiredReservationIDs(time.Now())
	if err != nil {
		log.Printf("[Trigger] DB tarama hatası: %v\n", err)
		return 0, err
	}

	if len(expiredIDs) == 0 {
		return 0, nil
	}

	log.Printf("[Trigger] %d adet süresi dolmuş rezervasyon bulundu, işçilere gönderiliyor...\n", len(expiredIDs))

	for _, id := range expiredIDs {
		expiredJobs <- id
	}

	return len(expiredIDs), nil
}
