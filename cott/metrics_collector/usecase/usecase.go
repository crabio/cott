package usecase

import (
	"time"

	container_launcher "github.com/iakrevetkho/components-tests/cott/container_launcher/usecase"
	"github.com/iakrevetkho/components-tests/cott/domain"
	"github.com/sirupsen/logrus"
)

const DEFAULT_NETWORK = "eth0"

type MetricsCollectorUsecase interface {
	CalcStepDuration(stepFunc func() error, stepName string) error
	GetStepContainerStats(stepFunc func() error, stepName string, containerID string) error
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

// TODO Create global method

func (mcuc *metricsCollectorUsecase) CalcStepDuration(stepFunc func() error, stepName string) error {
	start := time.Now()
	if err := stepFunc(); err != nil {
		logrus.WithError(err).WithField("stepName", stepName).Warn("error on step execution")
		mcuc.tcra.AddError(stepName + ". " + err.Error())
		return err
	}
	duration := time.Since(start)
	logrus.WithFields(logrus.Fields{"duration": duration, "stepName": stepName}).Debug("step finished")
	mcuc.tcra.AddMetric(stepName+"Duration", domain.UnitOfMeasurePrefix_Micro, domain.UnitOfMeasure_Second, float64(duration.Microseconds()))
	return nil
}

func (mcuc *metricsCollectorUsecase) GetStepContainerStats(stepFunc func() error, stepName string, containerID string) error {
	stats, err := mcuc.cluc.GetContainerStats(containerID)
	if err != nil {
		return err
	}

	startCpuTotalUsage := stats.CPUStats.CPUUsage.TotalUsage
	startStorageReadUsage := stats.StorageStats.ReadSizeBytes
	startStorageWriteUsage := stats.StorageStats.WriteSizeBytes
	startNetworkRxUsage := stats.Networks[DEFAULT_NETWORK].RxBytes
	startNetworkTxUsage := stats.Networks[DEFAULT_NETWORK].TxBytes

	if err := stepFunc(); err != nil {
		logrus.WithError(err).WithField("stepName", stepName).Warn("error on step execution")
		mcuc.tcra.AddError(stepName + ". " + err.Error())
		return err
	}

	stats, err = mcuc.cluc.GetContainerStats(containerID)
	if err != nil {
		return err
	}

	resCpuTotalUsage := stats.CPUStats.CPUUsage.TotalUsage - startCpuTotalUsage
	resMemUsage := stats.MemoryStats.Usage
	resMaxMemUsage := stats.MemoryStats.MaxUsage
	resStorageReadUsage := stats.StorageStats.ReadSizeBytes - startStorageReadUsage
	resStorageWriteUsage := stats.StorageStats.WriteSizeBytes - startStorageWriteUsage
	resNetworkRxUsage := stats.Networks[DEFAULT_NETWORK].RxBytes - startNetworkRxUsage
	resNetworkTxUsage := stats.Networks[DEFAULT_NETWORK].TxBytes - startNetworkTxUsage

	logrus.WithFields(logrus.Fields{
		"cpuUsage":             resCpuTotalUsage,
		"resMemUsage":          resMemUsage,
		"resMaxMemUsage":       resMaxMemUsage,
		"resStorageReadUsage":  resStorageReadUsage,
		"resStorageWriteUsage": resStorageWriteUsage,
		"resNetworkRxUsage":    resNetworkRxUsage,
		"resNetworkTxUsage":    resNetworkTxUsage,
	}).Debug("result container usage")

	return nil
}
