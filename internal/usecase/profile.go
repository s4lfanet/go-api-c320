package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/gosnmp/gosnmp"
	"github.com/rs/zerolog/log"
	"github.com/s4lfanet/go-api-c320/config"
	apperrors "github.com/s4lfanet/go-api-c320/internal/errors"
	"github.com/s4lfanet/go-api-c320/internal/model"
	"github.com/s4lfanet/go-api-c320/internal/repository"
	"golang.org/x/sync/singleflight"
)

// ProfileUseCaseInterface defines the interface for traffic profile operations
type ProfileUseCaseInterface interface {
	GetAllTrafficProfiles(ctx context.Context) ([]*model.TrafficProfile, error)
	GetTrafficProfile(ctx context.Context, profileID int) (*model.TrafficProfile, error)
	GetAllVlanProfiles(ctx context.Context) ([]*model.VlanProfile, error)
}

// profileUsecase implements traffic profile usecase
type profileUsecase struct {
	snmpRepository  repository.SnmpRepositoryInterface
	redisRepository repository.OnuRedisRepositoryInterface
	cfg             *config.Config
	sg              singleflight.Group
}

// NewProfileUsecase creates a new profile usecase instance
func NewProfileUsecase(
	snmpRepository repository.SnmpRepositoryInterface,
	redisRepository repository.OnuRedisRepositoryInterface,
	cfg *config.Config,
) ProfileUseCaseInterface {
	return &profileUsecase{
		snmpRepository:  snmpRepository,
		redisRepository: redisRepository,
		cfg:             cfg,
		sg:              singleflight.Group{},
	}
}

// GetAllTrafficProfiles retrieves all traffic profiles
func (u *profileUsecase) GetAllTrafficProfiles(ctx context.Context) ([]*model.TrafficProfile, error) {
	log.Info().Msg("Getting all traffic profiles")

	profiles := make(map[int]*model.TrafficProfile)

	// Walk the profile name OID to get all profile IDs
	// OID: .3.26.1.1.2.{profile_id}
	nameOID := fmt.Sprintf("%s.3.26.1.1.2", u.cfg.OltCfg.BaseOID1)

	err := u.snmpRepository.Walk(nameOID, func(pdu gosnmp.SnmpPDU) error {
		// Extract profile ID from OID
		oidParts := strings.Split(pdu.Name, ".")
		if len(oidParts) < 1 {
			return nil
		}

		profileIDStr := oidParts[len(oidParts)-1]
		profileID := 0
		_, _ = fmt.Sscanf(profileIDStr, "%d", &profileID)

		if profileID == 0 {
			return nil
		}

		// Get profile name
		name := ""
		switch v := pdu.Value.(type) {
		case string:
			name = v
		case []byte:
			name = string(v)
		default:
			name = fmt.Sprintf("%v", v)
		}

		profiles[profileID] = &model.TrafficProfile{
			ProfileID: profileID,
			Name:      name,
		}

		return nil
	})

	if err != nil {
		log.Error().Err(err).Msg("Failed to walk traffic profile names")
		return nil, apperrors.NewSNMPError("failed to get traffic profiles", err)
	}

	// For each profile, get CIR, PIR, and MaxBW
	for profileID, profile := range profiles {
		// Get CIR (.3.26.1.1.3)
		cirOID := fmt.Sprintf("%s.3.26.1.1.3.%d", u.cfg.OltCfg.BaseOID1, profileID)
		if result, err := u.snmpRepository.Get([]string{cirOID}); err == nil && len(result.Variables) > 0 {
			if cir, ok := result.Variables[0].Value.(int); ok {
				profile.CIR = cir
			}
		}

		// Get PIR (.3.26.1.1.4)
		pirOID := fmt.Sprintf("%s.3.26.1.1.4.%d", u.cfg.OltCfg.BaseOID1, profileID)
		if result, err := u.snmpRepository.Get([]string{pirOID}); err == nil && len(result.Variables) > 0 {
			if pir, ok := result.Variables[0].Value.(int); ok {
				profile.PIR = pir
			}
		}

		// Get MaxBW (.3.26.1.1.5)
		maxBWOID := fmt.Sprintf("%s.3.26.1.1.5.%d", u.cfg.OltCfg.BaseOID1, profileID)
		if result, err := u.snmpRepository.Get([]string{maxBWOID}); err == nil && len(result.Variables) > 0 {
			if maxBW, ok := result.Variables[0].Value.(int); ok {
				profile.MaxBW = maxBW
			}
		}
	}

	// Convert map to slice
	result := make([]*model.TrafficProfile, 0, len(profiles))
	for _, profile := range profiles {
		result = append(result, profile)
	}

	log.Info().Int("count", len(result)).Msg("Successfully retrieved traffic profiles")
	return result, nil
}

