package repository

import (
	"context"
	"fmt"
	"net"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/s4lfanet/go-api-c320/config"
	"github.com/s4lfanet/go-api-c320/internal/model"
	"github.com/ziutek/telnet"
)

// TelnetRepository defines the interface for Telnet operations
type TelnetRepository interface {
	// Connection management
	Connect() error
	Close() error
	IsConnected() bool
	Reconnect() error

	// Command execution
	Execute(ctx context.Context, command string) (*model.TelnetResponse, error)
	ExecuteMulti(ctx context.Context, commands []string) (*model.TelnetBatchResponse, error)
	ExecuteWithExpect(ctx context.Context, command, expectPattern string) (*model.TelnetResponse, error)

	// Mode management
	EnterEnableMode() error
	EnterConfigMode() error
	ExitConfigMode() error
	GetCurrentMode() string

	// Helper methods
	SaveConfig() error
	ShowRunningConfig() (string, error)
	GetConnectionInfo() *model.TelnetConnectionInfo
}

// telnetRepository implements TelnetRepository interface
type telnetRepository struct {
	config       *config.TelnetConfig
	conn         *telnet.Conn
	currentMode  string
	connected    bool
	connectedAt  time.Time
	lastActivity time.Time
	mu           sync.Mutex
	commandMu    sync.Mutex // Separate mutex for command execution
}

// NewTelnetRepository creates a new telnet repository instance
func NewTelnetRepository(cfg *config.TelnetConfig) TelnetRepository {
	return &telnetRepository{
		config:      cfg,
		currentMode: "disconnected",
		connected:   false,
	}
}

// Connect establishes a telnet connection and authenticates
func (r *telnetRepository) Connect() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.connected {
		log.Info().Msg("Telnet already connected")
		return nil
	}

	// Establish connection
	address := fmt.Sprintf("%s:%d", r.config.Host, r.config.Port)
	log.Info().Str("address", address).Msg("Connecting to OLT via Telnet")

	conn, err := telnet.DialTimeout("tcp", address, r.config.ConnectTimeout)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to OLT")
		return model.NewTelnetError(model.ErrCodeConnectionFailed,
			fmt.Sprintf("failed to connect: %v", err), true)
	}

	r.conn = conn
	r.connectedAt = time.Now()
	r.lastActivity = time.Now()

	// Set timeouts
	_ = r.conn.SetReadDeadline(time.Now().Add(r.config.ReadTimeout))
	_ = r.conn.SetWriteDeadline(time.Now().Add(r.config.WriteTimeout))

	// Login sequence
	if err := r.login(); err != nil {
		r.conn.Close()
		r.conn = nil
		return err
	}

	r.connected = true
	r.currentMode = "user"

	log.Info().Str("mode", r.currentMode).Msg("Telnet connection established")
	return nil
}

// login performs the login sequence
func (r *telnetRepository) login() error {
	// Wait for username prompt
	log.Debug().Msg("Waiting for username prompt")
	if err := r.expectString("Username:", r.config.Timeout); err != nil {
		return model.NewTelnetError(model.ErrCodeAuthFailed,
			"username prompt not received", false)
	}

	// Send username
	log.Debug().Str("username", r.config.Username).Msg("Sending username")
	if err := r.sendCommand(r.config.Username); err != nil {
		return err
	}

	// Wait for password prompt
	log.Debug().Msg("Waiting for password prompt")
	if err := r.expectString("Password:", r.config.Timeout); err != nil {
		return model.NewTelnetError(model.ErrCodeAuthFailed,
			"password prompt not received", false)
	}

	// Send password
	log.Debug().Msg("Sending password")
	if err := r.sendCommand(r.config.Password); err != nil {
		return err
	}

	// Wait for user prompt
	log.Debug().Str("prompt", r.config.PromptUser).Msg("Waiting for user prompt")
	if err := r.expectString(r.config.PromptUser, r.config.Timeout); err != nil {
		return model.NewTelnetError(model.ErrCodeAuthFailed,
			"login failed - user prompt not received", false)
	}

	log.Info().Msg("Login successful")
	return nil
}

