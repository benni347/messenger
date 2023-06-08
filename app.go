package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	utils "github.com/benni347/messengerutils"
	"github.com/joho/godotenv"
	"github.com/pusher/pusher-http-go/v5"
	amqp "github.com/rabbitmq/amqp091-go"
)

// App struct
type App struct {
	ctx     context.Context
	verbose bool
	config  Config
	user    User
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

type User struct {
	queueName string
	ch        *amqp.Channel
}

type Config struct {
	AppId            string `json:"appId"`
	AppSecret        string `json:"appSecret"`
	AppKey           string `json:"appKey"`
	ClusterId        string `json:"clusterId"`
	SupaBaseApiKey   string `json:"supaBaseApiKey"`
	SupaBaseUrl      string `json:"supaBaseUrl"`
	RabbitMqAdmin    string `json:"rabbitMqAdmin"`
	RabbitMqPassword string `json:"rabbitMqPassword"`
	RabbitMqHost     string `json:"rabbitMqHost"`
}

// RetrieveEnvValues retrieves the values from the .env file
// and saves them in the config struct
// Returns the values of the .env file
// in the following order:
// app_id, app_secret, app_key, cluster_id
func (a *App) RetrieveEnvValues() Config {
	m := &utils.MessengerUtils{
		Verbose: a.verbose,
	}
	m.PrintInfo("Loading .env file")
	err := godotenv.Load()
	if err != nil {
		utils.PrintError("Error loading .env file", err)
		return Config{}
	}
	a.config.AppId = os.Getenv("APP_ID")
	a.config.AppSecret = os.Getenv("SECRET")
	a.config.AppKey = os.Getenv("KEY")
	a.config.ClusterId = os.Getenv("CLUSTER")
	a.config.SupaBaseApiKey = os.Getenv("SUPABASE_API_KEY")
	a.config.SupaBaseUrl = os.Getenv("SUPABASE_URL")
	a.config.RabbitMqAdmin = os.Getenv("RABBITMQ_ADMIN")
	a.config.RabbitMqPassword = os.Getenv("RABBITMQ_PASSWORD")
	a.config.RabbitMqHost = os.Getenv("RABBITMQ_HOST")
	return a.config
}

func (a *App) ValidateEmail(email string) bool {
	return utils.ValidateEmailRegex(email)
}

func (a *App) GetAppId() string {
	return a.config.AppId
}

func (a *App) GetAppSecret() string {
	return a.config.AppSecret
}

func (a *App) GetAppKey() string {
	return a.config.AppKey
}

func (a *App) GetClusterId() string {
	return a.config.ClusterId
}

func (a *App) GetSupaBaseApiKey() string {
	return a.config.SupaBaseApiKey
}

func (a *App) GetSupaBaseUrl() string {
	return a.config.SupaBaseUrl
}

func (a *App) GetRabbitMqAdmin() string {
	return a.config.RabbitMqAdmin
}

func (a *App) GetRabbitMqPassword() string {
	return a.config.RabbitMqPassword
}

func (a *App) GetRabbitMqHost() string {
	return a.config.RabbitMqHost
}

type Message struct {
	Message string `json:"message"`
	Sender  string `json:"sender"`
	Time    string `json:"time"`
}

// SendMessage sends a message from a specific sender to a specified chat room.
//
// This function takes in the ID of the chat room (chatRoomId), the sender's identifier (sender), and the content of the message (message) as parameters.
// It then generates a timestamp, creates a Message object with the sender, message, and timestamp, and converts this object into a JSON format.
// Afterward, it gets the necessary Pusher credentials for this application, creates a new Pusher client using these credentials,
// and triggers a new Pusher event with the JSON message.
//
// The format of the message from the server should be: {"message": "message", "sender": "sender", "time": "time"}
//
// If there's an error while marshalling the Message object to JSON, this function prints the error and returns immediately.
//
// Note that SendMessage does not return any values.
//
// Usage:
// app := NewApp()
// app.SendMessage("chatRoomId", "senderId", "Hello, world!")
func (a *App) SendMessage(chatRoomId string, sender string, message string) {
	// The format from the server should be: {"message": "message", "time": "time"}
	currentTime := time.Now().UnixNano()
	currentTimeString := strconv.FormatInt(currentTime, 10)

	// Create a new Message object
	msg := Message{
		Message: message,
		Time:    currentTimeString,
		Sender:  sender,
	}

	// Convert Message object to JSON
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		utils.PrintError("marshalling JSON", err)
		return
	}

	// Get the Pusher credentials
	appId := a.GetAppId()
	appSecret := a.GetAppSecret()
	appKey := a.GetAppKey()
	clusterId := a.GetClusterId()

	// Create a new Pusher Client
	pusherClient := pusher.Client{
		AppID:   appId,
		Key:     appKey,
		Secret:  appSecret,
		Cluster: clusterId,
		Secure:  true,
	}

	// Create a new Pusher trigger
	pusherClient.Trigger(chatRoomId, "message", string(msgJSON))
}

