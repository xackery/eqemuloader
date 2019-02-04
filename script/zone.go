package script

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ZoneRun handles zone service
func (s *Script) ZoneRun(args []string) (err error) {
	count := 1
	var tmpCount int64
	if len(args) > 0 {
		tmpCount, err = strconv.ParseInt(args[0], 10, 64)
		if err == nil {
			count = int(tmpCount)
		}
	}
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
		for zoneID := 7000; zoneID < 7101; zoneID++ {
			if !strings.Contains(out, fmt.Sprintf("zone%d", zoneID)) {
				continue
			}
			s.commandRun(fmt.Sprintf("docker stop zone%d", zoneID))
			s.commandRunParse(fmt.Sprintf("docker rm zone%d", zoneID))
		}
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "logs" {
		s.commandRun(fmt.Sprintf("docker logs zone%d", 7000+count-1))
		return
	}

	fmt.Println("Zone", count)

	for i := 0; i < int(count); i++ {
		zoneID := 7000 + i
		if strings.Contains(out, fmt.Sprintf("zone%d", zoneID)) {
			continue
		}
		s.commandRunParse(fmt.Sprintf("docker rm zone%d", zoneID))
		useValgrind := ""
		if s.IsValgrind {
			useValgrind = "valgrind --tool=memcheck --leak-check=yes "
		}
                if s.IsDebug {
			err = s.commandRunAttached(fmt.Sprintf("docker run -i --cap-add=SYS_PTRACE --security-opt seccomp=unconfined --privileged -v %s:/src --ulimit core=10000000 --network=%s -p %d:%d/udp -e LD_LIBRARY_PATH=/src/ --name=zone%d eqemu/server gdb --args ./zone dynamic_zone%d:%d", s.Bin.Directory, s.Docker.Network, zoneID, zoneID, zoneID, zoneID, zoneID))
                	return
		}
		_, err = s.commandRunDetached(fmt.Sprintf("docker run -v %s:/src --ulimit core=10000000 --network=%s -p %d:%d/udp -e LD_LIBRARY_PATH=/src/ --name=zone%d eqemu/server %s./zone dynamic_zone%d:%d", s.Bin.Directory, s.Docker.Network, zoneID, zoneID, zoneID, useValgrind, zoneID, zoneID))
		if err != nil {
			return
		}
	}

	//fmt.Printf("use `docker logs zone7000` to `docker logs zone%d` to see logs\n", 7000+count-1)
	return
}
