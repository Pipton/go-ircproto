// Copyright (c) 2016, Aelita Styles
//
// This file is part of the go-ircproto project, and is licensed under a
// BSD-like license. A copy of the license should have been included with this
// source code, but in the event that it was not, you can find it here:
// https://github.com/AelitaStyles/go-ircproto/blob/master/LICENSE.md

package ircproto

import (
	"fmt"
	"reflect"
	"testing"
)

// ----------------------------------------------------------------------------
// PARSE RAW TESTS ------------------------------------------------------------
// ----------------------------------------------------------------------------

func testGoodParseRaw(testCmd string, expectedSource IrcUserMask,
	expectedRawType string, expectedRawArguments []string) (IrcCommand,
	error) {
	// Fetch parsed command
	var parsedCmd IrcCommand
	var err error
	parsedCmd, err = ParseRaw(testCmd)

	// Check if it errored unexpectedly
	if err != nil {
		return parsedCmd, fmt.Errorf("ParseRaw failed unexpectedly with error:"+
			" %s", err.Error())
	}

	// Check that it contains everything we expected
	if !reflect.DeepEqual(parsedCmd.Source, expectedSource) {
		return parsedCmd, fmt.Errorf("ParseRaw returned an unexpected source "+
			"of '%+v'. We were expecting it to be '%+v'.", parsedCmd.Source,
			expectedSource)
	} else if parsedCmd.RawType != expectedRawType {
		return parsedCmd, fmt.Errorf("ParseRaw returned an unexpected raw "+
			"type of '%s'. We were expecting it to be '%s'.",
			parsedCmd.RawType, expectedRawType)
	} else if len(parsedCmd.RawArguments) != len(expectedRawArguments) {
		return parsedCmd, fmt.Errorf("ParseRaw returned an unexpected number "+
			"of arguments. We were expecting %d, but got %d.",
			len(expectedRawArguments), len(parsedCmd.RawArguments))
	} else if len(parsedCmd.RawArguments) > 15 {
		return parsedCmd, fmt.Errorf("ParseRaw returned an unexpected number "+
			"of arguments. We were expecting 15 or less, but got %d.",
			len(parsedCmd.RawArguments))
	} else if parsedCmd.Type != "" {
		return parsedCmd, fmt.Errorf("ParseRaw returned an unexpected type of "+
			"'%s'. We were expecting it to be an empty string.", parsedCmd.Type)
	} else if parsedCmd.Data != nil {
		return parsedCmd, fmt.Errorf("ParseRaw returned an unexpected data "+
			"value of '%+v'. We were expecting it to be nil.", parsedCmd.Data)
	}

	for key, val := range expectedRawArguments {
		if val != parsedCmd.RawArguments[key] {
			return parsedCmd, fmt.Errorf("ParseRaw returned an unexpected "+
				"value at RawArguments[%d]. We were expecting '%s', but got "+
				"'%s'.", key, val, parsedCmd.RawArguments[key])
		}
	}

	return parsedCmd, nil
}

