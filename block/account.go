package block

type Account struct {
	AccountKey string `edn:"AccountKey"`
	//Balance    int
	//Name string //if introduce names for accounts should think about prefixes to enable hierarchies
}

func AccountFromString(key string) Account {
	return Account{AccountKey: key}
}
