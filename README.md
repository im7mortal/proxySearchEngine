# proxySearchEngine
The proxy server allow to do from the browsers input lines to search from different search engines dependently on input query.
Current logic do queries to [Yandex](https://yandex.com) if there are any cyrilic letters present(at least one) and to [Qwant](https://www.qwant.com).

### set up
Go to the https://{HOST}/discovery page and then follow the instructions for you browser **[FIREFOX](https://support.mozilla.org/en-US/kb/add-or-remove-search-engine-firefox)**

### Local usage

For better privacy is possible to run the proxy on your local machine instead of the server.

## Dev 

> PORT=":3456" HOST=1.2.3.4 go run github.com/im7mortal/proxySearchEngine/cmd/proxy