func TestParseRawValid1(t *testing.T) {
	parsedCmd, err := testGoodParseRaw(":server.test 001 TestUser :Welcome "+
		"to the TestNet IRC Network TestUser!test@user.client.test\r\n",
		IrcUserMask{Type: "Server", Host: "server.test"}, "001",
		[]string{"TestUser", "Welcome to the TestNet IRC Network " +
		"TestUser!test@user.client.test"})

	if err != nil {
		t.Error(err)
	}
	t.Logf("ParseRaw returned structure: %+v", parsedCmd)
}
func TestParseRawValid2(t *testing.T) {
	parsedCmd, err := testGoodParseRaw(":server.test NOTICE TestUser :This "+
		"is a notice. Boo!\r\n", IrcUserMask{Type: "Server",
		Host: "server.test"}, "NOTICE", []string{"TestUser", "This is a " +
		"notice. Boo!"})

	if err != nil {
		t.Error(err)
	}
	t.Logf("ParseRaw returned structure: %+v", parsedCmd)
}
func TestParseRawValid3(t *testing.T) {
	parsedCmd, err := testGoodParseRaw(":OtherUser!foo@second.client.test "+
		"PRIVMSG TestUser :This is a message. Boo!\r\n",
		IrcUserMask{Type: "User", Nick: "OtherUser", Username: "foo",
		Host: "second.client.test"}, "PRIVMSG", []string{"TestUser", "This is "+
		"a message. Boo!"})

	if err != nil {
		t.Error(err)
	}
	t.Logf("ParseRaw returned structure: %+v", parsedCmd)
}
func TestParseRawValid4(t *testing.T) {
	parsedCmd, err := testGoodParseRaw(":server.test 005 TestUser CAP1 CAP2 "+
		"CAP3 CAP4 CAP5 CAP6 CAP7 CAP8 CAP9 CAP10 CAP11 CAP12 CAP13 are "+
		"supported by this server\r\n", IrcUserMask{Type: "Server",
		Host: "server.test"}, "005", []string{"TestUser", "CAP1", "CAP2",
		"CAP3", "CAP4", "CAP5", "CAP6", "CAP7", "CAP8", "CAP9", "CAP10",
		"CAP11", "CAP12", "CAP13", "are supported by this server"})

	if err != nil {
		t.Error(err)
	}
	t.Logf("ParseRaw returned structure: %+v", parsedCmd)
}
func TestParseRawValid5(t *testing.T) {
	parsedCmd, err := testGoodParseRaw("PING :BOOP\r\n",
		IrcUserMask{Type: "None"}, "PING", []string{"BOOP"})

	if err != nil {
		t.Error(err)
	}
	t.Logf("ParseRaw returned structure: %+v", parsedCmd)
}
func TestParseRawValid6(t *testing.T) {
	parsedCmd, err := testGoodParseRaw("TEST BOOP\r\n",
		IrcUserMask{Type: "None"}, "TEST", []string{"BOOP"})

	if err != nil {
		t.Error(err)
	}
	t.Logf("ParseRaw returned structure: %+v", parsedCmd)
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
// PARSE RAW BENCHMARKS -------------------------------------------------------
// ----------------------------------------------------------------------------

func BenchmarkParseRawNumeric(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := ParseRaw(":server.test 001 TestUser :Welcome to the " +
			"TestNet IRC Network TestUser!test@user.client.test\r\n")
		if err != nil {
			b.FailNow()
		}
	}
}
func BenchmarkParseRawPrivmsg(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := ParseRaw(":OtherUser!foo@second.client.test PRIVMSG TestUser" +
			" :This is a message. Boo!\r\n")
		if err != nil {
			b.FailNow()
		}
	}
}
func BenchmarkParseRawCaps(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := ParseRaw(":server.test 005 TestUser CAP1 CAP2 CAP3 CAP4 " +
			"CAP5 CAP6 CAP7 CAP8 CAP9 CAP10 CAP11 CAP12 CAP13 are " +
			"supported by this server\r\n")
		if err != nil {
			b.FailNow()
		}
	}
}
func BenchmarkParseRawPing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := ParseRaw("PING :BOOP\r\n")
		if err != nil {
			b.FailNow()
		}
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
func TestParseUserMaskValidAmbigiousMask1(t *testing.T) {
	unparsedmask := "test"
	var parsedmask IrcUserMask
	var err error

	parsedmask, err = ParseUserMask(unparsedmask)
	t.Logf("Returned structure is: %+v", parsedmask)
	if err != nil {
		t.Fatalf("ParseUserMask failed with error '%s'", err)
	}

	if parsedmask.Type != "Unknown" {
		t.Fatalf("Expected Type field to be \"Unknown\", got \"%s\".",
			parsedmask.Type)
	} else if parsedmask.Nick != "test" {
		t.Fatalf("Expected Nick field to be \"test\", got \"%s\".",
			parsedmask.Nick)
	} else if parsedmask.Username != "" {
		t.Fatalf("Expected Username field to be empty, got \"%s\".",
			parsedmask.Username)
	} else if parsedmask.Host != "test" {
		t.Fatalf("Expected Host field to be \"test\", got \"%s\".",
			parsedmask.Host)
	}
}
func TestParseUserMaskInvalid1(t *testing.T) {
	unparsedmask := "test!x!x@client.test"
	var parsedmask IrcUserMask
	var err error

	parsedmask, err = ParseUserMask(unparsedmask)
	t.Logf("Returned error is: %+v", err)
	if err == nil {
		t.Fatalf("ParseUserMask should have failed an with error '%s'", err)
	}

	if parsedmask.Type != "" || parsedmask.Nick != "" || parsedmask.Username !=
		"" || parsedmask.Host != "" {
		t.Fatalf("ParseUserMark returned a non-empty IrcUserMask after an " +
			"error.")
	}

	if err.Error() != "User mask has multiple nick/user seperators." {
		t.Fatalf("Expected error to be \"User mask has multiple nick/user "+
			"seperators, got \"%s\"", err)
	}
}
func TestParseUserMaskInvalid2(t *testing.T) {
	unparsedmask := "te.st!x@client.test"
	var parsedmask IrcUserMask
	var err error

	parsedmask, err = ParseUserMask(unparsedmask)
	t.Logf("Returned error is: %+v", err)
	if err == nil {
		t.Fatalf("ParseUserMask should have failed an with error '%s'", err)
	}

	if parsedmask.Type != "" || parsedmask.Nick != "" || parsedmask.Username !=
		"" || parsedmask.Host != "" {
		t.Fatalf("ParseUserMark returned a non-empty IrcUserMask after an " +
			"error.")
	}

	if err.Error() != "Nickname contains dots, which isn't permitted." {
		t.Fatalf("Expected error to be \"Nickname contains dots, which isn't "+
			"permitted.\", got \"%s\"", err)
	}
}
func TestParseUserMaskInvalid3(t *testing.T) {
	unparsedmask := "test!x@client@client.test"
	var parsedmask IrcUserMask
	var err error

	parsedmask, err = ParseUserMask(unparsedmask)
	t.Logf("Returned error is: %+v", err)
	if err == nil {
		t.Fatalf("ParseUserMask should have failed an with error '%s'", err)
	}

	if parsedmask.Type != "" || parsedmask.Nick != "" || parsedmask.Username !=
		"" || parsedmask.Host != "" {
		t.Fatalf("ParseUserMark returned a non-empty IrcUserMask after an " +
			"error.")
	}

	if err.Error() != "User mask has multiple host seperators." {
		t.Fatalf("Expected error to be \"User mask has multiple host "+
			"seperators.\", got \"%s\"", err)
	}
}
func TestParseUserMaskInvalid4(t *testing.T) {
	unparsedmask := "test!x.boo@client.test"
	var parsedmask IrcUserMask
	var err error

	parsedmask, err = ParseUserMask(unparsedmask)
	t.Logf("Returned error is: %+v", err)
	if err == nil {
		t.Fatalf("ParseUserMask should have failed an with error '%s'", err)
	}

	if parsedmask.Type != "" || parsedmask.Nick != "" || parsedmask.Username !=
		"" || parsedmask.Host != "" {
		t.Fatalf("ParseUserMark returned a non-empty IrcUserMask after an " +
			"error.")
	}

	if err.Error() != "Username contains dots, which isn't permitted." {
		t.Fatalf("Expected error to be \"Username contains dots, which isn't "+
			"permitted.\", got \"%s\"", err)
	}
}
func TestParseUserMaskInvalid5(t *testing.T) {
	unparsedmask := "test!x boo@client.test"
	var parsedmask IrcUserMask
	var err error

	parsedmask, err = ParseUserMask(unparsedmask)
	t.Logf("Returned error is: %+v", err)
	if err == nil {
		t.Fatalf("ParseUserMask should have failed an with error '%s'", err)
	}

	if parsedmask.Type != "" || parsedmask.Nick != "" || parsedmask.Username !=
		"" || parsedmask.Host != "" {
		t.Fatalf("ParseUserMark returned a non-empty IrcUserMask after an " +
			"error.")
	}

	if err.Error() != "User mask contains reserved characters." {
		t.Fatalf("Expected error to be \"User mask contains reserved "+
			"characters.\", got \"%s\"", err)
	}
}

// ----------------------------------------------------------------------------
// PARSE USER MASK BENCHMARKS -------------------------------------------------
// ----------------------------------------------------------------------------

func BenchmarkParseUserMaskUserLong(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := ParseUserMask("TestUser!test@user.client.test")
		if err != nil {
			b.FailNow()
		}
	}
}
func BenchmarkParseUserMaskUserShort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := ParseUserMask("TestUser@user.client.test")
		if err != nil {
			b.FailNow()
		}
	}
}
func BenchmarkParseUserMaskServer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := ParseUserMask("server.test")
		if err != nil {
			b.FailNow()
		}
	}
}
func BenchmarkParseUserMaskAmbigious(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := ParseUserMask("test")
		if err != nil {
			b.FailNow()
		}
	}
}
