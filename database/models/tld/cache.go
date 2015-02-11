package tlds

var byID map[int32]TLD
var byName map[string]TLD

func init() {
	tldList, err := GetAll().List()
	if err != nil {
		panic(err)
	}

	byID = make(map[int32]TLD)
	byName = make(map[string]TLD)

	for _, tld := range tldList {
		byID[tld.ID] = tld
		byName[tld.Name] = tld
	}
}
