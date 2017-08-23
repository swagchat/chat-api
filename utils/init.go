package utils

import (
	"fmt"
	"log"

	yaml "gopkg.in/yaml.v2"
)

func init() {
	log.SetFlags(log.Llongfile)
	setupConfig()
	if IsShowVersion {
		return
	}
	setupLogger()

	yaml, err := yaml.Marshal(&Cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	aa := `
███████╗██╗    ██╗ █████╗  ██████╗  ██████╗██╗  ██╗ █████╗ ████████╗     ██████╗██╗  ██╗ █████╗ ████████╗     █████╗ ██████╗ ██╗
██╔════╝██║    ██║██╔══██╗██╔════╝ ██╔════╝██║  ██║██╔══██╗╚══██╔══╝    ██╔════╝██║  ██║██╔══██╗╚══██╔══╝    ██╔══██╗██╔══██╗██║
███████╗██║ █╗ ██║███████║██║  ███╗██║     ███████║███████║   ██║       ██║     ███████║███████║   ██║       ███████║██████╔╝██║
╚════██║██║███╗██║██╔══██║██║   ██║██║     ██╔══██║██╔══██║   ██║       ██║     ██╔══██║██╔══██║   ██║       ██╔══██║██╔═══╝ ██║
███████║╚███╔███╔╝██║  ██║╚██████╔╝╚██████╗██║  ██║██║  ██║   ██║       ╚██████╗██║  ██║██║  ██║   ██║       ██║  ██║██║     ██║
╚══════╝ ╚══╝╚══╝ ╚═╝  ╚═╝ ╚═════╝  ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝        ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝       ╚═╝  ╚═╝╚═╝     ╚═╝
`
	fmt.Println(aa)
	pointDown := string([]byte{0xF0, 0x9F, 0x91, 0x87})
	fmt.Println(pointDown, " Chat API is running with the following setting values")
	fmt.Printf("\n%s\n", string(yaml))
}
