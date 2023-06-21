package server

import (
	"fmt"
	"net"
)

// Fechar o listener é responsabilidade de quem chama
// para garantir que a alocação de multiplas portas aleatórias não colidam
func reservaPorta(porta int) (*net.TCPListener, int, error) {
	endpoint := fmt.Sprintf("localhost:%d", porta)
	addr, err := net.ResolveTCPAddr("tcp", endpoint)
	if err != nil {
		return nil, 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, 0, err
	}
	return l, l.Addr().(*net.TCPAddr).Port, nil
}

func reservaPortaLivre() (*net.TCPListener, int, error) {
	return reservaPorta(0)
}

func ReservaPorta(maxTries, fallback int) (*net.TCPListener, int, error) {
	for i := 0; i < maxTries; i++ {
		l, port, err := reservaPortaLivre()

		if err == nil {
			return l, port, nil
		}
	}

	// ultima tentativa com a porta fallback
	return reservaPorta(fallback)
}
