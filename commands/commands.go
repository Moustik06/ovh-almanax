package commands

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/go-resty/resty/v2"
	"github.com/robfig/cron/v3"
	"net/http"
	"strconv"
	"strings"
)

type AlmanaxBonus struct {
	Bonus struct {
		Description string `json:"description"`
	} `json:"bonus"`
	Date    string `json:"date"`
	Tribute struct {
		Item     Item `json:"item"`
		Quantity int  `json:"quantity"`
	} `json:"tribute"`
}

type Item struct {
	Name      string `json:"name"`
	ImageURLs struct {
		SD string `json:"sd"`
	} `json:"image_urls"`
}

func createEmbed(bonus *AlmanaxBonus) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("Almanax du %s", bonus.Date),
		Color:       0x00ff00, // Couleur verte
		Description: bonus.Bonus.Description,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Offrande",
				Value: bonus.Tribute.Item.Name,
			},
			{
				Name:  "Quantité",
				Value: strconv.Itoa(bonus.Tribute.Quantity),
			},
		},

		Footer: &discordgo.MessageEmbedFooter{
			Text: "Date : " + bonus.Date,
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: bonus.Tribute.Item.ImageURLs.SD,
		},
	}

	return embed
}
func getAlmanaxBonus() (*AlmanaxBonus, *[]AlmanaxBonus, error) {
	url := "https://api.dofusdu.de/dofus2/fr/almanax" // Remplace {language} par la langue souhaitée

	// Crée un client HTTP
	client := resty.New()
	// Effectue la requête GET à l'API
	response, err := client.R().Get(url)
	if err != nil {
		return nil, nil, err
	}

	if response.StatusCode() != http.StatusOK {
		return nil, nil, fmt.Errorf("La requête GET a retourné un statut d'erreur : %d", response.StatusCode())
	}

	var bonuses []AlmanaxBonus
	err = json.Unmarshal(response.Body(), &bonuses)
	if err != nil {
		return nil, nil, err
	}

	// Récupère le bonus du jour actuel (premier élément du tableau)
	bonusDuJour := bonuses[0]

	return &bonusDuJour, &bonuses, nil
}

func RegisterCommands(s *discordgo.Session) {
	s.AddHandler(handlePingCommand)
}

func handlePingCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if strings.HasPrefix(m.Content, "!ping") {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Pong!")
		bonus, _, _ := getAlmanaxBonus()
		println(bonus)
		_, _ = s.ChannelMessageSendEmbed(m.ChannelID, createEmbed(bonus))
		return
	}
	if strings.HasPrefix(m.Content, "!bonus") {
		args := strings.Split(m.Content, " ")
		if len(args) == 1 {
			bonus, _, _ := getAlmanaxBonus()
			_, _ = s.ChannelMessageSendEmbed(m.ChannelID, createEmbed(bonus))
			return
		}
		dayAhead, err := strconv.Atoi(args[1])
		if err != nil || dayAhead > 4 || dayAhead < 1 {
			_, _ = s.ChannelMessageSend(m.ChannelID, "Jour invalide")
			return
		}
		_, bonuses, _ := getAlmanaxBonus()
		_, _ = s.ChannelMessageSendEmbed(m.ChannelID, createEmbed(&(*bonuses)[dayAhead]))
		return
	}
}
func StartCronScheduler(s *discordgo.Session) {
	c := cron.New()
	_, err := c.AddFunc("@midnight", func() {
		runDailyBonus(s)
	})
	if err != nil {
		println(err)
	}
	c.Start()

}

func runDailyBonus(s *discordgo.Session) {
	bonus, _, _ := getAlmanaxBonus()
	_, _ = s.ChannelMessageSendEmbed("1088206530245038140", createEmbed(bonus))
}
