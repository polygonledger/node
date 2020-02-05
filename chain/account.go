package chain

type Account struct {
	AccountKey string
	//Balance    int
}

func AccountFromString(key string) Account {
	return Account{AccountKey: key}
}
