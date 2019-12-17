// The chat package provides a simple chat client web page.
package chat

import (
	"tour/static"
)

// MainPage serves a web page that provides a simple chat client interface.
var ChatPage = static.Serve("chat/chat.html")
