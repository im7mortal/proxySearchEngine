package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"text/template"
	"time"
)

func getEnv(name string, variable *string) {
	var exist bool
	*variable, exist = os.LookupEnv(name)
	if !exist {
		panic("ENV: " + name + " doesn't exist.")
	}
}

// Environment variable names
const (
	portEnv = "PORT"
	hostEnv = "HOST"
)

// Environment variable
var (
	port string
	host string
)
var searchEnginePlugin string
var searchEngineDiscovery string

func init() {
	flag.Parse()
	getEnv(hostEnv, &host)
	getEnv(portEnv, &port)
	var searchURL = url.URL{
		Scheme: "https",
		Host:   host + port,
		Path:   searchPath,
	}
	var searchEngineConfig = searchEngine{
		SearchEnginePluginURL: searchPluginPath,

		ShortName:     "SearchProxy",
		Description:   "Yandex for cyrilic letters, google for latin letters only and duckduckgo on weekends",
		InputEncoding: "utf-8",
		URL:           searchURL.String(),
		KeyName:       querySearchKey,
	}
	searchEnginePlugin = stringFromTemplate(openSearchTemplate, searchEngineConfig, "searchEnginePlugin")
	searchEngineDiscovery = stringFromTemplate(discoveryPlugin, searchEngineConfig, "searchEngineDiscovery")
}

var openSearchTemplate = `<OpenSearchDescription xmlns="http://a9.com/-/spec/opensearch/1.1/"
                       xmlns:moz="http://www.mozilla.org/2006/browser/search/">
  <ShortName>{{.ShortName}}</ShortName>
  <Description>{{.Description}}</Description>
  <InputEncoding>{{.InputEncoding}}</InputEncoding>
  <Url type="text/html" template="{{.URL}}">
    <Param name="{{.KeyName}}" value="{searchTerms}"/>
  </Url>
</OpenSearchDescription>`

var discoveryPlugin = ` 
<html>
<header>
<title>Register</title>
<link rel="search"
      type="application/opensearchdescription+xml"
      title="{{.ShortName}}"
      href="{{.SearchEnginePluginURL}}">
</header>
<body>
</body>
</html>
`

var searchPluginPath = "/search_plugin.xml"
var searchPath = "/proxysearchengine"
var querySearchKey = "proxyText"

func stringFromTemplate(tpl string, srch searchEngine, name string) string {
	var err error
	searchEnginePlugin, err := template.New(name).
		Parse(tpl)
	if err != nil {
		panic(err)
	}
	buff := bytes.NewBuffer([]byte{})
	err = searchEnginePlugin.Execute(buff, srch)
	if err != nil {
		panic(err)
	}
	return buff.String()
}

func main() {
	var err error
	println(searchEnginePlugin)
	println(searchEngineDiscovery)
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/discover", discovery)
	r.GET(searchPluginPath, searchPluginHandler)
	r.GET("/proxysearchengine", proxysearchengine)
	err = r.Run(port)
	if err != nil {
		log.Fatal(err)
	}
}

func discovery(c *gin.Context) {
	_, err := fmt.Fprintf(c.Writer, searchEngineDiscovery)
	if err != nil {
		println(err.Error())
	}
}
func searchPluginHandler(c *gin.Context) {
	c.Header("Content-Type", "application/xml")
	c.Status(http.StatusOK)
	_, err := c.Writer.Write([]byte(searchEnginePlugin))
	if err != nil {
		println(err.Error())
	}
}

const (
	google     = "https://www.google.com/search?q="
	yandex     = "https://yandex.ru/search/?text="
	duckduckgo = "https://duckduckgo.com/?q="
)

var cyrilicR *regexp.Regexp

func init() {
	var err error
	cyrilicR, err = regexp.Compile("[а-яА-Я]")
	if err != nil {
		panic(err.Error())
	}
}

func russianLetters(s string) bool {
	return cyrilicR.Match([]byte(s))
}

func proxysearchengine(c *gin.Context) {
	searchUrl := google
	q := c.Query(querySearchKey)
	weekday := time.Now().Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		searchUrl = duckduckgo
	}
	if russianLetters(q) {
		searchUrl = yandex
	}
	c.Redirect(http.StatusPermanentRedirect, searchUrl+url.QueryEscape(q))
	c.Status(http.StatusOK)
}

type searchEngine struct {
	// plugin discovery
	SearchEnginePluginURL string
	// plugin
	ShortName, Description, InputEncoding, URL, KeyName string
}
