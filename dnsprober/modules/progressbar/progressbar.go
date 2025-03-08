package progressbar

import (
	"fmt"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/RevoltSecurities/Dnsprober/dnsprober/modules/logger"
)

const (
	tcl = "\r\x1b[2K"
)

var Logger = logger.New(true)

type ProgressBar struct {
	ReqCount   int64
	ReqTotal   int64
	ErrorCount int64
	StartedAt  time.Time
	BarWidth   int 
}

func New(total int64) *ProgressBar {
	return &ProgressBar{
		ReqTotal:  total,
		StartedAt: time.Now(),
		BarWidth:  50, //returned as default limit
	}
}

func (pb *ProgressBar) Increment(success, errs int64) {
	atomic.AddInt64(&pb.ReqCount, success)
	atomic.AddInt64(&pb.ErrorCount, errs)
}

func (pb *ProgressBar) Render() {
	reqCount := atomic.LoadInt64(&pb.ReqCount)
	reqTotal := atomic.LoadInt64(&pb.ReqTotal)
	errorCount := atomic.LoadInt64(&pb.ErrorCount)

	// Calculate percentage complete.
	var percentage float64
	if reqTotal > 0 {
		percentage = float64(reqCount) / float64(reqTotal) * 100
	}

	// Calculate how many blocks to fill based on the desired bar width.
	barWidth := pb.BarWidth
	blockCounts := int(percentage / (100.0 / float64(barWidth)))
	if blockCounts < 0 {
		blockCounts = 0
	} else if blockCounts > barWidth {
		blockCounts = barWidth
	}
	bar := strings.Repeat("█", blockCounts) + strings.Repeat(" ", barWidth-blockCounts)

	// Calculate elapsed time.
	elapsed := time.Since(pb.StartedAt).Seconds()
	var rate float64
	if elapsed > 0 {
		rate = float64(reqCount) / elapsed
	}

	// Calculate ETA.
	var etaSeconds float64
	if reqCount > 0 && reqCount < reqTotal {
		estimatedTotal := elapsed * (float64(reqTotal) / float64(reqCount))
		etaSeconds = estimatedTotal - elapsed
	} else {
		etaSeconds = 0
	}

	// Format elapsed time.
	hours := int64(elapsed) / 3600
	minutes := (int64(elapsed) % 3600) / 60
	seconds := int64(elapsed) % 60

	// Format ETA.
	etaHours := int64(etaSeconds) / 3600
	etaMinutes := (int64(etaSeconds) % 3600) / 60
	etaRem := int64(etaSeconds) % 60

	pb.Printer(bar, percentage, hours, minutes, seconds, reqCount, reqTotal, rate, errorCount, etaHours, etaMinutes, etaRem)
}


func (pb *ProgressBar) Printer(bar string, percentage float64, hours, minutes, seconds, reqCount, reqTotal int64, rate float64, errorCount, etaHours, etaMinutes, etaSeconds int64) {
	output := fmt.Sprintf("%s[%s] %6.2f%% [%02d:%02d:%02d] ETA [%02d:%02d:%02d] [%d/%d] %.1f req/sec Errors: %d\r",
		tcl, bar, percentage, hours, minutes, seconds, etaHours, etaMinutes, etaSeconds, reqCount, reqTotal, rate, errorCount)
	formattedOutput := Logger.Bolder(output)
	fmt.Fprint(os.Stderr, formattedOutput)
}
