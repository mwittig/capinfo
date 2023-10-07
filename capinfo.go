package main

import (
	"archive/zip"
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/skythen/cap"
)

var aidMap = map[string]string{
	"A0000000620201": "javacardx.crypto",
	"A0000000620202": "javacardx.biometry",
	"A0000000620203": "javacardx.external",
	"A0000000620204": "javacardx.biometry1toN",
	"A0000000620205": "javacardx.security",

	"A000000062020801":   "javacardx.framework.util",
	"A00000006202080101": "javacardx.framework.util.intx",
	"A000000062020802":   "javacardx.framework.math",
	"A000000062020803":   "javacardx.framework.tlv",
	"A000000062020804":   "javacardx.framework.string",

	"A0000000620209":   "javacardx.apdu",
	"A000000062020901": "javacardx.apdu.util",

	"A0000000620101":   "javacard.framework",
	"A000000062010101": "javacard.framework.service",
	"A0000000620102":   "javacard.security",

	"A0000000620001": "java.lang",
	"A0000000620002": "java.io",
	"A0000000620003": "java.rmi",

	"A00000015100":     "org.globalplatform",
	"A0000000030000":   "visa.openplatform",
	"E804007F00070308": "NXP JCOP 5.2 CSP",
}

func stringify(bytes []byte) string {
	return strings.ToUpper(hex.EncodeToString(bytes))
}

func aid2PackageName(aidHex string) string {
	if v, ok := aidMap[aidHex]; ok {
		return v
	}

	return aidHex
}

func main() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	if len(os.Args) < 2 {
		log.Fatalln("Missing parameter, provide file name.")
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln("unable to read file:", os.Args[1])
	}

	bReader := bytes.NewReader(data)

	zReader, err := zip.NewReader(bReader, bReader.Size())
	if err != nil {
		log.Fatalln("unable to open zip reader from file provided")
	}

	capFile, err := cap.Parse(zReader)
	if err != nil || capFile == nil {
		log.Fatalln("failed to parse CAP file: ", err)
	}

	w := new(tabwriter.Writer)
	defer w.Flush()

	w.Init(os.Stdout, 0, 8, 1, ' ', 0)

	cardAPIVersion := "-"
	if capFile.GPCardAPIVersion != nil {
		cardAPIVersion = fmt.Sprintf("%d.%d", capFile.GPCardAPIVersion.Major, capFile.GPCardAPIVersion.Minor)
	}

	packageName := "-"
	if capFile.Header.PackageName != "" {
		packageName = capFile.Header.PackageName
	}

	fmt.Fprintf(w, "CAP Filname:\t%s\n", filepath.Base(os.Args[1]))
	fmt.Fprintf(w, "Package Name:\t%s\n", packageName)
	fmt.Fprintf(w, "Package AID:\t%s\n", stringify(capFile.Header.Package.AID))
	fmt.Fprintf(w, "Package Version:\t%d.%d\n", capFile.Header.Package.Version.Major, capFile.Header.Package.Version.Major)
	fmt.Fprintf(w, "JavaCard Version:\t%d.%d.%d\n", capFile.JavaCardVersion.Major, capFile.JavaCardVersion.Minor, capFile.JavaCardVersion.Patch)
	fmt.Fprintf(w, "GP Card API Version:\t%s\n", cardAPIVersion)
	fmt.Fprintf(w, "Contains Applet:\t%t\n", capFile.Header.ContainsApplet)
	fmt.Fprintf(w, "Contains Export:\t%t\n", capFile.Header.ContainsExport)
	fmt.Fprintf(w, "Imported Packages:\t%d\n", len(capFile.Imports))

	for i := 0; i < len(capFile.Imports); i++ {
		fmt.Fprintf(w, "- %s:\t%d.%d\n", aid2PackageName(stringify(capFile.Imports[i].AID)), capFile.Imports[i].Version.Major, capFile.Imports[i].Version.Minor)
	}
	fmt.Fprintf(w, "Load Bytes:\t%d\n", len(capFile.LoadBytes(false)))
	fmt.Fprintf(w, "CAP File Version:\t%d.%d\n", capFile.Header.CapFileVersion.Major, capFile.Header.CapFileVersion.Major)
}
