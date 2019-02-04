package script

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// NATSRun represents NATS
func (s *Script) NATSRun(args []string) (err error) {
	err = s.prepare()
	if err != nil {
		return
	}
	out, err := s.commandRunParse("docker ps")
	if err != nil {
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "stop" {
		if !strings.Contains(out, "nats") {
			return
		}
		s.commandRun("docker stop nats")
		s.commandRunParse("docker rm nats")
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "logs" {
		s.commandRun("docker logs nats")
		return
	}
	fmt.Println("NATS")
	if s.isPortListening("0.0.0.0:4222") {
		return
	}
	s.commandRunParse("docker rm nats")
	err = s.commandRun(fmt.Sprintf("docker run -p 4222:4222 -p 6222:6222 -p 8222:8222 --network=%s --name=nats --detach=true nats:latest", s.Docker.Network))
	if err != nil {
		return
	}
	for counter := 0; counter < 15; counter++ {
		<-time.After(1 * time.Second)
		if s.isPortListening("0.0.0.0:4222") {
			return
		}
	}
	err = fmt.Errorf("timeout waiting for nats to be up")
	return
}
