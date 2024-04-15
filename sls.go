package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "math"
  "os"
  "os/user"
  "path/filepath"
  "sort"
  "strconv"
  "strings"
  "syscall"

  flag "github.com/spf13/pflag"
)


type Metadata struct {
  Permission string
  Links      uint64
  Owner      string
  Group      string
  Size       string
  ModTime    string
  Path       string
}


type Icons struct {
  IconByFilename            map[string]string `json:"icons_by_filename"`
  IconByFileExtension       map[string]string `json:"icons_by_file_extension"`
  IconByOperatingSystem     map[string]string `json:"icons_by_operating_system"`
  IconByDesktopEnvironment  map[string]string `json:"icons_by_desktop_environment"`
  IconByWindowManager       map[string]string `json:"icons_by_window_manager"`
}

var (
  showHidden    *bool
  humanReadable *bool
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


func fileType(fileMode os.FileMode) string {
  switch {
  case fileMode.IsRegular():
    return "-"
  case fileMode.IsDir():
    return "d"
  case fileMode&os.ModeSymlink != 0:
    return "l"
  default:
    return "?"
  }
}


func fFileMode(fileMode os.FileMode) string {
  perms := []string{"---", "--x", "-w-", "-wx", "r--", "r-x", "rw-", "rwx"}
  permBits := fileMode.Perm()

  var sb strings.Builder
  sb.WriteString(fileType(fileMode))
  sb.WriteString(perms[(permBits>>6)&0x7])
  sb.WriteString(perms[(permBits>>3)&0x7])
  sb.WriteString(perms[permBits&0x7])

  return sb.String()
}


func colorPrint(message string, color string) {
  colorCode, ok := colors[color]
  if !ok {
    colorCode = colors["reset"]
  }

  fmt.Printf("%s%s%s\n", colorCode, message, colors["reset"])
}


func getFilePerms(filePath string) (string, error) {
  fileInfo, err := os.Stat(filePath)
  if err != nil {
    return "", err
  }

  fileMode := fileInfo.Mode()
  permissions := fFileMode(fileMode)

  return permissions, nil
}


func getOwnerAndGroup(info os.FileInfo) (string, string, error) {
	owner, err := user.LookupId(fmt.Sprintf("%d", info.Sys().(*syscall.Stat_t).Uid))
	if err != nil {
		return "", "", err
	}

	group, err := user.LookupGroupId(fmt.Sprintf("%d", info.Sys().(*syscall.Stat_t).Gid))
	if err != nil {
		return "", "", err
	}

	return owner.Username, group.Name, nil
}


func getFileMeta(path string) (string, uint64, string, string, int64, string, string, error) {
  perms, err := getFilePerms(path)
  if err != nil {
    fstr := fmt.Sprintf("Error: %v", err)
    colorPrint(fstr, "red")
    return "", 0, "", "", 0, "", "", err
  }
      
  fileInfo, err := os.Lstat(path)
  if err != nil {
    fstr := fmt.Sprintf("Error getting file info: %v", err)
    colorPrint(fstr, "red")
    return "", 0, "", "", 0, "", "", err
  }

  links := fileInfo.Sys().(*syscall.Stat_t).Nlink

  owner, group, err := getOwnerAndGroup(fileInfo)
  if err != nil {
    fstr := fmt.Sprintf("Error getting owner and group: %v", err)
    colorPrint(fstr, "red")
  }

  size := fileInfo.Size()
  modTime := fileInfo.ModTime().Format("Jan 02 15:04")

  fileName := filepath.Base(path)
  if fileInfo.Mode()&os.ModeSymlink != 0 {
    target, err := os.Readlink(path)
    if err != nil {
      fstr := fmt.Sprintf("Error resolving symlink: %v", err)
      colorPrint(fstr, "red")
    } else {
      fileName = fmt.Sprintf("%s -> %s", fileName, target)
    }
  }

  return perms, links, owner, group, size, modTime, fileName, nil
}


func convertSize(size int64) string {
  suffixes := []string{"B", "K", "M", "G", "T", "P", "E", "Z", "Y"}

  if size == 0 {
    return "0"
  }

  if size <= 1024 {
    return fmt.Sprintf("%d", size)
  }

  i := math.Floor(math.Log(float64(size)) / math.Log(1024))
  suffix := float64(size) / math.Pow(1024, i)

  return fmt.Sprintf("%.1f%s", suffix, suffixes[int(i)])
}


func sls(path string) error {
  if _, err := os.Stat(path); err == nil {
    info, err := os.Stat(path)
    if err != nil {
      fstr := fmt.Sprintf("Error: %v", err)
      colorPrint(fstr, "red")
      return err
    }

    exePath, err := os.Executable()
    if err != nil {
      fstr := fmt.Sprintf("Error getting executable path: %v", err)
      colorPrint(fstr, "red")
      return err
    }

    file, err := os.Open(filepath.Join(filepath.Dir(exePath), "..", "icons.json"))
    if err != nil {
      fstr := fmt.Sprintf("Error opening file: %v", err)
      colorPrint(fstr, "red")
      return err
    }

    defer file.Close()

    var icons Icons

    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&icons); err != nil {
      fstr := fmt.Sprintf("Error decoding JSON: %v", err)
      colorPrint(fstr, "red")
      return err
    }

    if info.IsDir() {
      dir, err := os.Open(path)
      if err != nil {
        fstr := fmt.Sprintf("Error opening directory: %v", err)
        colorPrint(fstr, "red")
        return err
      }

      entries, err := dir.Readdir(0)
      if err != nil {
        fstr := fmt.Sprintf("Error reading directory: %v", err)
        colorPrint(fstr, "red")
        return err
      }

      var fileList []Metadata
      for _, entry := range entries {
        name := entry.Name()
        if !*showHidden && (entry.IsDir() || strings.HasPrefix(name, ".")) {
          continue
        }

        perms, links, owner, group, size, modTime, fileName, err := getFileMeta(filepath.Join(path, name))
        if err != nil {
          fstr := fmt.Sprintf("Error: %v", err)
          colorPrint(fstr, "yellow")
          continue
        }

        var strSize string
        if *humanReadable {
          strSize = convertSize(size)
        } else {
          strSize = fmt.Sprintf("%d", size)
        }

        fileList = append(fileList, Metadata{
          Permission: perms,
          Links:      links,
          Owner:      owner,
          Group:      group,
          Size:       strSize,
          ModTime:    modTime,
          Path:       fileName,
        })
      }

      if *showHidden {
        dirs := []string{".", ".."}
        for _, dir_ := range dirs {
          perms, links, owner, group, size, modTime, fileName, err := getFileMeta(path + string(filepath.Separator) + dir_)
          if err != nil {
            fstr := fmt.Sprintf("Error: %v", err)
            colorPrint(fstr, "yellow")
          }
        
          var strSize string
          if *humanReadable {
            strSize = convertSize(size)
          } else {
            strSize = fmt.Sprintf("%d", size)
          }


          fileList = append(fileList, Metadata{
            Permission: perms,
            Links:      links,
            Owner:      owner,
            Group:      group,
            Size:       strSize,
            ModTime:    modTime,
            Path:       fileName,
          })
        }
      }

      if err != nil {
        fstr := fmt.Sprintf("Error: %v", err)
        colorPrint(fstr, "red")
        return err
      }

      sort.Slice(fileList, func(i, j int) bool {
        return fileList[i].Path < fileList[j].Path
      })

      max := func(a, b int) int {
        if a > b {
          return a
        }
        return b
      }

      var maxPerm, maxLinks, maxOwner, maxGroup, maxSize, maxMod, maxPath int
      for _, fileInfo := range fileList {
        maxPerm = max(maxPerm, len(fileInfo.Permission))
        maxLinks = max(maxLinks, len(strconv.FormatUint(fileInfo.Links, 10)))
        maxOwner = max(maxOwner, len(fileInfo.Owner))
        maxGroup = max(maxGroup, len(fileInfo.Group))
        maxSize = max(maxSize, len(fileInfo.Size))
        maxMod = max(maxMod, len(fileInfo.ModTime))
        maxPath = max(maxPath, len(fileInfo.ModTime))
      }

      for _, fileInfo := range fileList {
        clr := ""
        var match bool
        var icon string

        if fileInfo.Permission[3] == 'x' || fileInfo.Permission[6] == 'x' || fileInfo.Permission[9] == 'x' {
          clr = colors["green"]
          match = true
          icon = "\uead3 "
        }

        if fileInfo.Permission[0] == 'd' {
          clr = colors["blue"]
          match = true
          files, err := ioutil.ReadDir(filepath.Join(path, fileInfo.Path))
          if err != nil {
            fstr := fmt.Sprintf("Error reading directory: %v", err)
            colorPrint(fstr, "yellow")
          }

          if len(files) == 0 {
            icon = "\uea83 "
          } else {
            icon = "\ue6ad "
          }
        }

        fileExt := filepath.Ext(fileInfo.Path)

        if ico, ok := getIcon(icons.IconByFilename, fileInfo.Path); ok {
          match = true
          icon = ico + " "
        }

        if !match {
          for ext, ico := range icons.IconByFileExtension {
            if strings.HasSuffix(fileInfo.Path, "." + ext) {
              match = true
              icon = ico + " "
              break
            }
          }
        }

        if !match {
          if ico, ok := getIcon(icons.IconByOperatingSystem, fileExt); ok {
            match = true
            icon = ico + " "
          }
        }

        if !match {
          if ico, ok := getIcon(icons.IconByDesktopEnvironment, fileExt); ok {
            match = true
            icon = ico + " "
          }
        }

        if !match {
          if ico, ok := getIcon(icons.IconByWindowManager, fileExt); ok {
            match = true
            icon = ico + " "
          }
        }

        if !match {
          icon = "󰈚 "
        }

        fmt.Printf("%-*s  %*d  %-*s %-*s  %*s  %-*s  %s%s%s%s\n",
          maxPerm, fileInfo.Permission,
          maxLinks, fileInfo.Links,
          maxOwner, fileInfo.Owner,
          maxGroup, fileInfo.Group,
          maxSize, fileInfo.Size,
          maxMod, fileInfo.ModTime,
          clr, icon, fileInfo.Path, colors["reset"],
        )
      }
    } else {
      perms, links, owner, group, size, modTime, fileName, err := getFileMeta(path)
      if err != nil {
        fstr := fmt.Sprintf("Error: %v", err)
        colorPrint(fstr, "yellow")
      }
      fmt.Printf("%s  %d  %s %s  %d  %s  %s\n", perms, links, owner, group, size, modTime, fileName)
    }
  } else if os.IsNotExist(err) {
    fstr := fmt.Sprintf("%s does not exist", path)
    colorPrint(fstr, "red")
    return err
  } else {
    fstr := fmt.Sprintf("Error checking path: %v", err)
    colorPrint(fstr, "red")
    return err
  }
  return nil
}


