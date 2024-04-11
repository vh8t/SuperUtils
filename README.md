# SuperUtils

SuperUtils is a collection of enhanced Linux commands developed by vh8t. It aims to provide improved functionality and features compared to standard command-line tools.

## Supported Commands

### sls (Super ls)

```sh
sls
```

`sls` is an enhanced version of the `ls` command with icons and more readable output. It uses Nerd Font icons to enhance file representation in the terminal.

### ccat (Colorful cat)

```sh
ccat <file>
```

`ccat` is an improved version of the `cat` command that supports syntax highlighting using the Pygments library. It can display the contents of files with syntax colors for better readability.

More commands will be added to SuperUtils in future updates.

## Requirements

- Python 3.11+
- Pygments library (`pip install pygments`)
- Nerd Font as the terminal font for optimal display of icons

## Installation

To install SuperUtils, run the following script in your terminal:

```sh
curl -o- https://raw.githubusercontent.com/vh8t/SuperUtils/main/install.sh | bash
```

This script will set up the necessary dependencies and configurations for SuperUtils on your system.

## Usage

Once installed, you can use the supported commands (`sls`, `ccat`, etc.) directly from your terminal with enhanced features.

## License

This project is licensed under the [GNU General Public License v3.0](LICENSE).

---

**Note:** SuperUtils is a work in progress. Contributions and feedback are welcome!
