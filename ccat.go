package main

import (
  "bufio"
  "fmt"
  "os"
  "strings"
  "github.com/alecthomas/chroma/v2"
  "github.com/alecthomas/chroma/v2/lexers"
  "github.com/alecthomas/chroma/v2/styles"
  "github.com/alecthomas/chroma/v2/formatters"
  flag "github.com/spf13/pflag"
)


var (
  showEnds      *bool
  numberLines   *bool
  squeezeBlanks *bool
  showTabs      *bool
  showHelp      *bool
)


var colors = map[string] string {
  "red":     "\033[31m",
  "green":   "\033[32m",
  "yellow":  "\033[33m",
  "blue":    "\033[34m",
  "magenta": "\033[35m",
  "cyan":    "\033[36m",
  "white":   "\033[37m",
  "reset":   "\033[0m",
}


func colorPrint(message string, color string) {
  colorCode, ok := colors[color]
  if !ok {
    colorCode = colors["reset"]
  }

  fmt.Printf("%s%s%s\n", colorCode, message, colors["reset"])
}


func ccat(path string) error {
  fileInfo, err := os.Stat(path)
  if err != nil {
    fstr := fmt.Sprintf("Error accesing file: %v", err)
    colorPrint(fstr, "red")
    return err
  }

  if fileInfo.IsDir() {
    fstr := fmt.Sprintf("%s must be a file, not a directory", path)
    colorPrint(fstr, "red")
    return err
  }

  file, err := os.Open(path)
  if err != nil {
    fstr := fmt.Sprintf("Error opening file: %v", err)
    colorPrint(fstr, "red")
    return err
  }
  defer file.Close()

  var code string
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    code += scanner.Text() + "\n"
  }
  if err := scanner.Err(); err != nil {
    fstr := fmt.Sprintf("Error reading file: %v", err)
    colorPrint(fstr, "red")
    return nil
  }

  lexer := lexers.Match(path)
  if lexer == nil {
    lexer = lexers.Fallback
  }
  lexer = chroma.Coalesce(lexer)

  style := styles.Get("base16-snazzy")
  if style == nil {
    style = styles.Fallback
  }

  formatter := formatters.Get("terminal")
  if formatter == nil {
    formatter = formatters.Fallback
  }

  iterator, err := lexer.Tokenise(nil, code)
  if err != nil {
    fstr := fmt.Sprintf("Tokenise error: %v", err)
    colorPrint(fstr, "red")
    return nil
  }

  var fCode strings.Builder

  err = formatter.Format(&fCode, style, iterator)
  if err != nil {
    fstr := fmt.Sprintf("Format error: %v", err)
    colorPrint(fstr, "red")
    return nil
  }

  lines := strings.Split(fCode.String(), "\n")
  pureLines := strings.Split(code, "\n")
  padding := len(fmt.Sprintf("%d", len(lines)))
  for index, line := range lines {
    if *showEnds {
      line = fmt.Sprintf("%s$", line)
    }

    if *numberLines {
      line = fmt.Sprintf("%-*d %s", padding, index + 1, line)
    }

    if *squeezeBlanks {
      if !(index + 2 > len(pureLines)) {
        if strings.TrimSpace(pureLines[index]) == "" && strings.TrimSpace(pureLines[index + 1]) == "" {
          continue
        }
      }
    }

    if *showTabs {
      line = strings.ReplaceAll(line, "\t", "^I")
    }
    
    if !(index + 1 == len(pureLines)) {
      fmt.Println(line)
    }
  }

  return nil
}


func printHelp() {
  fmt.Println("Usage: ccat [flags] [file]")
  fmt.Println("\nShow the contents of a file with syntax highlighting")
  fmt.Println("\npositional arguments:")
  fmt.Println("  file                 file path")
  fmt.Println("\noptions:")
  fmt.Println("  -e, --show-ends      display $ at the end of each line")
  fmt.Println("  -n, --number         number all output lines")
  fmt.Println("  -s, --squeeze-blank  supress repeated empty output lines")
  fmt.Println("  -t, --show-tabs      display TAB characters as ^I")
  fmt.Println("      --help            show this help message and exit")
  fmt.Println("\nThis command is part of the SuperUtils collection (ccat - colorful cat)")
  fmt.Println("Copyright (C) 2024 vh8t\nGitHub: https://github.com/vh8t\nWebsite: https://vh8t.xyz\n")
}


func main() {
  showEnds = flag.BoolP("show-ends", "e", false, "display $ at the end of each line")
  numberLines = flag.BoolP("number", "n", false, "number all output lines")
  squeezeBlanks = flag.BoolP("squeeze-blanks", "s", false, "supress repeated empty output lines")
  showTabs = flag.BoolP("show-tabs", "t", false, "display TAB characters as ^I")
  showHelp = flag.BoolP("help", "", false, "show this help message and exit")
  flag.Parse()

  if *showHelp {
    printHelp()
    os.Exit(0)
  }

  if flag.NArg() == 0 {
    colorPrint("Usage: ccat [flags] [file]\nMissing file argument", "red")
    os.Exit(1)
  }

  if err := ccat(flag.Arg(0)); err != nil {
    os.Exit(1)
  }
}
