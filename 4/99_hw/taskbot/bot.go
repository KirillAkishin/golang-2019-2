package main

// сюда писать код

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

var (
	// @BotFather в телеграме даст вам это
	BotToken = "742388923:AAEp0Vhz1hkXe45iZNNibiXbM4iYwvFgEMc"

	// урл выдаст вам игрок или хероку
	WebhookURL = "https://telegram-bot-test1908.herokuapp.com"
)

// type Сommands struct {
// 	Сommands []string
// }

type TaskDescriotion struct {
	TaskName    string
	HolderName  string
	HolderID    int
	CreatorName string
	CreatorID   int
	IsActive    bool
}

type TasksData struct {
	Tasks         *[]TaskDescriotion
	CountOfTasks  int
	NumberOfTasks int
	Bot           *tgbotapi.BotAPI
}

func (tsData *TasksData) TasksCommand(update tgbotapi.Update) {
	currentID := update.Message.Chat.ID
	answer := ""

	if tsData.NumberOfTasks == 0 {
		answer = "Нет задач"
	} else {
		stringsOfTasks := make([]string, 0)
		for idOfTask, task := range *tsData.Tasks {
			if task.IsActive {
				idOfTask++
				strOfTask := fmt.Sprintf("%d. %s by @%s\n", idOfTask, task.TaskName, task.CreatorName)
				switch task.HolderID {
				case 0:
					strOfTask += fmt.Sprintf("/assign_%d", idOfTask)
				case int(currentID):
					strOfTask += fmt.Sprintf("assignee: я\n/unassign_%d /resolve_%d", idOfTask, idOfTask)
				default:
					strOfTask += fmt.Sprintf("assignee: @%s", task.HolderName)
				}
				stringsOfTasks = append(stringsOfTasks, strOfTask)
			}
		}
		answer = strings.Join(stringsOfTasks, "\n\n")
	}
	tsData.Bot.Send(tgbotapi.NewMessage(currentID, answer))
}

func (tsData *TasksData) NewCommand(update tgbotapi.Update, command string) {
	currentID := update.Message.Chat.ID
	answer := ""
	if len(command) <= 5 {
		log.Println("syntax error:", command)
		answer = fmt.Sprintf("syntax error")
		tsData.Bot.Send(tgbotapi.NewMessage(currentID, answer))
		return
	}
	command = command[5:]

	tsData.CountOfTasks++
	tsData.NumberOfTasks++
	task := TaskDescriotion{
		TaskName: command,
		// NameOfHolder:  update.Message.From.UserName,
		HolderID:    0,
		CreatorName: update.Message.From.UserName,
		CreatorID:   update.Message.From.ID,
		IsActive:    true,
	}
	*tsData.Tasks = append(*tsData.Tasks, task)
	answer = fmt.Sprintf("Задача \"%s\" создана, id=%d", command, tsData.CountOfTasks)
	tsData.Bot.Send(tgbotapi.NewMessage(currentID, answer))
}

func (tsData *TasksData) AssignCommand(update tgbotapi.Update, command string) {
	currentID := update.Message.Chat.ID
	answer := ""

	taskID, err := strconv.Atoi(command[8:])
	if (err != nil) || (tsData.CountOfTasks < taskID) {
		log.Println("syntax error:", command)
		answer = fmt.Sprintf("syntax error")
		tsData.Bot.Send(tgbotapi.NewMessage(currentID, answer))
		return
	}
	taskID--

	tasks := *tsData.Tasks
	taskName := tasks[taskID].TaskName
	creatorID := tasks[taskID].CreatorID
	currentName := update.Message.Chat.UserName

	answer = fmt.Sprintf("Задача \"%s\" назначена на вас", taskName)
	if int(currentID) != creatorID {
		answerToPrevHolder := fmt.Sprintf("Задача \"%s\" назначена на @%s", taskName, currentName)
		prevHolderID := int64(tasks[taskID].HolderID)
		if prevHolderID == 0 {
			prevHolderID = int64(tasks[taskID].CreatorID)
		}
		tsData.Bot.Send(tgbotapi.NewMessage(prevHolderID, answerToPrevHolder))
	}
	tasks[taskID].HolderID = int(currentID)
	tasks[taskID].HolderName = currentName
	tsData.Bot.Send(tgbotapi.NewMessage(currentID, answer))
}

