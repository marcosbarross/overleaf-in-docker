package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	backupDir := "./backups"
	retentionDays := 7
	timestamp := time.Now().Format("20060102_150405")

	// 1. Cria o diretório de backup se não existir
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		log.Fatalf("Erro ao criar diretório de backup: %v", err)
	}

	fmt.Printf("[%s] Iniciando backup do Overleaf...\n", timestamp)

	// 2. Backup do MongoDB
	mongoFile := filepath.Join(backupDir, fmt.Sprintf("mongo_%s.gz", timestamp))
	fmt.Println("-> Fazendo dump do MongoDB...")
	if err := runMongoDump(mongoFile); err != nil {
		log.Printf("Erro no backup do MongoDB: %v\n", err)
	}

	// 3. Backup do Volume de Dados
	fmt.Println("-> Compactando arquivos do volume (uploads e imagens)...")
	tarFile := fmt.Sprintf("overleaf_data_%s.tar.gz", timestamp)
	if err := runVolumeBackup(backupDir, tarFile); err != nil {
		log.Printf("Erro no backup do volume: %v\n", err)
	}

	// 4. Limpeza de Backups Antigos
	fmt.Printf("-> Removendo backups mais antigos que %d dias...\n", retentionDays)
	if err := cleanOldBackups(backupDir, retentionDays); err != nil {
		log.Printf("Erro na limpeza de backups antigos: %v\n", err)
	}

	fmt.Printf("[%s] Backup concluído com sucesso!\n", time.Now().Format("20060102_150405"))
}

// Executa o dump do mongo e joga a saída direto para o arquivo
func runMongoDump(outFile string) error {
	cmd := exec.Command("docker", "compose", "exec", "-T", "mongo", "mongodump", "--archive", "--gzip")
	
	outFileHandle, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer outFileHandle.Close()

	cmd.Stdout = outFileHandle // Redireciona o stdout para o arquivo .gz
	cmd.Stderr = os.Stderr     // Mantém os erros no terminal

	return cmd.Run()
}

// Sobe um container alpine temporário para ler o volume do Overleaf e criar o tar.gz
func runVolumeBackup(backupDir, tarFile string) error {
	// Pega o ID do container do sharelatex dinamicamente
	out, err := exec.Command("docker", "compose", "ps", "-q", "sharelatex").Output()
	if err != nil {
		return fmt.Errorf("falha ao obter ID do container sharelatex: %v", err)
	}
	
	containerID := strings.TrimSpace(string(out))
	if containerID == "" {
		return fmt.Errorf("container sharelatex não está rodando")
	}

	// Pega o caminho absoluto do diretório de backup para montar no container
	absBackupDir, err := filepath.Abs(backupDir)
	if err != nil {
		return err
	}

	dockerArgs := []string{
		"run", "--rm",
		"--volumes-from", containerID,
		"-v", fmt.Sprintf("%s:/backup", absBackupDir),
		"alpine", "tar", "czf", fmt.Sprintf("/backup/%s", tarFile), "/var/lib/overleaf",
	}

	cmd := exec.Command("docker", dockerArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Varre a pasta e deleta arquivos que passaram do tempo de retenção
func cleanOldBackups(dir string, days int) error {
	cutoff := time.Now().AddDate(0, 0, -days)

	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Verifica se é um arquivo, se é mais antigo que o cutoff e se termina em .gz
		if !info.IsDir() && info.ModTime().Before(cutoff) && strings.HasSuffix(info.Name(), ".gz") {
			fmt.Printf("   Removendo backup antigo: %s\n", info.Name())
			return os.Remove(path)
		}
		return nil
	})
}