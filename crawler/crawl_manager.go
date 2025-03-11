package crawler

/*
	CrawlManager is a construct developed to
		1. Co-ordinate between crawl workers,
		2. Aggregate metrics and
		3. Manage the crawl process (Shutdown and system limits).
*/
import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"sync"
)

type CrawlManager struct {
	mu                     sync.Mutex
	AvailableGoroutines    int
	runningAverage         float64
	ReceivedShutdownSignal bool
	Failures               int
	Success                int
}

func (cm *CrawlManager) RegisterGoroutine() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.AvailableGoroutines <= 0 {
		log.Warn().Msg("not enough available goroutines")
		return errors.New("not enough available goroutines")
	}
	cm.AvailableGoroutines--
	//log.Debug().Int("available-goroutines", cm.AvailableGoroutines).Msg("goroutine allocated")
	return nil
}

func (cm *CrawlManager) DeregisterGoroutine() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.AvailableGoroutines++
	//log.Debug().Int("available-goroutines", cm.AvailableGoroutines).Msg("goroutine de-allocated")
	return nil
}

func (cm *CrawlManager) IsGoroutineAvailable() bool {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	return cm.AvailableGoroutines > 0
}

func (cm *CrawlManager) UpdateRunningAverage(responseTime int64) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if cm.runningAverage == 0 {
		cm.runningAverage = float64(responseTime)
	} else {
		cm.runningAverage = (cm.runningAverage + float64(responseTime)) / 2.0
	}
}

func (cm *CrawlManager) RecordFailure() {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.Failures++
}

func (cm *CrawlManager) RecordSuccess() {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.Success++
}

func (cm *CrawlManager) GenerateReport(printToConsole bool) map[string]any {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	var totalRequests int
	var errorRate float64
	if cm.Failures > 0 || cm.Success > 0 {
		totalRequests = cm.Failures + cm.Success
		errorRate = float64(cm.Failures) / float64(totalRequests)
	}
	report := map[string]any{
		"averageResponseTime": cm.runningAverage,
		"totalProcessed":      totalRequests,
		"errorRate":           errorRate,
	}

	if printToConsole {
		fmt.Println("#####################################################")
		fmt.Println("Run report:")
		fmt.Printf("Total processed URLs:\t%v\n", report["totalProcessed"])
		fmt.Printf("Error rate:\t%v\n", report["errorRate"])
		fmt.Printf("Average response time:\t%vs\n", cm.runningAverage/1000.0)
		fmt.Println("#####################################################")
	}
	return report
}
func (cm *CrawlManager) ShutdownCrawls() {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.ReceivedShutdownSignal = true
}
