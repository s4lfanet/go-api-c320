package usecase

import (
	"context"
	"fmt"

	"github.com/gosnmp/gosnmp"
	"github.com/rs/zerolog/log"
	"github.com/s4lfanet/go-api-c320/config"
	apperrors "github.com/s4lfanet/go-api-c320/internal/errors"
	"github.com/s4lfanet/go-api-c320/internal/model"
	"github.com/s4lfanet/go-api-c320/internal/repository"
	"golang.org/x/sync/singleflight"
)

// PonUseCaseInterface defines the interface for PON port operations
type PonUseCaseInterface interface {
	GetPonPortInfo(ctx context.Context, boardID, ponID int) (*model.PonPortInfo, error)
}

// ponUsecase implements PON port usecase
type ponUsecase struct {
	snmpRepository  repository.SnmpRepositoryInterface
	redisRepository repository.OnuRedisRepositoryInterface
	cfg             *config.Config
	sg              singleflight.Group
}

// NewPonUsecase creates a new PON usecase instance
func NewPonUsecase(
	snmpRepository repository.SnmpRepositoryInterface,
	redisRepository repository.OnuRedisRepositoryInterface,
	cfg *config.Config,
) PonUseCaseInterface {
	return &ponUsecase{
		snmpRepository:  snmpRepository,
		redisRepository: redisRepository,
		cfg:             cfg,
		sg:              singleflight.Group{},
	}
}

// GetPonPortInfo retrieves PON port information
func (u *ponUsecase) GetPonPortInfo(ctx context.Context, boardID, ponID int) (*model.PonPortInfo, error) {
	log.Info().Msgf("Getting PON port info for Board %d PON %d", boardID, ponID)

	// Generate PON index based on board and PON ID
	// For V2.1: 268500992 + (ponID * 256) for Board 1
	// For Board 2: 268509184 + (ponID * 256)
	var ponIndex int
	if boardID == 1 {
		ponIndex = 268500992 + (ponID * 256)
	} else if boardID == 2 {
		ponIndex = 268509184 + (ponID * 256)
	} else {
		return nil, apperrors.NewConfigError("invalid board ID", nil)
	}

	ponInfo := &model.PonPortInfo{
		Board: boardID,
		PON:   ponID,
	}

	// Get admin status from .3.11.3.1.1.{pon_index}
	adminStatusOID := fmt.Sprintf("%s.3.11.3.1.1.%d", u.cfg.OltCfg.BaseOID1, ponIndex)
	if result, err := u.snmpRepository.Get([]string{adminStatusOID}); err == nil && len(result.Variables) > 0 {
		if status, ok := result.Variables[0].Value.(int); ok {
			if status == 1 {
				ponInfo.AdminStatus = "enabled"
			} else {
				ponInfo.AdminStatus = "disabled"
			}
		}
	}

	// Get operational status and distance from .3.11.5.1.{col}.{pon_index}
	distanceOID := fmt.Sprintf("%s.3.11.5.1.3.%d", u.cfg.OltCfg.BaseOID1, ponIndex)
	if result, err := u.snmpRepository.Get([]string{distanceOID}); err == nil && len(result.Variables) > 0 {
		if distance, ok := result.Variables[0].Value.(int); ok {
			ponInfo.Distance = distance
		}
	}

	// Get oper status from .3.11.5.1.4.{pon_index}
	operStatusOID := fmt.Sprintf("%s.3.11.5.1.4.%d", u.cfg.OltCfg.BaseOID1, ponIndex)
	if result, err := u.snmpRepository.Get([]string{operStatusOID}); err == nil && len(result.Variables) > 0 {
		if status, ok := result.Variables[0].Value.(int); ok {
			if status == 2 {
				ponInfo.OperStatus = "up"
			} else {
				ponInfo.OperStatus = "down"
			}
		}
	}

	// Count ONUs by doing SNMP walk on ONU table
	onuCountOID := fmt.Sprintf("%s.3.13.3.1.5.%d", u.cfg.OltCfg.BaseOID1, ponIndex)
	onuCount := 0
	err := u.snmpRepository.Walk(onuCountOID, func(pdu gosnmp.SnmpPDU) error {
		onuCount++
		return nil
	})
	if err != nil {
		log.Warn().Err(err).Msg("Failed to count ONUs")
	}
	ponInfo.OnuCount = onuCount

	return ponInfo, nil
}
