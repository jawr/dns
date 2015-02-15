package domains

var byUUID map[string]Domain

func init() {
	byUUID = make(map[string]Domain)
}

func addToCache(domain Domain) {
	if len(byUUID) > 1000 {
		i := 0
		for k, _ := range byUUID {
			if i > 100 {
				break
			}
			delete(byUUID, k)
			i++
		}
	}
	byUUID[domain.UUID.String()] = domain
}
