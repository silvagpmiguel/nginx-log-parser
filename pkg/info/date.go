package info

// Date represents the date of a log entry
type Date struct {
	Day      [2]byte
	Month    [2]byte
	Year     [4]byte
	DateTime string
}

// StringToMonth transforms a valid string to the correspondent month
func StringToMonth(str string) [2]byte {
	switch str {
	case "Jan":
		return [2]byte{'0', '1'}
	case "Feb":
		return [2]byte{'0', '2'}
	case "Mar":
		return [2]byte{'0', '3'}
	case "Apr":
		return [2]byte{'0', '4'}
	case "May":
		return [2]byte{'0', '5'}
	case "Jun":
		return [2]byte{'0', '6'}
	case "Jul":
		return [2]byte{'0', '7'}
	case "Aug":
		return [2]byte{'0', '8'}
	case "Sep":
		return [2]byte{'0', '9'}
	case "Oct":
		return [2]byte{'1', '0'}
	case "Nov":
		return [2]byte{'1', '1'}
	case "Dec":
		return [2]byte{'1', '2'}
	}

	return [2]byte{'0', '0'}
}

//CompareDay does the comparison of 2 days
func (d Date) CompareDay(snd [2]byte) int {
	ret := 0
	fst := d.Day

	if fst[0] > snd[0] {
		ret = 1
	} else if fst[0] < snd[0] {
		ret = -1
	} else if fst[0] == snd[0] && fst[1] < snd[1] {
		ret = -1
	} else if fst[0] == snd[0] && fst[1] > snd[1] {
		ret = 1
	}

	return ret
}

//CompareMonth does the comparison of 2 months
func (d Date) CompareMonth(snd [2]byte) int {
	ret := 0
	fst := d.Month

	if fst[0] > snd[0] {
		ret = 1
	} else if fst[0] < snd[0] {
		ret = -1
	} else if fst[0] == snd[0] && fst[1] < snd[1] {
		ret = -1
	} else if fst[0] == snd[0] && fst[1] > snd[1] {
		ret = 1
	}

	return ret
}

//CompareYear does the comparison of two years represented by 4 bytes
func (d Date) CompareYear(snd [4]byte) int {
	ret := 0
	fst := d.Year

	for i := 0; i < 4; i++ {
		if fst[i] > snd[i] {
			ret = 1
			break
		} else if fst[i] < snd[i] {
			ret = -1
			break
		}
	}

	return ret
}
