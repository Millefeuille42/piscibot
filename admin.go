package main

import (
	"github.com/bwmarrin/discordgo"
	"os"
)

func everyoneUnregister(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID != os.Getenv("DISCORDADMIN") {
		return
	}
	res, err := session.Guild(message.GuildID)
	logError(err)
	for _, user := range res.Members {
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
	res, err := session.Guild(message.GuildID)
	logError(err)
	for _, user := range res.Members {
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
