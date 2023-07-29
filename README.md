# 🤖 Discord Reply Bot

[![License: GPL-3.0](https://img.shields.io/badge/License-GPL%203.0-blue.svg)](https://opensource.org/licenses/GPL-3.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/Tobias-Pe/discord-reply-bot)](https://goreportcard.com/report/github.com/Tobias-Pe/discord-reply-bot)
[![Gitmoji](https://img.shields.io/badge/gitmoji-%20😜%20😍-FFDD67.svg)](https://gitmoji.dev)

> ✨ A feature-rich Discord bot that allows users to set up automatic replies based on message matches.

## 📝 Description

🤖 Discord Reply Bot is a feature-rich Discord bot built with Go and powered by the [discordgo](https://github.com/bwmarrin/discordgo) library. The bot allows users on the server to set up exact or occurrence matches for messages and associate them with specific responses. When a match is triggered, one of the possible responses is chosen randomly. This makes it easy for server members to customize and maintain automatic replies for certain patterns or keywords in the chat.

Additionally, the bot provides a convenient feature where it suggests the last 10 messages in the chat, helping users to craft their reply setups effectively.

## ❗ Known Problems

- Suggestions / options-choices are limited to 100 characters, but commands can still be executed even if not suggested for certain values.
- [Issue with Discord API caching](https://github.com/discord/discord-api-docs/discussions/5121): Discord client caches the selected choices, so if a mistake is made and something is changed in other options, it won't be updated in later steps. One must rewrite the command to clear the cache.

## 🚀 Features

- User-defined exact or occurrence matches with custom responses.
- Randomized selection of responses for matches.
- Simple Discord bot commands for easy setup and maintenance.
- Suggestions for the last 10 messages in the chat to aid in reply configuration.
- In-memory caching of keys for faster query responses from the Redis database.

## 🛠️ Installation and Configuration

### Using Docker Compose

1. Clone the repository: `git clone https://github.com/Tobias-Pe/discord-reply-bot.git`
2. Navigate to the project folder: `cd discord-reply-bot`
3. Create a `.env` file in the root directory with the following content:

    ```
    DISCORD_BOT_TOKEN=YOUR_DISCORD_BOT_TOKEN
    DISCORD_GUILD_ID=YOUR_DISCORD_GUILD_ID
    REDIS_URL=YOUR_REDIS_URL
    ```
    Replace `YOUR_DISCORD_BOT_TOKEN`, `YOUR_DISCORD_GUILD_ID`, and `YOUR_REDIS_URL` with appropriate values.

4. Deploy the bot using Docker Compose: `docker-compose up -d`

### Standalone Launch

1. Clone the repository: `git clone https://github.com/Tobias-Pe/discord-reply-bot.git`
2. Navigate to the project folder: `cd discord-reply-bot`
3. Build the bot: `go build -o bot`
4. Run the bot with the necessary arguments:

    ```bash
    ./bot -token YOUR_DISCORD_BOT_TOKEN -guild YOUR_DISCORD_GUILD_ID -redis YOUR_REDIS_URL
    ```
    
    Replace `YOUR_DISCORD_BOT_TOKEN`, `YOUR_DISCORD_GUILD_ID`, and `YOUR_REDIS_URL` with appropriate values.

### Logging Level

You can customize the logger in the `internal/logger/logger.go` file. By default, the logger is set to `zap.NewProduction()` which provides a production-ready configuration. However, you can adjust it according to your desired logging settings by modifying the logger initialization.

## ⚙️ Commands

The bot supports the following commands:

- `add_reply`: Adds a new reply for a specific key.
- `remove_key`: Removes a key and its associated replies from the database.
- `remove_reply`: Removes a reply from a specific key.
- `list_replies`: Lists all the replies associated with a specific key.
- `edit_reply`: Edits an existing reply for a specific key.
- `edit_key`: Edits an existing key and its associated replies.

## 🤝 Contributing

👍 Contributions are welcome! If you find any issues or have ideas for improvements, feel free to open an issue or submit a pull request.

## 📄 License

This project is licensed under the [GNU General Public License v3.0](LICENSE).

