package ca

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func addToAgent(cert *ssh.Certificate, key *ecdsa.PrivateKey) {
	con, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))

	if err != nil {
		fmt.Println("can't connect to SSH agent: ", err)
	}

	sshAgent := agent.NewClient(con)

	if err = sshAgent.Add(agent.AddedKey{
		PrivateKey:  key,
		Certificate: cert,
	}); err != nil {
		fmt.Println("ssh-agent failure: ", err)
	}
}

func genPrivateKey() *ecdsa.PrivateKey {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Println("generate failed: ", err)
	}

	return privateKey
}

// SignCert signs certs
func SignCert() {
	// user key pair
	userKey := genPrivateKey()

	// we sign this public key with the CA
	userKeyPub, e := ssh.NewPublicKey(&userKey.PublicKey)

	if e != nil {
		fmt.Println("user public key failed: ", e)
	}

	//read in private key to act as CA
	caFile, readErr := ioutil.ReadFile("./ca.private")

	if readErr != nil {
		fmt.Println("error reading private key: ", readErr)
	}

	ca, caParseErr := ssh.ParsePrivateKey(caFile)

	if caParseErr != nil {
		fmt.Println("ca parse error: ", caParseErr)
	}

	expireTime, _ := time.ParseDuration("2h")

	//we create this Cert struct using the user's public key
	//https://github.com/ejcx/sshcert/blob/1c64826f1a45d87777103946575701b0a062623a/sshcert.go#L82
	// SignCert is called to sign an ssh public key and produce an ssh certificate.
	certInstance := &ssh.Certificate{
		Key:             userKeyPub,
		Serial:          400,
		CertType:        ssh.UserCert,
		KeyId:           "using key from file",
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
	certErr := certInstance.SignCert(rand.Reader, ca)

	if certErr != nil {
		fmt.Println("Error signing Certificate: ", certErr)
	}

	addToAgent(certInstance, userKey)
}
