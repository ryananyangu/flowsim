package models

import (
	"fmt"
	"strconv"

	"github.com/zohari-tech/flowsim/utils"
)

type IScreen interface {
	Display() string           // Display the current screen
	NextPage(input string) int // Gets the input from user and return index of what to display next
	ScreenType() utils.ScreenType        // type of the struct that has implemented the interface
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

// NOTE: List screen methods
func (ls ListScreen) Display() string {
	message := fmt.Sprintf("CON %s", ls.Header)
	if ls.IsEnd {
		message = fmt.Sprintf("END %s", ls.Header)
	}
	for i, value := range ls.Options {
		message += fmt.Sprintf("\n%d.%s", i, value)
	}
	return message
}

func (ls ListScreen) NextPage(input string) int {
	location, err := strconv.Atoi(input)
	if err != nil || len(ls.Options) < location {
		return -1
	}
	// selection := ls.Options[location]
	// FIXME: Do something with the selected item probably cache them
	return int(ls.NextLocation)
}

func (ls ListScreen) ScreenType() utils.ScreenType {
	return utils.LIST_SCREEN
}

// NOTE: Route screen methods
func (rs RouteScreen) Display() string {
	message := fmt.Sprintf("CON %s", rs.Header)
	if rs.IsEnd {
		message = fmt.Sprintf("END %s", rs.Header)
	}

	for i, route := range rs.Routes {
		message += fmt.Sprintf("\n%d.%s", i, route.Value)
	}
	return message
}

func (rs RouteScreen) NextPage(input string) int {
	location, err := strconv.Atoi(input)
	if err != nil {
		return -1
	}
	selection := rs.Routes[location]
	return int(selection.NextLocation)
}

func (rs RouteScreen) ScreenType() utils.ScreenType {
	return utils.ROUTE_SCREEN
}

// NOTE: Rawinput method
func (ris RawInputScreen) Display() string {
	message := fmt.Sprintf("CON %s", ris.Header)
	if ris.IsEnd {
		message = fmt.Sprintf("END %s", ris.Header)
	}
	return message
}

func (ris RawInputScreen) NextPage(input string) int {
	return int(ris.NextLocation)
}

func (ris RawInputScreen) ScreenType() utils.ScreenType {
	return utils.RAW_INPUT_SCREEN
}

// // FIXME: External screen methods
// func (es *ExternalScreen) SetDisplay(message string) {
// 	es.Message = message
// }

// func (es *ExternalScreen) SetLocation(location int, routes []Route) (err error) {
// 	// NOTE: Must provide atleast routes or next location for this service
// 	if location < 0 || len(routes) <= 0 {
// 		err = fmt.Errorf("failed to next screen location")
// 		return
// 	}
// 	es.NextLocation = uint(location)
// 	es.Routes = routes
// 	return

// }
// func (es *ExternalScreen) Display(lang string) string {
// 	return fmt.Sprintf("%s\n%s", es.Header, es.Message)
// }

// func (es *ExternalScreen) NextPage(input string) int {
// 	location, err := strconv.Atoi(input)
// 	if err != nil {
// 		return -1
// 	}
// 	if location == int(es.NextLocation) {
// 		return int(es.NextLocation)
// 	}
// 	if len(es.Routes) >= location {
// 		return int(es.Routes[location].NextLocation)
// 	}
// 	return -1
// }

// func (es *ExternalScreen) ScreenType() utils.ScreenType {
// 	return utils.EXTERNAL_SCREEN
// }
