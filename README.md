PGDays'17 Hyperledger Fabric Demo
===

In this repository implemented demo chaincodes with primary goal to demonstrate
capabilities of Hyperledger Fabric chaincode development process to be
presented at PGDays'17 conference.


Use case
---

An example introduces simple bussines logic for accounts management, presented
two basic workflows:

1. **Adding and managing person information in the ledger**
2. **Opening new account, which has to be linked to real person**


Person management chaincode
---

Chaincode capable to handle opeations of maintaining records of persons in the
ledger, supports following operations:

1. Adding new person record to ledger
2. Retrieving person record from the ledger
3. Removing personal record
5. Update person address information
4. Providing history overview of person record change

For demo purpose person record simplified down to:

```
type Person struct {
    ID          string
    Name        string
    Phone       string
    Address     string
}
```

Account management chaincode
---

Chaincode which allows to create "bank" accounts, while doing this operation it
has to  verify that account binded to real person, e.g. has to leverage person
management chaincode to verify the information before creating new account.

Account struct is:

```
type Account struct {
    PersonID string
    Number   string
    Balance  float64
}

```


Commands to run for install and chaincode instantiation:

1. Install persons chaincode
```
peer chaincode install -n prsnMgmt -v 1.0 -p github.com/C0rwin/pgdays/chaincode/personsV3
```

2. Instantiate persons chaincode
```
./build/bin/peer chaincode instantiate -n prsnMgmt -v 1.0 -C mychannel -c '{"Args": ["init"]}' -o localhost:705
```

3. Install account chaincode
```
peer chaincode install -n accMgmt -v 1.0 -p github.com/C0rwin/pgdays/chaincode/account
```

4. Instantiate account chaincode
```
peer chaincode instantiate -n accMgmt -v 1.0 -C mychannel -c '{"Args": ["init", "prsnMgmt", "mychannel"]}' -o localhost:7050
```

5. Adding new person record

```
peer chaincode invoke -n loanApp -v 1.0 -C mychannel -c '{"Args": ["addPerson", "{ \"id\": \"1111111\", \"Name\": \"Vasily Terkin\", \"Phone\": \"0544444444\", \"address\": \"Moscow, RU\"}"]}' -o localhost:7050
```
