package script

func defaultEqemuConfig() string {
	return `
	{
		"server": {
			"discord" : {
                                "channelid" : "{{.Discord.ChannelID}}",
                                "itemurl" : "{{.Discord.ItemUrl}}",
                                "refreshrate" : "{{.Discord.RefreshRate}}",
                                "clientid": "{{.Discord.ClientID}}",
                                "serverid" : "{{.Discord.ServerID}}",
                                "username" : "{{.Discord.UserName}}",
                                "commandchannelid" : "{{.Discord.CommandChannelID}}"
                        },
			"chatserver": {
				"host": "{{.ChatServer.Host}}",
				"port": "{{.ChatServer.Port}}"
			},
			"nats": {
				"host": "{{.NATS.Host}}",
				"port": "{{.NATS.Port}}"
			},
			"world": {
				"locked": "{{.World.Locked}}",
				"localaddress": "{{.World.LocalAddress}}",
				"shortname": "{{.World.ShortName}}",
				"loginserver": {
					"password": "",
					"host": "login.eqemulator.net",
					"legacy": "1",
					"port": "5998",
					"account": ""
				},
				"telnet": {
					"host": "world",
					"port": "9000",
					"telnet": "enabled"
				},
				"http": {
					"port": "9080",
					"enabled": "false",
					"mimefile": "mime.types"
				},
				"longname": "{{.World.LongName}}",
				"tcp": {
					"host": "world",
					"port": "9000",
					"telnet": "enabled"
				},
				"key": "oasijfoaisfjoij23%j2o3ij5o23i222"
			},
			"mailserver": {
				"host": "127.0.0.1",
				"port": "7778"
			},
			"zones": {
				"defaultstatus": "20",
				"ports": {
					"low": "7000",
					"high": "7400"
				}
			},
			"database": {
				"host": "mariadb",
				"port": "3306",
				"username": "{{.Database.Username}}",
				"password": "{{.Database.Password}}",
				"db": "{{.Database.Name}}"
			},
			"qsdatabase": {
				"host": "mariadb",
				"port": "3306",
				"username": "{{.Database.Username}}",
				"password": "{{.Database.Password}}",
				"db": "{{.Database.Name}}"
			}
		}
	}
	`
}
