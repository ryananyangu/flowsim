package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

type IScreen interface {
	display(lang string) string // Display the current screen
	nextPage(input string) int  // Gets the input from user and return index of what to display next
	screenType() string         // type of the struct that has implemented the interface
	// previousPage(input string) int // Get Previous page
}

type CoreScreen struct {
	Name        string
	Header      string
	IsEnd       bool
	BackEnabled bool
	ExitEnabled bool
}

type ListScreen struct {
	CoreScreen
	NextLocation uint
	Options      []string
}

type Route struct {
	Value        string
	NextLocation uint
	IsEnd        bool
}

type RouteScreen struct {
	CoreScreen
	Routes []Route
}

type RawInputScreen struct {
	CoreScreen
	NextLocation uint
}

type ExternalScreen struct {
	CoreScreen
	NextLocation uint
	Message      string
	Routes       []Route
}

type Menu []IScreen

// NOTE: List screen methods
func (ls *ListScreen) display(lang string) string {
	message := ls.Header
	for i, value := range ls.Options {
		message += fmt.Sprintf("\n%d.%s", i, value)
	}
	return message
}

func (ls *ListScreen) nextPage(input string) int {
	location, err := strconv.Atoi(input)
	if err != nil || len(ls.Options) < location {
		return -1
	}
	// selection := ls.Options[location]
	// FIXME: Do something with the selected item probably cache them
	return int(ls.NextLocation)
}

func (ls *ListScreen) screenType() string {
	return "LS"
}

// NOTE: Route screen methods
func (rs *RouteScreen) display(lang string) string {
	message := rs.Header
	for i, route := range rs.Routes {
		message += fmt.Sprintf("\n%d.%s", i, route.Value)
	}
	return message
}

func (rs *RouteScreen) nextPage(input string) int {
	location, err := strconv.Atoi(input)
	if err != nil {
		return -1
	}
	selection := rs.Routes[location]
	return int(selection.NextLocation)
}

func (rs *RouteScreen) screenType() string {
	return "RS"
}

// NOTE: Rawinput method
func (ris *RawInputScreen) display(lang string) string {
	return ris.Header
}

func (ris *RawInputScreen) nextPage(input string) int {
	return int(ris.NextLocation)
}

func (ris *RawInputScreen) screenType() string {
	return "RIS"
}

// NOTE: External screen methods
func (es *ExternalScreen) SetDisplay(message string) {
	es.Message = message
}

func (es *ExternalScreen) SetLocation(location int, routes []Route) (err error) {
	// NOTE: Must provide atleast routes or next location for this service
	if location < 0 || len(routes) <= 0 {
		err = fmt.Errorf("failed to next screen location")
		return
	}
	es.NextLocation = uint(location)
	es.Routes = routes
	return

}
func (es *ExternalScreen) display(lang string) string {
	return fmt.Sprintf("%s\n%s", es.Header, es.Message)
}

func (es *ExternalScreen) nextPage(input string) int {
	location, err := strconv.Atoi(input)
	if err != nil {
		return -1
	}
	if location == int(es.NextLocation) {
		return int(es.NextLocation)
	}
	if len(es.Routes) >= location {
		return int(es.Routes[location].NextLocation)
	}
	return -1
}

func (es *ExternalScreen) screenType() string {
	return "ES"
}

// Navigation
func (mn Menu) navigate(sessionId string, input string) (displayText string) {

	ihistory, active := CacheInstance.Get(sessionId)
	history, ok := ihistory.([]uint)

	if !active || !ok {
		//NOTE: No value for session in cache starting from the top
		CacheInstance.Set(sessionId, []uint{0}, time.Minute)
		return mn[0].display("en")
	}
	inputForScreen := history[len(history)-1]
	nextpage := mn[inputForScreen].nextPage(input)
	history = append(history, uint(nextpage))
	CacheInstance.Set(sessionId, history, time.Minute)
	return mn[nextpage].display("en")

}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	sessionId := "ryantest" // NOTE: Auto generated
	dial := 0
	menu := Menu{}
	fmt.Println("Enter ussd shortcode : ")
	for {
		// reads user input until \n by default
		scanner.Scan()
		// Holds the string that was scanned
		text := scanner.Text()
		if len(text) != 0 {
			if dial == 0 {
				menu = GetMenu(text)
				fmt.Println(menu.navigate(sessionId, ""))
			} else {
				fmt.Println(menu.navigate(sessionId, text))
			}
		} else {
			// exit if user entered an empty string
			break
		}
		dial += 1
		fmt.Println("***********************************************************")
	}

}

func GetMenu(shortcode string) (menu Menu) {

	ls := &ListScreen{
		CoreScreen: CoreScreen{
			Header: "Please select youre bank selection",
			Name:   "HomePage",
		},
		NextLocation: 1,
		Options:      []string{"KCB", "NCBA", "COOP"},
	}
	rs := &RouteScreen{
		CoreScreen: CoreScreen{
			Name:   "RoutedPage",
			Header: "Show routed selection",
		},
		Routes: []Route{{Value: "Selected", NextLocation: 2}},
	}
	rsi := &RawInputScreen{
		CoreScreen: CoreScreen{
			Header: "Please enter your name",
			Name:   "RawInputTest",
		},
		NextLocation: 3,
	}

	es := &ExternalScreen{
		CoreScreen: CoreScreen{
			Header: "Please enter your name",
			Name:   "RawInputTest",
		},
	}

	es.SetDisplay("1.YTUTY0\n2.KJ8789")
	es.SetLocation(-1, []Route{{Value: "Account number", NextLocation: 4}})
	menu = append(menu, ls, rs, rsi, es)
	return
}
