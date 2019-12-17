package ajax

import (
	"net/http"
	"tour/static"
)

// LogPage serves a web page that shows log messages sent to Chan.
var LogPage http.Handler = static.Serve("ajax/logpage.html")
