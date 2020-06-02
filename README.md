# funcv [![GoDoc](https://godoc.org/moshenahmias/funcv?status.svg)](https://godoc.org/github.com/moshenahmias/funcv)

**funcv** helps you create CLI tools with Go.

**funcv** offers a different approach for dealing with command line arguments and flags.

**funcv** supplies an easy to use command builder, you use that builder to build your set of commands, each such command can be tested against a slice of string arguments, if the arguments are compatible with the command, a given action function is called, the parameters for that function are the extracted and parsed variables and flags.

Let's see how it works with a simple example:

```go
func main() {
	cmd := funcv.NewCommand("delete a file").
		AddConstant("delete", false).
		AddVariable("filename", "file to delete", new(funcv.StringConverter)).
		MustCompile()

	if _, err := cmd.Execute(os.Args[1:], func(name string) {
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

|          | Comment                                                |
| -------- | ------------------------------------------------------ |
| Constant | Static word (allowed characters: 0-9, A-Z, a-z, _, -). |
| Variable | With default value or without.                         |
| Flag     | -x or -x..x, with parameter or without.                |
| Variadic | Closing the list of arguments.                         |

The list of supported arguments is extendable via the `funcv.Argument` interface.



### Converters

Arguments that translates to function parameters (ex: not constant) require a `func.Converter`.

The package includes converters for Strings, Integers and Booleans.

### Groups

It is possible to group different commands together using a `funcv.Group`:

```go
func main() {
	var grp funcv.Group

	err := funcv.NewCommand("delete a file").
		AddConstant("example", false).
		AddConstant("delete", false).
		AddParameterlessFlag("r", "move to recycle bin", new(funcv.BooleanConverter), false, true).
		AddVariable("filename", "file to delete", new(funcv.StringConverter)).
		ToGroup(&grp, func(recycle bool, name string) {
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
		ToGroup(&grp, func() {
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
	if grp.ExecuteAll(append([]string{"example"}, os.Args[1:]...)) == 0 {
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



### Variadic Functions

Variadic action functions support is available. Example:  

```go
func main() {
	var grp funcv.Group

   	converter := new(funcv.IntegerConverter)
    
	if err := funcv.NewCommand("add two numbers").
		AddConstant("calc", false).
		AddConstant("add", false).
		AddVariable("1st", "first operand", converter).
		AddVariable("2nd", "second operand", converter).
		ToGroup(&grp, func(x, y int) {
			fmt.Println(x, "+", y, "=", x+y, "(I)")
		}); err != nil {
		panic(err)
	}

	if err := funcv.NewCommand("add two or more numbers").
		AddConstant("calc", false).
		AddConstant("add", false).
		AddVariable("1st", "first operand", converter).
		AddVariable("2nd", "second operand", converter).
		AddVariadic("operands", "list of operands", converter).
		ToGroup(&grp, func(operands ...int) {
			var sb strings.Builder
			var sum int

			for i, op := range operands {
				sum += op
				sb.WriteString(fmt.Sprint(op))

				if i+1 < len(operands) {
					sb.WriteString(" + ")
				}
			}

			sb.WriteString(fmt.Sprintf(" = %d (II)", sum))
			fmt.Println(sb.String())

		}); err != nil {
		panic(err)
	}

	if grp.ExecuteAll(append([]string{"calc"}, os.Args[1:]...)) == 0 {
		fmt.Fprintln(os.Stderr, "invalid command:", strings.Join(os.Args[1:], " "))
	}
}
```

```bash
$ calc add 1
invalid command: add 1

$ calc add 1 2
1 + 2 = 3 (I)
1 + 2 = 3 (II)

$ calc add 1 2 3
1 + 2 + 3 = 6 (II)
```

Use `ExecuteFirst` if you want to stop executing commands after the first successful executed command, the method returns the index of the executed command within the group, [0, `len(group)`), or a negative value if no command was found: 

```go
func main() {
	// ... same commands as before ...

	if grp.ExecuteFirst(append([]string{"calc"}, os.Args[1:]...)) < 0 {
		fmt.Fprintln(os.Stderr, "invalid command:", strings.Join(os.Args[1:], " "))
	}
}
```

```bash
$ calc add 1 2
1 + 2 = 3 (I)
```



### Download

```bash
$ go get -u github.com/moshenahmias/funcv
```
