# Daniel Schmerber MD5-Hash comparison 

This is a small CLI tool, to create MD5 hashes of either Strings or Files.

Usage: 
    -s <string> -> returns the MD5 hash of the passed string
    -f <path> -> returns the MD5 hash of a given file

# Gleam

## Requirements:
Gleam (https://gleam.run/getting-started/installing/)
ErlangVM

## Run Project:
Download Dependencies (if missing): gleam deps download
Run : gleam run

## Test Project:
Test: gleam test

## Build Project
gleam build
gleam run -m gleescript
Run Standalone escript md5

## Benchmark Project

Run the cli with the --benchmark flag

# Go

## Requirements:
Go (https://go.dev/dl/)

## Run Project:
Download Dependencies (if missing): go mod download
Run: go run .

## Test Project:
go test ./...

## Build Project:
go build
Run Binary: ./<binary-name>

## Benchmark Project:
go test -bench .