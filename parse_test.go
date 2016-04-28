// Copyright (c) 2016, Aelita Styles
//
// This file is part of the go-ircproto project, and is licensed under a
// BSD-like license. A copy of the license should have been included with this
// source code, but in the event that it was not, you can find it here:
// https://github.com/AelitaStyles/go-ircproto/blob/master/LICENSE.md

package ircproto

import "testing"

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
