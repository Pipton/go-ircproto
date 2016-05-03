# ircproto

[![Build Status](https://img.shields.io/travis/AelitaStyles/go-ircproto/master.svg)](https://travis-ci.org/AelitaStyles/go-ircproto)
[![Current Version](https://img.shields.io/github/release/AelitaStyles/go-ircproto.svg)](https://github.com/AelitaStyles/go-ircproto/releases)
[![License](https://img.shields.io/github/license/aelitastyles/go-ircproto.svg)](https://github.com/AelitaStyles/go-ircproto/blob/master/LICENSE.md)

Package ircproto provides a parser, generator and state tracker for the IRC
protocol. How you get the IRC command and what you do with the information
ircproto gives you is irrelevant; just drop a string into the parser, and
it'll hand back a neatly arranged structure (or an error).

When you want to send data to the IRC server, simply give the relevent ircproto
function what you want to send, and it'll hand back a string ready to go. Oh,
and it'll keep track of what channels you are subscribed to and who's in them.

Neat!

**NOT CURRENTLY IN A FUNCTIONING STATE;** come back later!
