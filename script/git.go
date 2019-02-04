package script

import (
	"fmt"
	"io/ioutil"
	"net"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	gogitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

func getAuth(authType string, authUsername string, authPassword string, authKey string) (auth transport.AuthMethod, err error) {
	var sshKey []byte
	var signer ssh.Signer
	switch strings.ToLower(authType) {
	case "http":
		auth = &http.BasicAuth{Username: authUsername, Password: authPassword}
	case "ssh":
		hostKeyCallback := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		}
		auth = &gogitssh.Password{User: authUsername, Password: authPassword, HostKeyCallbackHelper: gogitssh.HostKeyCallbackHelper{
			HostKeyCallback: hostKeyCallback,
		}}
	case "ssh_key":
		sshKey, err = ioutil.ReadFile(authKey)
		if err != nil {
			err = errors.Wrapf(err, "failed to read ssh key at %s", authKey)
			return
		}
		signer, err = ssh.ParsePrivateKey([]byte(sshKey))
		if err != nil {
			err = errors.Wrapf(err, "failed to sign private key %s", authKey)
			return
		}
		auth = &gogitssh.PublicKeys{User: "git", Signer: signer}
	default:
		err = fmt.Errorf("unknown bin auth type: %s", authKey)
		return
	}
	return
}
