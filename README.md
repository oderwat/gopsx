# gopsx

This is a helper which aims to make the usage of [gops](https://github.com/google/gops) easier by allowing the usage of globbing process name patterns instead of pids. The script also can run for multiple targets at once. See below for examples.

It was implemented in a way that it can either being build as executable or run like a shell script by making it executable and put it into your path. I created a `Makefile` which lets me either install the binary, the script, or a linked version of the script and calling them all with `gopsx` from the shell.

The installation can be done with

* Using `go get https://github.com/oderwat/gopsx` to install the binary.
* Downloading the raw code of the script, place it into the path, and make it executable (and maybe do more, see the `Makefile`).
* Get the source and use the Makefile like `make linked` or `make script`. Both use `~/bin/` as the installation path. So you may need to adapt it.

All of them need a working Go compiler, which you pretty much have installed when you are interested in `gops` anyway.

The calling syntax is the same as `gops` but you can replace the `pid` parameter with a globbing pattern. If there are multiple matches `goplx` will list them and you can either select one of them by using the ordinal number like `s*@1` or `s*@2`. This will call `gops` for this process. Or you can use another `*` like in `s*@*` to run `gops` once for each PID. As not every process has the agent running and most commands need the agent you can also use `+` to just select everything which runs the agent like `s*@+` which would only run `gops` on processes which name start with `s` and have the agent running.

Examples:

```bash
# listing all Go processes as with gops
gopsx

# listing all Go processes and present a list if you have more as one
gopsx \*

# listing all Go processes and run gops as if you called it for every running go process pid
gopsx \*@*

# running gops on the second instance of syncthing
gopsx syncthing@2

# listing all Go processes which have a gops agent included and run gops as if you called it for every running go process pid
gopsx \*@+

# running the gops commnand stats on all Go processes which have an agent
gopsx stats \*@+

# listing pids of all your minio servers
gopsx pid minio@*

# showing a stack trace for your app named "simpleserver"
gopsx trace simpleserver
```

The last example is the reason why I wrote this helper. It was annoying to get the PID for my process on each new test run. So I wrote a shell script to get the PID. I found that useful and wanted to share it with others. But a bash script is not a lot of fun, so I decided to write an enhanced version and also try out what needs to be done for using Go as a scripting tool for my needs.
