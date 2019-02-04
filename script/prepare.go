package script

import (
	"os/exec"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/pkg/errors"
)

// Prepare is used to prepare an environment.
func (s *Script) prepare() (err error) {
	_, err = exec.LookPath("docker")
	if err != nil {
		err = fmt.Errorf("docker not found. Please install before using this")
		return
	}

	err = s.ensureImage("nats:latest")
	if err != nil {
		return
	}
	err = s.ensureImage("eqemu/server:latest")
	if err != nil {
		return
	}
	err = s.ensureImage("bitnami/mariadb:latest")
	if err != nil {
		return
	}

	s.commandRunParse("docker network create " + s.Docker.Network)

	err = s.binCheck()
	if err != nil {
		return
	}
	_, err = os.Stat("./bin/logs")
        if err != nil {
                err = os.Mkdir("./bin/logs", 0777)
                if err != nil {
                        err = errors.Wrap(err, "failed to make bin/logs directory")
                        return
                }
        }
	s.checkEqemuConfig()
	return
}

// Check for existing eqemu_config.json and generate one if needed.
func (s *Script) checkEqemuConfig() (err error) {
	_, err = os.Stat("./bin/eqemu_config.json")
        if err != nil {
			var tmpl *template.Template
			data := defaultEqemuConfig()
			tmpl, err = template.New("EqemuConfig").Parse(data)
			if err != nil {
				err = errors.Wrap(err, "failed to parse eqemu_config.json template")
				return
			}
			tmp := []byte{}
			buf := bytes.NewBuffer(tmp)
			err = tmpl.Execute(buf, s)
			if err != nil {
				err = errors.Wrap(err, "failed to execute eqemu_config.json template")
				return
			}
			err = ioutil.WriteFile("./bin/eqemu_config.json", buf.Bytes(), 0744)
			if err != nil {
				err = errors.Wrap(err, "failed to write bin/eqemu_config.json")
				return
			}
		}
	return
}
