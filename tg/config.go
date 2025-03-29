package TG

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

type ConfigManager struct {
	configs map[string]string
	lock    sync.RWMutex
	changed bool // Track whether any changes have been made
}

// Default configurations
var defaultConfig = map[string]string{
	"pluginmanager": "default",
}

func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		configs: make(map[string]string),
	}
}

func (cm *ConfigManager) Load() error {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	// Check if config file exists
	if _, err := os.Stat("config"); os.IsNotExist(err) {
		// If the config file does not exist, create it and write defaults
		if err := cm.createDefaultConfig(); err != nil {
			return err
		}
	}
	file, err := os.Open("config")
	if err != nil {
		return fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Split into key and value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			cm.configs[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	return nil
}

func (cm *ConfigManager) Save() error {
	if !cm.changed {
		return nil // No changes to save
	}

	cm.lock.Lock()
	defer cm.lock.Unlock()

	file, err := os.Create("config")
	if err != nil {
		return fmt.Errorf("failed to create config file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for key, value := range cm.configs {
		_, err := writer.WriteString(fmt.Sprintf("%s=%s\n", key, value))
		if err != nil {
			return fmt.Errorf("error writing config file: %v", err)
		}
	}

	cm.changed = false // Reset the flag after saving
	return writer.Flush()
}

func (cm *ConfigManager) Get(key string) (string, bool) {
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	value, exists := cm.configs[key]
	return value, exists
}

func (cm *ConfigManager) Set(key string, value string) {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	cm.configs[key] = value
	cm.changed = true // Mark as changed
}

func (cm *ConfigManager) createDefaultConfig() error {
	file, err := os.Create("config")
	if err != nil {
		return fmt.Errorf("failed to create config file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Write default configurations
	for key, value := range defaultConfig {
		_, err := writer.WriteString(fmt.Sprintf("%s=%s\n", key, value))
		if err != nil {
			return fmt.Errorf("error writing to config file: %v", err)
		}
	}

	// Flush to ensure data is written
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("error flushing config file: %v", err)
	}

	return nil
}
