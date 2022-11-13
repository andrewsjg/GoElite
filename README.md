# Go Elite

[![Go Report Card](https://goreportcard.com/badge/github.com/andrewsjg/GoElite)](https://goreportcard.com/report/github.com/andrewsjg/GoElite)
![tests](https://github.com/andrewsjg/GoElite/workflows/tests/badge.svg)

TxtElite Implemented in Go. Including a terminal UI for playing the trading game.

Based off of [Ian Bell's Text Elite](http://www.iancgbell.clara.net/elite/text/)

## Installing

**NOTE:** txtelite looks best when run in a terminal with a dark background. In future I will add colour schemes to suit light terminals as well.

### From Source

This version of txtelite has been tested on macOS versions 13 (Ventura) and 12 (Monterey) and on Ubuntu Linux 22.04.

It will also compile and run on Windows, however the terminal UI doesnt work well in the Windows shell. There is bound to be a fix for this, and I will get to it eventually.

To build from source, ensure the latest version of `go` is installed. Any versions > `go v1.18` will work.

- Clone this repo: `git clone https://github.com/andrewsjg/GoElite.git`
- Run `make build`
- The binary is: `.\bin\txtelite`

### Download the latest release build

Go to the [releases](https://github.com/andrewsjg/GoElite/releases) page in this repo, download the binary for your system, place it in the path and run `txtelite`

### Via Homebrew on macOS

- Add my TAP to homebrew: `brew tap andrewsjg/tap`
- Install `txtelite`: `brew install txtelite`



## Todo

- ~~Market Command~~
- ~~Local Command~~
- ~~Buy Command~~
- ~~Sell Command~~
- ~~Fuel Command~~
- ~~Jump Command~~
- ~~Info Command~~
- ~~Local Command~~
- ~~Hyperspace command~~
- ~~Help Command~~
- ~~TUI Status bar~~
- ~~Add ship info panel~~
- ~~Add ships hold table~~
- ~~Tidy up basic TUI~~
- ~~Fix bug that means market isnt generated on hyperspace jump~~
- ~~Fix fuel bug that allows a player to buy more than max fuel~~
- ~~Fix buy/sell commands where the commodity has a space in its name~~
- ~~Check Fuel maths. Strange things happen when buying fuel~~
- Tidy up the game title. Styling with Lipgloss? - Sort of done. Not sure I like it yet
- ~~Fix variadic buy functions. Used incorrectly.~~
- Basic Commander rank info
- Basic Commander name function
- Styling for light terminals
- Check for terminal width
- Fix Local output display
- Add webpage for homebrew
- ~~Add version option. Do I need Cobra?~~
- Fix TUI when running on windows

## Todo - Sometime

- Improved TUI
- Game saves
- Score
- "Improved" economy

## Aspirational

- 2D game engine
- Basic combat
