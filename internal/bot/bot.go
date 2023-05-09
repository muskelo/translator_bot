package bot

import (
	"fmt"
	"strings"
	"time"

	lt "github.com/muskelo/translator_bot/internal/libretranslate"
	"github.com/muskelo/translator_bot/internal/sql"
	tele "gopkg.in/telebot.v3"
	"xorm.io/xorm"
)

// -- Bot --

type Defaults struct {
	PrimaryLanguage   string
	SecondaryLanguage string
}

type Settings struct {
	DatabaseConnStr   string
	LibreTranslateUrl string
	BotToken          string
	Defaults          Defaults
}

type Bot struct {
	*tele.Bot
	Engine      *xorm.Engine
	Translation *lt.Translation
	Defaults    Defaults
}

func New(settings Settings) (*Bot, error) {
	engine, err := sql.New(settings.DatabaseConnStr)
	if err != nil {
		return nil, err
	}

	translation, err := lt.New(settings.LibreTranslateUrl)
	if err != nil {
		return nil, err
	}
	err = translation.LoadSupportedLanguages()
	if err != nil {
		return nil, err
	}

	pref := tele.Settings{
		Token:  settings.BotToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	telebot, err := tele.NewBot(pref)
	if err != nil {
		return nil, err
	}

	return &Bot{
		Bot:         telebot,
		Engine:      engine,
		Translation: translation,
		Defaults:    settings.Defaults,
	}, nil
}

func (bot *Bot) Init() {
	bot.Use(bot.errorHandlerMiddlware)
	bot.Handle("/start", bot.handlerStart)
	bot.Handle("/supported", bot.handlerSupported)
	bot.Handle("/primary", bot.handlerSetPrimaryLang)
	bot.Handle("/secondary", bot.handlerSetSecondaryLang)
	bot.Handle("/t", bot.handlerTranslate)
	bot.Handle(tele.OnText, bot.handlerAutoTranslate)
}

//
// -- Bot handlers --
//

func (bot *Bot) handlerStart(ctx tele.Context) error {
	return sql.CreateUser(bot.Engine, sql.User{
		TelegramId:        ctx.Sender().ID,
		PrimaryLanguage:   bot.Defaults.PrimaryLanguage,
		SecondaryLanguage: bot.Defaults.SecondaryLanguage,
	})
}

func (bot *Bot) handlerSetPrimaryLang(ctx tele.Context) error {
	args := ctx.Args()
	if len(args) == 0 {
		return ctx.Send("No language specified")
	}

	langCode := args[0]
	if !bot.Translation.IsSupported(langCode) {
		return ctx.Send("Language not supported")
	}

	return sql.UpdateUser(bot.Engine, sql.User{
		TelegramId:      ctx.Sender().ID,
		PrimaryLanguage: langCode,
	})
}

func (bot *Bot) handlerSetSecondaryLang(ctx tele.Context) error {
	args := ctx.Args()
	if len(args) == 0 {
		return ctx.Send("No language specified")
	}

	langCode := args[0]
	if !bot.Translation.IsSupported(langCode) {
		return ctx.Send("Language not supported")
	}

	return sql.UpdateUser(bot.Engine, sql.User{
		TelegramId:        ctx.Sender().ID,
		SecondaryLanguage: langCode,
	})
}

func (bot *Bot) handlerAutoTranslate(ctx tele.Context) error {
	confidence, sourseLangCode, targetLangCode, err := bot.determineLangs(ctx)
	if err != nil {
		return err
	}

	targetText, err := bot.Translation.Translate(sourseLangCode, targetLangCode, ctx.Text())
	if err != nil {
		return err
	}

	msgText := fmt.Sprintf("%v\n<i>%v -> %v</i>", targetText, sourseLangCode, targetLangCode)
	if confidence < 0.5 {
		msgText = fmt.Sprintf("%v\n%v", msgText, "<i>Not confident in sourse language</i>")
	}
	msgOption := &tele.SendOptions{
		ReplyTo:   ctx.Message(),
		ParseMode: tele.ModeHTML,
	}
	return ctx.Send(msgText, msgOption)
}

func (bot *Bot) handlerTranslate(ctx tele.Context) error {
	if ctx.Message().ReplyTo == nil {
		return ctx.Send("Translation command must reply to the message being translated")
	}
	sourseText := ctx.Message().ReplyTo.Text

	args := ctx.Args()
	if len(args) != 2 {
		return ctx.Send("Translation command require 2 args.")
	}
	sourseLangCode := args[0]
	targetLangCode := args[1]

	targetText, err := bot.Translation.Translate(sourseLangCode, targetLangCode, sourseText)
	if err != nil {
		return err
	}

	msgText := fmt.Sprintf("%v\n<i>%v -> %v</i>", targetText, sourseLangCode, targetLangCode)
	msgPption := &tele.SendOptions{
		ReplyTo:   ctx.Message().ReplyTo,
		ParseMode: tele.ModeHTML,
	}
	return ctx.Send(msgText, msgPption)
}

func (bot *Bot) handlerSupported(ctx tele.Context) error {
	elems := []string{"Langs:"}
	for i, lang := range bot.Translation.Langs {
		str := fmt.Sprintf("%v: %v (%v)", i+1, lang.Name, lang.Code)
		elems = append(elems, str)
	}
	return ctx.Send(strings.Join(elems, "\n"))
}

//
// -- Bot middleware --
//

func (bot *Bot) errorHandlerMiddlware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(ctx tele.Context) error {
		unknowErr := next(ctx)
		if unknowErr == nil {
			return nil
		}

		switch unknowErr {
		case sql.ErrUserNotFound:
			return ctx.Send("Please, run /start command")
		}

		switch err := unknowErr.(type) {
		case lt.NotSupportedError:
			return ctx.Send(fmt.Sprintf("%v not supported", err.LangCode))
		case lt.NoTargetError:
			return ctx.Send(fmt.Sprintf("%v cannot be translated into %v", err.SourceLangCode, err.TragerLangCode))
		}

		ctx.Send("Sory someting went wrong")
		return unknowErr
	}
}

//
// -- Bot other funcs --
//

// determine the languages by input and user information
func (bot *Bot) determineLangs(ctx tele.Context) (confidence float32, sourseLangCode string, targetLangCode string, err error) {
	user, err := sql.GetUserByID(bot.Engine, ctx.Sender().ID)
	if err != nil {
		return
	}

	confidence, sourseLangCode, err = bot.Translation.Detect(ctx.Text())
	if err != nil {
		return
	}

	if sourseLangCode == user.PrimaryLanguage {
		targetLangCode = user.SecondaryLanguage
	} else {
		targetLangCode = user.PrimaryLanguage
	}
	return
}