// Close closes the telnet connection
func (r *telnetRepository) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.connected || r.conn == nil {
		return nil
	}

	// Try to exit gracefully
	if r.currentMode == "config" {
		_ = r.sendCommand("end")
		_ = r.expectString(r.config.PromptEnable, 2*time.Second)
	}

	_ = r.sendCommand("exit")

	// Close connection
	err := r.conn.Close()
	r.conn = nil
	r.connected = false
	r.currentMode = "disconnected"

	log.Info().Msg("Telnet connection closed")
	return err
}

// IsConnected checks if the connection is active
func (r *telnetRepository) IsConnected() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.connected && r.conn != nil
}

// Reconnect closes and reopens the connection
func (r *telnetRepository) Reconnect() error {
	log.Info().Msg("Reconnecting to OLT")
	r.Close()
	return r.Connect()
}

// EnterEnableMode enters enable (privileged) mode
func (r *telnetRepository) EnterEnableMode() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.currentMode == "enable" || r.currentMode == "config" {
		return nil
	}

	if r.currentMode != "user" {
		return fmt.Errorf("cannot enter enable mode from %s mode", r.currentMode)
	}

	log.Debug().Msg("Entering enable mode")

	// Send enable command
	if err := r.sendCommand("enable"); err != nil {
		return err
	}

	// Check if password prompt appears
	output, _ := r.readUntilPrompt(r.config.Timeout)
	if strings.Contains(output, "Password:") {
		// Send enable password
		if err := r.sendCommand(r.config.EnablePassword); err != nil {
			return err
		}
		// Wait for enable prompt
		if err := r.expectString(r.config.PromptEnable, r.config.Timeout); err != nil {
			return model.NewTelnetError(model.ErrCodeAuthFailed,
				"failed to enter enable mode", false)
		}
	}

	r.currentMode = "enable"
	log.Info().Msg("Entered enable mode")
	return nil
}

// EnterConfigMode enters configuration mode
func (r *telnetRepository) EnterConfigMode() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.currentMode == "config" {
		return nil
	}

	// Must be in enable mode first
	if r.currentMode != "enable" {
		r.mu.Unlock()
		if err := r.EnterEnableMode(); err != nil {
			return err
		}
		r.mu.Lock()
	}

	log.Debug().Msg("Entering configuration mode")

	// Send configure terminal command
	if err := r.sendCommand("configure terminal"); err != nil {
		return err
	}

	// Wait for config prompt
	if err := r.expectString(r.config.PromptConfig, r.config.Timeout); err != nil {
		return model.NewTelnetError(model.ErrCodeCommandFailed,
			"failed to enter config mode", true)
	}

	r.currentMode = "config"
	log.Info().Msg("Entered configuration mode")
	return nil
}

// ExitConfigMode exits configuration mode back to enable mode
func (r *telnetRepository) ExitConfigMode() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.currentMode != "config" {
		return nil
	}

	log.Debug().Msg("Exiting configuration mode")

	// Send end command
	if err := r.sendCommand("end"); err != nil {
		return err
	}

	// Wait for enable prompt
	if err := r.expectString(r.config.PromptEnable, r.config.Timeout); err != nil {
		return model.NewTelnetError(model.ErrCodeCommandFailed,
			"failed to exit config mode", true)
	}

	r.currentMode = "enable"
	log.Info().Msg("Exited configuration mode")
	return nil
}

// GetCurrentMode returns the current mode
func (r *telnetRepository) GetCurrentMode() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.currentMode
}

