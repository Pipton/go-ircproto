// Copyright (c) 2016, Aelita Styles
//
// This file is part of the go-ircproto project, and is licensed under a
// BSD-like license. A copy of the license should have been included with this
// source code, but in the event that it was not, you can find it here:
// https://github.com/AelitaStyles/go-ircproto/blob/master/LICENSE.md

package ircproto

// The IrcCommand structure represents a single IRC command/message.
type IrcCommand struct {
	Source       IrcUserMask // The source of the message (if given)
	Type         string      // The type of command, if ircproto understood it
	Data         interface{} // The parsed data, if iroproto understood the command
	RawType      string      // The command type as given in the command
	RawArguments []string    // The raw arguments as an array
}

// The IrcUserMask structure represents a IRC user mask
type IrcUserMask struct {
	Type     string // Whether this is a user, server or neither
	Nick     string // Nickname of the user (if a user)
	Username string // Username of the user (if a user)
	Host     string // Host of the user or server (if given)
}
