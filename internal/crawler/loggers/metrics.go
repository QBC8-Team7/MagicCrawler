package loggers

import (
	"fmt"

	"github.com/QBC8-Team7/MagicCrawler/pkg/db/sqlc"
	"github.com/QBC8-Team7/MagicCrawler/pkg/logger"
	"github.com/QBC8-Team7/MagicCrawler/pkg/utils"
)

func MetricLog(metricLogger logger.AppLogger, err error, usage utils.Usage, crawlJob sqlc.CrawlJob) {

	status := "succeed"
	if err != nil {
		status = "failed"
	}

	line := fmt.Sprintf("jobid: %d | type: %s | status: %s | time: %v | RAM: %.2f MB | CPU_Before: %.2f%% CPU_After %.2f%%", crawlJob.ID, crawlJob.PageType, status, usage.Time, usage.RAM, usage.CPU_Before, usage.CPU_After)

	if err != nil {
		metricLogger.Error(line)
	} else {
		metricLogger.Info(line)
	}
}
