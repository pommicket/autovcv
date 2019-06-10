# AutoVCV

Randomly-generated [VCV Rack](https://vcvrack.com/) patches

You can run AutoVCV by downloading one of the
[releases](https://github.com/pommicket/autovcv/releases),
and running the executable file.

**IMPORTANT:** You might have one or more of the following problems
when you open the output of AutoVCV in VCV Rack:
1. The audio device is set to "(No device)" on all Audio modules.
   To fix this, you will need to first determine the correct audio device.
   You can do this by clicking on the "(No device)" text, and
   looking at the options VCV Rack gives you. If you would like to use,
   for example, "Speakers (Realtek High Definition Audio) (1-2 out)" as
   the device for Audio modules, you can use these command line arguments:
   ```bash
   autovcv -device 'Speakers (Realtek High Definition Audio)'
   ```

   (You do not need the \*-\* out).

2. (not very likely)
   VCV Rack might give an error about the version number in the `.vcv` file.
   This has never happened to me, but it's possible. If this happens,
   you should determine your "Core" and "Fundamental" plugin versions.
   When you search for a module which is part of a plugin, it'll show you
   the plugin version next to the name of the module (e.g. "Fundamental 0.6.2").
   To pass in the correct plugin versions, use something like:

   ```bash
   autovcv -version-core "0.6.2c" -version-fundamental "0.6.2"
   ```
3. This might not work on some very new versions of VCV Rack (i.e. ones which don't have VCOs, LFOs, etc.)

Also, VCV Rack might take a while to load files with a large number of modules/wires.

## Building from source

To build AutoVCV from source, you can use

```bash
go get github.com/pommicket/autovcv
cd $GOPATH/github.com/pommicket/autovcv
mkdir -p bin
go build -o bin/autovcv
```

You will need to install Go, which you can either do by going to
https://golang.org, or by using your package manager, e.g.
```bash
sudo apt install golang
```
on Debian/Ubuntu.

## Command-line arguments


You can see a full list of command-line arguments for AutoVCV with:

```bash
autovcv -help
```

These include things like setting the output file and changing the number
of modules or wires

## License

This program is licensed under the [GNU General Public License, Version 3](https://www.gnu.org/licenses/gpl-3.0.en.html).
