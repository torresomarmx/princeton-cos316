/*****************************************************************************
 * router.go
 * Name:
 * NetId:
 *****************************************************************************/

package http_router

import (
	"errors"
	"log"
	"net/http"
	"regexp"
	"strings"
)

const R int = 256

// Student defined types or constants go here
type RouteNode struct {
	Handler     http.HandlerFunc
	HTTPMethod  string
	Next        *[R]RouteNode
	CaptureName string
}

type AddRouteNodeInput struct {
	Node       RouteNode
	CharIdx    int
	HTTPMethod string
	Pattern    string
	Handler    http.HandlerFunc
}

type RouteInfo struct {
	Pattern     string
	HTTPMethod  string
	Handler     http.HandlerFunc
	QueryString string
}

// HTTPRouter stores the information necessary to route HTTP requests
type HTTPRouter struct {
	// Place anything you'd like here
	RoutesTrieRoot RouteNode
}

// NewRouter creates a new HTTP Router, with no initial routes
func NewRouter() *HTTPRouter {
	return new(HTTPRouter)
}

// AddRoute adds a new route to the router, associating a given method and path
// pattern with the designated http handler.
func (router *HTTPRouter) AddRoute(method string, pattern string, handler http.HandlerFunc) {
	directories := strings.Split(pattern, "/")
	if len(directories) == 1 || len(directories[0]) != 0 {
		// pattern contains no "/" separators || pattern doesn't begin with "/"
		log.Fatal("Invalid pattern")
	}

	validDirectoryOrCapture := regexp.MustCompile("^[:]?[a-zA-Z0-9_.-]+$")
	for _, v := range directories[1:] {
		if len(v) == 0 || !validDirectoryOrCapture.MatchString(v) {
			log.Fatal("Invalid directory/capture in pattern")
		}
	}
	router.RoutesTrieRoot = router.addRouteNode(AddRouteNodeInput{
		Node:       router.RoutesTrieRoot,
		CharIdx:    0,
		HTTPMethod: method,
		Pattern:    pattern,
		Handler:    handler,
	})
	return
}

func (router *HTTPRouter) addRouteNode(input AddRouteNodeInput) RouteNode {
	if input.Node.Next == nil {
		input.Node.Next = new([R]RouteNode)
	}
	if input.CharIdx == len(input.Pattern) {
		input.Node.Handler = input.Handler
		input.Node.HTTPMethod = input.HTTPMethod
		return input.Node
	}

	var captureName string
	newCharIdx := input.CharIdx + 1
	// look for capture name
	if string(input.Pattern[input.CharIdx]) == ":" {
		captureName, newCharIdx = getDirectoryName(input.Pattern, newCharIdx)
	}

	input.Node.Next[input.Pattern[input.CharIdx]] = router.addRouteNode(AddRouteNodeInput{
		Node:       input.Node.Next[input.Pattern[input.CharIdx]],
		CharIdx:    newCharIdx,
		HTTPMethod: input.HTTPMethod,
		Pattern:    input.Pattern,
		Handler:    input.Handler,
	})

	input.Node.Next[input.Pattern[input.CharIdx]].CaptureName = captureName

	return input.Node
}

func getDirectoryName(pattern string, initCharIdx int) (string, int) {
	var directoryName string
	endCharIdx := initCharIdx
	for ; endCharIdx < len(pattern); endCharIdx++ {
		// look for end of directory name
		if string(pattern[endCharIdx]) == "/" {
			break
		}
	}

	directoryName = string(pattern[initCharIdx:endCharIdx])

	return directoryName, endCharIdx
}

func (router *HTTPRouter) GetRouteInfo(pattern string) (RouteInfo, error) {
	currentIdx := 0
	currentNode := router.RoutesTrieRoot

	if currentNode.Next == nil {
		return RouteInfo{}, errors.New("No routes in router")
	}

	queryBuilder := []string{}
	for currentIdx < len(pattern) {
		currentChar := pattern[currentIdx]
		if currentNode.Next[currentChar].Next != nil {
			currentNode = currentNode.Next[currentChar]
			currentIdx += 1
		} else if currentNode.Next[":"[0]].Next != nil {
			// find where capture arg ends
			captureNode := currentNode.Next[":"[0]]
			captureArg, captureArgEndIdx := getDirectoryName(pattern, currentIdx)
			queryBuilder = append(queryBuilder, captureNode.CaptureName+"="+captureArg)
			currentIdx = captureArgEndIdx
			currentNode = captureNode
		} else {
			return RouteInfo{}, errors.New("No route found for pattern")
		}
	}

	routeInfo := RouteInfo{Pattern: pattern,
		HTTPMethod:  currentNode.HTTPMethod,
		Handler:     currentNode.Handler,
		QueryString: strings.Join(queryBuilder, "&"),
	}

	return routeInfo, nil
}

func (router *HTTPRouter) Contains(pattern string) bool {
	_, err := router.GetRouteInfo(pattern)
	return err == nil
}

// ServeHTTP writes an HTTP response to the provided response writer
// by invoking the handler associated with the route that is appropriate
// for the provided request.
func (router *HTTPRouter) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	routeInfo, err := router.GetRouteInfo(request.URL.Path)
	if err != nil || routeInfo.HTTPMethod != request.Method {
		http.NotFound(response, request)
	}

	request.URL.RawQuery = routeInfo.QueryString
	routeInfo.Handler(response, request)
}
