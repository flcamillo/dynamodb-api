package logs

import (
	"fmt"
	"time"
)

const (
	// Nível de log: Erro
	LevelError = 1
	// Nível de log: Aviso
	LevelWarning = 2
	// Nível de log: Informação
	LevelInfo = 3
	// Nível de log: Depuração
	LevelDebug = 4
)

// StdoutLog implementa a interface de log escrevendo no stdout.
type Stdout struct {
	// nível de log
	Level int
}

// Cria uma nova instância do log para stdout.
func NewStdoutLog() *Stdout {
	return &Stdout{
		Level: LevelInfo,
	}
}

// Escreve uma mensagem de informação no log.
func (p *Stdout) Info(format string, a ...any) {
	if p.Level >= LevelInfo {
		fmt.Printf("%s [INFO] %s\n", time.Now().Format(time.RFC3339), fmt.Sprintf(format, a...))
	}
}

// Escreve uma mensagem de erro no log.
func (p *Stdout) Error(format string, a ...any) {
	if p.Level >= LevelError {
		fmt.Printf("%s [ERROR] %s\n", time.Now().Format(time.RFC3339), fmt.Sprintf(format, a...))
	}
}

// Escreve uma mensagem de aviso no log.
func (p *Stdout) Warn(format string, a ...any) {
	if p.Level >= LevelWarning {
		fmt.Printf("%s [WARNING] %s\n", time.Now().Format(time.RFC3339), fmt.Sprintf(format, a...))
	}
}

// Escreve uma mensagem de depuração no log.
func (p *Stdout) Debug(format string, a ...any) {
	if p.Level >= LevelDebug {
		fmt.Printf("%s [DEBUG] %s\n", time.Now().Format(time.RFC3339), fmt.Sprintf(format, a...))
	}
}
