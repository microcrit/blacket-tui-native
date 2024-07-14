package proxies

import (
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"slices"
	"sync"
	"time"

	"github.com/gbin/goncurses"
)

var URLS = [...]string{
	"https://api.proxyscrape.com/v3/free-proxy-list/get?request=getproxies&proxy_format=ipport&format=text",
	"https://www.proxy-list.download/api/v1/get?type=http",
	"https://raw.githubusercontent.com/TheSpeedX/PROXY-List/master/http.txt",
	"https://raw.githubusercontent.com/sunny9577/proxy-scraper/master/generated/http_proxies.txt",
	"https://raw.githubusercontent.com/proxifly/free-proxy-list/main/proxies/protocols/http/data.txt",
	"https://free-proxy-list.net", // not the safest way of doing this, but this extracts from the "raw list" textarea. :thumbsup:
	"https://raw.githubusercontent.com/ErcinDedeoglu/proxies/main/proxies/http.txt",
	"https://github.com/zloi-user/hideip.me/raw/main/http.txt",
	"https://raw.githubusercontent.com/vakhov/fresh-proxy-list/master/http.txt",
}

func CheckProxy(ip string, port string, shouldStop *bool) bool {
	proxyUrl, err := url.Parse("http://" + ip + ":" + port)
	if err != nil {
		return false
	}
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	client.Timeout = 3 * time.Second
	request := &http.Request{Method: "GET", URL: &url.URL{Host: "httpbin.org", Scheme: "http"}}
	if *shouldStop {
		return false
	}
	_, err = client.Do(request)
	return err == nil
}

var IP_PORT_REGEX = `(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}):(\d{1,5})`

func Log(text string, currentY int, maxY int, maxX int) int {
	if currentY >= maxY-3 {
		currentY = 0
		for i := 0; i < maxY-2; i++ {
			stdscr.MovePrint(i, 0, strings.Repeat(" ", maxX))
		}
	}
	stdscr.MovePrint(currentY+1, 0, text)
	stdscr.Refresh()
	return currentY + 1
}

func LogWorking(maxY int, workingCount int, failedCount int, total int) {
	stdscr.MovePrint(maxY-1, 0, strconv.Itoa(workingCount)+"[✓] / "+strconv.Itoa(failedCount)+"[✗] - "+strconv.Itoa(workingCount+failedCount)+"/"+strconv.Itoa(total))
	stdscr.Refresh()
}

var stdscr *goncurses.Window

func scrapeProxy(chunk []string, channel chan string, logChannel chan string, _ chan int, failedChannel chan int, max int, stopChannel chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	checked := 0
	working := 0
	shouldStop := false
	go func(shouldStop *bool) {
		for range stopChannel {
			logChannel <- "Stopping thread... - Checked " + strconv.Itoa(checked) + " proxies"
			*shouldStop = true
			return
		}
	}(&shouldStop)
	for _, ip := range chunk {
		if shouldStop || working >= max {
			return
		}
		i, p := strings.Split(ip, ":")[0], strings.Split(ip, ":")[1]
		if !shouldStop {
			logChannel <- "- Checking " + ip + " [...]"
		} else {
			return
		}
		good := false
		if !shouldStop {
			good = CheckProxy(i, p, &shouldStop)
		} else {
			return
		}
		if shouldStop || working >= max {
			return
		}
		if good && !shouldStop {
			logChannel <- " ~> " + ip + " is working [✓]"
			channel <- ip
		} else if !shouldStop {
			logChannel <- " ~> " + ip + " is not working [✗]"
			failedChannel <- 1
		} else {
			return
		}
		checked++
	}
}

func chunkSlice(slice []string, chunkSize int64) [][]string {
	var divided [][]string
	groupSize := int64(len(slice)) / chunkSize
	for i := int64(0); i < chunkSize; i++ {
		chunk := slice[i*groupSize : (i+1)*groupSize]
		divided = append(divided, chunk)
	}
	if int64(len(slice)) > chunkSize*groupSize {
		remainder := slice[chunkSize*groupSize:]
		divided[len(divided)-1] = append(divided[len(divided)-1], remainder...)
	}
	if int64(len(divided)) > chunkSize {
		panic("Too many chunks: " + strconv.Itoa(len(divided)))
	}
	return divided
}

func Handler(stdscrx *goncurses.Window, limit int64, threads int64) []string {
	stdscrx.Clear()

	maxY, maxX := stdscrx.MaxYX()
	currentY := 0
	workingCount := 0
	failedCount := 0

	stdscr = stdscrx

	ips := []string{}

	compiledRgx := regexp.MustCompile(IP_PORT_REGEX)

	wg := sync.WaitGroup{}

	for _, url := range URLS {
		resp, err := http.Get(url)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}
		ipsInner := compiledRgx.FindAllString(string(body), -1)
		ips = append(ips, ipsInner...)
		for _, ip := range ipsInner {
			currentY = Log(ip, currentY, maxY, maxX)
		}
	}

	slices.Sort(ips)
	ips = slices.Compact(ips)

	results := []string{}
	chunked := chunkSlice(ips, threads)
	currentY = Log("Starting to check "+strconv.Itoa(len(ips))+" proxies...", currentY, maxY, maxX)
	currentY = Log("Chunks: "+strconv.Itoa(len(chunked))+" | Limit: "+strconv.Itoa(int(limit)), currentY, maxY, maxX)
	ipLen := 0
	stdscrx.MovePrint(maxY-2, 0, strconv.Itoa(len(chunked))+" chunks")
	for _, chunk := range chunked {
		ipLen += len(chunk)
	}

	channel := make(chan string)
	logChannel := make(chan string)
	workingChannel := make(chan int)
	failedChannel := make(chan int)

	stopChannel := make(chan bool)

	for i, chunk := range chunked {
		Log("Starting thread "+strconv.Itoa(i), currentY, maxY, maxX)
		go scrapeProxy(chunk, channel, logChannel, workingChannel, failedChannel, int(limit), stopChannel, &wg)
	}
	wg.Add(len(chunked))

	stopped := 0

	go func(stopped *int) {
		for {
			select {
			case log := <-logChannel:
				currentY = Log(log, currentY, maxY, maxX)

			case failed := <-failedChannel:
				failedCount += failed
				LogWorking(maxY, workingCount, failedCount, ipLen)

				if workingCount >= int(limit) {
					stdscrx.MovePrint(0, 0, "[Event] Stopping threads...")
					stdscrx.Refresh()
					stopChannel <- true
				}

			case <-stopChannel:
				*stopped++
				if *stopped >= len(chunked)+1 {
					stdscrx.MovePrint(0, 0, "[Event] Threads stopped.")
					stdscrx.Refresh()
				} else {
					stdscrx.MovePrint(0, 0, "[Event] Stopping threads... ("+strconv.Itoa(*stopped)+"/"+strconv.Itoa(len(chunked)+1)+")")
					stdscrx.Refresh()
				}

			case resultsX := <-channel:
				results = append(results, resultsX)
				workingCount++
				LogWorking(maxY, workingCount, failedCount, ipLen)

				if workingCount >= int(limit) {
					stdscrx.MovePrint(0, 0, "[Event] Stopping threads...")
					stdscrx.Refresh()
					stopChannel <- true
				}
			}
		}
	}(&stopped)

	wg.Wait()

	return results[:limit]
}
