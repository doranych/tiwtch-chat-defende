## Twitch Chat Defender
Created for streamers protection from "shit" raids.

Bot executes specified command with list of bots. Reminds moderator permissions 

Based on [hackbolt/twitchbotsnbigots](https://github.com/hackbolt/twitchbotsnbigots) collecting list of bots accounts.

### How to use?
1. Download binary for your system from release page
2. Place `config.yaml` right near executable
3. Fill `config.yaml` using any suitable text editor
4. Run it

### Configuration
```yaml
twitch:
  username: <bot username>
  token: '<bot token>'
  channel: <channel>
defender:
  command: /ban {username}
  repoUrl: https://raw.githubusercontent.com/hackbolt/twitchbotsnbigots/main/
  syncInterval: 10m
```
Fields description:
```
  username     - username for bot used to authorize via twitch chat.
  token        - authorization token. Can be acquired from https://twitchapps.com/tmi/
  channel      - channel to moderate
  command      - command, that will be executed. Replace `{username}` with actual bot's username
  repoUrl      - RAW link to repository with list of bot names. 
                 Use this https://raw.githubusercontent.com/hackbolt/twitchbotsnbigots/main/
  syncInterval - list reload interval. Show how ofter it will load data from repository. 
                 Must be valid formatted. Required format: `number[s|m|h]`.
```

## Contribution
Contribution is allowed and appreciated
