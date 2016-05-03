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
func testBadParseRaw(testCmd string, expectedErrorMsg string) (error, error) {
	var parsedCmd IrcCommand
	var err error
	parsedCmd, err = ParseRaw(testCmd)

	if !reflect.DeepEqual(parsedCmd, IrcCommand{}) {
		return err, fmt.Errorf("ParseRaw returned a non-empty IrcCommand "+
			"object, which wasn't expected. The returned object is: %+v",
			parsedCmd)
	} else if err == nil {
		return err, fmt.Errorf("ParseRaw returned an empty error, which "+
			"wasn't expected.")
	} else if err.Error() != expectedErrorMsg {
		return err, fmt.Errorf("ParseRaw returned an unexpected error of '%s'."+
			" We were expecting an error of '%s'.", err.Error(),
			expectedErrorMsg)
	}

	return err, nil
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
	prErr, testErr := testBadParseRaw(":x!x!foo@test PRIVMSG TestUser "+
		":Hello\r\n", "User mask has multiple nick/user seperators.")

	if testErr != nil {
		t.Error(testErr)
	}
	t.Logf("ParseRaw returned error: %+v", prErr)
}
func TestParseRawInvalid2(t *testing.T) {
	prErr, testErr := testBadParseRaw(":server.test 0a0 TestUser :Hello\r\n",
		"Command type is not a valid numeric.")

	if testErr != nil {
		t.Error(testErr)
	}
	t.Logf("ParseRaw returned error: %+v", prErr)
}
func TestParseRawInvalid3(t *testing.T) {
	prErr, testErr := testBadParseRaw(":server.test M30W TestUser :Hello\r\n",
		"Command type contains invalid characters.")

	if testErr != nil {
		t.Error(testErr)
	}
	t.Logf("ParseRaw returned error: %+v", prErr)
}
func TestParseRawInvalid4(t *testing.T) {
	prErr, testErr := testBadParseRaw(":server.test 001 TestUser :Hello",
		"Command does not end with CRLF.")

	if testErr != nil {
		t.Error(testErr)
	}
	t.Logf("ParseRaw returned error: %+v", prErr)
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

func testGoodParseUserMask(testMask, expectedType, expectedNick, expectedUser,
	expectedHost string) (IrcUserMask, error){
	// Parse the test mask
	var parsedMask IrcUserMask
	var err error
	parsedMask, err = ParseUserMask(testMask)

	// Check for unexpected error
	if err != nil {
		return parsedMask, fmt.Errorf("ParseUserMask failed unexpected with "+
			"error: %s", err.Error())
	}

	// Check if we got expected values
	if parsedMask.Type != expectedType {
		return parsedMask, fmt.Errorf("ParseUserMask returned unexpected type "+
			"value of '%s'. We were expecting '%s'.", parsedMask.Type,
			expectedType)
	} else if parsedMask.Nick != expectedNick {
		return parsedMask, fmt.Errorf("ParseUserMask returned unexpected nick "+
			"value of '%s'. We were expecting '%s'.", parsedMask.Nick,
			expectedNick)
	} else if parsedMask.Username != expectedUser {
		return parsedMask, fmt.Errorf("ParseUserMask returned unexpected user "+
			"value of '%s'. We were expecting '%s'.", parsedMask.Username,
			expectedUser)
	} else if parsedMask.Host != expectedHost {
		return parsedMask, fmt.Errorf("ParseUserMask returned unexpected host "+
			"value of '%s'. We were expecting '%s'.", parsedMask.Host,
			expectedHost)
	}

	return parsedMask, nil
}

func TestParseUserMaskValidUser1(t *testing.T) {
	parsedMask, err := testGoodParseUserMask("TestUser!test@user.client.test",
		"User", "TestUser", "test", "user.client.test")

	if err != nil {
		t.Error(err)
	}
	t.Logf("Returned structure is: %+v", parsedMask)
}
func TestParseUserMaskValidUser2(t *testing.T) {
	parsedMask, err := testGoodParseUserMask("TestUser@user.client.test",
		"User", "TestUser", "", "user.client.test")

	if err != nil {
		t.Error(err)
	}
	t.Logf("Returned structure is: %+v", parsedMask)
}
func TestParseUserMaskValidServer1(t *testing.T) {
	parsedMask, err := testGoodParseUserMask("server.test",
		"Server", "", "", "server.test")

	if err != nil {
		t.Error(err)
	}
	t.Logf("Returned structure is: %+v", parsedMask)
}
func TestParseUserMaskValidAmbigiousMask1(t *testing.T) {
	parsedMask, err := testGoodParseUserMask("test",
		"Unknown", "test", "", "test")

	if err != nil {
		t.Error(err)
	}
	t.Logf("Returned structure is: %+v", parsedMask)
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
