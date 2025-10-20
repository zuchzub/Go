#!/bin/bash

# Script to deploy TG Go Music Bot to Heroku

echo "=========================================="
echo "  TG Go Music Bot - Heroku Deployment"
echo "=========================================="
echo ""

# Check if Heroku CLI is installed
if ! command -v heroku &> /dev/null; then
    echo "❌ Heroku CLI is not installed."
    echo "   Please install it from: https://devcenter.heroku.com/articles/heroku-cli"
    exit 1
fi

echo "✅ Heroku CLI detected"
echo ""

# Check if user is logged in to Heroku
if ! heroku auth:whoami &> /dev/null; then
    echo "❌ You are not logged in to Heroku."
    echo "   Please run: heroku login"
    exit 1
fi

echo "✅ Logged in to Heroku as: $(heroku auth:whoami)"
echo ""

# Ask for app name
read -p "Enter your Heroku app name (or press Enter to create a new one): " APP_NAME

if [ -z "$APP_NAME" ]; then
    echo "Creating a new Heroku app..."
    heroku create
    APP_NAME=$(heroku apps:info -j | grep -o '"name":"[^"]*' | cut -d'"' -f4)
else
    # Check if app exists
    if heroku apps:info --app "$APP_NAME" &> /dev/null; then
        echo "✅ Using existing app: $APP_NAME"
    else
        echo "Creating app: $APP_NAME"
        heroku create "$APP_NAME"
    fi
fi

echo ""
echo "=========================================="
echo "  Setting up buildpacks..."
echo "=========================================="

# Add buildpacks
heroku buildpacks:clear --app "$APP_NAME"
heroku buildpacks:add --index 1 https://github.com/heroku/heroku-buildpack-apt --app "$APP_NAME"
heroku buildpacks:add --index 2 heroku/go --app "$APP_NAME"

echo ""
echo "=========================================="
echo "  Configuring environment variables..."
echo "=========================================="
echo ""

# Check if .env file exists
if [ -f ".env" ]; then
    echo "Found .env file. Setting config vars..."
    
    # Read .env and set config vars
    while IFS='=' read -r key value; do
        # Skip empty lines and comments
        if [ -z "$key" ] || [[ "$key" =~ ^#.* ]]; then
            continue
        fi
        
        # Remove quotes from value
        value=$(echo "$value" | sed -e 's/^"//' -e 's/"$//' -e "s/^'//" -e "s/'$//")
        
        echo "Setting $key..."
        heroku config:set "$key=$value" --app "$APP_NAME" > /dev/null
    done < .env
    
    echo "✅ Environment variables configured"
else
    echo "⚠️  No .env file found. You'll need to set environment variables manually."
    echo "   Run: heroku config:set KEY=VALUE --app $APP_NAME"
fi

echo ""
echo "=========================================="
echo "  Deploying to Heroku..."
echo "=========================================="
echo ""

# Add Heroku remote if it doesn't exist
if ! git remote | grep -q heroku; then
    heroku git:remote --app "$APP_NAME"
fi

# Deploy
echo "Pushing code to Heroku..."
git push heroku master

echo ""
echo "=========================================="
echo "  Scaling worker dyno..."
echo "=========================================="
echo ""

heroku ps:scale worker=1 --app "$APP_NAME"

echo ""
echo "=========================================="
echo "  ✅ Deployment Complete!"
echo "=========================================="
echo ""
echo "App URL: https://$APP_NAME.herokuapp.com"
echo ""
echo "Useful commands:"
echo "  - View logs: heroku logs --tail --app $APP_NAME"
echo "  - Restart app: heroku restart --app $APP_NAME"
echo "  - Scale worker: heroku ps:scale worker=1 --app $APP_NAME"
echo "  - Stop worker: heroku ps:scale worker=0 --app $APP_NAME"
echo ""
