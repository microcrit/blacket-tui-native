package packs

import (
	"crypto/rand"
	"math/big"
	"strings"
	"sync"

	"github.com/gbin/goncurses"

	"crit.rip/blacket-tui/api/objects/user"
)

var stdscr *goncurses.Window

func openPack(wg *sync.WaitGroup, packName string, name string, password string, proxy interface{}, tokenEndChannel chan bool, userInfoChannel chan []string) {
	defer wg.Done()

	realProxy := ""
	if proxy != nil {
		realProxy = proxy.(string)
	}

	token := user.LoginProxy(name, password, realProxy)

	ranOutOfTokens := false
	for !ranOutOfTokens {
		resp := user.OpenPackProxy(token, packName, realProxy)
		if resp.Error {
			ranOutOfTokens = true
			tokenEndChannel <- true
			break
		}
		userInfoChannel <- []string{name, resp.Blook}
	}
}

func PackSelectMenu() string {
	stdscr.Clear()
	stdscr.MovePrint(0, 0, "Enter a pack name: ")
	stdscr.Refresh()
	result := ""
	for {
		ch := stdscr.GetChar()
		if ch == goncurses.KEY_ENTER || ch == 10 {
			break
		}
		if ch == goncurses.KEY_BACKSPACE || ch == 127 {
			if len(result) > 0 {
				result = result[:len(result)-1]
			}
		} else {
			result += string(ch)
		}
		stdscr.MovePrint(0, 0, "Enter a pack name: "+result)
		stdscr.Refresh()
	}
	stdscr.Clear()
	stdscr.Refresh()
	return result
}

func LogAccountInfo(accountInfo map[string]string) {
	_, maxX := stdscr.MaxYX()
	stdscr.Clear()
	i := 0
	for i := 0; i < len(accountInfo); i++ {
		stdscr.MovePrint(i+1, 0, strings.Repeat(" ", maxX))
	}
	for k, v := range accountInfo {
		stdscr.MovePrint(i+1, 0, k+" - "+v)
		i++
	}
	stdscr.Refresh()
}

func Handler(stdscrx *goncurses.Window, accounts []map[string]string, threads int64, proxies []string) []string {
	stdscrx.Clear()
	stdscr = stdscrx

	pack := PackSelectMenu()

	wg := sync.WaitGroup{}

	tokenEndChannel := make(chan bool)
	userInfoChannel := make(chan []string)

	for _, account := range accounts {
		wg.Add(1)
		r, err := rand.Int(rand.Reader, big.NewInt(int64(len(proxies))))
		if err != nil {
			panic(err)
		}
		username, password := account["Username"], account["Password"]
		go openPack(&wg, pack, username, password, proxies[r.Int64()], tokenEndChannel, userInfoChannel)
	}

	userInfoM := map[string]string{}

	go func() {
		for {
			select {
			case <-tokenEndChannel:
				wg.Done()
			case userInfo := <-userInfoChannel:
				userInfoM[userInfo[0]] = userInfo[1]
				LogAccountInfo(userInfoM)
			}
		}
	}()
	wg.Add(1)

	wg.Wait()

	return []string{}
}
