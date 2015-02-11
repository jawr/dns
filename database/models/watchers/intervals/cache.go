package intervals

var byID map[int32]Interval
var byValue map[string]Interval

func init() {
	intervalList, err := GetAll().List()
	if err != nil {
		panic(err)
	}

	byID = make(map[int32]Interval)
	byValue = make(map[string]Interval)

	for _, interval := range intervalList {
		byID[interval.ID] = interval
		byValue[interval.Value] = interval
	}
}
