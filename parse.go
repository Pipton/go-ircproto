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

// ParseRaw will check the IRC command is valid, split it up and return the
// contents as an IrcCommand structure, without doing any additional parsing
// or state tracking.
func ParseRaw(cmd string) (IrcCommand, error) {
	cmdlen := len(cmd)
	prefixEnd := 1
	cmdTypeEnd := 0
	var arguments []string
	var parsedcmd IrcCommand

	// Check if we have a prefix
	if cmd[:1] == ":" {
		for cmd[prefixEnd-1:prefixEnd] != " " && prefixEnd < cmdlen {
			prefixEnd++
		}

		parsedprefix, err := ParseUserMask(cmd[1 : prefixEnd-1])
		if err != nil {
			return IrcCommand{}, err
		}
		parsedcmd.Source = parsedprefix
	} else {
		parsedcmd.Source = IrcUserMask{Type: "None"}
		prefixEnd = 0
	}

	// Parse the command type
	if cmd[prefixEnd] > 47 && cmd[prefixEnd] < 58 {
		if cmd[prefixEnd+1] > 47 && cmd[prefixEnd+1] < 58 &&
			cmd[prefixEnd+2] > 47 && cmd[prefixEnd+2] < 58 &&
			cmd[prefixEnd+3] == ' ' {
			cmdTypeEnd = prefixEnd + 3
		} else {
			return IrcCommand{}, fmt.Errorf("Command type is not a valid " +
				"numeric.")
		}
	} else {
		for i, v := range cmd[prefixEnd:] {
			if v == ' ' {
				cmdTypeEnd = i + prefixEnd
				break
			} else if (v < 65 || v > 122) && !(v > 90 && v < 97) {
				return IrcCommand{}, fmt.Errorf("Command type contains invalid" +
					"characters.")
			}
		}
	}
	parsedcmd.RawType = cmd[prefixEnd:cmdTypeEnd]

	// Fetch arguments
	argStart := cmdTypeEnd + 1
	argEnd := 0
	argCount := 0
	for i, v := range cmd[cmdTypeEnd+1:] {
		if argCount == 14 {
			break
		} else if v == ':' {
			argEnd = cmdlen
			arguments = append(arguments, cmd[argStart+2:cmdlen-2])
			break
		} else if v == ' ' {
			argEnd = i + cmdTypeEnd + 1
			arguments = append(arguments, cmd[argStart:argEnd])
			argCount++
			argStart = argEnd
		} else if v == '\r' {
			argEnd = i + cmdTypeEnd + 1
			arguments = append(arguments, cmd[argStart:argEnd])
			break
		}
	}
	parsedcmd.RawArguments = arguments

	if cmd[cmdlen-2:cmdlen] != "\r\n" {
		return IrcCommand{}, fmt.Errorf("Command does not end with CRLF.")
	}

	return parsedcmd, nil
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
		parsedmask.Host = mask[hostsep+1 : len(mask)]
		if usersep != 0 {
			parsedmask.Username = mask[usersep+1 : hostsep]
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