// Execute executes a single command
func (r *telnetRepository) Execute(ctx context.Context, command string) (*model.TelnetResponse, error) {
	r.commandMu.Lock()
	defer r.commandMu.Unlock()

	if !r.IsConnected() {
		return nil, model.NewTelnetError(model.ErrCodeDisconnected,
			"not connected to OLT", true)
	}

	startTime := time.Now()
	log.Debug().Str("command", command).Msg("Executing command")

	// Send command
	if err := r.sendCommand(command); err != nil {
		return &model.TelnetResponse{
			Command:   command,
			Success:   false,
			Error:     err.Error(),
			Timestamp: time.Now().Format(time.RFC3339),
		}, err
	}

	// Read response
	output, err := r.readUntilPrompt(r.config.Timeout)
	if err != nil {
		return &model.TelnetResponse{
			Command:   command,
			Output:    output,
			Success:   false,
			Error:     err.Error(),
			Timestamp: time.Now().Format(time.RFC3339),
		}, err
	}

	// Clean output
	output = r.cleanOutput(output, command)

	duration := time.Since(startTime)
	log.Debug().
		Str("command", command).
		Dur("duration", duration).
		Msg("Command executed successfully")

	return &model.TelnetResponse{
		Command:   command,
		Output:    output,
		Success:   true,
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}

// ExecuteMulti executes multiple commands in sequence
func (r *telnetRepository) ExecuteMulti(ctx context.Context, commands []string) (*model.TelnetBatchResponse, error) {
	startTime := time.Now()
	responses := make([]model.TelnetResponse, 0, len(commands))
	allSuccess := true

	for _, cmd := range commands {
		resp, err := r.Execute(ctx, cmd)
		if err != nil {
			allSuccess = false
		}
		responses = append(responses, *resp)

		// Check context cancellation
		select {
		case <-ctx.Done():
			return &model.TelnetBatchResponse{
				Responses: responses,
				Success:   false,
				TotalTime: time.Since(startTime).String(),
			}, ctx.Err()
		default:
		}
	}

	return &model.TelnetBatchResponse{
		Responses: responses,
		Success:   allSuccess,
		TotalTime: time.Since(startTime).String(),
	}, nil
}

// ExecuteWithExpect executes a command and waits for a specific pattern
func (r *telnetRepository) ExecuteWithExpect(ctx context.Context, command, expectPattern string) (*model.TelnetResponse, error) {
	r.commandMu.Lock()
	defer r.commandMu.Unlock()

	if !r.IsConnected() {
		return nil, model.NewTelnetError(model.ErrCodeDisconnected,
			"not connected to OLT", true)
	}

	log.Debug().
		Str("command", command).
		Str("expect", expectPattern).
		Msg("Executing command with expect")

	// Send command
	if err := r.sendCommand(command); err != nil {
		return &model.TelnetResponse{
			Command:   command,
			Success:   false,
			Error:     err.Error(),
			Timestamp: time.Now().Format(time.RFC3339),
		}, err
	}

	// Read until expected pattern
	if err := r.expectString(expectPattern, r.config.Timeout); err != nil {
		return &model.TelnetResponse{
			Command:   command,
			Success:   false,
			Error:     fmt.Sprintf("expected pattern '%s' not found", expectPattern),
			Timestamp: time.Now().Format(time.RFC3339),
		}, err
	}

	return &model.TelnetResponse{
		Command:   command,
		Success:   true,
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}

// SaveConfig saves the current configuration
func (r *telnetRepository) SaveConfig() error {
	log.Info().Msg("Saving configuration")

	// Exit to enable mode if in config mode
	if r.GetCurrentMode() == "config" {
		if err := r.ExitConfigMode(); err != nil {
			return err
		}
	}

	resp, err := r.Execute(context.Background(), "write")
	if err != nil {
		return model.NewTelnetError(model.ErrCodeConfigSaveFailed,
			fmt.Sprintf("failed to save config: %v", err), true)
	}

	if !resp.Success {
		return model.NewTelnetError(model.ErrCodeConfigSaveFailed,
			"config save command failed", true)
	}

	log.Info().Msg("Configuration saved successfully")
	return nil
}

// ShowRunningConfig retrieves the running configuration
func (r *telnetRepository) ShowRunningConfig() (string, error) {
	resp, err := r.Execute(context.Background(), "show running-config")
	if err != nil {
		return "", err
	}
	return resp.Output, nil
}

// GetConnectionInfo returns connection information
func (r *telnetRepository) GetConnectionInfo() *model.TelnetConnectionInfo {
	r.mu.Lock()
	defer r.mu.Unlock()

	uptime := "0s"
	if r.connected {
		uptime = time.Since(r.connectedAt).String()
	}

	return &model.TelnetConnectionInfo{
		Host:       r.config.Host,
		Port:       r.config.Port,
		Connected:  r.connected,
		Mode:       r.currentMode,
		LastActive: r.lastActivity,
		Uptime:     uptime,
	}
}

// Helper methods

// sendCommand sends a command with newline
func (r *telnetRepository) sendCommand(command string) error {
	r.lastActivity = time.Now()
	_ = r.conn.SetWriteDeadline(time.Now().Add(r.config.WriteTimeout))

	_, err := r.conn.Write([]byte(command + "\n"))
	if err != nil {
		log.Error().Err(err).Str("command", command).Msg("Failed to send command")
		return model.NewTelnetError(model.ErrCodeCommandFailed,
			fmt.Sprintf("failed to send command: %v", err), true)
	}
	return nil
}

// expectString waits for a specific string to appear
func (r *telnetRepository) expectString(pattern string, timeout time.Duration) error {
	r.lastActivity = time.Now()
	_ = r.conn.SetReadDeadline(time.Now().Add(timeout))

	data := make([]byte, 4096)
	buffer := ""

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		n, err := r.conn.Read(data)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				return model.NewTelnetError(model.ErrCodeTimeout,
					fmt.Sprintf("timeout waiting for '%s'", pattern), true)
			}
			return model.NewTelnetError(model.ErrCodeCommandFailed,
				fmt.Sprintf("read error: %v", err), true)
		}

		buffer += string(data[:n])
		if strings.Contains(buffer, pattern) {
			return nil
		}
	}

	return model.NewTelnetError(model.ErrCodeTimeout,
		fmt.Sprintf("timeout waiting for '%s'", pattern), true)
}

