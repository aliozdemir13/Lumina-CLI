package internal

import "time"

func FormatToLocal(isoString string) (string, error) {
	const layout = "2006-01-02T15:04Z"

	// Parse using the custom layout
	t, err := time.Parse(layout, isoString)
	if err != nil {
		return "", err
	}

	// Format: "short month, 2-digit day, hour:min" (German style)
	// Layout: "Jan 02, 15:04"
	return t.Format("02. Jan, 15:04"), nil
}
