// Copyright (c) 2016, Aelita Styles
//
// This file is part of the go-ircproto project, and is licensed under a
// BSD-like license. A copy of the license should have been included with this
// source code, but in the event that it was not, you can find it here:
// https://github.com/AelitaStyles/go-ircproto/blob/master/LICENSE.md

package ircproto

import (
	"fmt"
	"testing"
)

// ----------------------------------------------------------------------------
// PARSE RAW TESTS ------------------------------------------------------------
// ----------------------------------------------------------------------------

func TestParseRawValid1(t *testing.T) {
	unparsedcmd := ":server.test 001 TestUser :Welcome to the TestNet IRC " +
		"Network TestUser!test@user.client.test\r\n"
	var parsedcmd IrcCommand
	var err error

	parsedcmd, err = ParseRaw(unparsedcmd)
	t.Logf("Returned structure is: %+v", parsedcmd)
	if err != nil {
		t.Fatalf("ParseRaw failed with error '%s'", err)
	}

	if parsedcmd.Source.Type != "Server" || parsedcmd.Source.Host !=
		"server.test" {
		t.Fatalf("Source has type \"%s\" and host \"%s\", expected type "+
			"\"Server\" and host \"server.test\".", parsedcmd.Source.Type,
			parsedcmd.Source.Host)
	} else if parsedcmd.RawType != "001" {
		t.Fatalf("Expected RawType to be \"001\", got \"%s\".",
			parsedcmd.RawType)
	} else if parsedcmd.Type != "" {
		t.Fatalf("Expected Type to be unset, it was set to \"%s\"",
			parsedcmd.Type)
	} else if len := len(parsedcmd.RawArguments); len != 2 {
		t.Fatalf("Expected RawArguments array length to be 2, it was %d.", len)
	} else if parsedcmd.RawArguments[0] != "TestUser" {
		t.Fatalf("Expected RawArguments[0] to be set to \"TestUser\", it was "+
			"\"%s\".", parsedcmd.RawArguments[0])
	} else if parsedcmd.RawArguments[1] != "Welcome to the TestNet IRC "+
		"Network TestUser!test@user.client.test" {
		t.Fatalf("Expected RawArguments[1] to be set to \"Welcome to the "+
			"TestNet IRC Network TestUser!test@user.client.test\", it was "+
			"\"%s\".", parsedcmd.RawArguments[1])
	} else if parsedcmd.Data != nil {
		t.Fatalf("Expected Data to be nil, it was \"%+v\".", parsedcmd.Data)
	}
}
func TestParseRawValid2(t *testing.T) {
	unparsedcmd := ":server.test NOTICE TestUser :This is a notice. Boo!\r\n"
	var parsedcmd IrcCommand
	var err error

	parsedcmd, err = ParseRaw(unparsedcmd)
	t.Logf("Returned structure is: %+v", parsedcmd)
	if err != nil {
		t.Fatalf("ParseRaw failed with error '%s'", err)
	}

	if parsedcmd.Source.Type != "Server" || parsedcmd.Source.Host !=
		"server.test" {
		t.Fatalf("Source has type \"%s\" and host \"%s\", expected type "+
			"\"Server\" and host \"server.test\".", parsedcmd.Source.Type,
			parsedcmd.Source.Host)
	} else if parsedcmd.RawType != "NOTICE" {
		t.Fatalf("Expected RawType to be \"NOTICE\", got \"%s\".",
			parsedcmd.RawType)
	} else if parsedcmd.Type != "" {
		t.Fatalf("Expected Type to be unset, it was set to \"%s\"",
			parsedcmd.Type)
	} else if len := len(parsedcmd.RawArguments); len != 2 {
		t.Fatalf("Expected RawArguments array length to be 2, it was %d.", len)
	} else if parsedcmd.RawArguments[0] != "TestUser" {
		t.Fatalf("Expected RawArguments[0] to be set to \"TestUser\", it was "+
			"\"%s\".", parsedcmd.RawArguments[0])
	} else if parsedcmd.RawArguments[1] != "This is a notice. Boo!" {
		t.Fatalf("Expected RawArguments[1] to be set to \"This is a notice. "+
			"Boo!\", it was \"%s\".", parsedcmd.RawArguments[1])
	} else if parsedcmd.Data != nil {
		t.Fatalf("Expected Data to be nil, it was \"%+v\".", parsedcmd.Data)
	}
}
func TestParseRawValid3(t *testing.T) {
	unparsedcmd := ":OtherUser!foo@second.client.test PRIVMSG TestUser :This " +
		"is a message. Boo!\r\n"
	var parsedcmd IrcCommand
	var err error

	parsedcmd, err = ParseRaw(unparsedcmd)
	t.Logf("Returned structure is: %+v", parsedcmd)
	if err != nil {
		t.Fatalf("ParseRaw failed with error '%s'", err)
	}

	if parsedcmd.Source.Type != "User" || parsedcmd.Source.Nick !=
		"OtherUser" || parsedcmd.Source.Username != "foo" ||
		parsedcmd.Source.Host != "second.client.test" {
		t.Fatalf("Source has type \"%s\", nick \"%s\", username \"%s\" and "+
			"host \"%s\", expected type \"User\", nick \"OtherUser\", "+
			"username \"foo\" and host \"second.client.test\".",
			parsedcmd.Source.Type, parsedcmd.Source.Nick,
			parsedcmd.Source.Username, parsedcmd.Source.Host)
	} else if parsedcmd.RawType != "PRIVMSG" {
		t.Fatalf("Expected RawType to be \"PRIVMSG\", got \"%s\".",
			parsedcmd.RawType)
	} else if parsedcmd.Type != "" {
		t.Fatalf("Expected Type to be unset, it was set to \"%s\"",
			parsedcmd.Type)
	} else if len := len(parsedcmd.RawArguments); len != 2 {
		t.Fatalf("Expected RawArguments array length to be 2, it was %d.", len)
	} else if parsedcmd.RawArguments[0] != "TestUser" {
		t.Fatalf("Expected RawArguments[0] to be set to \"TestUser\", it was "+
			"\"%s\".", parsedcmd.RawArguments[0])
	} else if parsedcmd.RawArguments[1] != "This is a message. Boo!" {
		t.Fatalf("Expected RawArguments[1] to be set to \"This is a message. "+
			"Boo!\", it was \"%s\".", parsedcmd.RawArguments[1])
	} else if parsedcmd.Data != nil {
		t.Fatalf("Expected Data to be nil, it was \"%+v\".", parsedcmd.Data)
	}
}
func TestParseRawValid4(t *testing.T) {
	unparsedcmd := ":server.test 005 TestUser CAP1 CAP2 CAP3 CAP4 CAP5 CAP6 " +
		"CAP7 CAP8 CAP9 CAP10 CAP11 CAP12 CAP13 are supported by this " +
		"server\r\n"
	var parsedcmd IrcCommand
	var err error

	parsedcmd, err = ParseRaw(unparsedcmd)
	t.Logf("Returned structure is: %+v", parsedcmd)
	if err != nil {
		t.Fatalf("ParseRaw failed with error '%s'", err)
	}

	if parsedcmd.Source.Type != "Server" || parsedcmd.Source.Host !=
		"server.test" {
		t.Fatalf("Source has type \"%s\" and host \"%s\", expected type "+
			"\"Server\" and host \"server.test\".", parsedcmd.Source.Type,
			parsedcmd.Source.Host)
	} else if parsedcmd.RawType != "005" {
		t.Fatalf("Expected RawType to be \"005\", got \"%s\".",
			parsedcmd.RawType)
	} else if parsedcmd.Type != "" {
		t.Fatalf("Expected Type to be unset, it was set to \"%s\"",
			parsedcmd.Type)
	} else if len := len(parsedcmd.RawArguments); len != 15 {
		t.Fatalf("Expected RawArguments array length to be 15, it was %d.", len)
	} else if parsedcmd.RawArguments[0] != "TestUser" {
		t.Fatalf("Expected RawArguments[0] to be set to \"TestUser\", it was "+
			"\"%s\".", parsedcmd.RawArguments[0])
	} else if parsedcmd.RawArguments[14] != "are supported by this server" {
		t.Fatalf("Expected RawArguments[14] to be set to \"are supported by "+
			"this server\", it was \"%s\".", parsedcmd.RawArguments[14])
	} else if parsedcmd.Data != nil {
		t.Fatalf("Expected Data to be nil, it was \"%+v\".", parsedcmd.Data)
	}

	for i := 1; i < 14; i++ {
		if parsedcmd.RawArguments[i] != fmt.Sprintf("CAP%d", i) {
			t.Fatalf("Expected RawArguments[%d] to be \"CAP%d\", but it was "+
				"\"%s\".", i, i, parsedcmd.RawArguments[i])
		}
	}
}
func TestParseRawValid5(t *testing.T) {
	unparsedcmd := "PING :BOOP\r\n"
	var parsedcmd IrcCommand
	var err error

	parsedcmd, err = ParseRaw(unparsedcmd)
	t.Logf("Returned structure is: %+v", parsedcmd)
	if err != nil {
		t.Fatalf("ParseRaw failed with error '%s'", err)
	}

	if parsedcmd.Source.Type != "None" {
		t.Fatalf("Source type expected to be \"None\", got \"%s\".",
			parsedcmd.Source.Type)
	} else if parsedcmd.RawType != "PING" {
		t.Fatalf("Raw type texpected to be \"PING\", got \"%s\".",
			parsedcmd.RawType)
	} else if parsedcmd.Type != "" {
		t.Fatalf("Parsed type not meant to be set here, it was set to \"%s\".",
			parsedcmd.Type)
	} else if parsedcmd.Data != nil {
		t.Fatalf("Parsed data not meant to be set here, it was set to \"%+v\".",
			parsedcmd.Data)
	} else if len := len(parsedcmd.RawArguments); len != 1 {
		t.Fatalf("Raw arguments array expected to be of length 1, length is "+
			"%d.", len)
	} else if parsedcmd.RawArguments[0] != "BOOP" {
		t.Fatalf("Expected RawArguments[0] to be set to \"BOOP\", it is "+
			"\"%s\".", parsedcmd.RawArguments[0])
	}
}
func TestParseRawValid6(t *testing.T) {
	unparsedcmd := "TEST BOOP\r\n"
	var parsedcmd IrcCommand
	var err error

	parsedcmd, err = ParseRaw(unparsedcmd)
	t.Logf("Returned structure is: %+v", parsedcmd)
	if err != nil {
		t.Fatalf("ParseRaw failed with error '%s'", err)
	}

	if parsedcmd.Source.Type != "None" {
		t.Fatalf("Source type expected to be \"None\", got \"%s\".",
			parsedcmd.Source.Type)
	} else if parsedcmd.RawType != "TEST" {
		t.Fatalf("Raw type texpected to be \"TEST\", got \"%s\".",
			parsedcmd.RawType)
	} else if parsedcmd.Type != "" {
		t.Fatalf("Parsed type not meant to be set here, it was set to \"%s\".",
			parsedcmd.Type)
	} else if parsedcmd.Data != nil {
		t.Fatalf("Parsed data not meant to be set here, it was set to \"%+v\".",
			parsedcmd.Data)
	} else if len := len(parsedcmd.RawArguments); len != 1 {
		t.Fatalf("Raw arguments array expected to be of length 1, length is "+
			"%d.", len)
	} else if parsedcmd.RawArguments[0] != "BOOP" {
		t.Fatalf("Expected RawArguments[0] to be set to \"BOOP\", it is "+
			"\"%s\".", parsedcmd.RawArguments[0])
	}
}
func TestParseRawInvalid1(t *testing.T) {
	unparsedcmd := ":x!x!foo@test PRIVMSG TestUser :Hello\r\n"
	var parsedcmd IrcCommand
	var err error

	parsedcmd, err = ParseRaw(unparsedcmd)
	t.Logf("Returned error is: %s", err)
	if err == nil {
		t.Fatalf("ParseRaw should have failed with error, but didn't.")
	}

	if parsedcmd.Data != nil || parsedcmd.RawType != "" ||
		parsedcmd.Type != "" || len(parsedcmd.RawArguments) != 0 {
		t.Fatalf("ParseRaw returned a non-empty IrcCommand after an error.")
	}

	if err.Error() != "User mask has multiple nick/user seperators." {
		t.Fatalf("Expected err to be \"User mask has multiple nick/user "+
			"seperators.\", it was \"%+v\".", err)
	}
}
func TestParseRawInvalid2(t *testing.T) {
	unparsedcmd := ":server.test 0a0 TestUser :Hello\r\n"
	var parsedcmd IrcCommand
	var err error

	parsedcmd, err = ParseRaw(unparsedcmd)
	t.Logf("Returned error is: %s", err)
	if err == nil {
		t.Fatalf("ParseRaw should have failed with error, but didn't.")
	}

	if parsedcmd.Data != nil || parsedcmd.RawType != "" ||
		parsedcmd.Type != "" || len(parsedcmd.RawArguments) != 0 {
		t.Fatalf("ParseRaw returned a non-empty IrcCommand after an error.")
	}

	if err.Error() != "Command type is not a valid numeric." {
		t.Fatalf("Expected err to be \"Command type is not a valid numeric.\","+
			" it was \"%+v\".", err)
	}
}
func TestParseRawInvalid3(t *testing.T) {
	unparsedcmd := ":server.test M30W TestUser :Hello\r\n"
	var parsedcmd IrcCommand
	var err error

	parsedcmd, err = ParseRaw(unparsedcmd)
	t.Logf("Returned error is: %s", err)
	if err == nil {
		t.Fatalf("ParseRaw should have failed with error, but didn't.")
	}

	if parsedcmd.Data != nil || parsedcmd.RawType != "" ||
		parsedcmd.Type != "" || len(parsedcmd.RawArguments) != 0 {
		t.Fatalf("ParseRaw returned a non-empty IrcCommand after an error.")
	}

	if err.Error() != "Command type contains invalid characters." {
		t.Fatalf("Expected err to be \"Command type contains invalid "+
			"characters.\", it was \"%+v\".", err)
	}
}
func TestParseRawInvalid4(t *testing.T) {
	unparsedcmd := ":server.test 001 TestUser :Hello"
	var parsedcmd IrcCommand
	var err error

	parsedcmd, err = ParseRaw(unparsedcmd)
	t.Logf("Returned error is: %s", err)
	if err == nil {
		t.Fatalf("ParseRaw should have failed with error, but didn't.")
	}

	if parsedcmd.Data != nil || parsedcmd.RawType != "" ||
		parsedcmd.Type != "" || len(parsedcmd.RawArguments) != 0 {
		t.Fatalf("ParseRaw returned a non-empty IrcCommand after an error.")
	}

	if err.Error() != "Command does not end with CRLF." {
		t.Fatalf("Expected err to be \"Command does not end with CRLF.\", "+
			"it was \"%+v\".", err)
	}
}