// readUntilPrompt reads until a prompt is detected
func (r *telnetRepository) readUntilPrompt(timeout time.Duration) (string, error) {
	r.lastActivity = time.Now()
	_ = r.conn.SetReadDeadline(time.Now().Add(timeout))

	data := make([]byte, 4096)
	buffer := ""

	prompts := []string{r.config.PromptUser, r.config.PromptEnable, r.config.PromptConfig}

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		n, err := r.conn.Read(data)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// Timeout - check if we have a prompt
				for _, prompt := range prompts {
					if strings.Contains(buffer, prompt) {
						return buffer, nil
					}
				}
				return buffer, model.NewTelnetError(model.ErrCodeTimeout,
					"timeout reading response", true)
			}
			return buffer, model.NewTelnetError(model.ErrCodeCommandFailed,
				fmt.Sprintf("read error: %v", err), true)
		}

		buffer += string(data[:n])

		// Check for prompts
		for _, prompt := range prompts {
			if strings.Contains(buffer, prompt) {
				return buffer, nil
			}
		}
	}

	return buffer, model.NewTelnetError(model.ErrCodeTimeout,
		"timeout waiting for prompt", true)
}

// cleanOutput removes command echo and prompt from output
func (r *telnetRepository) cleanOutput(output, command string) string {
	// Remove command echo
	output = strings.Replace(output, command, "", 1)

	// Remove prompts
	prompts := []string{r.config.PromptUser, r.config.PromptEnable, r.config.PromptConfig}
	for _, prompt := range prompts {
		output = strings.ReplaceAll(output, prompt, "")
	}

	// Remove ANSI escape codes
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	output = ansiRegex.ReplaceAllString(output, "")

	// Trim whitespace
	output = strings.TrimSpace(output)

	return output
}
