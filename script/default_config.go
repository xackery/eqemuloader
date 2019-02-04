package script

func defaultConfig() string {
	return `
## Database represents a docker instance running MariaDB
## Note that if you change username/password, you MUST recreate the entire database
## in order for it to go into effect.
[database]

	## Docker name of database. (default: eqemu)
	name = "eqemu"

	## Username to log into database
	username = "eqemu"

	## Password to log into database
	password = "{{.Database.Password}}"

        ## Password to log into database as root
        rootpassword = "{{.Database.RootPassword}}"

	## Directory to save database data to, relative to loader's directory
	directory = "database"

	## For latest database snapshots, you need to set a release endpoint
	## two types are supported: github or gitea
	release_type = "{{.Database.ReleaseType}}"

	## Url for latest database. In most situations it is simply api.github.com
	release_url = "{{.Database.ReleaseURL}}"

	## Automatically pull down new releases as they become available.
	release_auto = true

	## user (or org) a repository exists under
	release_user = "{{.Database.ReleaseUser}}"
	## repository name to get releases
	release_repo = "{{.Database.ReleaseRepo}}"

	## personal access token, if not a public repository
	release_access_token = "{{.Database.ReleaseAccessToken}}"
[docker]
	network = "eqemu"

## Bin is the binary folder
[bin]
	## URL is where to fetch the bin repository from.
	## It will honor submodules, pulling them if any are found
	url = "git@git.rebuildeq.com:eq/bin.git"

	## Directory to git clone bin to. It is recommended to leave it the same
	directory = "bin"
	
	## There are 3 types of auth supported: http, ssh, or ssh_key
	auth_type = "ssh_key"

	## Used by all auth types, your credentials to log into bin url
	auth_username = "{{.Bin.AuthUsername}}"

	## Used by http and ssh auth types
	auth_password = "{{.Bin.AuthPassword}}"

	## Used by ssh_key only. absolute path to auth key
	auth_key = "{{.Bin.AuthKey}}"

	## For latest eqemu binaries, you need to set a release endpoint
	## two types are supported: github or gitea
	release_type = "{{.Bin.ReleaseType}}"

	## Url for latest binaries. In most situations it is simply api.github.com
	release_url = "{{.Bin.ReleaseURL}}"

	## Automatically pull down new releases as they become available.
	release_auto = true

	## user (or org) a repository exists under
	release_user = "{{.Bin.ReleaseUser}}"
	## repository name to get releases
	release_repo = "{{.Bin.ReleaseRepo}}"

	## personal access token, if not a public repository
	release_access_token = "{{.Bin.ReleaseAccessToken}}"

## EQEMU ChatServer configuration
[chatserver]
	host = "127.0.0.1"
	port = "7778"

[nats]
	host = "nats"
	port = "4222"

[world]
	long_name = "{{.World.LongName}}"
	short_name = "{{.World.ShortName}}"
	locked = "true"
	localaddress = ""

[web]
	port = "80"
	url = "http://rebuildeq.com"

[peqeditor]
	port = "81"

[discord]
	channelid = ""
	itemurl = ""
	refreshrate = "15"
	clientid = ""
	serverid = ""
	username = ""
	commandchannelid = ""
`
}
