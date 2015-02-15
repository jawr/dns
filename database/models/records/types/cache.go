package types

var byID map[int32]Type
var byName map[string]Type
var byNameAndTLD map[string]map[int32]Type

func init() {
	tList, err := GetAll().List()
	if err != nil {
		panic(err)
	}

	byID = make(map[int32]Type)
	byName = make(map[string]Type)
	byNameAndTLD = make(map[string]map[int32]Type)

	for _, t := range tList {
		byID[t.ID] = t
		byName[t.Name] = t
	}
}
