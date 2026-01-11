package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/s4lfanet/go-api-c320/internal/model"
	"github.com/rs/zerolog/log"
)

// RebootONU reboots/resets an ONU
func (m *TelnetSessionManager) RebootONU(ctx context.Context, req *model.ONURebootRequest) error {
	log.Info().
		Str("pon_port", req.PONPort).
		Int("onu_id", req.ONUID).
		Msg("Rebooting ONU via Telnet")

	// Enter PON interface configuration mode
	interfaceCmd := fmt.Sprintf("interface gpon-olt_%s", req.PONPort)
	if _, err := m.ExecuteCommand(ctx, interfaceCmd); err != nil {
		return fmt.Errorf("failed to enter interface mode: %w", err)
	}

	// Execute ONU reset command
	resetCmd := fmt.Sprintf("onu reset %d", req.ONUID)
	if _, err := m.ExecuteCommand(ctx, resetCmd); err != nil {
		// Check if error is due to ONU not existing
		if strings.Contains(err.Error(), "not exist") || strings.Contains(err.Error(), "invalid") {
			return fmt.Errorf("ONU %s/%d not found", req.PONPort, req.ONUID)
		}
		return fmt.Errorf("failed to reset ONU: %w", err)
	}

	// Exit interface mode
	if _, err := m.ExecuteCommand(ctx, "exit"); err != nil {
		log.Warn().Err(err).Msg("Failed to exit interface mode, session will auto-recover")
	}

	log.Info().
		Str("pon_port", req.PONPort).
		Int("onu_id", req.ONUID).
		Msg("ONU reboot command executed successfully")

	return nil
}

// BlockONU disables an ONU (sets state to disable)
func (m *TelnetSessionManager) BlockONU(ctx context.Context, req *model.ONUBlockRequest) error {
	log.Info().
		Str("pon_port", req.PONPort).
		Int("onu_id", req.ONUID).
		Bool("block", req.Block).
		Msg("Blocking ONU via Telnet")

	// Enter PON interface configuration mode
	interfaceCmd := fmt.Sprintf("interface gpon-olt_%s", req.PONPort)
	if _, err := m.ExecuteCommand(ctx, interfaceCmd); err != nil {
		return fmt.Errorf("failed to enter interface mode: %w", err)
	}

	// Execute ONU state disable command
	stateCmd := fmt.Sprintf("onu %d state disable", req.ONUID)
	if _, err := m.ExecuteCommand(ctx, stateCmd); err != nil {
		// Check if error is due to ONU not existing
		if strings.Contains(err.Error(), "not exist") || strings.Contains(err.Error(), "invalid") {
			return fmt.Errorf("ONU %s/%d not found", req.PONPort, req.ONUID)
		}
		return fmt.Errorf("failed to disable ONU: %w", err)
	}

	// Exit interface mode
	if _, err := m.ExecuteCommand(ctx, "exit"); err != nil {
		log.Warn().Err(err).Msg("Failed to exit interface mode, session will auto-recover")
	}

	log.Info().
		Str("pon_port", req.PONPort).
		Int("onu_id", req.ONUID).
		Msg("ONU blocked successfully")

	return nil
}

// UnblockONU enables an ONU (sets state to enable)
func (m *TelnetSessionManager) UnblockONU(ctx context.Context, req *model.ONUBlockRequest) error {
	log.Info().
		Str("pon_port", req.PONPort).
		Int("onu_id", req.ONUID).
		Bool("block", req.Block).
		Msg("Unblocking ONU via Telnet")

	// Enter PON interface configuration mode
	interfaceCmd := fmt.Sprintf("interface gpon-olt_%s", req.PONPort)
	if _, err := m.ExecuteCommand(ctx, interfaceCmd); err != nil {
		return fmt.Errorf("failed to enter interface mode: %w", err)
	}

	// Execute ONU state enable command
	stateCmd := fmt.Sprintf("onu %d state enable", req.ONUID)
	if _, err := m.ExecuteCommand(ctx, stateCmd); err != nil {
		// Check if error is due to ONU not existing
		if strings.Contains(err.Error(), "not exist") || strings.Contains(err.Error(), "invalid") {
			return fmt.Errorf("ONU %s/%d not found", req.PONPort, req.ONUID)
		}
		return fmt.Errorf("failed to enable ONU: %w", err)
	}

	// Exit interface mode
	if _, err := m.ExecuteCommand(ctx, "exit"); err != nil {
		log.Warn().Err(err).Msg("Failed to exit interface mode, session will auto-recover")
	}

	log.Info().
		Str("pon_port", req.PONPort).
		Int("onu_id", req.ONUID).
		Msg("ONU unblocked successfully")

	return nil
}

