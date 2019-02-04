package script

import (
	"fmt"
	"os"
	"strings"
)

// SharedRun runs shared_memory
func (s *Script) SharedRun(args []string) (err error) {
	err = s.prepare()
	if err != nil {
		return
	}
	out, err := s.commandRunParse("docker ps")
	if err != nil {
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "crash" {
		err = s.crash("zone")
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "stop" {
		if !strings.Contains(out, "shared") {
			return
		}
		s.commandRun("docker stop shared")
		s.commandRunParse("docker rm shared")
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "logs" {
		s.commandRun("docker logs shared_memory")
		return
	}
	fmt.Println("Shared Memory")
	s.commandRunParse("docker rm shared_memory")
	useValgrind := ""
	if s.IsValgrind {
		useValgrind = "valgrind --tool=memcheck --leak-check=yes "
	}
	if s.IsDebug {
		err = s.commandRunAttached(fmt.Sprintf("docker run -i --cap-add=SYS_PTRACE --security-opt seccomp=unconfined --privileged -v %s:/src --network=%s --ulimit core=10000000 -e LD_LIBRARY_PATH=/src/ --name=shared_memory eqemu/server gdb ./shared_memory", s.Bin.Directory, s.Docker.Network))
		return
	}
	err = s.commandRun(fmt.Sprintf("docker run -v %s:/src --network=%s --ulimit core=10000000 -e LD_LIBRARY_PATH=/src/ --name=shared_memory eqemu/server %s./shared_memory", s.Bin.Directory, s.Docker.Network, useValgrind))
	if err != nil {
		return
	}
	return
}
