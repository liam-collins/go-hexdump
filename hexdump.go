package main

/*
	hexdump outputs a hex and ASCII dump, of a io stream, to STDOUT.

	The program will either read in from STDIN or take 1 or more
	REGULAR files and process them as an IO stream. If a given file
	is not a REGULAR file or the user does not have persmission to
	the file then it is skipped.

	Edits:

		2020-08-26		lc 		Created from scratch

	Copyright (c) 2020 NOVA Industries Limited

	No warrenty implied or otherwise

*/

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

const (
	bufferSize = 4096

	normalWidth    = 16
	wideWidth      = 32
	extraWideWidth = 64

	maxUint16 = ^uint16(0)
	maxUint32 = ^uint32(0)

	hex8Bits  = "%2.2X"
	hex16Bits = "%4.4X"
	hex32Bits = "%8.8X"
	hex64Bits = "%16.16X"

	chSpace = 0x20
	chDel   = 0x7F
)

func main() {

	var displayWidth int
	wide := flag.Bool("w", false, "32 byte wide display (cannot use with '-x')")
	extraWide := flag.Bool("x", false, "64 byte wide display (cannot use with '-w'")

	flag.Parse()
	args := flag.Args()

	if *wide && *extraWide {
		fmt.Fprintf(os.Stderr, "Error: Wide and Extra wide options are mutually exclusive\n")
		os.Exit(1)
	}

	switch {
	case *wide:
		displayWidth = wideWidth
	case *extraWide:
		displayWidth = extraWideWidth
	default:
		displayWidth = normalWidth
	}

	numberOfFiles := flag.NArg()

	if numberOfFiles == 0 {
		hexdump(os.Stdin, hex64Bits, displayWidth)
	} else {
		for i := range args {
			file := args[i]

			if fh, fileScale, err := openRegularFile(file); err != nil {
				fmt.Fprintf(os.Stderr, "\nWarning: Skipping file: %s\n", err)
			} else {
				defer fh.Close()
				hexdump(fh, fileScale, displayWidth)
			}
		}
	}
}

// hexdump dump the content of an IO stream in hex and ASCII format
//		The function reads io from an inout stream and writes the
//		hex & ASCII characters to STDOUT.
//		The output format is:
//			<File Offset>   <hex> ... <hex>  : <printable ASCII chars>

func hexdump(fh *os.File, fileScale string, displayWidth int) {

	buffer := make([]byte, bufferSize)
	var offset uint64

	for {
		if bufferRead, err := fh.Read(buffer); err == nil {
			offset = formatBuffer(buffer, bufferRead, fileScale, offset, displayWidth)
		} else {
			if err != io.EOF {
				fmt.Println("Error:", err)
			}
			return
		}
	}
}

// formatBuffer takes the content of a buffer and prints. The code produces
// 	an output formatted as follows:
//
//	<offset hex address>   <16 hex bytes>  <ASCii characters>

func formatBuffer(buffer []byte,
	bytesInBuffer int,
	fileScale string,
	position uint64,
	displayWidth int) uint64 {

	var hexDigits string
	var chrDigits string
	linePosition := position

	// Dynamically build the output format string for "Printf"
	outputFormat := fileScale + " : " + fmt.Sprintf("%%-%ds", 3*displayWidth) + "  : %s\n"
	for i := 0; i < bytesInBuffer; i++ {

		widthCounter := position % uint64(displayWidth)
		if widthCounter == 0 && i > 0 {
			fmt.Printf(outputFormat, linePosition, hexDigits, chrDigits)
			hexDigits = ""
			chrDigits = ""
			linePosition = position
		}

		hexDigits = fmt.Sprintf("%s %2.2x", hexDigits, buffer[i])

		if isPrintable(buffer[i]) {
			chrDigits = fmt.Sprintf("%s%c", chrDigits, buffer[i])
		} else {
			chrDigits = chrDigits + "."
		}

		position++
	}

	fmt.Printf(outputFormat, linePosition, hexDigits, chrDigits)
	return position
}

// isPrintable checks to see if a character is a printable
// character. This is based on the "C" code. It will probably
// be converted into a lambda function - wish this had "macros"

func isPrintable(ch byte) (printable bool) {

	if ch >= chSpace && ch < chDel {
		return true
	}

	return false
}

// openRegularFile will only allow a regular file to be opened for reading.
// 		The function returns a file handle to a requested file only
// 		if the following conditions are cleared:
//
//		1. The file is a regular file (links are allowed to regular files)
//		2. The user has permissions to read the file

func openRegularFile(filename string) (fh *os.File, fileSizeScale string, err error) {

	err = nil
	fh = nil

	fileInfo, err := os.Stat(filename)
	if err != nil {
		return
	}

	fileMode := fileInfo.Mode()
	if !fileMode.IsRegular() {
		errorMsg := fmt.Sprintf("open %s: It's not a regular file", filename)
		err = errors.New(errorMsg)
		return
	}

	switch {
	case fileInfo.Size() < int64(maxUint16):
		fileSizeScale = hex16Bits
	case fileInfo.Size() < int64(maxUint32):
		fileSizeScale = hex32Bits
	default:
		fileSizeScale = hex64Bits
	}

	fh, err = os.Open(filename)
	return
}