// ----------------------------------------------------------------------------
// PARSE USER MASK TESTS ------------------------------------------------------
// ----------------------------------------------------------------------------

func TestParseUserMaskValidUser1(t *testing.T) {
	unparsedmask := "TestUser!test@user.client.test"
	var parsedmask IrcUserMask
	var err error

	parsedmask, err = ParseUserMask(unparsedmask)
	t.Logf("Returned structure is: %+v", parsedmask)
	if err != nil {
		t.Fatalf("ParseUserMask failed with error '%s'", err)
	}

	if parsedmask.Type != "User" {
		t.Fatalf("Expected Type field to be \"User\", got \"%s\".",
			parsedmask.Type)
	} else if parsedmask.Nick != "TestUser" {
		t.Fatalf("Expected Nick field to be \"TestUser\", got \"%s\".",
			parsedmask.Nick)
	} else if parsedmask.Username != "test" {
		t.Fatalf("Expected Username field to be \"test\", got \"%s\".",
			parsedmask.Username)
	} else if parsedmask.Host != "user.client.test" {
		t.Fatalf("Expected Host field to be \"user.client.test\", got \"%s\".",
			parsedmask.Host)
	}
}
func TestParseUserMaskValidUser2(t *testing.T) {
	unparsedmask := "TestUser@user.client.test"
	var parsedmask IrcUserMask
	var err error

	parsedmask, err = ParseUserMask(unparsedmask)
	t.Logf("Returned structure is: %+v", parsedmask)
	if err != nil {
		t.Fatalf("ParseUserMask failed with error '%s'", err)
	}

	if parsedmask.Type != "User" {
		t.Fatalf("Expected Type field to be \"User\", got \"%s\".",
			parsedmask.Type)
	} else if parsedmask.Nick != "TestUser" {
		t.Fatalf("Expected Nick field to be \"TestUser\", got \"%s\".",
			parsedmask.Nick)
	} else if parsedmask.Username != "" {
		t.Fatalf("Expected Username field to be empty, got \"%s\".",
			parsedmask.Username)
	} else if parsedmask.Host != "user.client.test" {
		t.Fatalf("Expected Host field to be \"user.client.test\", got \"%s\".",
			parsedmask.Host)
	}
}
func TestParseUserMaskValidServer1(t *testing.T) {
	unparsedmask := "server.test"
	var parsedmask IrcUserMask
	var err error

	parsedmask, err = ParseUserMask(unparsedmask)
	t.Logf("Returned structure is: %+v", parsedmask)
	if err != nil {
		t.Fatalf("ParseUserMask failed with error '%s'", err)
	}

	if parsedmask.Type != "Server" {
		t.Fatalf("Expected Type field to be \"Server\", got \"%s\".",
			parsedmask.Type)
	} else if parsedmask.Nick != "" {
		t.Fatalf("Expected Nick field to be empty, got \"%s\".",
			parsedmask.Nick)
	} else if parsedmask.Username != "" {
		t.Fatalf("Expected Username field to be empty, got \"%s\".",
			parsedmask.Username)
	} else if parsedmask.Host != "server.test" {
		t.Fatalf("Expected Host field to be \"server.test\", got \"%s\".",
			parsedmask.Host)
	}
}
