package main

import (
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/gbin/goncurses"
	"github.com/pelletier/go-toml"

	"crit.rip/blacket-tui/api/objects/user"
	confParser "crit.rip/blacket-tui/config"
	"crit.rip/blacket-tui/providers/chat"
	"crit.rip/blacket-tui/providers/packs"
	"crit.rip/blacket-tui/providers/proxies"

	"crit.rip/blacket-tui/ui"
)

func LayoutLogin(stdscr *goncurses.Window) (uname string, pwd string) {
	totalY, totalX := stdscr.MaxYX()
	centerY, centerX := totalY/2, totalX/2

	uLength, pLength := len("Username: "), len("Password: ")
	tLength := len("Enter a username or ID and password to login.")

	stdscr.MovePrint(centerY-3, centerX-(tLength/2), "Enter a username or ID and password to login.")
	username := ui.TextBox(stdscr, centerY-2, centerX-(uLength/2), "Username: ")
	password := ui.PasswordBox(stdscr, centerY-1, centerX-(pLength/2), "Password: ")

	stdscr.Clear()

	return username, password
}

type Choice struct {
	Text   string
	Action func()
}

const ESCAPE = 27
const UP_ARROW = 259
const DOWN_ARROW = 258

func SelectLayout(stdscr *goncurses.Window, title string, choices []Choice) int {
	choicesTexts := []string{}
	for _, choice := range choices {
		choicesTexts = append(choicesTexts, choice.Text)
	}
	choicesTexts = append(choicesTexts, "Exit")

	selected := RawSelect(stdscr, title, choicesTexts)
	if selected == -1 {
		return 0
	} else if selected == len(choicesTexts)-1 {
		return -1
	}
	choices[selected].Action()
	return 1
}

func RawSelect(stdscr *goncurses.Window, title string, choices []string) int {
	_, totalX := stdscr.MaxYX()

	stdscr.MovePrint(1, 2, title)
	for i, choice := range choices {
		stdscr.MovePrint(2+i, 2, "  "+choice)
	}

	stdscr.Refresh()

	selected := 0
	for {
		for i, choice := range choices {
			stdscr.MovePrint(2+i, 2, strings.Repeat(" ", totalX))
			if i == selected {
				stdscr.MovePrint(2+i, 2, "> "+choice)
			} else {
				stdscr.MovePrint(2+i, 2, "  "+choice)
			}
		}
		stdscr.Refresh()

		ch := stdscr.GetChar()
		if ch == 'q' || ch == ESCAPE {
			return -1
		} else if ch == UP_ARROW {
			if selected > 0 {
				selected--
			}
		} else if ch == DOWN_ARROW {
			if selected < len(choices)-1 {
				selected++
			}
		} else if ch == '\n' {
			break
		}
	}

	stdscr.Clear()
	stdscr.Refresh()

	return selected
}

