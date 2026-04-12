package internal

import "fmt"

const (
	ColorIndigo    = "\033[38;5;99m"
	BgIndigo       = "\033[48;5;99m"
	ColorWhiteBold = "\033[1;37m"
	ColorReset     = "\033[0m"
	ColorDim       = "\033[38;5;245m"
	ColorGreen     = "\033[92m"
	ColorYellow    = "\033[38;2;223;142;29m"
	ColorRed       = "\033[38;2;210;15;57m"
	ColorLavender  = "\033[38;2;114;135;253m"
	BgIndigoLight  = "\033[48;5;105m" // A slightly lighter indigo for the accent
	TextIndigo     = "\033[38;5;141m" // A bright "glowing" indigo
)

// Function: Dimming the text visibility using darker color
func Dim(text string) string {
	return ColorDim + text + ColorReset
}

// Function: Coloring text indigo
func Indigo(text string) string {
	return ColorIndigo + text + ColorReset
}

func StyledBar(text string) string {
	return fmt.Sprintf("%s %s %s%s", BgIndigo, ColorWhiteBold+text, ColorReset, "\n")
}

func FancyBar(title string, version string) string {
	iconPart := BgIndigo + ColorWhiteBold + " ⚡ " + ColorReset
	titlePart := BgIndigoLight + ColorWhiteBold + " " + title + " " + ColorReset
	versionPart := "\033[38;5;239m" + " " + version + ColorReset // Dark gray version
	return iconPart + titlePart + versionPart + "\n"
}

func MegaLogo() string {
	return ColorIndigo + `
	██╗     ██╗   ██╗███╗   ███╗██╗███╗   ██╗ █████╗ 
	██║     ██║   ██║████╗ ████║██║████╗  ██║██╔══██╗
	██║     ██║   ██║██╔████╔██║██║██╔██╗ ██║███████║
	██║     ██║   ██║██║╚██╔╝██║██║██║╚██╗██║██╔══██║
	███████╗╚██████╔╝██║ ╚═╝ ██║██║██║ ╚████║██║  ██║
	╚══════╝ ╚═════╝ ╚═╝     ╚═╝╚═╝╚═╝  ╚═══╝╚═╝  ╚═╝
                                                 ` + ColorReset + "\n"
}
