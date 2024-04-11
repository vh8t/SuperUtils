#!/usr/bin/env python
from argparse import ArgumentParser
from os import listdir, path, stat
from pwd import getpwuid
from grp import getgrgid
from datetime import datetime
import stat as _stat

FILE_MAP = {
    'c': '\ue61e',
    'cpp': '\ue61d',
    'java': '\ue738',
    'py': '\ue73c',
    'js': '\uf2ef',
    'ts': '\ue69d',
    'html': '\ue736',
    'css': '\ue749',
    'php': '\ue73d',
    'rb': '\ue739',
    'swift': '\ue755',
    'go': '\ue65e',
    'rust': '\ue7a8',
    'rs': '\ue7a8',
    'dart': '\ue64c',
    'kt': '\ue634',
    'csharp': '\ue648',
    'cs': '\ue648',
    'lua': '\ue620',
    'perl': '\ue769',
    'sh': '\ue691',
    'ps1': '\uebc7',
    'asm': '\ue637',
    'json': '\ueb0f',
    'xml': '\ue619',
    'toml': '\ue6b2',
    'md': '\ue73e',
    'txt': '\uf15c',
    'rst': '\uf15c',
    'tex': '\uf15c',
    'csv': '\ueefc',
    'ini': '\ue615',
    'cfg': '\ue615',
    'conf': '\ue615',
    'properties': '\ue615',
    'env': '\ue615',
    'sql': '\uf1c0',
    'db': '\uf1c0',
    'sqlite': '\ue7c4',
    'xls': '\ue6a6',
    'xlsx': '\ue6a6',
    'gitignore': '\ue702',
    'gitattributes': '\ue702',
    'png': '\uf03e',
    'jpg': '\uf03e',
    'gif': '\uf03e',
    'svg': '\uf03e',
    'mp4': '\uf1c8',
    'mp3': '\uf1c7',
    'makefile': '\ue673',
    'gradle': '\ue660',
    'maven': '\ue674',
    'pkg': '\uf487',
    'deb': '\ue77d',
    'pdf': '\uf1c1',
    'rtf': '\uf15c',
    'zip': '\uf1c6',
    'tar': '\uf1c6',
    'gz': '\uf1c6',
    'bz2': '\uf1c6',
    'dockerfile': '\uf21f',
    'ejs': '\ue618',
    'twig': '\ue61c',
    'pug': '\ue686',
    'vue': '\ue6a0',
    'psd': '\ue7b8',
    'ai': '\ue7b4',
    'sketch': '\uef64',
    'unity': '\ue721',
    'prefab': '\ue721',
    'editorconfig': '\ue652',
    'npmrc': '\ue71e',
    'yarnrc': '\ue6a7',
    'babelrc': '\ue639',
    'eslintrc': '\ue655',
    'lock': '\uf023',
    'key': '\ueb11',
    'ttf': '\uf031',
    'otf': '\uf031',
    'woff': '\uf031',
    'sav': '\uf0c7',
    'zsh': '\ue691',
    'vim': '\ue7c5'
}

RED = '\033[31m'
GREEN = '\033[32m'
BLUE = '\033[34m'
RESET = '\033[0m'


