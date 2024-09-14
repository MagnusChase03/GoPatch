package patch

import (
    "fmt"
    "os"
    "strconv"
    "strings"
    "regexp"
    "encoding/binary"
)

// Definitions of a binary patch
type PatchSection struct {
    Address uint64
    Data []byte
}

type Patch struct {
    Patches []PatchSection
}

/*
 * Returns a new patch given a patchfile.
 *
 * Arguments:
 *   - patchfile (string): The file path to the patchfile.
 *
 * Returns:
 *   - *Patch: The defined patches to be made.
 *   - error: An error if any occurred.
 */
func NewPatch(patchfile string) (*Patch, error) {
    f, err := os.Open(patchfile);
    if err != nil {
        return nil, fmt.Errorf("[ERROR] Failed to open patchfile %s. %w", patchfile, err);
    }

    // Read file into memory
    patchBuffer := make([]byte, 0);
    readBuffer := make([]byte, 1);
    for true {
        bytesRead, err := f.Read(readBuffer);
        if bytesRead == 0 {
            break;
        }

        if err != nil {
            return nil, fmt.Errorf("[ERROR] Failed to read patchfile %s. %w", patchfile, err);
        }
        patchBuffer = append(patchBuffer, readBuffer[0]);
    }

    // Parse file into patch sections
    addressPattern, err := regexp.Compile("^0x[0-9a-fA-F]+:$");
    dataPattern, err := regexp.Compile("^0x[0-9a-fA-F]+$");

    lines := strings.Split(string(patchBuffer[:]), "\n");
    patchSections := make([]PatchSection, 0);
    var currentAddress uint64;
    for _, line := range lines {
        if line == "" {
            continue;
        }

        // Match <Address>:
        if addressPattern.MatchString(line) {
            end := len(line) - 1;
            currentAddress, err = strconv.ParseUint(line[2:end], 16, 64);
            if err != nil {
                return nil, fmt.Errorf("[ERROR] Failed to parse address %s. %w", line, err);
            }

        // Match <Bytes>
        } else if dataPattern.MatchString(line) {
            currentBytes, err := strconv.ParseUint(line[2:], 16, 64);
            if err != nil {
                return nil, fmt.Errorf("[ERROR] Failed to parse bytes %s. %w", line, err);
            }
            data := make([]byte, 8);
            binary.BigEndian.PutUint64(data, currentBytes);
            start := 0;
            for data[start] == 0 {
                start += 1;
            }

            patchSections = append(patchSections, PatchSection{
                Address: currentAddress,
                Data: data[start:],
            });

        // Unrecognized pattern
        } else {
            return nil, fmt.Errorf("[ERROR] Failed to parse line %s.", line);
        }
    }

    f.Close();
    return &Patch{patchSections}, nil;
}

/*
 * Writes specified patches to a file.
 *
 * Arguments:
 *   - binaryFile (string): The file to be patched.
 *
 * Returns:
 *   - error: An error if any occurred.
 */
func (p *Patch) WritePatch(binaryFile string) error {
    f, err := os.OpenFile(binaryFile, os.O_RDWR, 0644);
    if err != nil {
        return fmt.Errorf("[ERROR] Failed to open binary %s. %w", binaryFile, err);
    }

    for _, patch := range p.Patches {
        _, err = f.Seek(int64(patch.Address), 0);
        if err != nil {
            return fmt.Errorf("[ERROR] Failed to seek to 0x%x in binary %s. %w", patch.Address, binaryFile, err);
        }

        _, err := f.Write(patch.Data);
        if err != nil {
            return fmt.Errorf("[ERROR] Failed to write to 0x%x in binary %s. %w", patch.Address, binaryFile, err);
        }
    }

    f.Close();
    return nil;
}
