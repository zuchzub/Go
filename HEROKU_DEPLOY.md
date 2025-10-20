# Deploying TG Go Music Bot to Heroku

This guide will help you deploy the TG Go Music Bot to Heroku.

## Prerequisites

1. **Heroku Account**: Sign up at [heroku.com](https://heroku.com)
2. **Heroku CLI**: Install from [heroku.com/cli](https://devcenter.heroku.com/articles/heroku-cli)
3. **Git**: Ensure Git is installed on your system
4. **Required Credentials**:
   - Telegram Bot Token (from [@BotFather](https://t.me/BotFather))
   - API ID & API Hash (from [my.telegram.org](https://my.telegram.org))
   - Pyrogram String Session (generate using `python3 generate_pyrogram_session.py`)
   - MongoDB URL (free tier from [MongoDB Atlas](https://www.mongodb.com/cloud/atlas))

## Quick Deploy

### Method 1: One-Click Deploy

Click the button below to deploy directly to Heroku:

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)

Then fill in all the required environment variables.

### Method 2: Manual Deploy with Script

1. **Login to Heroku**:
   ```bash
   heroku login
   ```

2. **Make deploy script executable**:
   ```bash
   chmod +x deploy_heroku.sh
   ```

3. **Run the deployment script**:
   ```bash
   ./deploy_heroku.sh
   ```

The script will:
- Create a new Heroku app (or use existing one)
- Set up buildpacks
- Configure environment variables from `.env` file
- Deploy the code
- Scale the worker dyno

### Method 3: Manual Deploy Step by Step

1. **Login to Heroku**:
   ```bash
   heroku login
   ```

2. **Create a new Heroku app**:
   ```bash
   heroku create your-app-name
   ```

3. **Add buildpacks**:
   ```bash
   heroku buildpacks:add --index 1 https://github.com/heroku/heroku-buildpack-apt
   heroku buildpacks:add --index 2 heroku/go
   ```

4. **Set environment variables**:
   ```bash
   heroku config:set TOKEN="your_bot_token"
   heroku config:set API_ID="your_api_id"
   heroku config:set API_HASH="your_api_hash"
   heroku config:set STRING_SESSION="your_string_session"
   heroku config:set MONGODB_URL="your_mongodb_url"
   heroku config:set LOGGER_ID="your_telegram_id"
   heroku config:set OWNER_ID="your_telegram_id"
   ```

5. **Deploy the code**:
   ```bash
   git push heroku master
   ```
   
   Or if you're on a different branch:
   ```bash
   git push heroku your-branch:master
   ```

6. **Scale the worker dyno**:
   ```bash
   heroku ps:scale worker=1
   ```

## Configuration

### Required Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `TOKEN` | Telegram Bot Token from @BotFather | `1234567890:ABCdefGHIjklMNOpqrsTUVwxyz` |
| `API_ID` | API ID from my.telegram.org | `12345678` |
| `API_HASH` | API Hash from my.telegram.org | `abcdef1234567890abcdef1234567890` |
| `STRING_SESSION` | Pyrogram String Session | `BQHGAdUA...` |
| `MONGODB_URL` | MongoDB connection string | `mongodb+srv://user:pass@cluster.mongodb.net/` |
| `LOGGER_ID` | Telegram Chat ID for logs | `123456789` |
| `OWNER_ID` | Bot owner's Telegram User ID | `123456789` |

### Optional Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `API_URL` | API URL for YouTube search | `https://api.safone.dev` |
| `API_KEY` | API Key for YouTube search | `SAF_ONE` |
| `PPROF_PORT` | Port for pprof profiling | `5068` |
| `DOWNLOADS_DIR` | Directory for downloads | `/tmp/downloads` |

## Generate String Session

Before deploying, you need to generate a Pyrogram String Session:

```bash
python3 generate_pyrogram_session.py
```

Follow the prompts to:
1. Enter your API ID
2. Enter your API Hash
3. Enter your phone number (with country code)
4. Enter the verification code from Telegram
5. Enter your 2FA password (if enabled)

Save the generated string session for use in Heroku config.

## MongoDB Setup

1. Go to [MongoDB Atlas](https://www.mongodb.com/cloud/atlas)
2. Create a free cluster
3. Create a database user
4. Whitelist all IP addresses (0.0.0.0/0) for Heroku
5. Get your connection string and replace `<password>` with your database user password

## Useful Heroku Commands

### View Logs
```bash
heroku logs --tail --app your-app-name
```

### Restart App
```bash
heroku restart --app your-app-name
```

### Stop Bot
```bash
heroku ps:scale worker=0 --app your-app-name
```

### Start Bot
```bash
heroku ps:scale worker=1 --app your-app-name
```

### View Config Vars
```bash
heroku config --app your-app-name
```

### Set Config Var
```bash
heroku config:set KEY=VALUE --app your-app-name
```

### Open Heroku Dashboard
```bash
heroku open --app your-app-name
```

## Troubleshooting

### Bot not starting?
- Check logs: `heroku logs --tail --app your-app-name`
- Verify all environment variables are set correctly
- Ensure worker dyno is scaled to 1: `heroku ps --app your-app-name`

### "No such file or directory" errors?
- Make sure Aptfile includes all required dependencies (ffmpeg, python3)
- Verify buildpacks are in correct order

### Database connection issues?
- Check MongoDB connection string is correct
- Ensure IP whitelist includes 0.0.0.0/0 in MongoDB Atlas
- Verify database user credentials

### Out of memory?
- Upgrade to a higher dyno plan
- Monitor memory usage: `heroku ps --app your-app-name`

## Important Notes

⚠️ **Free Dyno Limitations**:
- Free dynos sleep after 30 minutes of inactivity
- You get 550-1000 free dyno hours per month
- Free dynos restart every 24 hours
- Consider upgrading to Basic or higher for production use

⚠️ **Files Storage**:
- Heroku uses ephemeral filesystem
- Downloaded files are deleted on dyno restart
- Use `/tmp` directory for temporary files

⚠️ **String Session Security**:
- Never commit `.env` file to Git
- Keep your string session private
- It has full access to your Telegram account

## Support

For issues or questions:
- Open an issue on GitHub
- Check the main [README.md](README.md)
- Review Heroku logs for error messages

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
