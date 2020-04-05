package main

import (
	"bytes"
	"flag"
	"fmt"
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
	http.HandleFunc("/discover", discovery)
	http.HandleFunc(searchPluginPath, searchPluginHandler)
	http.HandleFunc(searchPath, proxysearchengine)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func discovery(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, searchEngineDiscovery)
	if err != nil {
		println(err.Error())
	}
}
func searchPluginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(searchEnginePlugin))
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

func proxysearchengine(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	var q string
	if val, exist := queries[querySearchKey]; exist {
		q = val[0]
	}
	searchUrl := google
	weekday := time.Now().Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		searchUrl = duckduckgo
	}
	if russianLetters(q) {
		searchUrl = yandex
	}
	http.Redirect(w, r, searchUrl+url.QueryEscape(q), http.StatusPermanentRedirect)
}

type searchEngine struct {
	// plugin discovery
	SearchEnginePluginURL string
	// plugin
	ShortName, Description, InputEncoding, URL, KeyName string
}
