package httpx

import (
	"context"
	"crypto/tls"

	utls "github.com/refraction-networking/utls"
)

type UTLSConnectionAdapter struct {
	*utls.UConn
}

func (adapter *UTLSConnectionAdapter) ConnectionState() tls.ConnectionState {
	state := adapter.UConn.ConnectionState()
	return tls.ConnectionState{
		Version:                     state.Version,
		HandshakeComplete:           state.HandshakeComplete,
		DidResume:                   state.DidResume,
		CipherSuite:                 state.CipherSuite,
		NegotiatedProtocol:          state.NegotiatedProtocol,
		NegotiatedProtocolIsMutual:  state.NegotiatedProtocolIsMutual,
		ServerName:                  state.ServerName,
		PeerCertificates:            state.PeerCertificates,
		VerifiedChains:              state.VerifiedChains,
		SignedCertificateTimestamps: state.SignedCertificateTimestamps,
		OCSPResponse:                state.OCSPResponse,
		TLSUnique:                   state.TLSUnique,
	}
}

func (adapter *UTLSConnectionAdapter) HandshakeContext(ctx context.Context) error {
	handshakeErrCh := make(chan error, 1)
	go func() {
		handshakeErrCh <- adapter.UConn.Handshake()
	}()

	select {
	case err := <-handshakeErrCh:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