func (tsData *TasksData) UnassignCommand(update tgbotapi.Update, command string) {
	currentID := update.Message.Chat.ID
	answer := ""

	taskID, err := strconv.Atoi(command[10:])
	if (err != nil) || (tsData.CountOfTasks < taskID) {
		log.Println("syntax error:", command)
		answer = fmt.Sprintf("syntax error")
		tsData.Bot.Send(tgbotapi.NewMessage(currentID, answer))
		return
	}
	taskID--

	tasks := *tsData.Tasks
	if tasks[taskID].HolderID == update.Message.From.ID {
		tasks[taskID].HolderID = 0
		answer = "Принято"
		answerToCreator := fmt.Sprintf("Задача \"%s\" осталась без исполнителя", tasks[taskID].TaskName)
		creatorID := int64(tasks[taskID].CreatorID)
		tsData.Bot.Send(tgbotapi.NewMessage(creatorID, answerToCreator))
	} else {
		answer = "Задача не на вас"
	}
	tsData.Bot.Send(tgbotapi.NewMessage(currentID, answer))
}

func (tsData *TasksData) ResolveCommand(update tgbotapi.Update, command string) {
	currentID := update.Message.Chat.ID
	answer := ""
	taskID, err := strconv.Atoi(command[9:])
	if (err != nil) || (tsData.CountOfTasks < taskID) {
		log.Println("syntax error:", command)
		answer = fmt.Sprintf("syntax error")
		tsData.Bot.Send(tgbotapi.NewMessage(currentID, answer))
		return
	}
	taskID--

	tasks := *tsData.Tasks
	answer = fmt.Sprintf("Задача \"%s\" выполнена", tasks[taskID].TaskName)
	tsData.Bot.Send(tgbotapi.NewMessage(currentID, answer))

	creatorID := int64(tasks[taskID].CreatorID)
	currentName := update.Message.From.UserName
	answerToCreator := fmt.Sprintf("Задача \"%s\" выполнена @%s", tasks[taskID].TaskName, currentName)
	tsData.Bot.Send(tgbotapi.NewMessage(creatorID, answerToCreator))

	tsData.NumberOfTasks--
	tasks[taskID].IsActive = false
}

func (tsData *TasksData) MyCommand(update tgbotapi.Update) {
	currentID := update.Message.Chat.ID
	answer := ""
	stringsToAnswer := make([]string, 0)
	for taskID, task := range *tsData.Tasks {
		if (task.HolderID == int(currentID)) && (task.IsActive) {
			taskID++
			strOfTask := fmt.Sprintf("%d. %s by @%s\n/unassign_%d /resolve_%d", taskID, task.TaskName, task.CreatorName, taskID, taskID)
			stringsToAnswer = append(stringsToAnswer, strOfTask)
		}
	}
	answer = strings.Join(stringsToAnswer, "\n\n")
	tsData.Bot.Send(tgbotapi.NewMessage(currentID, answer))
}

func (tsData *TasksData) OwnerCommand(update tgbotapi.Update) {
	currentID := update.Message.Chat.ID
	answer := ""
	stringsToAnswer := make([]string, 0)
	for taskID, task := range *tsData.Tasks {
		if (task.CreatorID == int(currentID)) && (task.IsActive) {
			taskID++
			strOfTask := fmt.Sprintf("%d. %s by @%s\n/assign_%d", taskID, task.TaskName, task.CreatorName, taskID)
			stringsToAnswer = append(stringsToAnswer, strOfTask)
		}
	}
	answer = strings.Join(stringsToAnswer, "\n\n")
	tsData.Bot.Send(tgbotapi.NewMessage(currentID, answer))
}

func startTaskBot(ctx context.Context) error {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		return errors.New("Can't connect to Telegram")
	}
	log.Println("Authorized on account", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(WebhookURL))
	if err != nil {
		return errors.New("Can't connect to WebhookURL")
	}
	updates := bot.ListenForWebhook("/")

	// port := os.Getenv("PORT")
	port := "8081"
	go http.ListenAndServe(":"+port, nil)
	log.Println("start listen:", port)

	tsData := &TasksData{}
	tsData.Tasks = new([]TaskDescriotion)
	tsData.CountOfTasks = 0
	tsData.NumberOfTasks = 0
	tsData.Bot = bot
	for update := range updates {
		command := update.Message.Text
		switch {
		case strings.HasPrefix(command, "/tasks"):
			tsData.TasksCommand(update)
		case strings.HasPrefix(command, "/new"):
			tsData.NewCommand(update, command)
		case strings.HasPrefix(command, "/assign_"):
			tsData.AssignCommand(update, command)
		case strings.HasPrefix(command, "/unassign_"):
			tsData.UnassignCommand(update, command)
		case strings.HasPrefix(command, "/resolve_"):
			tsData.ResolveCommand(update, command)
		case strings.HasPrefix(command, "/my"):
			tsData.MyCommand(update)
		case strings.HasPrefix(command, "/owner"):
			tsData.OwnerCommand(update)
		}
	}
	return nil
}

func main() {
	ctx := context.Background()
	err := startTaskBot(ctx)
	if err != nil {
		log.Println("FATAL ERROR:", err)
	}
}
