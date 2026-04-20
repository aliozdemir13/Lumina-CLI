// Package internal is managing the API logic, this class handles styling
package internal

import "fmt"

const (
	// ColorIndigo renders indigo color
	ColorIndigo = "\033[38;5;99m"
	// BgIndigo renders different shade of indigo color for background
	BgIndigo = "\033[48;5;99m"
	// ColorWhiteBold renders white color
	ColorWhiteBold = "\033[1;37m"
	// ColorReset resets added style
	ColorReset = "\033[0m"
	// ColorDim renders dimmed colors
	ColorDim = "\033[38;5;245m"
	// ColorGreen renders green color
	ColorGreen = "\033[92m"
	// ColorYellow renders yellow color
	ColorYellow = "\033[38;2;223;142;29m"
	// ColorRed renders red color
	ColorRed = "\033[38;2;210;15;57m"
	// ColorLavender renders lavender color
	ColorLavender = "\033[38;2;114;135;253m"
	// CBgIndigoLight slightly lighter indigo for the accent
	BgIndigoLight = "\033[48;5;105m"
	// TextIndigo bright "glowing" indigo
	TextIndigo = "\033[38;5;141m"
)

// Dim - Dimming the text visibility using darker color
func Dim(text string) string {
	return ColorDim + text + ColorReset
}

// Indigo - Coloring text indigo
func Indigo(text string) string {
	return ColorIndigo + text + ColorReset
}

// StyledBar styles the header bars
func StyledBar(text string) string {
	return fmt.Sprintf("%s %s %s%s", BgIndigo, ColorWhiteBold+text, ColorReset, "\n")
}

// FancyBar renders sub header with logos
func FancyBar(title string, version string) string {
	iconPart := BgIndigo + ColorWhiteBold + " ⚡ " + ColorReset
	titlePart := BgIndigoLight + ColorWhiteBold + " " + title + " " + ColorReset
	versionPart := "\033[38;5;239m" + " " + version + ColorReset // Dark gray version
	return iconPart + titlePart + versionPart + "\n"
}

// MegaLogo prinst the logo
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
