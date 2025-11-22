package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/matutetandil/autonode/internal/core"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// ConfigCommand implements the config command for setting local Node.js version and npm profile
// Single Responsibility Principle: Only responsible for managing .autonode.yml configuration
type ConfigCommand struct {
	nodeVersion string
	npmProfile  string
	show        bool
	remove      bool
}

// autonodeConfig represents the structure of .autonode.yml file
type autonodeConfig struct {
	NodeVersion string `yaml:"nodeVersion,omitempty"`
	NpmProfile  string `yaml:"npmProfile,omitempty"`
}

// init registers this command automatically when the package is imported
func init() {
	Register(&ConfigCommand{})
}

// GetCobraCommand returns the cobra command for this command
func (c *ConfigCommand) GetCobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configure local Node.js version and npm profile for current directory",
		Long: `Configure a specific Node.js version and/or npm profile for the current directory.

This creates or updates a .autonode.yml file in the current directory.
The .autonode.yml configuration has the highest priority (before .nvmrc).

Examples:
  autonode config --node 20           # Set Node.js version to 20
  autonode config --node 18.17.0      # Set specific version
  autonode config --profile work      # Set npm profile
  autonode config --node 20 --profile work  # Set both
  autonode config --show              # Show current configuration
  autonode config --remove            # Remove .autonode.yml file
  autonode config --node ""           # Remove only nodeVersion
  autonode config --profile ""        # Remove only npmProfile`,
		RunE: c.run,
	}

	cmd.Flags().StringVarP(&c.nodeVersion, "node", "n", "", "Node.js version to use (empty string to remove)")
	cmd.Flags().StringVarP(&c.npmProfile, "profile", "p", "", "npm profile to use (empty string to remove)")
	cmd.Flags().BoolVarP(&c.show, "show", "s", false, "Show current configuration")
	cmd.Flags().BoolVarP(&c.remove, "remove", "r", false, "Remove .autonode.yml configuration file")

	return cmd
}

// run executes the config command
func (c *ConfigCommand) run(cmd *cobra.Command, args []string) error {
	logger := core.NewConsoleLogger()

	// Get current working directory
	projectPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	configPath := filepath.Join(projectPath, ".autonode.yml")

	// Handle --show flag
	if c.show {
		return c.showConfig(configPath, logger)
	}

	// Handle --remove flag
	if c.remove {
		return c.removeConfig(configPath, logger)
	}

	// Check if any configuration flag was provided
	nodeChanged := cmd.Flags().Changed("node")
	profileChanged := cmd.Flags().Changed("profile")

	if !nodeChanged && !profileChanged {
		// No flags provided, show help
		return cmd.Help()
	}

	// Load existing config or create new one
	config, err := c.loadConfig(configPath)
	if err != nil {
		return err
	}

	// Update configuration based on flags
	if nodeChanged {
		if c.nodeVersion == "" {
			config.NodeVersion = ""
			logger.Info("Removed nodeVersion from configuration")
		} else {
			config.NodeVersion = c.nodeVersion
			logger.Success(fmt.Sprintf("Set nodeVersion to '%s'", c.nodeVersion))
		}
	}

	if profileChanged {
		if c.npmProfile == "" {
			config.NpmProfile = ""
			logger.Info("Removed npmProfile from configuration")
		} else {
			config.NpmProfile = c.npmProfile
			logger.Success(fmt.Sprintf("Set npmProfile to '%s'", c.npmProfile))
		}
	}

	// If both fields are empty, remove the file
	if config.NodeVersion == "" && config.NpmProfile == "" {
		if _, err := os.Stat(configPath); err == nil {
			if err := os.Remove(configPath); err != nil {
				return fmt.Errorf("failed to remove config file: %w", err)
			}
			logger.Info("Removed .autonode.yml (no configuration left)")
		}
		return nil
	}

	// Save configuration
	return c.saveConfig(configPath, config, logger)
}

// showConfig displays the current configuration
func (c *ConfigCommand) showConfig(configPath string, logger core.Logger) error {
	config, err := c.loadConfig(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Info("No local configuration found.")
			logger.Info("Use --node <version> or --profile <name> to configure.")
			return nil
		}
		return err
	}

	// Check if config file exists but is empty
	if config.NodeVersion == "" && config.NpmProfile == "" {
		logger.Info("No local configuration found.")
		logger.Info("Use --node <version> or --profile <name> to configure.")
		return nil
	}

	logger.Info("Current configuration (.autonode.yml):")
	if config.NodeVersion != "" {
		logger.Info(fmt.Sprintf("  nodeVersion: %s", config.NodeVersion))
	}
	if config.NpmProfile != "" {
		logger.Info(fmt.Sprintf("  npmProfile: %s", config.NpmProfile))
	}

	return nil
}

// removeConfig removes the .autonode.yml file
func (c *ConfigCommand) removeConfig(configPath string, logger core.Logger) error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		logger.Info("No .autonode.yml file found")
		return nil
	}

	if err := os.Remove(configPath); err != nil {
		return fmt.Errorf("failed to remove config file: %w", err)
	}

	logger.Success("Removed .autonode.yml configuration")
	return nil
}

// loadConfig loads the configuration from .autonode.yml file
func (c *ConfigCommand) loadConfig(configPath string) (*autonodeConfig, error) {
	config := &autonodeConfig{}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return config, nil // Return empty config if file doesn't exist
		}
		return nil, err
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse .autonode.yml: %w", err)
	}

	return config, nil
}

// saveConfig saves the configuration to .autonode.yml file
func (c *ConfigCommand) saveConfig(configPath string, config *autonodeConfig, logger core.Logger) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to serialize configuration: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write .autonode.yml: %w", err)
	}

	logger.Info(fmt.Sprintf("Configuration saved to %s", configPath))
	return nil
}