def list_dir(pth='.', show_hidden=False, human_readable=False) -> None:
    try:
        if path.isdir(pth):    
            entries = listdir(pth)

            entries.insert(0, '.')
            entries.insert(1, '..')

            if not show_hidden:
                entries = [entry for entry in entries if not entry.startswith('.')]

            entries = sorted(entries, key = lambda x: x.lower())
        else:
            entries = [pth]

        metadata = []
        for entry in entries:
            if path.isdir(pth):
                full_path = path.join(pth, entry)
            else:
                full_path = pth
            status = stat(full_path)
        
            mode = status.st_mode
            permissions = _get_permissions(mode)
            nlink = status.st_nlink
            user = getpwuid(status.st_uid).pw_name
            group = getgrgid(status.st_gid).gr_name
        
            size = status.st_size
            if human_readable:
                size = _format_size(size)

            modified_time = datetime.fromtimestamp(status.st_mtime).strftime('%b %d %H:%M')

            metadata.append((permissions, str(nlink), user, group, str(size), modified_time, entry))

        max_lengths = [max(len(meta[i]) for meta in metadata) for i in range(7)]
        for meta in metadata:
            line = ''
            is_exe = False
            is_dir = False
            for i, m in enumerate(meta):
                if i == 0:
                    if m[3] == 'x' or m[6] == 'x' or m[9] == 'x':
                        is_exe = True
                    if m[0] != '-':
                        is_dir = True

                if i == 6:
                    try:
                        if m.startswith('.bash'):
                            line += '   \uebca'
                        elif m.startswith('.zsh'):
                            line += '   \ue691'
                        elif m.startswith('.vim'):
                            line += '   \ue7c5'
                        else:
                            if is_exe and not is_dir:
                                line += f'   {GREEN}{FILE_MAP[m.split(".")[-1]]}'
                            else:
                                line += f'   {FILE_MAP[m.split(".")[-1]]}'
                    except KeyError:
                        if is_dir:
                            if len(listdir(m)) != 0:
                                line += f'   {BLUE}\uf07b'
                            else:
                                line += f'   {BLUE}\uf114'
                        else:
                            line += '   \ue64e'

                if i != 4:
                    if i == 6:
                        line += f' {m.ljust(max_lengths[i])}{RESET}'
                    elif i != 0:
                        line += f'   {m.ljust(max_lengths[i])}{RESET}'
                    else:
                        line += f'{m.ljust(max_lengths[i])}{RESET}'
                else:
                    line += f'   {m.rjust(max_lengths[i])}{RESET}' if i != 0 else f'{m.rjust(max_lengths[i])}{RESET}'
            print(line)
    except Exception as e:
        print(f'{RED}sls: {e}{RESET}')
        return


def _get_permissions(mode) -> str:
    if _stat.S_ISDIR(mode):
        file_type = 'd'
    elif _stat.S_ISREG(mode):
        file_type = '-'
    elif _stat.S_ISLNK(mode):
        file_type = 'l'
    elif _stat.S_ISFIFO(mode):
        file_type = '|'
    elif _stat.S_ISSOCK(mode):
        file_type = 's'
    elif _stat.S_ISCHR(mode):
        file_type = 'c'
    elif _stat.S_ISBLK(mode):
        file_type = 'b'
    else:
        file_type = '?'

    perms = [
        ('r' if mode & _stat.S_IRUSR else '-'),
        ('w' if mode & _stat.S_IWUSR else '-'),
        ('x' if mode & _stat.S_IXUSR else '-'),
        ('r' if mode & _stat.S_IRGRP else '-'),
        ('w' if mode & _stat.S_IWGRP else '-'),
        ('x' if mode & _stat.S_IXGRP else '-'),
        ('r' if mode & _stat.S_IROTH else '-'),
        ('w' if mode & _stat.S_IWOTH else '-'),
        ('x' if mode & _stat.S_IXOTH else '-')
    ]

    return f'{file_type}{"".join(perms)}'


def _format_size(size) -> str:
    units = ['B', 'K', 'M', 'G', 'T', 'P']
    unit_index = 0
    
    while size >= 1024 and unit_index < len(units) - 1:
        size /= 1024.0
        unit_index += 1

    if size == 0:
        return '0'

    return f'{size:.1f}{units[unit_index]}'


def print_help() -> None:
    print('''Usage: sls [flags] [file]

List information about file (the current directory by default)

positional arguments:
  file                  directory/file path (default: current directory)

options:
  -a, --all             show hidden files
  -h, --human-readable  print filesize in human-readable format
      --help            show this help message and exit
      --version         output version information and exit

This command is part of the SuperUtils collection (sls - super ls)
Copyright (C) 2024 vh8t
Author: vh8t
GitHub: https://github.com/vh8t
Website: https://vh8t.xyz
    ''')
    exit(0)


def main() -> None:
    parser = ArgumentParser(add_help=False)

    parser.add_argument('file', nargs='?', default='.')
    parser.add_argument('--help', action='store_true')
    parser.add_argument('--version', action='store_true')
    parser.add_argument('-a', '--all', action='store_true')
    parser.add_argument('-h', '--human-readable', action='store_true')

    args = parser.parse_args()
    if args.help:
        print_help()
    elif args.version:
        print('sls (SuperUtils) 1.0\nCopyright (C) 2024 vh8t')
        exit(0)
    else:
        list_dir(args.file, args.all, args.human_readable)


main()
