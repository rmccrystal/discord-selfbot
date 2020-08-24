# Discord Selfbot
A Discord selfbot written in GoLang

## Getting your token
* On your Discord, press Ctrl + Shift + I to open devtools
* Go to the network tab on the top
* Click on any request on the left
* Under `Request Headers` find the `authorization` key
* Copy the value

[Screenshot](https://prnt.sc/u4pfns)

## Installation
### Using `go get`
If not already installed, download [Go](https://golang.org/dl/)

Download and install the repo
```
go get github.com/rmccrystal/discord-selfbot/cmd/selfbot
```
Run the bot
```
selfbot
```
Enter your auth key when prompted. It will automatically generate a config
file with your token so you don't have to enter it every time. Make sure
you are in the same directory every time you run it.

## List of commands
### `delete`

Description: `Deletes the past n commands`

Usage: `delete [messageCount]`

Aliases: `d`

### `ascii`

Description: `Turns your text to ascii art and prints it to the chat`

Usage: `ascii <text>`

### `react`

Description: `Reacts custom text to the most recent message`

Usage: `react <text>`

Aliases: `r`

### `clear`

Description: `Sends x blank lines to the channel. Defaults to 60`

Usage: `clear [lines]`

Aliases: `c`

### `help`

Description: `Prints help about all commands or about a specific command`

Usage: `help [command]`

Aliases: `h`

## Using multiple configs
If you would like to host multiple users at once, change your config.json
file to an array of config objects. For example:

```json
[
  {
    "token": "token1"
  },
  {
    "token": "token2",
    "prefix": "!"
  },
  {
    "token": "token3",
  }
]
```

This would create three selfbot instances with the three specified configs.