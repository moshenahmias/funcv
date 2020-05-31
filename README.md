# funcv [![GoDoc](https://godoc.org/moshenahmias/funcv?status.svg)](https://godoc.org/github.com/moshenahmias/funcv)

**funcv** helps you create CLI tools with Go.

It offers a different approach for dealing with command line arguments and flags.

**funcv** supplies an easy to use command builder, you use that builder to build your set of commands, each such command can be tested against a slice of string arguments, if the arguments are compatible with the command, a given action function is called, the parameters for that function are the extracted and parsed variables and flags input values.

Let's see how it works with a simple example:

```go
func main() {
	cmd := funcv.NewCommand("delete a file").
		AddConstant("delete", false).
		AddStrVar("filename", "file to delete").
		MustCompile()

	if err := cmd.Execute(os.Args[1:], func(name string) {
		fmt.Println("deleting", name, "...")
		// ...
	}); err != nil {
        fmt.Fprintln(os.Stderr, "invalid command:", strings.Join(os.Args[1:], " "))
	}
}
```

```bash
$ example delete song.mp3 
deleting song.mp3 ...
```

First, we called the `funcv.NewCommand` function with a command description, then, we used the returned builder (`funcv.Builder`) to add our command components, a constant text ("delete") and a string variable ("filename"). The call for `MustCompile()` finished the building process and returned the new command (`failure.Command`).

Next, we called the command's `Execute`  method with a slice of string arguments (`os.Args[1:]`) and an action function (`func(name string){...}`).

The `Execute` method tests the given arguments slice (`[]string{"delete", "song.mp3"}`)  and finds that it contains two arguments, the first argument equals "delete", therefore, the action function is called with the second argument as a parameter (`name string`).



### Arguments

Currently supported list of arguments:

|                               | Type   | Comment                                                      |
| ----------------------------- | ------ | ------------------------------------------------------------ |
| Constant                      | -      | Static word (allowed characters: 0-9, A-Z, a-z, _, -)        |
| String variable               | string |                                                              |
| Integer variable              | int64  | Every integer type is supported as the action function parameter (converted from int64 with a possible data loss) |
| String variable with default  | string | No other arguments allowed after that argument except other variables with default value |
| Integer variable with default | int64  | No other arguments allowed after that argument except other variables with default value |
| String flag                   | string | -x <value> or --x..x <value>                                 |
| Integer flag                  | int64  | -x <value> or --x..x <value>, every integer type is supported as the action function parameter (converted from int64 with a possible data loss) |
| Boolean flag                  | bool   | -x / -x <false/true> or --x..x / --x..x <false/true>         |

The list of supported arguments is extendable via the `funcv.Arg` interface.



### Groups

It is possible to group different commands together using a `funcv.Group` and test a slice of arguments against all grouped commands via a single call:

```go
func main() {
	grp := funcv.NewGroup()

    err := funcv.NewCommand("delete a file").
		AddConstant("example", false).
		AddConstant("delete", false).
		AddBoolFlag("r", "move to recycle bin").
		AddStrVar("filename", "file to delete").
		ToGroup(grp, func(recycle bool, name string) {
            // the count, order and type of params must match the count, order
            // and type of flags and variables in the command (excluding constants)	
            if recycle {
				fmt.Println("recycling", name, "...")
			} else {
				fmt.Println("deleting", name, "...")
			}
			// ...
		})

	if err != nil {
		panic(err)
	}

	err = funcv.NewCommand("print this help").
		AddConstant("example", false).
		AddConstant("help", false).
		ToGroup(grp, func() {
            // groups, commands and arguments implement io.WriterTo
            // and will write their informative usage text into the
            // writer given to WriteTo(io.Writer)
			w := new(tabwriter.Writer)
			w.Init(os.Stdout, 18, 8, 0, '\t', 0)
			defer w.Flush()
			if _, err := grp.WriteTo(w); err != nil {
				panic(err)
			}
		})

	if err != nil {
		panic(err)
	}

    // test against all commands in grp
    // returns the number of executed commands
	if grp.Execute(append([]string{"example"}, os.Args[1:]...)) == 0 {
		fmt.Fprintln(os.Stderr, "invalid command:", strings.Join(os.Args[1:], " "))
	}
}
```

```bash
$ example delete song.mp3 
deleting song.mp3 ...

$ example delete -r song.mp3 
recycling song.mp3 ...

$ example delete -r false song.mp3 
deleting song.mp3 ...

$ example delete -r true song.mp3 
recycling song.mp3 ...

$ example help
delete a file:          > example delete [-r] <filename>

                        -r                      move to recycle bin (default: false)
                        filename                file to delete

print this help:        > example help

$ example typo
invalid command: typo
```



### Download

```bash
$ go get -u github.com/moshenahmias/funcv
```