func Layout(config map[string]interface{}) string {
	stdscr, err := goncurses.Init()
	if err != nil {
		panic(err)
	}

	defer goncurses.End()

	stdscr.Clear()
	stdscr.Keypad(true)
	stdscr.Timeout(-1)

	maxY, _ := stdscr.MaxYX()

	selectedNew := true
	selectedAccount := -1

	var accounts []interface{}
	if config["Accounts"] != nil {
		var a []interface{}
		if len(config["Accounts"].([]interface{})) > 0 {
			a = config["Accounts"].([]interface{})
		}
		accounts = a
	}

	var account interface{}
	resulting := []string{}
	for selectedNew {
		accountNames := []string{}
		for _, account := range accounts {
			rau := account.(map[string]interface{})["Username"]

			if rau != "" {
				accountNames = append(accountNames, rau.(string))
			}
		}

		accountNames = append(accountNames, "+ New")
		resulting = append(accountNames, "Exit")
		selectedAccount = RawSelect(stdscr, "Select or configure your controller (master) account.", resulting)
		if selectedAccount == -1 {
			continue
		}
		if selectedAccount == len(resulting)-2 {
			username, password := LayoutLogin(stdscr)
			var ax map[string]interface{} = map[string]interface{}{
				"Username": username,
				"Password": password,
			}
			accounts = append(accounts, interface{}(ax))
			a := config["Accounts"]
			if a == nil {
				a = []interface{}{}
			}
			for i, account := range a.([]interface{}) {
				if account.(map[string]interface{})["Username"] == username {
					stdscr.MovePrint(maxY, 0, "Account already exists.")
					stdscr.Refresh()
					selectedAccount = i
					selectedNew = false
					break
				}
			}
			if !selectedNew {
				break
			}
			config["Accounts"] = append(a.([]interface{}), ax)
			str, err := toml.Marshal(config)
			if err != nil {
				panic(err)
			}
			err = os.WriteFile("config.toml", []byte(str), fs.ModePerm)
			if err != nil {
				panic(err)
			}
			selectedNew = true
		} else if selectedAccount == len(resulting)-1 {
			return "Exiting..."
		} else {
			selectedNew = false
		}
	}
	if selectedAccount == -1 || selectedAccount == len(resulting)-1 {
		return "Exiting..."
	}
	account = (accounts)[selectedAccount]

	raux := account.(map[string]interface{})["Username"]

	rawx := account.(map[string]interface{})["Username"]

	if raux == "" || rawx == "" {
		panic("Invalid account")
	}

	token := user.Login(raux.(string), rawx.(string))
	user := user.GetUser(stdscr, token)

	defaultFile := "proxies.txt"
	defaultMax := 10
	defaultThreads := 5

	var ps map[string]interface{} = map[string]interface{}{
		"File":    &defaultFile,
		"Max":     &defaultMax,
		"Threads": &defaultThreads,
	}
	if config["ProxyScraper"] != nil {
		ps = config["ProxyScraper"].(map[string]interface{})
	}

	for {
		stdscr.Clear()
		stdscr.Refresh()

		action := SelectLayout(stdscr, "Select an action", []Choice{
			{
				Text: "Scrape Proxies - Controller account",
				Action: func() {
					maxProxies := ps["Max"].(int64)
					proxyPath := ps["File"].(string)
					maxThreads := ps["Threads"].(int64)
					proxyList := proxies.Handler(stdscr, maxProxies, maxThreads)

					fi, err := os.OpenFile(proxyPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fs.ModePerm)

					if err != nil {
						panic(err)
					}

					defer fi.Close()

					fi.Write([]byte(strings.Join(proxyList, "\n")))
				},
			},
			{
				Text: "Open Chat - Controller account",
				Action: func() {
					chat.Handler(token, user)
				},
			},
			{
				Text: "Open Packs - Sub accounts",
				Action: func() {
					op := config["PackOpener"].(map[string]interface{})
					accounts := config["Accounts"].([]interface{})
					maxThreads := 5
					if op["Threads"] != nil {
						maxThreads = int(op["Threads"].(int64))
					}
					proxyFile := config["ProxyScraper"].(map[string]interface{})["File"].(string)
					content, err := os.ReadFile(proxyFile)
					if err != nil {
						panic(err)
					}
					proxies := strings.Split(string(content), "\n")
					realAccounts := []map[string]string{}
					for _, account := range accounts {
						if account.(map[string]interface{})["Username"] != raux {
							realAccounts = append(realAccounts, map[string]string{
								"Username": account.(map[string]interface{})["Username"].(string),
								"Password": account.(map[string]interface{})["Password"].(string),
							})
						}
					}
					packs.Handler(stdscr, realAccounts, int64(maxThreads), proxies)
				},
			},
		})

		if action == -1 || action == 0 {
			break
		}
	}

	stdscr.Refresh()

	return ""
}

func main() {
	config := confParser.ParseConfig("config.toml")
	resp := Layout(config)
	if resp != "" {
		log.Fatalln(resp)
	}
}
