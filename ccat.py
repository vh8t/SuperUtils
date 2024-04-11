#!/usr/bin/env python
from argparse import ArgumentParser
from os import path

from pygments.lexers import get_lexer_for_filename
from pygments.formatters import TerminalFormatter
from pygments import highlight


RED = '\033[31m'
RESET = '\033[0m'


def cat_file(pth, show_ends, number, squeeze_blank, show_tabs) -> None:
    try:
        lexer = get_lexer_for_filename(pth, stripall=True)

        with open(pth, 'r') as f:
            content = f.read()

        highlighted = highlight(content, lexer, TerminalFormatter())

        if show_ends:
            highlighted = highlighted.replace('\n', '$\n')
        if show_tabs:
            highlighted = highlighted.replace('\t', '^I')

        lines = highlighted.split('\n')
        uneditted = content.split('\n')

        if number:
            padding = len(str(len(lines)))
            for i, line in enumerate(lines):
                lines[i] = str(i).ljust(padding + 3) + line

        if squeeze_blank:
            highlighted = []
            previous = False

            for i, line in enumerate(lines):
                try:
                    if uneditted[i].strip() == '' and uneditted[i + 1].strip() == '':
                        continue
                    highlighted.append(line)
                except IndexError:
                    continue
            lines = highlighted

        highlighted = '\n'.join(lines)

        print(highlighted)
    except Exception as e:
        print(f'{RED}ccat: {e}{RESET}')
        return




def print_help() -> None:
    print('''Usage: ccat [flags] [file]

Show the contents of a file with syntax highlighting

positional arguments:
  file                 file path

options:
  -e, --show-ends      display $ at the end of each line
  -n, --number         number all output lines
  -s, --squeeze-blank  supress repeated empty output lines
  -t, --show-tabs      display TAB characters as ^I
      --help           show this help message and exit
      --version        output version information and exit

This command is part of the SuperUtils collection (ccat - colorful cat)
Copyright (C) 2024 vh8t
Author: vh8t
GitHub: https://github.com/vh8t
Website: https://vh8t.xyz
    ''')
    exit(0)


def main() -> None:
    parser = ArgumentParser(add_help=False)

    parser.add_argument('file', nargs='?')
    parser.add_argument('--help', action='store_true')
    parser.add_argument('--version', action='store_true')
    parser.add_argument('-e', '--show-ends', action='store_true')
    parser.add_argument('-n', '--number', action='store_true')
    parser.add_argument('-s', '--squeeze-blank', action='store_true')
    parser.add_argument('-t', '--show-tabs', action='store_true')

    args = parser.parse_args()
    if args.help:
        print_help()
    elif args.version:
        print('ccat (SuperUtils) 1.0\nCopyright (C) 2024 vh8t')
        exit(0)
    else:
        cat_file(args.file, args.show_ends, args.number, args.squeeze_blank, args.show_tabs)


main()