func failOnError(err error, msg string) {
	if err != nil {
		utils.PrintError(msg, err)
	}
}

func (a *App) Send(message, chatRoomId string) {
	var amqpHost string
	var amqpPort string
	var amqpUser string
	var amqpPassword string

	amqpHost = a.GetRabbitMqHost()
	amqpPort = "5672"
	amqpUser = a.GetRabbitMqAdmin()
	amqpPassword = a.GetRabbitMqPassword()

	amqpUrl := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		amqpUser,
		url.QueryEscape(amqpPassword),
		amqpHost,
		amqpPort,
	)

	conn, err := amqp.Dial(amqpUrl)
	failOnError(err, "Failed to connect to RabbitMQ")

	channel, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	queueName := chatRoomId
	queue, err := channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	body := message

	err = channel.Publish(
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")

	defer func() {
		if err := conn.Close(); err != nil {
			utils.PrintError("closing connection", err)
		}
		if err := channel.Close(); err != nil {
			utils.PrintError("closing channel", err)
		}
	}()

	time.Sleep(1 * time.Second)
}

func (a *App) Receive(channelId string) <-chan string {
	var amqpHost string
	var amqpPort string
	var amqpUser string
	var amqpPassword string

	amqpHost = a.GetRabbitMqHost()
	amqpPort = "5672"
	amqpUser = a.GetRabbitMqAdmin()
	amqpPassword = a.GetRabbitMqPassword()

	amqpUrl := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		amqpUser,
		url.QueryEscape(amqpPassword),
		amqpHost,
		amqpPort,
	)

	conn, err := amqp.Dial(amqpUrl)
	failOnError(err, "Failed to connect to RabbitMQ")

	channel, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	queueName := channelId
	queue, err := channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to declare a queue")

	msgs, err := channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to register a consumer")

	out := make(chan string)

	go func() {
		for d := range msgs {
			out <- string(d.Body)
		}
		close(out)
	}()

	return out
}

func (a *App) ReciveFormatForJs(chatRoomId string) string {
	messages := a.Receive(chatRoomId)

	for message := range messages {
		return message
	}
	return ""
}

func intToSpecificBaseToString(num, base int) string {
	const charset = "0123456789abcdefghijklmnopqrstuvwxyz"
	if base < 2 || base > 36 {
		return ""
	}

	if num == 0 {
		return "0"
	}

	result := ""
	isNegative := num < 0

	if isNegative {
		num = -num
	}

	for num > 0 {
		result = string(charset[num%base]) + result
		num /= base
	}

	if isNegative {
		result = "-" + result
	}

	return result
}

func (a *App) GenerateUserName(length int) string {
	rand.Seed(time.Now().UnixNano())
	const charset = "0123456789abcdefghijklmnopqrstuvwxyz"
	basis := 36

	random_number := rand.Intn(int(math.Pow(float64(basis), float64(length))))

	randomNumberStringInCorrectBase := intToSpecificBaseToString(random_number, basis)

	if len(randomNumberStringInCorrectBase) < length {
		// If the random string is too short, pad it with random characters
		for len(randomNumberStringInCorrectBase) < length {
			randomChar := charset[rand.Intn(len(charset))]
			randomNumberStringInCorrectBase = randomNumberStringInCorrectBase + string(randomChar)
		}
	}

	randomStr := randomNumberStringInCorrectBase[len(randomNumberStringInCorrectBase)-length:]

	user_name_string := "user_" + randomStr
	return user_name_string
}

func (a *App) CreateChatRoomId(otherId, currentId string) string {
	var chatRoomId string
	if otherId < currentId {
		chatRoomId = otherId + currentId
	} else {
		chatRoomId = currentId + otherId
	}
	return chatRoomId
}

func (a *App) SetQueuName(queueName string) {
	a.user.queueName = queueName
}

func (a *App) getQueueName() string {
	return a.user.queueName
}

func (a *App) GetOtherUserId(concatenatedUUIDs string, myUUID string) string {
	// Remove dashes from the UUIDs
	myUUID = strings.ReplaceAll(myUUID, "-", "")
	concatenatedUUIDs = strings.ReplaceAll(concatenatedUUIDs, "-", "")

	// Retrieve the other UUID
	otherUUID := strings.Replace(concatenatedUUIDs, myUUID, "", 1)

	// Add dashes back to the other UUID
	otherUUID = otherUUID[0:8] + "-" + otherUUID[8:12] + "-" + otherUUID[12:16] + "-" + otherUUID[16:20] + "-" + otherUUID[20:]

	return otherUUID
}
