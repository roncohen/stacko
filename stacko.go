// Copyright 2014 Christoffer Hallas.

// Package stacko provides the ability to generate a structured stacktrace.
package stacko

import (
  "os"
  "io/ioutil"
  "path/filepath"
  "strings"
  "runtime"
)

// The Stacktrace type is a slice of frames.
type Stacktrace []Frame

// The Frame type is a struct that hold a structured stacktrace frame.
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

// NewStacktrace generates a complete stacktrace except for those initial skip
// frames that are skipped.
func NewStacktrace (skip int) (Stacktrace, error) {

  // Create the actual stacktrace as a slice of frames.
  var stacktrace Stacktrace

  // Loop from skip and forever. We therefor rely on the execution of the loop
  // to provide termination.
  for i := skip; ; i++ {

    // Get the program counter, path and line number for a the frame i.
    pc, path, lineNumber, ok := runtime.Caller(i)

    // If not ok, we break and subsequently return the generated stacktrace.
    if !ok {
      break
    }

    // Call our own API to get the package and function names.
    packageName, functionName := FunctionInfo(pc)

    // We extract the context of a frame, e.g. the line it self, preceding and
    // proceding lines.
    gopath := os.Getenv("GOPATH")
    fileName, err := filepath.Rel(gopath, path)
    if err != nil {
      return nil, err
    }

    // Get the actual context, a slice of strings.
    context, offset, err := ContextInfo(path, lineNumber)
    if err != nil {
      return nil, err
    }

    // If this is the first frame or the frame has the same package as the first
    // frame then mark it as in domain.
    InDomain := i == skip || stacktrace[0].PackageName == packageName

    // Append the frame to the stacktrace.
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

// FunctionInfo takes a program counter and returns the function and package
// name for the frame at that counter.
func FunctionInfo (pc uintptr) (string, string) {

  // Get the function.
  function := runtime.FuncForPC(pc)
  if function == nil {
    return "", ""
  }

  // We take the name which at this point is a complete address of the function,
  // including path and file, then we take out the last part which is the
  // function and package name seperated by a dot.
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

// ContextInfo takes a path and a line number and returns a slice of strings
// that represent the line in self, the preceding and the proceding lines. It
// also returns the offset for the actual line in the context.
func ContextInfo (path string, lineNumber int) ([]string, int, error) {
  file, err := os.Open(path)
  if err != nil {
    return nil, -1, err
  }

  data, err := ioutil.ReadAll(file)
  if err != nil {
    return nil, -1, err
  }

  // We split on linebreak to get the lines.
  lines := strings.Split(string(data), "\n")

  // We try to get at most the 7 preceding lines.
  start := lineNumber - 7
  if start < 0 {
    start = 0
  }

  // Similar to the preceding lines, we also try for at most the 7 proceding
  // lines.
  end := lineNumber + 7
  if end >= len(lines) {
    end = len(lines) - 1
  }

  return lines[start:end], lineNumber - start, nil
}
