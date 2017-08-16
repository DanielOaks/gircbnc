// Copyright (c) 2016-2017 Daniel Oaks <daniel@danieloaks.net>
// released under the MIT license

package ircbnc

import (
	"github.com/goshuirc/irc-go/ircmsg"
)

// ClientCommands holds all commands executable by a client connected to a listener.
var ClientCommands map[string]ClientCommand

func init() {
	ClientCommands = make(map[string]ClientCommand)
	loadClientCommands()
}

// ClientCommand represents a command accepted on a listener.
type ClientCommand struct {
	handler      func(listener *Listener, msg ircmsg.IrcMessage)
	usablePreReg bool
	minParams    int
}

// Run runs this command with the given listener/message.
func (cmd *ClientCommand) Run(listener *Listener, msg ircmsg.IrcMessage) {
	if !listener.Registered && !cmd.usablePreReg {
		// command silently ignored
		return
	}
	if len(msg.Params) < cmd.minParams {
		listener.Send(nil, "", "461", listener.ClientNick, msg.Command, "Not enough parameters")
		return
	}
	cmd.handler(listener, msg)

	// after each command, see if we can send registration to the listener
	if !listener.Registered {
		listener.tryRegistration()
	}
}
