# Hexdump in Go

## Introduction:

This is a Golang version of the standard (monitor style) hexdump program.

The aim of the program was to learn how to do the following in Golang:

    1. Handle parameters using the default ("flags") package
    2. Handle IO from files and STDIN
    3. Use the core strings package library
    4. Explore the limitations of the in built arrays
    5. See how Golang pointers differ from C

## Usage:

    $ hexdump [-x | -w ] <files>

or:

    $ cat <file> | hexdump [ -x | -w ]

Options:

    -x      64 byte wide display ("extra wide")
    -w      32 byte wide display ("wide")

By default the display is 16 bytes wide

The format of the output is:

    <Address> : <Hex bytes> : <ASCII bytes>

The size of the "\<Address\>" field is dependant on the size of the file or if the file is streamed from STDIN. STDIN streams always use a 64bit wide hex address value. If the file is less than 64KiB long a 16bit address is used for the \<Address\> value. If the length of the file is less than MaxUint32 (2^32bytes) then a 32bit address is used for the \<Address\>. If neither of the above two states are true then the program will default to a 64bit address for \<Address\>. 

The \<Address\> value is always in upper case.

The \<Hex bytes\> vales are always in lower case.

The \<ASCII bytes\> will only display **printable ASCII** characters (range 0x20 to 0x7E). The output is in ASCII and **not** in UTF-8
