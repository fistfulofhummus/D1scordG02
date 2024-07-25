Discord Botnet :)

# Why ?

Decided to bounce off of https://github.com/emmaunel/DiscordGo.git and modify it into something similar that fits my specific needs.
I also found that a discord c2 would be kinda funny as well.


# Installation

To use DiscordGo, you need to create a Discord bot and a Discord server. After that, invite the bot to your server.

Click [here](https://support.discord.com/hc/en-us/articles/204849977-How-do-I-create-a-server-) to learn how to create a server and [here](https://discordjs.guide/preparations/setting-up-a-bot-application.html#creating-your-bot) to create a bot. And finally, learn to invite the bot to your server with [this.](https://discordjs.guide/preparations/adding-your-bot-to-servers.html#bot-invite-links)

When creating the bot, you need it give it some permission. For testing, I gave the bot full `administrative` permission. But the required permission are as follow:

* Send Messages
* Read Messages
* Attach Files
* Manage Server

# Usage

Edit this file `pkg/util/variables.go` with a url that gets a json file containing the serverID and botToken remotely. This is done so that the malware does not contain any hard coded secrets.

The bot token can be found on discord developer dashboard where you created the bot. To get your server ID, go to your server setting and click on `widget`. On the right pane, you see the your ID.

After that is done, all you have to do is run `make`. This will create binaries in the /bin folder. The windows binary is built with garble and normally as well for the sake of variety.

## Organizer Bot

When you have target connecting back to your discord server, channels are created by their ip addresses. This can quickly get hard to manage. Solution: Another bot to organize the targets channels.

To use the organizer bot, run the csv generator script in the scripts folder:
```
$ pip3 install -r requirements.txt
$ python3 csv_generator.py
```

This will create a csv like this:

```
192168185200,team01,hostname1,windows
192168185201,team02,hostname2,linux
```

To start the organizer bot: `go run cmd/organizer/main.go -f <csv_filename>.csv`

Run `clean` in any channel to organize bots into their respective categories.

# WIP (Work in Progress)

- [x] Use windows syscalls
- [x] DDOS capability (GET and POST)
- [x] Encrypted File Transfer
- [ ] Dumping LSASS
- [ ] Keylogging

# Disclamers
The author is in no way responsible for any illegal use of this software. It is provided purely as an educational proof of concept. I am also not responsible for any damages or mishaps that may happen in the course of using this software. Use at your own risk.

Every message on discord are saved on Discord's server, so be careful and not upload any sensitive or confidential documents.

# Used Libraries
* [discordgo](https://github.com/bwmarrin/discordgo)
* [garble](https://github.com/burrowers/garble.git)


Inspired by/mostly ripped from [emmaunel](https://github.com/emmaunel)
