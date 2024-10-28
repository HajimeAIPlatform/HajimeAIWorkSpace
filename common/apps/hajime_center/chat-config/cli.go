package chat_config

var CLI struct {
	Verbose bool   `help:"Verbose mode."`
	Config  string `help:"Config file." name:"chat-config" type:"file" default:"chat-config.json"`
}
