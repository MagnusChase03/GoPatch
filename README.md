# GoPatch

*Tool to easily patch binary files.*

## Build

To build the tool:

```
$ go build main.go
$ strip main
$ mv main /usr/bin/gopatch
```

## Usage

To patch a binary:

```
$ gopatch <binary> <patchfile>
```

## Patchfile

A patchfile should follow the below syntax:

```
1 <address>:
2 0x<bytes>
```

Example:

```
1 0x12345678:
2 0x5058
3 0xe9
4
5 0x23456789:
6 0x45
```

This patches will happen at specified binary addresses, and the bytes will be
written in order.

## Example

Here is an example of gopatch being used:

`patch.txt`
```
1 0x0:
2 0x67
3 0x6869
```

```
$ xxd binary
00000000: 6162 630a 6465 660a                      abc.def.

$ gopatch binary patch.txt

$ xxd binary
00000000: 6768 690a 6465 660a                      ghi.def.
```
