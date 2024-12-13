package util

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

const (
	BYTE = 1 << (10 * iota)
	KILOBYTE
	MEGABYTE
	GIGABYTE
	TERABYTE
	PETABYTE
	EXABYTE
)

var invalidByteQuantityError = errors.New("byte quantity must be a positive integer with a unit of measurement like M, MB, MiB, G, GiB, or GB")

// ByteSize returns a human-readable byte string of the form 10M, 12.5K, and so forth.  The following units are available:
//
//	E: Exabyte
//	P: Petabyte
//	T: Terabyte
//	G: Gigabyte
//	M: Megabyte
//	K: Kilobyte
//	B: Byte
//
// The unit that results in the smallest number greater than or equal to 1 is always chosen.
func ByteSize(bytes uint64) string {
	unit := ""
	value := float64(bytes)

	switch {
	case bytes >= EXABYTE:
		unit = "Ei"
		value = value / EXABYTE
	case bytes >= PETABYTE:
		unit = "Pi"
		value = value / PETABYTE
	case bytes >= TERABYTE:
		unit = "Ti"
		value = value / TERABYTE
	case bytes >= GIGABYTE:
		unit = "Gi"
		value = value / GIGABYTE
	case bytes >= MEGABYTE:
		unit = "Mi"
		value = value / MEGABYTE
	case bytes >= KILOBYTE:
		unit = "Ki"
		value = value / KILOBYTE
	case bytes >= BYTE:
		unit = "B"
	case bytes == 0:
		return "0B"
	}

	result := strconv.FormatFloat(value, 'f', 1, 64)
	result = strings.TrimSuffix(result, ".0")
	return result + unit
}

func ByteSizeInMi(bytes uint64) string {
	value := float64(bytes)
	value = value / MEGABYTE
	if bytes < MEGABYTE {
		result := strconv.FormatFloat(value, 'f', 4, 64)

		result = strings.TrimRight(result, "0")

		return result + "Mi"
	}

	result := strconv.FormatFloat(value, 'f', 1, 64)
	result = strings.TrimSuffix(result, ".0")
	return result + "Mi"
}

// ToMegabytes parses a string formatted by ByteSize as megabytes.
func ToMegabytes(s string) (uint64, error) {
	bytes, err := ToBytes(s)
	if err != nil {
		return 0, err
	}

	return bytes / MEGABYTE, nil
}

// ToBytes parses a string formatted by ByteSize as bytes. Note binary-prefixed and SI prefixed units both mean a base-2 units
// KB = K = KiB	= 1024
// MB = M = MiB = 1024 * K
// GB = G = GiB = 1024 * M
// TB = T = TiB = 1024 * G
// PB = P = PiB = 1024 * T
// EB = E = EiB = 1024 * P
func ToBytes(s string) (uint64, error) {
	s = strings.TrimSpace(s)
	s = strings.ToUpper(s)

	i := strings.IndexFunc(s, unicode.IsLetter)

	if i == -1 {
		return 0, invalidByteQuantityError
	}

	bytesString, multiple := s[:i], s[i:]
	bytes, err := strconv.ParseFloat(bytesString, 64)
	if err != nil || bytes < 0 {
		return 0, invalidByteQuantityError
	}

	switch multiple {
	case "E", "EB", "EI":
		return uint64(bytes * EXABYTE), nil
	case "P", "PB", "PI":
		return uint64(bytes * PETABYTE), nil
	case "T", "TB", "TI":
		return uint64(bytes * TERABYTE), nil
	case "G", "GB", "GI":
		return uint64(bytes * GIGABYTE), nil
	case "M", "MB", "MI":
		return uint64(bytes * MEGABYTE), nil
	case "K", "KB", "KI":
		return uint64(bytes * KILOBYTE), nil
	case "B":
		return uint64(bytes), nil
	default:
		return 0, invalidByteQuantityError
	}
}

func CoreSize(milli uint64) string {
	base := milli / 1000
	if base >= 1 {
		r := milli - base*1000

		if r != 0 {
			point := float64(base) + float64(r)/float64(1000)

			return strings.TrimRight(fmt.Sprintf("%.3f", point), "0")
		}

		return strconv.FormatUint(base, 10)
	}

	return fmt.Sprintf("%dm", milli)
}

func ToMilliCores(s string) (uint64, error) {

	if strings.Contains(s, "m") {
		return strconv.ParseUint(strings.TrimRight(s, "m"), 10, 64)
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}

	return uint64(1000 * f), nil
}
