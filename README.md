# AI Shell Command Generator (aishell)

A command-line tool that converts natural language descriptions into shell commands using AI. Simply describe what you want to do, and aishell will generate the appropriate shell command.

## Features

- Convert natural language to shell commands
- Interactive mode with command prompt
- Command preview and confirmation before execution
- Command editing capability
- Verbose mode for detailed output
- Cross-platform support (with some limitations on Windows)

## Prerequisites

- Go 1.23.1 or higher
- OpenRouter API key (get it from https://openrouter.ai/)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/aishell.git
cd aishell
cp .env.example .env
```

2. Replace your-api-key-here from `OPENROUTER_API_KEY` in the `.env` file with your actual OpenRouter API key.
3. Build the project:
```bash
./scripts/build.sh
```
## Usage
Run the `aishell` command in your terminal:
```bash
./aishell
```
The tool will prompt you to enter a description of the command you want to generate. Simply type your description, and the tool will generate the corresponding shell command.
You can also use the interactive mode by running `aishell -i`. In this mode, you can enter multiple commands and see the generated commands one by one.
## Examples
```bash
./aishell
Describe the command to create a new directory: mkdir new_dir
Generated command: mkdir new_dir
Do you want to execute this command? (y/n): y
Command executed successfully.
```
## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details. 