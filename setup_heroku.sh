#!/bin/bash

# Quick setup script for preparing Heroku deployment

echo "=========================================="
echo "  TG Go Music Bot - Setup for Heroku"
echo "=========================================="
echo ""

# Check if .env exists
if [ ! -f ".env" ]; then
    echo "⚠️  .env file not found. Creating from sample.env..."
    if [ -f "sample.env" ]; then
        cp sample.env .env
        echo "✅ .env file created. Please edit it with your credentials."
    else
        echo "❌ sample.env not found!"
        exit 1
    fi
else
    echo "✅ .env file found"
fi

echo ""
echo "=========================================="
echo "  Checking Dependencies"
echo "=========================================="
echo ""

# Check Python
if command -v python3 &> /dev/null; then
    echo "✅ Python3: $(python3 --version)"
else
    echo "❌ Python3 not found. Please install Python 3.7+"
    exit 1
fi

# Check if pyrogram is installed
if python3 -c "import pyrogram" &> /dev/null 2>&1; then
    echo "✅ Pyrogram installed"
else
    echo "⚠️  Pyrogram not installed. Installing..."
    pip install pyrogram tgcrypto -q
    echo "✅ Pyrogram installed"
fi

# Check Heroku CLI
if command -v heroku &> /dev/null; then
    echo "✅ Heroku CLI: $(heroku --version | head -n1)"
else
    echo "⚠️  Heroku CLI not found"
    echo "   Install from: https://devcenter.heroku.com/articles/heroku-cli"
fi

# Check Git
if command -v git &> /dev/null; then
    echo "✅ Git: $(git --version)"
else
    echo "❌ Git not found. Please install Git"
    exit 1
fi

echo ""
echo "=========================================="
echo "  Configuration Check"
echo "=========================================="
echo ""

# Function to check if variable is set in .env
check_var() {
    local var_name=$1
    local var_value=$(grep "^${var_name}=" .env 2>/dev/null | cut -d'=' -f2- | tr -d '"' | tr -d "'")
    
    if [ -z "$var_value" ] || [ "$var_value" = "your_${var_name,,}" ]; then
        echo "❌ $var_name: Not configured"
        return 1
    else
        echo "✅ $var_name: Configured"
        return 0
    fi
}

# Check required variables
all_configured=true

check_var "TOKEN" || all_configured=false
check_var "API_ID" || all_configured=false
check_var "API_HASH" || all_configured=false
check_var "STRING_SESSION" || all_configured=false
check_var "MONGODB_URL" || all_configured=false
check_var "LOGGER_ID" || all_configured=false
check_var "OWNER_ID" || all_configured=false

echo ""

if [ "$all_configured" = false ]; then
    echo "⚠️  Some required variables are not configured in .env"
    echo ""
    echo "To generate STRING_SESSION, run:"
    echo "  python3 generate_pyrogram_session.py"
    echo ""
    echo "For MongoDB, you can use:"
    echo "  - MongoDB Atlas (free tier): https://www.mongodb.com/cloud/atlas"
    echo ""
    echo "Please edit .env file and configure all required variables."
    echo ""
    read -p "Do you want to generate STRING_SESSION now? (y/n): " generate_session
    
    if [ "$generate_session" = "y" ] || [ "$generate_session" = "Y" ]; then
        if [ -f "generate_pyrogram_session.py" ]; then
            python3 generate_pyrogram_session.py
        else
            echo "❌ generate_pyrogram_session.py not found!"
        fi
    fi
else
    echo "✅ All required variables are configured!"
    echo ""
    echo "=========================================="
    echo "  Ready for Deployment!"
    echo "=========================================="
    echo ""
    echo "Next steps:"
    echo ""
    echo "1. Commit your changes (if needed):"
    echo "   git add ."
    echo "   git commit -m 'Prepare for Heroku deployment'"
    echo ""
    echo "2. Deploy to Heroku:"
    echo "   ./deploy_heroku.sh"
    echo ""
    echo "   Or manually:"
    echo "   heroku login"
    echo "   heroku create your-app-name"
    echo "   git push heroku master"
    echo ""
fi
