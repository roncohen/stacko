package stacko

import (
  "os"
  "io/ioutil"
  "path/filepath"
  "strings"
  "runtime"
)

type Stacktrace []Frame

type Frame struct {
  FileName     string
  FunctionName string
  PackageName  string
  Path         string
  LineNumber   int
  InDomain     bool
  PreContext   []string
  PostContext  []string
  Context      string
}

func NewStacktrace (skip int) (Stacktrace, error) {
  var stacktrace Stacktrace
  for i := skip; ; i++ {
    pc, path, lineNumber, ok := runtime.Caller(i)
    if !ok {
      break
    }

    packageName, functionName := FunctionInfo(pc)

    gopath := os.Getenv("GOPATH")
    fileName, err := filepath.Rel(gopath, path)
    if err != nil {
      return nil, err
    }

    context, offset, err := ContextInfo(path, lineNumber)
    if err != nil {
      return nil, err
    }

    InDomain := i == skip || stacktrace[0].PackageName == packageName

    stacktrace = append(stacktrace, Frame{
      fileName,
      functionName,
      packageName,
      path,
      lineNumber,
      InDomain,
      context[:offset - 1],
      context[offset:],
      context[offset - 1],
    })
  }

  return stacktrace, nil
}

func FunctionInfo (pc uintptr) (string, string) {
  function := runtime.FuncForPC(pc)
  if function == nil {
    return "", ""
  }

  name := function.Name()
  slash := strings.LastIndex(name, "/")
  if slash < 0 {
    slash = 0
  } else {
    slash++
  }

  info := strings.Split(name[slash:], ".")

  return info[0], info[1]
}

func ContextInfo (path string, lineNumber int) ([]string, int, error) {
  file, err := os.Open(path)
  if err != nil {
    return nil, -1, err
  }

  data, err := ioutil.ReadAll(file)
  if err != nil {
    return nil, -1, err
  }

  lines := strings.Split(string(data), "\n")

  start := lineNumber - 7
  if start < 0 {
    start = 0
  }

  end := lineNumber + 7
  if end >= len(lines) {
    end = len(lines) - 1
  }

  return lines[start:end], lineNumber - start, nil
}
