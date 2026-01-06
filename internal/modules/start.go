/*
  - This file is part of YukkiMusic.
*/

package modules

import (
	"github.com/Laky-64/gologging"
	tg "github.com/amarnathcjd/gogram/telegram"

	"main/internal/config"
	"main/internal/core"
	"main/internal/database"
	"main/internal/locales"
	"main/internal/utils"
)

func init() {
	helpTexts["/start"] = `<i>Start the bot and show main menu.</i>`
}

func startHandler(m *tg.NewMessage) error {

	// ===== GROUP CHAT =====
	if m.ChatType() != tg.EntityUser {
		database.AddServed(m.ChannelID())
		m.Reply(F(m.ChannelID(), "start_group"))
		return tg.ErrEndGroup
	}

	// ===== PRIVATE CHAT =====
	arg := m.Args()
	database.AddServed(m.ChannelID(), true)

	if arg != "" {
		gologging.Info(
			"Got Start parameter: " + arg + " in ChatID: " +
				utils.IntToStr(m.ChannelID()),
		)
	}

	switch arg {

	case "pm_help":
		gologging.Info("User requested help via start param")
		helpHandler(m)

	default:
		caption := F(m.ChannelID(), "start_private", locales.Arg{
			"user": utils.MentionHTML(m.Sender),
			"bot":  utils.MentionHTML(core.BUser),
		})

		// ---------- SEND IMAGE FIRST ----------
		if config.StartImage != "" {
			_, err := m.RespondMedia(config.StartImage, &tg.MediaOptions{
				NoForwards: true,
			})
			if err != nil {
				gologging.Error("[start] Image send failed: " + err.Error())
			}
		}

		// ---------- SEND TEXT + BUTTONS ----------
		_, err := m.Respond(caption, &tg.SendOptions{
			ReplyMarkup: core.GetStartMarkup(m.ChannelID()),
			NoForwards:  true,
		})
		if err != nil {
			return err
		}
	}

	// ===== LOGGER =====
	if config.LoggerID != 0 && isLogger() {
		uName := "N/A"
		if m.Sender.Username != "" {
			uName = "@" + m.Sender.Username
		}
		msg := F(m.ChannelID(), "logger_bot_started", locales.Arg{
			"mention":       utils.MentionHTML(m.Sender),
			"user_id":       m.SenderID(),
			"user_username": uName,
		})

		_, err := m.Client.SendMessage(config.LoggerID, msg)
		if err != nil {
			gologging.Error(
				"Failed to send logger_bot_started msg: " + err.Error(),
			)
		}
	}

	return tg.ErrEndGroup
}

func startCB(cb *tg.CallbackQuery) error {
	cb.Answer("")

	caption := F(cb.ChannelID(), "start_private", locales.Arg{
		"user": utils.MentionHTML(cb.Sender),
		"bot":  utils.MentionHTML(core.BUser),
	})

	// ---------- SEND IMAGE FIRST ----------
	if config.StartImage != "" {
		_, err := cb.Client.SendMessage(cb.ChannelID(), "", &tg.SendOptions{
			Media:      config.StartImage,
			NoForwards: true,
		})
		if err != nil {
			gologging.Error("[startCB] Image send failed: " + err.Error())
		}
	}

	// ---------- SEND TEXT + BUTTONS ----------
	_, err := cb.Client.SendMessage(cb.ChannelID(), caption, &tg.SendOptions{
		ReplyMarkup: core.GetStartMarkup(cb.ChannelID()),
		NoForwards:  true,
	})
	if err != nil {
		return err
	}

	return tg.ErrEndGroup
}
