package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	WinDesc   = "You Win!"
	DrawDesc  = "Draw!"
	LoseDesc  = "You Lose!"
	WinColor  = "43b581"
	DrawColor = "faa61a"
	LoseColor = "f04747"
)

var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)
	dg.AddHandler(messageReactionAdd)

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages | discordgo.IntentGuildMessageReactions)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "!rps") {
		hex, _ := strconv.ParseInt("43b581", 16, 64)
		embed := discordgo.MessageEmbed{Title: "Rock Paper Scissors", Description: "Select your hand", Color: int(hex)}
		message, _ := s.ChannelMessageSendEmbed(m.ChannelID, &embed)

		s.MessageReactionAdd(m.ChannelID, message.ID, "✊")
		s.MessageReactionAdd(m.ChannelID, message.ID, "✌")
		s.MessageReactionAdd(m.ChannelID, message.ID, "🖐")
	}
}

func messageReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.UserID == s.State.User.ID {
		return
	}
	hands := []string{"✊", "✌", "🖐"}

	rand.Seed(time.Now().UnixNano())
	rand := rand.Intn(len(hands))
	hand := hands[rand]

	var desc string
	var color string
	handFiled := discordgo.MessageEmbedField{Name: "Bot Hand", Value: hand, Inline: true}
	playerHandFiled := discordgo.MessageEmbedField{Name: "Your Hand", Value: r.Emoji.Name, Inline: true}

	switch hand {
	case "✊":
		switch r.Emoji.Name {
		case "✊":
			desc = DrawDesc
			color = DrawColor
		case "✌":
			desc = LoseDesc
			color = LoseColor
		case "🖐":
			desc = WinDesc
			color = WinColor
		}
	case "✌":
		switch r.Emoji.Name {
		case "✊":
			desc = WinDesc
			color = WinColor
		case "✌":
			desc = DrawDesc
			color = DrawColor
		case "🖐":
			desc = LoseDesc
			color = LoseColor
		}
	case "🖐":
		switch r.Emoji.Name {
		case "✊":
			desc = LoseDesc
			color = LoseColor
		case "✌":
			desc = WinDesc
			color = WinColor
		case "🖐":
			desc = DrawDesc
			color = DrawColor
		}
	}

	embed := discordgo.MessageEmbed{Title: "Rock Paper Scissors", Description: "Processing...", Color: 0}
	s.ChannelMessageEditEmbed(r.ChannelID, r.MessageID, &embed)
	time.Sleep(time.Second * 1)

	hex, _ := strconv.ParseInt(color, 16, 64)
	finalEmbed := discordgo.MessageEmbed{
		Title:       "Rock Paper Scissors",
		Description: desc,
		Color:       int(hex),
		Fields:      []*discordgo.MessageEmbedField{&handFiled, &playerHandFiled},
	}
	s.ChannelMessageEditEmbed(r.ChannelID, r.MessageID, &finalEmbed)
}