// GetTrafficProfile retrieves a specific traffic profile by ID
func (u *profileUsecase) GetTrafficProfile(ctx context.Context, profileID int) (*model.TrafficProfile, error) {
	log.Info().Int("profile_id", profileID).Msg("Getting traffic profile")

	profile := &model.TrafficProfile{
		ProfileID: profileID,
	}

	// Get profile name (.3.26.1.1.2)
	nameOID := fmt.Sprintf("%s.3.26.1.1.2.%d", u.cfg.OltCfg.BaseOID1, profileID)
	result, err := u.snmpRepository.Get([]string{nameOID})
	if err != nil {
		return nil, apperrors.NewNotFoundError("traffic profile", map[string]int{"profile_id": profileID})
	}

	if len(result.Variables) > 0 {
		switch v := result.Variables[0].Value.(type) {
		case string:
			profile.Name = v
		case []byte:
			profile.Name = string(v)
		default:
			profile.Name = fmt.Sprintf("%v", v)
		}
	}

	// Get CIR (.3.26.1.1.3)
	cirOID := fmt.Sprintf("%s.3.26.1.1.3.%d", u.cfg.OltCfg.BaseOID1, profileID)
	if result, err := u.snmpRepository.Get([]string{cirOID}); err == nil && len(result.Variables) > 0 {
		if cir, ok := result.Variables[0].Value.(int); ok {
			profile.CIR = cir
		}
	}

	// Get PIR (.3.26.1.1.4)
	pirOID := fmt.Sprintf("%s.3.26.1.1.4.%d", u.cfg.OltCfg.BaseOID1, profileID)
	if result, err := u.snmpRepository.Get([]string{pirOID}); err == nil && len(result.Variables) > 0 {
		if pir, ok := result.Variables[0].Value.(int); ok {
			profile.PIR = pir
		}
	}

	// Get MaxBW (.3.26.1.1.5)
	maxBWOID := fmt.Sprintf("%s.3.26.1.1.5.%d", u.cfg.OltCfg.BaseOID1, profileID)
	if result, err := u.snmpRepository.Get([]string{maxBWOID}); err == nil && len(result.Variables) > 0 {
		if maxBW, ok := result.Variables[0].Value.(int); ok {
			profile.MaxBW = maxBW
		}
	}

	log.Info().Int("profile_id", profileID).Str("name", profile.Name).Msg("Successfully retrieved traffic profile")
	return profile, nil
}

// GetAllVlanProfiles retrieves all VLAN profiles
func (u *profileUsecase) GetAllVlanProfiles(ctx context.Context) ([]*model.VlanProfile, error) {
	log.Info().Msg("Getting all VLAN profiles")

	vlanProfiles := make(map[string]*model.VlanProfile)

	// Walk the VLAN profile OID to get all VLAN names
	// OID: .3.50.20.15.1.{col}.{vlan_name_ascii}
	// The VLAN name is encoded as ASCII decimal values in the OID
	baseOID := fmt.Sprintf("%s.3.50.20.15.1", u.cfg.OltCfg.BaseOID1)

	err := u.snmpRepository.Walk(baseOID, func(pdu gosnmp.SnmpPDU) error {
		// Extract VLAN name from OID
		// OID format: baseOID.column.length.ascii_char1.ascii_char2...
		oidParts := strings.Split(pdu.Name, ".")

		// Find where our base OID ends
		baseOIDParts := strings.Split(baseOID, ".")
		if len(oidParts) <= len(baseOIDParts)+2 {
			return nil
		}

		// Get the column number (first part after base OID)
		column := 0
		if len(oidParts) > len(baseOIDParts) {
			_, _ = fmt.Sscanf(oidParts[len(baseOIDParts)], "%d", &column)
		}

		// Get the length (second part after base OID)
		length := 0
		if len(oidParts) > len(baseOIDParts)+1 {
			_, _ = fmt.Sscanf(oidParts[len(baseOIDParts)+1], "%d", &length)
		}

		// Get ASCII values (skip column and length)
		asciiParts := oidParts[len(baseOIDParts)+2:]

		// Only convert the number of characters specified by length
		vlanName := ""
		for i := 0; i < length && i < len(asciiParts); i++ {
			asciiVal := 0
		_, _ = fmt.Sscanf(asciiParts[i], "%d", &asciiVal)
		if asciiVal > 0 && asciiVal < 128 {
			vlanName += string(rune(asciiVal))
		}
	}

		if vlanName == "" {
			return nil
		}

		// Initialize or update profile
		if _, exists := vlanProfiles[vlanName]; !exists {
			vlanProfiles[vlanName] = &model.VlanProfile{
				Name: vlanName,
			}
		}

		// Parse value based on column
		profile := vlanProfiles[vlanName]

		switch column {
		case 2: // VLAN ID
			if vlanID, ok := pdu.Value.(int); ok {
				profile.VlanID = vlanID
			}
		case 3: // Priority
			if priority, ok := pdu.Value.(int); ok {
				profile.Priority = priority
			}
		case 4: // Mode
			mode := ""
			switch v := pdu.Value.(type) {
			case string:
				mode = v
			case []byte:
				mode = string(v)
			case int:
				if v == 1 {
					mode = "tag"
				} else if v == 2 {
					mode = "untag"
				} else {
					mode = fmt.Sprintf("mode_%d", v)
				}
			default:
				mode = fmt.Sprintf("%v", v)
			}
			profile.Mode = mode
		case 5: // Description
			switch v := pdu.Value.(type) {
			case string:
				profile.Description = v
			case []byte:
				profile.Description = string(v)
			default:
				profile.Description = fmt.Sprintf("%v", v)
			}
		}

		return nil
	})

	if err != nil {
		log.Error().Err(err).Msg("Failed to walk VLAN profiles")
		return nil, apperrors.NewSNMPError("failed to get VLAN profiles", err)
	}

	// Convert map to slice
	result := make([]*model.VlanProfile, 0, len(vlanProfiles))
	for _, profile := range vlanProfiles {
		result = append(result, profile)
	}

	log.Info().Int("count", len(result)).Msg("Successfully retrieved VLAN profiles")
	return result, nil
}
