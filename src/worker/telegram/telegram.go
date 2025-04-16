package telegram

import (
	"bmkg/src/db"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"log"
	"strconv"
)

// Bot represents a Telegram bot instance
type Bot struct {
	app core.App
	api *tgbotapi.BotAPI
}

// NewBot creates a new Bot instance
func NewBot(app core.App) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	api.Debug = false
	log.Printf("Authorized on account %s", api.Self.UserName)

	return &Bot{
		app: app,
		api: api,
	}, nil
}

// Start begins listening for and processing updates
func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		// Route message to appropriate handler
		if update.Message.Location != nil {
			b.handleLocationMessage(update.Message)
		} else if update.Message.Text == "Bantuan" {
			b.handleHelpMessage(update.Message)
		} else {
			b.showMainMenu(update.Message)
		}
	}
}

// handleLocationMessage processes location data and registers user for notifications
func (b *Bot) handleLocationMessage(message *tgbotapi.Message) {
	location := message.Location
	latitude := location.Latitude
	longitude := location.Longitude
	chatID := message.Chat.ID

	// Register user in database for earthquake notifications
	err := b.registerUserLocation(chatID, latitude, longitude)

	var reply tgbotapi.MessageConfig
	if err != nil {
		log.Printf("Error registering user: %v", err)
		reply = tgbotapi.NewMessage(chatID,
			"Terjadi kesalahan saat menyimpan lokasi Anda. Silakan coba lagi nanti.")
	} else {
		reply = tgbotapi.NewMessage(chatID,
			"Lokasi Anda telah disimpan! Anda akan menerima notifikasi gempa bumi yang terjadi di sekitar lokasi Anda.")
	}

	if _, err := b.api.Send(reply); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

// handleHelpMessage sends help information to the user
func (b *Bot) handleHelpMessage(message *tgbotapi.Message) {
	reply := tgbotapi.NewMessage(message.Chat.ID,
		"Selamat datang di bot ini!\n\n"+
			"Anda dapat mengirimkan lokasi dengan menekan tombol 'Kirim Lokasi'.")

	reply.ReplyMarkup = b.getMainMenuKeyboard()

	if _, err := b.api.Send(reply); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

// showMainMenu displays the main menu options
func (b *Bot) showMainMenu(message *tgbotapi.Message) {
	reply := tgbotapi.NewMessage(message.Chat.ID,
		"Silakan kirim lokasi Anda atau pilih opsi lainnya:")

	reply.ReplyMarkup = b.getMainMenuKeyboard()

	if _, err := b.api.Send(reply); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

// getMainMenuKeyboard creates the keyboard for main menu
func (b *Bot) getMainMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonLocation("Kirim Lokasi"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Bantuan"),
		),
	)
}

// registerUserLocation saves the user's location in the database
func (b *Bot) registerUserLocation(chatID int64, latitude, longitude float64) error {
	// Implement the logic to save the user's location in the database
	// This is a placeholder for actual database interaction
	fmt.Printf("Registering user with chat ID %d at location (%f, %f)\n", chatID, latitude, longitude)

	d, _ := db.NewProxy[db.UserNotify](b.app)
	//d := &db.UserNotify{}

	err := b.app.RecordQuery("user_notify").
		AndWhere(dbx.NewExp("identifier={:identifier}", dbx.Params{
			"identifier": fmt.Sprintf("%d", chatID), // convert int64 to string
		})).
		Limit(1).
		One(d)

	if err != nil {
		fmt.Printf("Failed to get user notify from db: %v\n", err)

	}

	d.SetIdentifier(strconv.FormatInt(chatID, 10))
	d.SetLintang(strconv.FormatFloat(latitude, 'f', -1, 64))
	d.SetBujur(strconv.FormatFloat(longitude, 'f', -1, 64))

	d.SetType(db.Telegram)

	err = b.app.Save(d)
	if err != nil {
		return err
	}

	// Example: Save to database (pseudo-code)
	// err := b.app.DB().SaveUserLocation(chatID, latitude, longitude)
	// return err

	return nil
}

// SendMessage sends a message to a specific chat ID
func (b *Bot) SendMessage(chatID int64, message string) error {
	msg := tgbotapi.NewMessage(chatID, message)
	_, err := b.api.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message to %d: %v", chatID, err)
	}
	return nil
}

// Run initializes and starts the bot
func Run(app core.App) {
	bot, err := NewBot(app)
	if err != nil {
		log.Panic(err)
	}

	go bot.Start()
}
