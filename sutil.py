#!/usr/bin/env python
from argparse import ArgumentParser
from os import path


def setup() -> None:
    command_aliases = {
        'sls': 'ls',
        'ccat': 'cat'
    }
    
    start_setup = input("Do you want to set up command aliases? (yes/no): ").lower()
    if start_setup != 'yes':
        print("Setup aborted.")
        return
    
    selected_commands = []
    for command, alias in command_aliases.items():
        response = input(f"Do you want to alias '{command}' to '{alias}'? (yes/no): ").lower()
        if response == 'yes':
            selected_commands.append((command, alias))
    
    if not selected_commands:
        print("No commands selected. Setup aborted.")
        return
    
    available_profiles = []
    profile_filenames = [
        ".bashrc",
        ".bash_profile",
        ".zshrc",
        ".profile",
        ".bash_aliases",
        ".config/fish/config.fish"
    ]
    
    for filename in profile_filenames:
        full_path = path.expanduser(f"~/{filename}")
        if path.exists(full_path):
            available_profiles.append(full_path)
    
    if not available_profiles:
        print("No terminal profile files found. Setup aborted.")
        return
    
    print("Available terminal profiles:")
    for i, profile in enumerate(available_profiles, start=1):
        print(f"{i}. {profile}")
    
    choice = input(f"Enter the number(s) of the profile(s) to modify (e.g., '1' or '1,2' for all): ")
    selected_profiles = []
    
    try:
        profile_indices = [int(index.strip()) - 1 for index in choice.split(",")]
        selected_profiles = [available_profiles[index] for index in profile_indices]
    except (ValueError, IndexError):
        print("Invalid input. Setup aborted.")
        return
    
    if not selected_profiles:
        print("No profile selected. Setup aborted.")
        return
    
    for profile in selected_profiles:
        with open(profile, 'a') as file:
            file.write("\n# Aliases added by SuperUtils\n")
            for command, alias in selected_commands:
                file.write(f"alias {alias}='{command}'\n")
        
        print(f"Aliases added to {profile}")

    print("Setup completed successfully.")


def print_help() -> None:
    print('''Usage: sutil [flags]

List information about file (the current directory by default)

options:
  --help     show this help message and exit
  --version  output version information and exit

This command is part of the SuperUtils collection (sutil - super util)
Copyright (C) 2024 vh8t
Author: vh8t
GitHub: https://github.com/vh8t
Website: https://vh8t.xyz
    ''')
    exit(0)


def main() -> None:
    parser = ArgumentParser(add_help=False)

    parser.add_argument('-s', '--setup', action='store_true')
    parser.add_argument('--help', action='store_true')
    parser.add_argument('--version', action='store_true')

    args = parser.parse_args()
    if args.setup:
        setup()
    elif args.version:
        print('SuperUtils 1.0\nCopyright (C) 2024 vh8t')
        exit(0)
    else:
        print_help()


main()
