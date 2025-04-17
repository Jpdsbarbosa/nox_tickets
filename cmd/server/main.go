package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"nox_tickets/internal/interfaces/http"
)

func main() {
	// cria servidor na porta 8080
	port := "8080" // Inicialmente, para testes - TODO: pegar da variavel de ambiente
	server := http.NewServer(port)

	// Canal par capturar sinais de interrupção (Ctrl+C)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Inicia o servidor em uma goroutine
	go func() {
		log.Printf("Servidor Iniciado na porta: %s\n", port)
		if err := server.Start(); err != nil {
			log.Fatalf("Erro ao iniciar o servidor: %v", err)
		}
	}()

	// Aguarda o sinal de interrupção
	<-stop

	// Desligar servidor
	log.Printf("Desligando servidor...")
	ctx := context.Background()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Erro ao desligar servidor", err)
	}

	log.Printf("Servidor desligado com sucesso")
}
