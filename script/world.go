package script

import (
	"fmt"
	"os"
	"strings"
)

// WorldRun handles world service
func (s *Script) WorldRun(args []string) (err error) {
	err = s.prepare()
	if err != nil {
		return
	}
	out, err := s.commandRunParse("docker ps")
	if err != nil {
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "crash" {
		err = s.crash("world")
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "pull" {
		err = s.binReleaseDownload()
		if err != nil {
			return
		}
		err = s.binPull()
		if err != nil {
			return
		}
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "stop" {
		if !strings.Contains(out, "world") {
			return
		}
		s.commandRun("docker stop world")
		s.commandRunParse("docker rm world")
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "logs" {
		s.commandRun("docker logs world")
		return
	}
	s.commandRunParse("docker rm world")

	fmt.Println("World")
	useValgrind := ""
	if s.IsValgrind {
		useValgrind = "valgrind --tool=memcheck --leak-check=yes "
	}
	if s.IsDebug {
		err = s.commandRunAttached(fmt.Sprintf("docker run -i --cap-add=SYS_PTRACE --security-opt seccomp=unconfined --privileged -v %s:/src --ulimit core=10000000 -e LD_LIBRARY_PATH=/src/ --network=%s --name=world -p 5998:5998/udp -p 5999:5999/udp -p 9000:9000 -p 9000:9000/udp -p 9001:9001 -p 9080:9080 eqemu/server gdb ./world", s.Bin.Directory, s.Docker.Network))
		return
	}
	_, err = s.commandRunDetached(fmt.Sprintf("docker run -v %s:/src --ulimit core=10000000 -e LD_LIBRARY_PATH=/src/ --network=%s --name=world -p 5998:5998/udp -p 5999:5999/udp -p 9000:9000 -p 9000:9000/udp -p 9001:9001 -p 9080:9080 eqemu/server %s./world", s.Bin.Directory, s.Docker.Network, useValgrind))
	if err != nil {
		return
	}
	//fmt.Println("started. use `docker logs world` to see logs")
	return
}
