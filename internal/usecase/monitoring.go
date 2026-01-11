package usecase

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gosnmp/gosnmp"
	"github.com/rs/zerolog/log"
	"github.com/s4lfanet/go-api-c320/config"
	apperrors "github.com/s4lfanet/go-api-c320/internal/errors"
	"github.com/s4lfanet/go-api-c320/internal/model"
	"github.com/s4lfanet/go-api-c320/internal/repository"
	"github.com/s4lfanet/go-api-c320/internal/utils"
)

// MonitoringUsecase handles real-time ONU monitoring operations
type MonitoringUsecase struct {
	snmp      *gosnmp.GoSNMP
	cfg       *config.Config
	onuRepo   *repository.OnuRepository
	telnetMgr *repository.TelnetSessionManager
}

// NewMonitoringUsecase creates a new MonitoringUsecase instance
func NewMonitoringUsecase(snmp *gosnmp.GoSNMP, cfg *config.Config, onuRepo *repository.OnuRepository, telnetMgr *repository.TelnetSessionManager) *MonitoringUsecase {
	return &MonitoringUsecase{
		snmp:      snmp,
		cfg:       cfg,
		onuRepo:   onuRepo,
		telnetMgr: telnetMgr,
	}
}

// GetONUMonitoring fetches real-time monitoring data for a single ONU
func (uc *MonitoringUsecase) GetONUMonitoring(ctx context.Context, ponPort string, onuID int) (*model.ONUMonitoringInfo, error) {
	log.Info().Str("pon_port", ponPort).Int("onu_id", onuID).Msg("Getting ONU monitoring data")

	// Get PON index from config
	boardPonKey := config.BoardPonKey{BoardID: 1, PonID: utils.ConvertStringToInt(ponPort)}
	_, exists := uc.cfg.BoardPonMap[boardPonKey]
	if !exists {
		return nil, apperrors.NewNotFoundError("PON port", ponPort)
	}

	ponIndex := repository.CalculatePonIndex(1, utils.ConvertStringToInt(ponPort))
	onuIndexStr := fmt.Sprintf("%d.%d", ponIndex, onuID)

	monitoring := &model.ONUMonitoringInfo{
		PonPort:    ponPort,
		OnuID:      onuID,
		LastUpdate: time.Now(),
	}

	// Get ONU basic info (serial, model, firmware)
	serialOID := fmt.Sprintf("1.3.6.1.4.1.3902.1012.3.13.3.1.5.%s", onuIndexStr)    // Device SN
	modelOID := fmt.Sprintf("1.3.6.1.4.1.3902.1012.3.13.3.1.10.%s", onuIndexStr)    // Model
	firmwareOID := fmt.Sprintf("1.3.6.1.4.1.3902.1012.3.13.3.1.11.%s", onuIndexStr) // Firmware
	statusOID := fmt.Sprintf("1.3.6.1.4.1.3902.1012.3.31.4.1.100.%s", onuIndexStr)  // Online status

	oids := []string{serialOID, modelOID, firmwareOID, statusOID}
	result, err := uc.snmp.Get(oids)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get ONU basic info")
		return nil, apperrors.NewSNMPError("get ONU basic info", err)
	}

	// Parse basic info
	for _, variable := range result.Variables {
		oidStr := variable.Name
		switch {
		case oidStr == serialOID:
			monitoring.SerialNumber = utils.ExtractStringValue(variable)
		case oidStr == modelOID:
			monitoring.Model = utils.ExtractStringValue(variable)
		case oidStr == firmwareOID:
			monitoring.FirmwareVer = utils.ExtractStringValue(variable)
		case oidStr == statusOID:
			monitoring.OnlineStatus = utils.ExtractIntValue(variable)
		}
	}

	// Get ONU statistics
	rxPacketsOID := fmt.Sprintf("1.3.6.1.4.1.3902.1012.3.31.4.1.3.%s", onuIndexStr) // RX packets
	rxBytesOID := fmt.Sprintf("1.3.6.1.4.1.3902.1012.3.31.4.1.6.%s", onuIndexStr)   // RX bytes

	statOIDs := []string{rxPacketsOID, rxBytesOID}
	statResult, err := uc.snmp.Get(statOIDs)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get ONU statistics, continuing without stats")
	} else {
		stats := &model.ONUStatistics{}
		for _, variable := range statResult.Variables {
			oidStr := variable.Name
			switch {
			case oidStr == rxPacketsOID:
				stats.RxPackets = utils.ExtractUint64Value(variable)
			case oidStr == rxBytesOID:
				stats.RxBytes = utils.ExtractUint64Value(variable)
			}
		}

		// Calculate rate (simplified - would need time-based calculation for accurate rate)
		if stats.RxBytes > 0 {
			stats.RxRate = utils.FormatBytesRate(stats.RxBytes)
		}
		monitoring.Statistics = stats
	}

	// Get optical info via Telnet (V2.1.0 doesn't have SNMP OIDs for optical power)
	if uc.telnetMgr != nil {
		opticalInfo, err := uc.telnetMgr.GetONUOpticalInfo(ctx, 1, utils.ConvertStringToInt(ponPort), onuID)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to get ONU optical info via Telnet, continuing without optical data")
		} else if opticalInfo != nil {
			monitoring.Optical = &model.OpticalInfo{
				RxPower:           opticalInfo.RxPower,
				TxPower:           opticalInfo.TxPower,
				OLTRxPower:        opticalInfo.OLTRxPower,
				Temperature:       opticalInfo.Temperature,
				Voltage:           opticalInfo.Voltage,
				BiasCurrent:       opticalInfo.BiasCurrent,
				RxPowerStatus:     opticalInfo.RxPowerStatus,
				TxPowerStatus:     opticalInfo.TxPowerStatus,
				TemperatureStatus: opticalInfo.TemperatureStatus,
			}
		}
	}

	log.Info().
		Str("serial", monitoring.SerialNumber).
		Str("model", monitoring.Model).
		Int("status", monitoring.OnlineStatus).
		Msg("Successfully retrieved ONU monitoring data")

	return monitoring, nil
}

