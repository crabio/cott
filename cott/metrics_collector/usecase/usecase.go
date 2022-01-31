package usecase

import (
	"time"

	container_launcher "github.com/iakrevetkho/components-tests/cott/container_launcher/usecase"
	"github.com/iakrevetkho/components-tests/cott/domain"
	"github.com/sirupsen/logrus"
)

const (
	DEFAULT_NETWORK = "eth0"
)

type MetricsCollectorUsecase interface {
	CollectStepMetrics(step *domain.TestCaseStep) error
}

type metricsCollectorUsecase struct {
	tcra *domain.TestCaseResultsAccumulator
	cluc container_launcher.ContainerLauncherUsecase
}

func NewMetricsCollectorUsecase(tcra *domain.TestCaseResultsAccumulator, cluc container_launcher.ContainerLauncherUsecase) MetricsCollectorUsecase {
	mcuc := new(metricsCollectorUsecase)
	mcuc.tcra = tcra
	mcuc.cluc = cluc
	return mcuc
}

func (mcuc *metricsCollectorUsecase) CollectStepMetrics(step *domain.TestCaseStep) error {
	tcsra := new(domain.TestCaseStepResultsAccumulator)

	stats, err := mcuc.cluc.GetContainerStats(step.ContainerID)
	if err != nil {
		return err
	}

	startCpuTotalUsage := stats.CPUStats.CPUUsage.TotalUsage
	startMemUsage := stats.MemoryStats.Usage
	startStorageReadUsage := stats.StorageStats.ReadSizeBytes
	startStorageWriteUsage := stats.StorageStats.WriteSizeBytes
	startNetworkRxUsage := stats.Networks[DEFAULT_NETWORK].RxBytes
	startNetworkTxUsage := stats.Networks[DEFAULT_NETWORK].TxBytes

	startTime := time.Now()
	if err := step.StepFunc(); err != nil {
		logrus.WithError(err).WithField("step", step).Warn("error on step execution")
		tcsra.AddError(err.Error())
		return err
	}
	tcsra.AddMetric(domain.MetricMeta_Duration, float64(time.Since(startTime).Microseconds()))

	stats, err = mcuc.cluc.GetContainerStats(step.ContainerID)
	if err != nil {
		return err
	}

	tcsra.AddMetric(domain.MetricMeta_CpuUsage, float64(stats.CPUStats.CPUUsage.TotalUsage-startCpuTotalUsage))
	tcsra.AddMetric(domain.MetricMeta_MemoryUsage, float64(stats.MemoryStats.Usage))
	tcsra.AddMetric(domain.MetricMeta_MemoryUsageDiff, float64(stats.MemoryStats.Usage-startMemUsage))
	tcsra.AddMetric(domain.MetricMeta_StorageReadUsage, float64(stats.StorageStats.ReadSizeBytes-startStorageReadUsage))
	tcsra.AddMetric(domain.MetricMeta_StorageWriteUsage, float64(stats.StorageStats.WriteSizeBytes-startStorageWriteUsage))
	tcsra.AddMetric(domain.MetricMeta_NetworkReceiveUsage, float64(stats.Networks[DEFAULT_NETWORK].RxBytes-startNetworkRxUsage))
	tcsra.AddMetric(domain.MetricMeta_NetworkSendUsage, float64(stats.Networks[DEFAULT_NETWORK].TxBytes-startNetworkTxUsage))

	return nil
}