func getIcon(m map[string]string, key string) (string, bool) {
  value, ok := m[key]
  return value, ok
}


func printHelp() {
  fmt.Println("Usage: sls [flags] [file]")
  fmt.Println("\nList information about file (the current directory by default)")
  fmt.Println("\npositional arguments:")
  fmt.Println("  file                  directory/file path (default: current directory)")
  fmt.Println("\noptions:")
  fmt.Println("  -a, --all             show hidden files")
  fmt.Println("  -h, --human-readable  print filesize in human-readable format")
  fmt.Println("      --help            show this help message and exit")
  fmt.Println("\nThis command is part of the SuperUtils collection (sls - super ls)")
  fmt.Println("Copyright (C) 2024 vh8t\nGitHub: https://github.com/vh8t\nWebsite: https://vh8t.xyz\n")
}


func main() {
  showHidden = flag.BoolP("all", "a", false, "show hidden files")
  humanReadable = flag.BoolP("human-readable", "h", false, "print filesize in human-readable format")
  showHelp = flag.BoolP("help", "", false, "show this help message and exit")
  flag.Parse()

  if *showHelp {
    printHelp()
    os.Exit(0)
  }

  path := "."
  if flag.NArg() > 0 {
    path = flag.Arg(0)
  }

  if err := sls(path); err != nil {
    os.Exit(1)
  }
}
