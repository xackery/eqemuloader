package script

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"text/template"

	toml "github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// Script wraps all variables
type Script struct {
	IsValgrind bool
	IsVerbose  bool
	IsDebug    bool
	Global     Global
	Database   Database
	Docker     Docker
	Bin        Bin
	ChatServer ChatServer `toml:"chatserver"`
	NATS       NATS
	World      World
	cache      cache `toml:"cache"`
	Web        Web
	PeqEditor  PeqEditor
	Discord    Discord
}

type cache struct {
	Entries map[string]string `toml:"cache"`
}

// Global Configuration
type Global struct {
	Stage string `toml:"stage"`
}

// Database Configuration
type Database struct {
	Username           string
	Password           string
	RootPassword       string
	Name               string
	Host               string
	Directory          string
	ReleaseType        string `toml:"release_type"`
	ReleaseURL         string `toml:"release_url"`
	ReleaseAuto        bool   `toml:"release_auto"`
	ReleaseUser        string `toml:"release_user"`
	ReleaseRepo        string `toml:"release_repo"`
	ReleaseAccessToken string `toml:"release_access_token"`
}

// Docker configuration
type Docker struct {
	Network string
}

// Web configuration
type Web struct {
	Port	string
	Url	string
}

// PeqEditor configuration
type PeqEditor struct {
	Port    string
}

// Docker configuration
type Discord struct {
        ChannelID        string
        ItemUrl          string
        RefreshRate      string
        ClientID         string
        ServerID         string
        UserName         string
        CommandChannelID string
}

// Bin Configuration
type Bin struct {
	Directory          string `toml:"directory"`
	URL                string `toml:"url"`
	AuthType           string `toml:"auth_type"`
	AuthUsername       string `toml:"auth_username"`
	AuthPassword       string `toml:"auth_password"`
	AuthKey            string `toml:"auth_key"`
	ReleaseType        string `toml:"release_type"`
	ReleaseURL         string `toml:"release_url"`
	ReleaseAuto        bool   `toml:"release_auto"`
	ReleaseUser        string `toml:"release_user"`
	ReleaseRepo        string `toml:"release_repo"`
	ReleaseAccessToken string `toml:"release_access_token"`
}

// ChatServer configuration
type ChatServer struct {
	Host string `toml:"host"`
	Port string `toml:"port"`
}

// NATS configuration
type NATS struct {
	Host string `toml:"host"`
	Port string `toml:"port"`
}

// World configuration
type World struct {
	ShortName    string `toml:"short_name"`
	LongName     string `toml:"long_name"`
	Locked       string
	LocalAddress string
}

// New creates a new script process
func New(cmd *cobra.Command) (s *Script, err error) {
	s = &Script{}

	var tree *toml.Tree
	f, err := os.Open("loader.conf")
	var tmpl *template.Template
	if err != nil {
		err = s.setup()
		if err != nil {
			return
		}

		data := defaultConfig()
		tmpl, err = template.New("config").Parse(data)
		if err != nil {
			err = errors.Wrap(err, "failed to parse config template")
			return
		}
		tmp := []byte{}
		buf := bytes.NewBuffer(tmp)
		err = tmpl.Execute(buf, s)
		if err != nil {
			err = errors.Wrap(err, "failed to execute template")
			return
		}
		err = ioutil.WriteFile("loader.conf", buf.Bytes(), 0744)
		if err != nil {
			err = errors.Wrap(err, "failed to write loader")
			return
		}

		f, err = os.Open("loader.conf")
		if err != nil {
			err = errors.Wrap(err, "failed to open loader.conf")
			return
		}
		os.Exit(0)
	}
	defer f.Close()
	tree, err = toml.LoadReader(f)
	if err != nil {
		err = errors.Wrap(err, "failed to load config")
		return
	}
	err = tree.Unmarshal(s)
	if err != nil {
		err = errors.Wrap(err, "failed to unmarshal config")
		return
	}
	s.IsValgrind, _ = cmd.Flags().GetBool("valgrind")
	s.IsVerbose, _ = cmd.Flags().GetBool("verbose")
	s.IsDebug, _ = cmd.Flags().GetBool("debug")
	if s.Database.Directory == "" || s.Database.Directory == "db" || s.Database.Directory == "db/" || s.Database.Directory == "./db" {
		err = fmt.Errorf("invalid database directory provided")
		return
	}
	wd, err := os.Getwd()
	if err != nil {
		err = errors.Wrap(err, "failed to get working directory")
	}
	if len(s.Bin.Directory) == 0 {
		err = fmt.Errorf("bin.directory is empty")
		return
	}
	s.Bin.Directory = fmt.Sprintf("%s/%s/", wd, s.Bin.Directory)
	if err != nil {
		err = fmt.Errorf("bin.directory is not a valid path: %s", s.Bin.Directory)
		return
	}
	if len(s.Database.Directory) == 0 {
		err = fmt.Errorf("database.directory is empty")
		return
	}
	s.Database.Directory = fmt.Sprintf("%s/%s/", wd, s.Database.Directory)
	if err != nil {
		err = fmt.Errorf("database.directory is not a valid path: %s", s.Database.Directory)
		return
	}
	if len(s.Docker.Network) < 3 {
		err = fmt.Errorf("docker.network must be at least 3 characters long")
		return
	}
	err = s.cacheLoad()
	if err != nil {
		err = errors.Wrap(err, "failed to load cache")
		return
	}
	return
}

func (s *Script) isPortListening(addr string) bool {
	if s.IsVerbose {
		fmt.Println("checking if addr", addr, "is listening")
	}
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return false
	}
	if conn != nil {
		return true
	}
	return false
}

func (s *Script) ensureImage(image string) (err error) {
	if s.IsVerbose {
		fmt.Println("ensuring image", image, "exists")
	}
	out, err := s.commandRunParse("docker images -q " + image)
	if err != nil {
		err = errors.Wrap(err, "failed to find docker image")
		return
	}
	if len(out) > 3 {
		return
	}

	fmt.Println("pulling image", image)
	err = s.commandRun("docker pull " + image)
	if err != nil {
		err = errors.Wrap(err, "failed to pull image")
		return
	}
	return
}
