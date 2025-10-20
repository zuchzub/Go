#!/usr/bin/env python3
"""
Script untuk generate Pyrogram String Session
"""

from pyrogram import Client
import sys

def main():
    print("=" * 60)
    print("Pyrogram String Session Generator")
    print("=" * 60)
    print()
    
    # Input API credentials
    try:
        api_id = int(input("Enter your API ID: "))
        api_hash = input("Enter your API Hash: ")
        phone = input("Enter your phone number (with country code, e.g., +628123456789): ")
    except ValueError:
        print("Error: API ID must be a number")
        sys.exit(1)
    
    print("\nStarting authentication process...")
    print("You will receive a verification code on Telegram.")
    print()
    
    # Create client and get string session
    with Client(
        name="my_account",
        api_id=api_id,
        api_hash=api_hash,
        phone_number=phone,
        in_memory=True
    ) as app:
        string_session = app.export_session_string()
        
        print("\n" + "=" * 60)
        print("Your Pyrogram String Session:")
        print("=" * 60)
        print(string_session)
        print("=" * 60)
        print("\n⚠️  IMPORTANT: Save this string session securely!")
        print("    Do not share it with anyone!")
        print("    It can be used to access your Telegram account.")
        print()

if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print("\n\nOperation cancelled by user.")
        sys.exit(0)
    except Exception as e:
        print(f"\n❌ Error: {e}")
        sys.exit(1)
