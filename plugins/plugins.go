package plugins

// Plugin gives ways for users to interact with a Bot
type Plugin struct {
	Names []string
	Usage string
	Help  string
	Run   func(name *string, args *string) (text string)
}

// Available commands
var Plugins = []Plugin{
	Images, // Fetches images from Google
}
