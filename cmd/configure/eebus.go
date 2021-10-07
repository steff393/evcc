package configure

import (
	"crypto/x509/pkix"
	"fmt"

	certhelper "github.com/evcc-io/eebus/cert"
	"github.com/evcc-io/eebus/communication"
	"github.com/evcc-io/evcc/server"
)

// setup EEBus
func (c *CmdConfigure) configureEEBus(conf map[string]interface{}) error {
	var err error
	if server.EEBusInstance, err = server.NewEEBus(conf); err == nil {
		go server.EEBusInstance.Run()
	}

	return nil
}

// setup EEBUS certificate
// returns privagte key, public key and error
func (c *CmdConfigure) eebusCertificate() (map[string]interface{}, error) {
	details := communication.ManufacturerDetails{
		DeviceName:    "EVCC",
		DeviceCode:    "EVCC_HEMS_01",
		DeviceAddress: "EVCC_HEMS",
		BrandName:     "EVCC",
	}

	subject := pkix.Name{
		CommonName:   details.DeviceCode,
		Country:      []string{"DE"},
		Organization: []string{details.BrandName},
	}

	var eebusConfig map[string]interface{}

	cert, err := certhelper.CreateCertificate(true, subject)
	if err != nil {
		return eebusConfig, fmt.Errorf("could not create certificate")
	}

	pubKey, privKey, err := certhelper.GetX509KeyPair(cert)
	if err != nil {
		return eebusConfig, fmt.Errorf("could not process generated certificate")
	}

	eebusConfig["certificate"] = map[string]interface{}{
		"public":  pubKey,
		"private": privKey,
	}

	return eebusConfig, nil
}