// GetPONMonitoring fetches aggregated monitoring data for a PON port
func (uc *MonitoringUsecase) GetPONMonitoring(ctx context.Context, ponPort string) (*model.PONMonitoringInfo, error) {
	log.Info().Str("pon_port", ponPort).Msg("Getting PON monitoring data")

	// Get PON config
	boardPonKey := config.BoardPonKey{BoardID: 1, PonID: utils.ConvertStringToInt(ponPort)}
	_, exists := uc.cfg.BoardPonMap[boardPonKey]
	if !exists {
		return nil, apperrors.NewNotFoundError("PON port", ponPort)
	}

	ponIndex := repository.CalculatePonIndex(1, utils.ConvertStringToInt(ponPort))

	monitoring := &model.PONMonitoringInfo{
		PonPort:    ponPort,
		PonIndex:   ponIndex,
		LastUpdate: time.Now(),
		ONUs:       []model.ONUMonitoringInfo{},
	}

	// Get all ONUs for this PON
	onus, err := uc.onuRepo.GetByBoardIDAndPonID(ctx, 1, utils.ConvertStringToInt(ponPort))
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get ONUs for PON, continuing with empty list")
	}

	// Fetch monitoring for each ONU
	monitoring.OnuCount = len(onus)
	for _, onu := range onus {
		onuMon, err := uc.GetONUMonitoring(ctx, ponPort, onu.ID)
		if err != nil {
			log.Warn().Err(err).Int("onu_id", onu.ID).Msg("Failed to get ONU monitoring")
			continue
		}

		monitoring.ONUs = append(monitoring.ONUs, *onuMon)
		if onuMon.OnlineStatus == 1 {
			monitoring.OnlineCount++
		} else {
			monitoring.OfflineCount++
		}
	}

	// Get PON port statistics
	ponIndexStr := fmt.Sprintf("%d.1", ponIndex)
	rxPacketsOID := fmt.Sprintf(".1.3.6.1.4.1.3902.1012.3.31.5.1.3.%s", ponIndexStr)
	rxBytesOID := fmt.Sprintf(".1.3.6.1.4.1.3902.1012.3.31.5.1.6.%s", ponIndexStr)

	statOIDs := []string{rxPacketsOID, rxBytesOID}
	statResult, err := uc.snmp.Get(statOIDs)
	if err == nil && len(statResult.Variables) > 0 {
		stats := &model.PONStatistics{}
		for _, variable := range statResult.Variables {
			oidStr := variable.Name
			switch {
			case strings.Contains(oidStr, ".3.31.5.1.3."):
				stats.RxPackets = utils.ExtractUint64Value(variable)
			case strings.Contains(oidStr, ".3.31.5.1.6."):
				stats.RxBytes = utils.ExtractUint64Value(variable)
			}
		}

		if stats.RxBytes > 0 {
			stats.RxRate = utils.FormatBytesRate(stats.RxBytes)
		}
		monitoring.Statistics = stats
	}

	log.Info().
		Int("onu_count", monitoring.OnuCount).
		Int("online", monitoring.OnlineCount).
		Int("offline", monitoring.OfflineCount).
		Msg("Successfully retrieved PON monitoring data")

	return monitoring, nil
}

// GetOLTMonitoring fetches overall OLT monitoring summary
func (uc *MonitoringUsecase) GetOLTMonitoring(ctx context.Context) (*model.OLTMonitoringSummary, error) {
	log.Info().Msg("Getting OLT monitoring summary")

	summary := &model.OLTMonitoringSummary{
		PONPorts:   []model.PONMonitoringInfo{},
		LastUpdate: time.Now(),
	}

	// Get all configured PON ports
	for key := range uc.cfg.BoardPonMap {
		ponPort := strconv.Itoa(key.PonID)
		ponMon, err := uc.GetPONMonitoring(ctx, ponPort)
		if err != nil {
			log.Warn().Err(err).Str("pon", ponPort).Msg("Failed to get PON monitoring")
			continue
		}

		summary.PONPorts = append(summary.PONPorts, *ponMon)
		summary.TotalONUs += ponMon.OnuCount
		summary.OnlineONUs += ponMon.OnlineCount
		summary.OfflineONUs += ponMon.OfflineCount
	}

	log.Info().
		Int("total_onus", summary.TotalONUs).
		Int("online", summary.OnlineONUs).
		Int("offline", summary.OfflineONUs).
		Int("pon_ports", len(summary.PONPorts)).
		Msg("Successfully retrieved OLT monitoring summary")

	return summary, nil
}
