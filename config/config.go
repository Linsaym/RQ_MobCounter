package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	LogPath    string `json:"log_path"`
	FilePrefix string `json:"file_prefix"`
}

const DefaultLogPath = `D:\B.A.S.E\Games\Royal Quest\chatlogs`
const DefaultFilePrefix = "exp"

func Load() (*Config, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("ошибка определения пути приложения: %w", err)
	}

	exeDir := filepath.Dir(exePath)
	configPath := filepath.Join(exeDir, "config.json")

	cfg := &Config{
		LogPath:    DefaultLogPath,
		FilePrefix: DefaultFilePrefix,
	}

	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("ошибка чтения конфига: %w", err)
		}

		if err := json.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("ошибка парсинга конфига: %w", err)
		}
	} else {
		fmt.Printf("⚠️ Предупреждение: файл config.json не найден в папке приложения!\n\n")
		fmt.Printf("Создайте файл config.json со следующим содержимым:\n")
		fmt.Printf("{\n")
		fmt.Printf("  \"log_path\": \"D:\\\\B.A.S.E\\\\Games\\\\Royal Quest\\\\chatlogs\",\n")
		fmt.Printf("  \"file_prefix\": \"exp\"\n")
		fmt.Printf("}\n\n")
		fmt.Printf("Использую значения по умолчанию.\n")
		fmt.Printf("LogPath: %s\n", cfg.LogPath)
		fmt.Printf("FilePrefix: %s\n\n", cfg.FilePrefix)
	}

	cfg.LogPath = filepath.FromSlash(cfg.LogPath)

	return cfg, nil
}

func (c *Config) Save() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("ошибка определения пути приложения: %w", err)
	}

	exeDir := filepath.Dir(exePath)
	configPath := filepath.Join(exeDir, "config.json")

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("ошибка сериализации конфига: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("ошибка сохранения конфига: %w", err)
	}

	return nil
}
