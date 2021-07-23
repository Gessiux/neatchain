package neatbyft

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"gopkg.in/urfave/cli.v1"

	cfg "github.com/Gessiux/go-config"
	tmcfg "github.com/Gessiux/neatchain/consensus/neatbyft/config/neatcon"
)

func GetTendermintConfig(chainId string, ctx *cli.Context) cfg.Config {
	datadir := ctx.GlobalString(DataDirFlag.Name)
	config := tmcfg.GetConfig(datadir, chainId)

	return config
}

func HomeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}

func DefaultDataDir() string {
	// Try to place the data folder in the user's home dir
	home := HomeDir()
	if home != "" {
		if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "neatchain")
		} else {
			return filepath.Join(home, ".neatchain")
		}
	}
	// As we cannot guess a stable location, return empty and handle later
	return ""
}

func ConcatCopyPreAllocate(slices [][]byte) []byte {
	var totalLen int
	for _, s := range slices {
		totalLen += len(s)
	}
	tmp := make([]byte, totalLen)
	var i int
	for _, s := range slices {
		i += copy(tmp[i:], s)
	}
	return tmp
}