// UpdateDescription updates the name/description of an ONU
func (m *TelnetSessionManager) UpdateDescription(ctx context.Context, req *model.ONUDescriptionRequest) error {
	log.Info().
		Str("pon_port", req.PONPort).
		Int("onu_id", req.ONUID).
		Str("description", req.Description).
		Msg("Updating ONU description via Telnet")

	// Enter PON interface configuration mode
	interfaceCmd := fmt.Sprintf("interface gpon-olt_%s", req.PONPort)
	if _, err := m.ExecuteCommand(ctx, interfaceCmd); err != nil {
		return fmt.Errorf("failed to enter interface mode: %w", err)
	}

	// Execute ONU name command
	// Escape quotes in description
	description := strings.ReplaceAll(req.Description, `"`, `\"`)
	nameCmd := fmt.Sprintf(`onu %d name "%s"`, req.ONUID, description)
	if _, err := m.ExecuteCommand(ctx, nameCmd); err != nil {
		// Check if error is due to ONU not existing
		if strings.Contains(err.Error(), "not exist") || strings.Contains(err.Error(), "invalid") {
			return fmt.Errorf("ONU %s/%d not found", req.PONPort, req.ONUID)
		}
		return fmt.Errorf("failed to update ONU description: %w", err)
	}

	// Exit interface mode
	if _, err := m.ExecuteCommand(ctx, "exit"); err != nil {
		log.Warn().Err(err).Msg("Failed to exit interface mode, session will auto-recover")
	}

	log.Info().
		Str("pon_port", req.PONPort).
		Int("onu_id", req.ONUID).
		Str("description", req.Description).
		Msg("ONU description updated successfully")

	return nil
}

// DeleteONU removes ONU configuration from the OLT
func (m *TelnetSessionManager) DeleteONU(ctx context.Context, req *model.ONUDeleteRequest) error {
	log.Info().
		Str("pon_port", req.PONPort).
		Int("onu_id", req.ONUID).
		Msg("Deleting ONU configuration via Telnet")

	// Enter PON interface configuration mode
	interfaceCmd := fmt.Sprintf("interface gpon-olt_%s", req.PONPort)
	if _, err := m.ExecuteCommand(ctx, interfaceCmd); err != nil {
		return fmt.Errorf("failed to enter interface mode: %w", err)
	}

	// Execute no onu command to delete ONU
	deleteCmd := fmt.Sprintf("no onu %d", req.ONUID)
	if _, err := m.ExecuteCommand(ctx, deleteCmd); err != nil {
		// Check if error is due to ONU not existing
		if strings.Contains(err.Error(), "not exist") || strings.Contains(err.Error(), "invalid") {
			return fmt.Errorf("ONU %s/%d not found", req.PONPort, req.ONUID)
		}
		return fmt.Errorf("failed to delete ONU: %w", err)
	}

	// Exit interface mode
	if _, err := m.ExecuteCommand(ctx, "exit"); err != nil {
		log.Warn().Err(err).Msg("Failed to exit interface mode, session will auto-recover")
	}

	log.Info().
		Str("pon_port", req.PONPort).
		Int("onu_id", req.ONUID).
		Msg("ONU deleted successfully")

	return nil
}
