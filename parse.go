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
	Source       IrcUserMask // The source of the message (if given).
	Type         int         // The type of command; see the IrcCommandType constants.
	Data         interface{} // The parsed data, if ircproto understood the command.
	RawType      string      // The command type as given in the command.
	RawArguments []string    // The raw arguments as an array.
}

// The list of IrcCommandType values
const (
	IrcCommandType_Unknown = iota // Used if ircproto doesn't understand the command
)

// The IrcUserMask structure represents a IRC user mask. If the mask is that
// of a user, the host and nickname fields should be set, and the username
// field may be set. If the source is a server, only the host field will be
// set. If the source can't be determined, both nick and host will be set
// with the same value.
//
// Occasionally ircproto will be unable to determine the source type, in
// such situations, you must work out the source yourself. If the type
// field is marked as empty, this means the source is the object you are
// directly connected to. If the field is marked as unknown, the meaning
// depends on whether you are a server or client.
//
// An unknown source type for a client means the source is either a server,
// or yourself (or the server acting on your behalf). For servers, it could
// either be a user or server.
//
// If you are using ircproto's state tracking capabilities, this should never
// be a problem.
type IrcUserMask struct {
	Type     int    // The type; see the IrcSourceType constants.
	Nick     string // Nickname of the user.
	Username string // Username of the user.
	Host     string // Host of the user or server.
}

// The list of IrcSourceType values
const (
	IrcSourceType_Empty = iota
	IrcSourceType_Unknown = iota
	IrcSourceType_User = iota
	IrcSourceType_Server = iota
)

// ParseRaw will check the IRC command is valid, split it up and return the
// contents as an IrcCommand structure, without doing any additional parsing
// or state tracking.
func ParseRaw(cmd string) (IrcCommand, error) {
	cmdlen := len(cmd)
	prefixEnd := 1
	cmdTypeEnd := 0
	var arguments []string
	parsedcmd := IrcCommand{Type: IrcCommandType_Unknown}

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
		parsedcmd.Source = IrcUserMask{Type: IrcSourceType_Empty}
		prefixEnd = 0
	}

	// Parse the command type
	if cmd[prefixEnd] > 47 && cmd[prefixEnd] < 58 {
		// Handle numeric types
		if cmd[prefixEnd+1] > 47 && cmd[prefixEnd+1] < 58 &&
			cmd[prefixEnd+2] > 47 && cmd[prefixEnd+2] < 58 &&
			cmd[prefixEnd+3] == ' ' {
			cmdTypeEnd = prefixEnd + 3
		} else {
			return IrcCommand{}, fmt.Errorf("Command type is not a valid " +
				"numeric.")
		}
	} else {
		// Handle named commands
		for i, v := range cmd[prefixEnd:] {
			if v == ' ' {
				cmdTypeEnd = i + prefixEnd
				break
			} else if (v < 65 || v > 122) && !(v > 90 && v < 97) {
				return IrcCommand{}, fmt.Errorf("Command type contains invalid" +
					" characters.")
			}
		}
	}
	parsedcmd.RawType = cmd[prefixEnd:cmdTypeEnd]

	// Fetch arguments
	argStart := cmdTypeEnd + 1
	argEnd := 0
	argCount := 0
	for i, v := range cmd[cmdTypeEnd+1:] {
		if v == ':' {
			argEnd = cmdlen - 2
			arguments = append(arguments, cmd[argStart+1:argEnd])
			break
		} else if v == '\r' {
			argEnd = i + cmdTypeEnd + 1
			arguments = append(arguments, cmd[argStart:argEnd])
			break
		} else if argCount == 14 {
			argEnd = cmdlen - 2
			arguments = append(arguments, cmd[argStart:argEnd])
			break
		} else if v == ' ' {
			argEnd = i + cmdTypeEnd + 1
			arguments = append(arguments, cmd[argStart:argEnd])
			argCount++
			argStart = argEnd + 1
		}
	}
	parsedcmd.RawArguments = arguments

	if cmd[cmdlen-2:cmdlen] != "\r\n" {
		return IrcCommand{}, fmt.Errorf("Command does not end with CRLF.")
	}

	return parsedcmd, nil
}

// ParseUserMask will parse a usermask string and return the appropriate
// IrcUserMask structure.
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
		// Handle usernames
		parsedmask.Type = IrcSourceType_User
		parsedmask.Host = mask[hostsep+1 : len(mask)]
		if usersep != 0 {
			parsedmask.Username = mask[usersep+1 : hostsep]
			parsedmask.Nick = mask[:usersep]
		} else {
			parsedmask.Nick = mask[:hostsep]
		}
	} else if usersep == 0 && hostsep == 0 && dotcount != 0 {
		// Handle servers
		parsedmask.Type = IrcSourceType_Server
		parsedmask.Host = mask
	} else if usersep == 0 && hostsep == 0 && dotcount == 0 {
		// Handle unknowns
		parsedmask.Type = IrcSourceType_Unknown
		parsedmask.Host = mask
		parsedmask.Nick = mask
	} else {
		return IrcUserMask{}, fmt.Errorf("Invalid user mask.")
	}

	return parsedmask, nil
}
