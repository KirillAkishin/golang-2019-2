package main

import (
	//"fmt"
	"strings"
)

type Person struct {
	Id        int
	Name      string
	Address   string
	Inventory string
}

var player = Person{}

func (p *Person) SetName(name string) {
	p.Name = name
}

type Room struct {
	Name             string
	Description      func(text string) (result string)
	ShortDescription string
	Neighbours       string
	Inventory        string
	ClosedDoors      string
}

func (r *Room) Move(strs []string) string {
	if len(strs) != 2 {
		return "не могу туда пройти"
	}
	destination := strs[1]
	if strings.Contains(r.Neighbours, destination) {
		if strings.Contains(r.ClosedDoors, destination) {
			return "дверь закрыта"
		}
		player.Address = destination
		return world[destination].ShortDescription
	} else {
		return "нет пути в " + destination
	}
}

func (r *Room) PutOn(strs []string) string {
	if len(strs) != 2 {
		return "не могу надеть это"
	}
	thing := strs[1]
	object := stuff[thing]
	if strings.Contains(r.Inventory, thing) {
		if !object.outfited {
			return "это нельзя надеть"
		}
		player.Inventory = player.Inventory + " " + thing
		r.Inventory = strings.Replace(r.Inventory, thing, "", -1)
		object.Address = player.Name
		return "вы надели: " + thing
	} else {
		return "не могу надеть " + thing
	}
}

func (r *Room) Take(strs []string) string {
	if len(strs) != 2 {
		return "нет такого"
	}
	thing := strs[1]
	object := stuff[thing]
	if strings.Contains(r.Inventory, thing) {
		if !strings.Contains(player.Inventory, "рюкзак") {
			return "некуда класть"
		}
		player.Inventory = player.Inventory + " " + thing
		r.Inventory = strings.Replace(r.Inventory, thing, "", -1)
		object.Address = player.Name
		return "предмет добавлен в инвентарь: " + thing
	} else {
		return "нет такого"
	}
}

func (r *Room) Apply(thing1 string, thing2 string) string {
	subject := stuff[thing1]
	object := stuff[thing2]
	if strings.Contains(player.Inventory, thing1) {
		if !strings.Contains(subject.applicable, thing2) {
			return "не к чему применить"
		}
		if object.closed != "" { //	если объект класса Дверь, то
			r.ClosedDoors = strings.Replace(r.ClosedDoors, object.closed, "", -1)
		}
		return object.appMessage
	} else {
		return "нет предмета в инвентаре - " + thing1
	}
}

var world = map[string]*Room{}

type Object struct {
	Name       string
	Address    string
	outfited   bool // предметы к которым можно применить "надеть"
	applicable string
	appMessage string
	closed     string
}

var stuff = map[string]*Object{}

func main() {
	return
}

func initGame() {
	world = map[string]*Room{
		"кухня": &Room{
			Name: "кухня",
			Description: func(string) (result string) {
				text := "собрать рюкзак и "
				if strings.Contains(player.Inventory, "рюкзак") {
					text = ""
				}
				result = "ты находишься на кухне, на столе: чай, надо " + text + "идти в универ. можно пройти - коридор"
				return result
			},
			ShortDescription: "кухня, ничего интересного. можно пройти - коридор",
			Neighbours:       "коридор",
		},
		"коридор": &Room{
			Name: "коридор",
			Description: func(string) string {
				return "Nil"
			},
			ShortDescription: "ничего интересного. можно пройти - кухня, комната, улица",
			Neighbours:       "кухня комната улица",
			ClosedDoors:      "улица",
		},
		"комната": &Room{
			Name: "комната",
			Description: func(invOfRoom string) (result string) {
				result = "на столе: "
				if !strings.Contains(invOfRoom, "конспекты") && !strings.Contains(invOfRoom, "рюкзак") && !strings.Contains(invOfRoom, "ключи") {
					return "пустая комната. можно пройти - коридор"
				}
				if strings.Contains(invOfRoom, "ключи") {
					result = result + "ключи"
				}
				if strings.Contains(invOfRoom, "ключи") && strings.Contains(invOfRoom, "конспекты") {
					result = result + ", "
				}
				if strings.Contains(invOfRoom, "конспекты") {
					result = result + "конспекты"
				}
				if (strings.Contains(invOfRoom, "конспекты") || strings.Contains(invOfRoom, "ключи")) && strings.Contains(invOfRoom, "рюкзак") {
					result = result + ", "
				}
				if strings.Contains(invOfRoom, "рюкзак") {
					result = result + "на стуле: рюкзак"
				}
				result = result + ". можно пройти - коридор"
				return result
			},
			ShortDescription: "ты в своей комнате. можно пройти - коридор",
			Neighbours:       "коридор",
			Inventory:        "ключи конспекты рюкзак",
		},
		"улица": &Room{
			Name: "улица",
			Description: func(text string) string {
				return "Nil"
			},
			ShortDescription: "на улице весна. можно пройти - домой",
			Neighbours:       "коридор",
		},
	}

	player = Person{
		Id:      0,
		Address: "кухня",
	}

	stuff = map[string]*Object{
		"рюкзак": &Object{
			Name:     "рюкзак",
			Address:  "комната",
			outfited: true,
		},
		"ключи": &Object{
			Name:       "ключи",
			Address:    "комната",
			outfited:   false,
			applicable: "дверь",
		},
		"конспекты": &Object{
			Name:     "конспекты",
			Address:  "комната",
			outfited: false,
		},
		"дверь": &Object{
			Name:       "дверь",
			Address:    "коридор",
			outfited:   false,
			appMessage: "дверь открыта",
			closed:     "улица",
		},
	}
}

func handleCommand(command string) string {
	strs := strings.Split(command, " ")
	//fmt.Println(">>>", command)
	switch strs[0] {
	case "осмотреться":
		{
			return world[player.Address].Description(world[player.Address].Inventory)
		}
	case "идти":
		{
			return world[player.Address].Move(strs)
		}
	case "применить":
		{
			thing1 := strs[1]
			thing2 := strs[2]
			return world[player.Address].Apply(thing1, thing2)
		}
	case "взять":
		{
			return world[player.Address].Take(strs)
		}
	case "надеть":
		{
			return world[player.Address].PutOn(strs)
		}
	}
	return "неизвестная команда"
}
