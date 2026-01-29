package main

import (
	"api/interfaces"
	"encoding/json"
	"os"
)

// Config representa a configuração da aplicação.
type Config struct {
	// arquivo de configuração
	File string `json:"-"`
	// cliente do DynamoDB
	DynamoDBClient interfaces.DynamoDBClient `json:"-"`
	// repositório de dados
	Repository interfaces.Repository `json:"-"`
	// log de aplicação
	Log interfaces.Log `json:"-"`
	// endereço para ativar o servidor
	Address string `json:"address"`
	// porta do servidor
	Port int `json:"port"`
	// tempo de expiração dos registros em minutos
	RecordTTLMinutes int64 `json:"record_ttl_minutes"`
}

// Cria uma instância da configuração da aplicação com valores padrão.
func NewConfig(file string) *Config {
	return &Config{
		File:             file,
		Address:          "0.0.0.0",
		Port:             7000,
		RecordTTLMinutes: 24 * 60,
	}
}

// Salva as configurações em um arquivo.
func (p *Config) Save() error {
	f, err := os.OpenFile(p.File, os.O_CREATE|os.O_WRONLY, 0750)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", " ")
	err = enc.Encode(p)
	if err != nil {
		return err
	}
	return nil
}

// Carrega as configurações de um arquivo ou retorna a configuração padrão se o arquivo não existir.
func LoadConfig(file string) (config *Config, err error) {
	f, err := os.OpenFile(file, os.O_RDONLY, 0750)
	if err != nil {
		if os.IsNotExist(err) {
			cfg := NewConfig(file)
			return cfg, cfg.Save()
		}
		return nil, err
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
