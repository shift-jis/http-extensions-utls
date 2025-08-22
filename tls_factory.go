package httpx

import (
	"crypto/tls"
	"net"

	oohttp "github.com/ooni/oohttp"
	utls "github.com/refraction-networking/utls"
)

type SecureTLSConnFactory struct {
	TLSClientProfile *utls.ClientHelloID
}

func (factory *SecureTLSConnFactory) CreateTLSConnection(conn net.Conn, config *tls.Config) oohttp.TLSConn {
	if factory.TLSClientProfile == nil {
		factory.TLSClientProfile = &utls.HelloChrome_Auto
	}

	tlsConfig := &utls.Config{
		RootCAs:                     config.RootCAs,
		NextProtos:                  config.NextProtos,
		ServerName:                  config.ServerName,
		DynamicRecordSizingDisabled: config.DynamicRecordSizingDisabled,
		InsecureSkipVerify:          config.InsecureSkipVerify,
		ClientSessionCache:          utls.NewLRUClientSessionCache(0),
	}

	return &UTLSConnectionAdapter{
		UConn: utls.UClient(conn, tlsConfig, *factory.TLSClientProfile),
	}
}
