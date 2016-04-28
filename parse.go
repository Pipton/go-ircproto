// Copyright (c) 2016, Aelita Styles
//
// This file is part of the go-ircproto project, and is licensed under a
// BSD-like license. A copy of the license should have been included with this
// source code, but in the event that it was not, you can find it here:
// https://github.com/AelitaStyles/go-ircproto/blob/master/LICENSE.md

package ircproto

import "fmt"

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

// ParseUserMask will parse a usermask string and return the appropriate
// IrcUserMask structure. The Type field will be set to either "None",
// "User", "Server" or "Unknown".
func ParseUserMask(mask string) (IrcUserMask, error) {
	var parsedmask IrcUserMask
	var usersep int
	var hostsep int
	dotcount := 0

	// Look for the user and host seperators
	for i, v := range mask {
		if v == '!' && usersep == 0 && dotcount == 0 {
			usersep = i
		} else if v == '!' && usersep != 0 {
			return IrcUserMask{}, fmt.Errorf("User mask has multiple " +
				"nick/user seperators.")
		} else if v == '!' && dotcount != 0 {
			return IrcUserMask{}, fmt.Errorf("Nickname contains dots, " +
				"which isn't permitted.")
		} else if v == '@' && hostsep == 0 && dotcount == 0 {
			hostsep = i
		} else if v == '@' && hostsep != 0 {
			return IrcUserMask{}, fmt.Errorf("User mask has multiple " +
				"host seperators.")
		} else if v == '@' && dotcount != 0 {
			return IrcUserMask{}, fmt.Errorf("Username contains dots, " +
				"which isn't permitted.")
		} else if v == '\r' || v == '\n' || v == ':' || v == ' ' {
			return IrcUserMask{}, fmt.Errorf("User mask contains " +
				"reserved characters.")
		} else if v == '.' {
			dotcount++
		}
	}

	// Check if hostsep is set
	if hostsep != 0 {
		parsedmask.Type = "User"
		parsedmask.Host = mask[hostsep+1:len(mask)]
		if usersep != 0 {
			parsedmask.Username = mask[usersep+1:hostsep]
			parsedmask.Nick = mask[:usersep]
		} else {
			parsedmask.Nick = mask[:hostsep]
		}
	} else if usersep == 0 && hostsep == 0 && dotcount != 0 {
		parsedmask.Type = "Server"
		parsedmask.Host = mask
	} else if usersep == 0 && hostsep == 0 && dotcount == 0 {
		parsedmask.Type = "Unknown"
		parsedmask.Host = mask
		parsedmask.Nick = mask
	} else {
		return IrcUserMask{}, fmt.Errorf("Invalid user mask.")
	}

	return parsedmask, nil
}
