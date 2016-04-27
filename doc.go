// Copyright (c) 2016, Aelita Styles
//
// This file is part of the go-ircproto project, and is licensed under a
// BSD-like license. A copy of the license should have been included with this
// source code, but in the event that it was not, you can find it here:
// https://github.com/AelitaStyles/go-ircproto/blob/master/LICENSE.md

// Package ircproto provides a parser, generator and state tracker for the IRC
// protocol. How you get the IRC command and what you do with the information
// ircproto gives you is irrelevant; just drop a string into the parser, and
// it'll hand back a neatly arranged structure (or an error).
//
// When you want to send data to the IRC server, simply give the relevent
// ircproto function what you want to send, and it'll hand back a string ready
// to go. Oh, and it'll keep track of what channels you are subscribed to and
// who's in them. Neat!
//
// NOT CURRENTLY IN A FUNCTIONING STATE; come back later!
package ircproto
