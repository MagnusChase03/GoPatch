package main;

import (
    "os"
    "fmt"

    "github.com/MagnusChase03/GoPatch/patch"
)

func main() {
    if len(os.Args) < 3 {
        fmt.Printf("[USAGE] gopatch <BINARY> <PATCHFILE>\n");
        return;
    }

    p, err := patch.NewPatch(os.Args[2]);
    if err != nil {
        fmt.Printf("[ERROR] Failed to load patchfile. %v\n", err);
        return;
    }
    p.WritePatch(os.Args[1]);
}
