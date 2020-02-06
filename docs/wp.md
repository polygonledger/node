# basic whitepaper

(alpha)

Polygon is based on an accounts based ledger.

The movement of value is formally defined by transactions. Traditionally extensbility of blockchains were
designed through programmable transactions starting with Bitcoin. The sender and receiver of Bitcoins are public keys
and the script language defines the possible transactions between these endpoints. The transactions are type-less
in that the code is defined in byte code expressions. Ethereum contracts can be much more complex, but 
similarly are byte code. In the case of Ethereum there are two distinct transactions: Ether and contracts. What a contract does can not be seen unless one understands the higher level code which was compiled into the byte code. 

In Polygon transactions come in various distinct forms i.e. types of transactions (in principle accounts could have types too, that is open for research)

account => transactionTypeSimple => account

account => transactionTypeComplex => account

This means we can contruct the system through iterations of complexity. The first iteration has only simple transaction
types i.e. sending and receiving electronic cash.