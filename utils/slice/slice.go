package sliceutil

import (
	"bytes"
	stringutil "dms/utils/string"
	"io"
)

func RemoveString(slice []string, remove func(item string) bool) []string {
	for i := 0; i < len(slice); i++ {
		if remove(slice[i]) {
			slice = append(slice[:i], slice[i+1:]...)
			i--
		}
	}
	return slice
}

func HasString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

func IOCopy(w io.Writer, slice []string) error {
	output := new(bytes.Buffer)
	for _, s := range slice {
		output.WriteString(stringutil.StripAnsi(s))
	}

	_, err := io.Copy(w, output)
	if err != nil {
		return err
	}

	return nil
}
