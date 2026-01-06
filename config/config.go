package config

import (
	"errors"
	"path/filepath"
	"rc_mavlink/tools"

	"github.com/spf13/viper"
)

const EnvPrefix = "control"

const (
	MavLinkAddr = "mavlink_addr"
	MavLinkPort = "mavlink_port"
	SerialPort  = "serial_port"
	SerialBaud  = "serial_baud"
)

type Config struct {
	MavLinkAddr string `json:"mavlink_addr,omitempty"`
	MavlinkPort string `json:"mavlink_port,omitempty"`
	SerialPort  string `json:"serial_port,omitempty"`
	SerialBaud  int    `json:"serial_baud,omitempty"`
}

// Valores de conf default
func init() {
	viper.SetDefault(MavLinkAddr, "192.168.2.15")
	viper.SetDefault(MavLinkPort, 8090)
	viper.SetDefault(SerialPort, "/dev/ttyUSB0")
	viper.SetDefault(SerialBaud, 460800)
}

func Load() (*Config, error) {

	viper.SetEnvPrefix(EnvPrefix)
	viper.AutomaticEnv()

	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	path, err := tools.GetCurrentFilePath()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(path, "/config")

	viper.AddConfigPath(configPath)
	err = viper.ReadInConfig()
	if err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return nil, err
		}
	}

	cfg := new(Config)

	cfg.MavLinkAddr = viper.GetString(MavLinkAddr)
	cfg.MavlinkPort = viper.GetString(MavLinkPort)
	cfg.SerialPort = viper.GetString(SerialPort)
	cfg.SerialBaud = viper.GetInt(SerialBaud)

	return cfg, nil
}
