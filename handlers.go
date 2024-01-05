package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	embed "github.com/Clinet/discordgo-embed"
	"github.com/alervaz/head-flipper/model"
	"github.com/bwmarrin/discordgo"
)

const (
	MAX = 12
)

func Flip(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Message.Content, "!help") {
		s.ChannelMessageSendEmbed(
			m.ChannelID,
			embed.NewGenericEmbed(
				"Help!",
				fmt.Sprintf(
					"!balance: to see your points\n !gamble {amount}: to gamble points from 1 to %d",
					MAX,
				),
			),
		)
	}

	if strings.HasPrefix(m.Message.Content, "!points") {
		points, err := model.GetPoints(m.Author.Username, m.GuildID)
		if err != nil {
			s.ChannelMessageSendEmbed(
				m.ChannelID,
				embed.NewErrorEmbed(
					"Unexpected error",
					err.Error(),
				),
			)
		}
		s.ChannelMessageSendEmbed(
			m.ChannelID,
			embed.NewGenericEmbed(
				m.Author.Username,
				fmt.Sprintf("Your balance is %d", points),
			),
		)
	}

	if strings.HasPrefix(m.Message.Content, "!gamble") {
		args := strings.Split(m.Message.Content, " ")
		if len(args) < 2 {
			s.ChannelMessageSendEmbed(
				m.ChannelID,
				embed.NewErrorEmbed(
					"Not enought arguments",
					fmt.Sprintf(
						"Please provide how many points you want to risk from 1 to %d",
						MAX,
					),
				),
			)
			return
		}

		pointsStr := args[1]
		points, err := strconv.Atoi(pointsStr)
		if err != nil {
			s.ChannelMessageSendEmbed(
				m.ChannelID,
				embed.NewErrorEmbed(
					"Invalid argument",
					fmt.Sprintf("Please specify an integer number between 1 to %d", MAX),
				),
			)
			return
		}

		if points > MAX || points < 1 {
			s.ChannelMessageSendEmbed(
				m.ChannelID,
				embed.NewErrorEmbed(
					"Invalid range",
					fmt.Sprintf(
						"Please provide how many points you want to risk from 1 to %d",
						MAX,
					),
				),
			)
			return
		}

		posibilities := []bool{true, false}
		won := posibilities[rand.Intn(len(posibilities))]
		userPoints, err := model.GetPoints(m.Author.Username, m.GuildID)
		if err != nil {
			s.ChannelMessageSendEmbed(
				m.ChannelID,
				embed.NewGenericEmbed(
					"Unexpected error",
					err.Error(),
				),
			)
			return
		}

		if won {
			err := model.Gamble(m.Author.Username, m.GuildID, userPoints+points)
			if err != nil {
				s.ChannelMessageSendEmbed(
					m.ChannelID,
					embed.NewGenericEmbed(
						"Unexpected error",
						err.Error(),
					),
				)
				return
			}

			s.ChannelMessageSendEmbed(
				m.ChannelID,
				embed.NewGenericEmbed(
					"You won it was heads!",
					fmt.Sprintf("You won %d points", points),
				),
			)
			return
		}

		result := userPoints - (points - 2)
		if result < 0 {
			result = 0
		}
		err = model.Gamble(m.Author.Username, m.GuildID, result)
		if err != nil {
			s.ChannelMessageSendEmbed(
				m.ChannelID,
				embed.NewGenericEmbed(
					"Unexpected error",
					err.Error(),
				),
			)
			return
		}

		until := time.Now().Add(time.Minute * time.Duration(points*5))
		err = s.GuildMemberTimeout(m.GuildID, m.Author.ID, &until)
		if err != nil {
			log.Println(err)
			s.ChannelMessageSendEmbed(
				m.ChannelID,
				embed.NewGenericEmbed(
					"Cannot mute a mod",
					"You cant play dirty mod",
				),
			)
			return
		}

		s.ChannelMessageSendEmbed(
			m.ChannelID,
			embed.NewGenericEmbed(
				"You lost",
				fmt.Sprintf("Enjoy your %d minutes out", points*5),
			),
		)
	}
}
