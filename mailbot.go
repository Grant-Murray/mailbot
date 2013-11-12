// Package mailbot supplies functions for a server to send emails
package mailbot

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/smtp"
)

type ServerConfig struct {
	Host      string
	Port      int
	User      string
	Password  string
	EmailFrom string
}

func getRootCAs() (rootCAs *x509.CertPool) {
	rootCAs = x509.NewCertPool()
	data, err := ioutil.ReadFile("/etc/ssl/certs/ca-certificates.crt")
	if err == nil {
		rootCAs.AppendCertsFromPEM(data)
	} else {
		panic(fmt.Sprintf("Error while reading certs file: %s", err))
	}
	return rootCAs
}

func verifyCert(vo x509.VerifyOptions) {
	data, err := ioutil.ReadFile("/tmp/test.mailbot.net.pem")
	if err != nil {
		panic(fmt.Sprintf("Error while reading cert file: %s", err))
	}
	asn1, _ := pem.Decode(data)
	cert, err := x509.ParseCertificate(asn1.Bytes)
	if err != nil {
		panic(fmt.Sprintf("Error parsing cert data: %s\n", err))
	}

	chains, err := cert.Verify(vo)
	if err != nil {
		panic(fmt.Sprintf("Error verifying: %s", err))
	}

	if len(chains) > 0 {
		fmt.Printf("Woohoo verified!\n")
	} else {
		panic("Sadface :( did not verify")
	}

}

// TLSSend makes the initial connection as TLS and starts the smtp conversation, this should be used with port 465. smtp.SendMail connects and sends HELO in the clear first and then upgrades to TLS if allowed or errors out.
func TLSSend(sc ServerConfig, emailFromOverride string, emailTo []string, body []byte) (err error) {

	addrport := fmt.Sprintf("%s:%d", sc.Host, sc.Port)

	var vo x509.VerifyOptions
	//vo.Roots = getRootCAs()

	verifyCert(vo)

	tlsConf := new(tls.Config)
	tlsConf.RootCAs = vo.Roots

	conn, err := tls.Dial("tcp", addrport, tlsConf)
	if err != nil {
		fmt.Printf("Dial error: %s", err)
		return err
	}

	smtpClient, err := smtp.NewClient(conn, "")
	if err != nil {
		fmt.Printf("NewClient err: %s", err)
		return err
	}

	from := emailFromOverride
	if from == "" {
		from = sc.EmailFrom
	}

	err = smtpClient.Hello("localhost")
	if err != nil {
		fmt.Printf("smtp Hello err: %s", err)
		return err
	}

	//auth := smtp.PlainAuth("", sc.User, sc.Password, sc.Host)

	// bug(Grant) finish the dialog
	panic("implementation of TLSSend not complete")

	return nil
}

// UpgradeTLSSend sends an email from ServerConfig.EmailFrom unless emailFromOverride is not nil and that is used instead. The body is sent to addresses in emailTo
func UpgradeTLSSend(sc ServerConfig, emailFromOverride string, emailTo []string, body []byte) (err error) {

	addrport := fmt.Sprintf("%s:%d", sc.Host, sc.Port)
	auth := smtp.PlainAuth("", sc.User, sc.Password, sc.Host)

	from := emailFromOverride
	if from == "" {
		from = sc.EmailFrom
	}

	err = smtp.SendMail(addrport, auth, from, emailTo, body)
	if err != nil {
		return err
	}

	return nil
}
