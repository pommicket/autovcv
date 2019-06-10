/*
Copyright (C) 2019 Leo Tenenbaum

This file is part of AutoVCV.

AutoVCV is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

AutoVCV is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with AutoVCV.  If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
    "autovcv/autovcv"
    "autovcv/vcv"
    "fmt"
    "os"
    "math/rand"
    "time"
    "flag"
    "runtime"
    "strings"
)

type versionNumber struct {
    plugin string
    version string
}

func main() {

    modules := flag.Int("modules", 50, "The number of modules to create")
    wires := flag.Int("wires", 100, "The maximum number of wires to create between modules")
    seed := flag.Int64("seed", time.Now().UTC().UnixNano(), "The seed to use for the random number generator")
    outputFile := flag.String("out", "out.vcv", "The name of the output VCV file")

    var versions []versionNumber = make([]versionNumber, 0, len(vcv.Versions))
    for key, val := range vcv.Versions {
        versions = append(versions, versionNumber{key, val})
        flag.StringVar(&versions[len(versions)-1].version,
            "version-" + strings.ToLower(key), val,
            fmt.Sprintf("The version of the %v plugin", key))

    }

    var defaultDevice string
    if runtime.GOOS == "windows" {
        defaultDevice = "Speakers (Realtek High Definition Audio)"
    } else {
        defaultDevice = "default"
    }
    deviceName := flag.String("device", defaultDevice, "The audio device to use")
    flag.Parse()

    // Put versions into vcv.Versions
    for _, vnum := range versions {
        vcv.Versions[vnum.plugin] = vnum.version
    }

    rand.Seed(*seed)
    fmt.Println("Using seed:", *seed)

    vcv := autovcv.RandomVCV(*modules, *wires)
    file, err := os.Create(*outputFile)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer file.Close()
    err = vcv.Write(file, *deviceName)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Printf("Successfully outputted to %v\n", *outputFile)
}
