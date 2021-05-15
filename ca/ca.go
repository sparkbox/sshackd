package ca

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"golang.org/x/crypto/ssh"
)

// writePemToFile writes keys to a file
func writeKeyToFile(keyBytes []byte, saveFileTo string) {
	err := ioutil.WriteFile(saveFileTo, keyBytes, 0600)
	if err != nil {
		fmt.Println("error saving: ", err)
	}

	log.Printf("Key saved to: %s", saveFileTo)
}

func genPrivateKey() *ecdsa.PrivateKey {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Println("generate failed: ", err)
	}

	return privateKey
}

func privateString(key *ecdsa.PrivateKey) []byte {
	privDer, err := x509.MarshalECPrivateKey(key)

	if err != nil {
		fmt.Println("Error generating string: ", err)
	}

	privBlock := pem.Block{
		Type:    "EC PRIVATE KEY",
		Headers: nil,
		Bytes:   privDer,
	}
	privatePEM := pem.EncodeToMemory(&privBlock)
	return privatePEM
}

func SignCert() string {
	// user key pair
	userKey := genPrivateKey()

	// we sign this public key with the CA
	userKeyPub, e := ssh.NewPublicKey(&userKey.PublicKey)

	if e != nil {
		fmt.Println("user public key failed: ", e)
	}

	ca := genPrivateKey()

	// caPub goes to the server to accept things signed by the CA
	caPub, e := ssh.NewPublicKey(&ca.PublicKey)

	if e != nil {
		fmt.Println("ca public key failed: ", e)
	}

	caSigner, err := ssh.ParsePrivateKey(privateString(ca))

	if err != nil {
		fmt.Println("error creating CA Signer: ", err)
	}

	expireTime, _ := time.ParseDuration("24h")

	//we create this Cert struct using the user's public key
	//https://github.com/ejcx/sshcert/blob/1c64826f1a45d87777103946575701b0a062623a/sshcert.go#L82
	// SignCert is called to sign an ssh public key and produce an ssh certificate.
	certInstance := &ssh.Certificate{
		Key:             userKeyPub,
		Serial:          333,
		CertType:        ssh.UserCert,
		KeyId:           "fool",
		ValidAfter:      uint64(time.Now().Unix()),
		ValidBefore:     uint64(time.Now().Add(expireTime).Unix()),
		ValidPrincipals: []string{"root"},
		Permissions: ssh.Permissions{
			CriticalOptions: map[string]string{},
			Extensions: map[string]string{
				"permit-X11-forwarding":   "",
				"permit-agent-forwarding": "",
				"permit-port-forwarding":  "",
				"permit-pty":              "",
				"permit-user-rc":          "",
			},
		},
	}

	//certInstance is now signed!
	certErr := certInstance.SignCert(rand.Reader, caSigner)

	if certErr != nil {
		fmt.Println("Error signing Certificate: ", certErr)
	}

	//need this + private key to login
	caTxt := privateString(ca)
	//need this + private key to login
	cert := ssh.MarshalAuthorizedKey(certInstance)
	//need this + cert to login
	userPrivate := privateString(userKey)
	//this is just for testing
	userPublic := ssh.MarshalAuthorizedKey(userKeyPub)
	//need this to be on any server hostj
	caPubText := ssh.MarshalAuthorizedKey(caPub)

	writeKeyToFile(cert, "./cert.pub")
	writeKeyToFile(caTxt, "./ca.private")
	writeKeyToFile(userPrivate, "./user.private")
	writeKeyToFile(userPublic, "./user.pub")
	writeKeyToFile(caPubText, "./ca.pub")
	return "yay"
}
