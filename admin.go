package main

import (
	"github.com/bwmarrin/discordgo"
	"os"
)

func everyoneUnregister(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID != os.Getenv("DISCORDADMIN") {
		return
	}
	users, _ := session.GuildMembers(message.GuildID, "", 100)
	for _, user := range users {
		_ = session.GuildMemberRoleAdd(message.GuildID, user.User.ID, os.Getenv("DISCORDUNREGISTEREDROLE"))
		_ = session.GuildMemberRoleRemove(message.GuildID, user.User.ID, os.Getenv("DISCORDREGISTEREDROLE"))
		_ = session.GuildMemberRoleRemove(message.GuildID, user.User.ID, os.Getenv("DISCORDSPECTATORROLE"))
	}
	_, _ = session.ChannelMessageSend(message.ChannelID, "Task successfully accomplished")
}

func everyoneSpectator(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID != os.Getenv("DISCORDADMIN") {
		return
	}
	users, _ := session.GuildMembers(message.GuildID, "", 100)
	for _, user := range users {
		registeredFlag := false
		for _, role := range user.Roles {
			if role == os.Getenv("DISCORDREGISTEREDROLE") {
				registeredFlag = true
			}
		}
		if !registeredFlag {
			_ = session.GuildMemberRoleAdd(message.GuildID, user.User.ID, os.Getenv("DISCORDSPECTATORROLE"))
			_ = session.GuildMemberRoleRemove(message.GuildID, user.User.ID, os.Getenv("DISCORDUNREGISTEREDROLE"))
		}
	}
	_, _ = session.ChannelMessageSend(message.ChannelID, "Task successfully accomplished")
}
