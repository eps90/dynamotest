PROJECT IS AN EXPERIMENT UNDER DEVELOPMENT, API MAY CHANGE ANYTIME. USE IT ON YOUR OWN RISK.
===

# Dynamotest

**Dynamotest** is my first side project written for my own needs to help me easily run integration tests **in parallel**.
Every new instance creates a table dedicated for a test case according to **migrations definitions** and run fixtures 
on it selectively by choosing **fixtures files to run**. 

Project's is **still under development** so you use it on your own risk. 
Feel free to discuss on new features but please take a look at issues list first to check if it's something I planned.
I'm open to any discussion.

Feel free to check [TODO List](#todos-and-other-plans) to check what are nearest things to do

## Table of contents

* [Requirements](#requirements)
* [Installation](#installation)
* [Usage](#usage)
    * [What happens under the hood?](#what-happens-under-the-hood)
* [Configuring and extending](#configuring-and-extending)
     * [Loading migration and fixtures files](#loading-migration-and-fixtures-files)                                                         
     * [Loading particular fixtures](#loading-particular-fixtures)
     * [Resolving table names](#resolving-table-names)
* [TODOs and other plans](#todos-and-other-plans)

## Requirements 
I started developing it in **Go 1.11**, currently I'm using **Go 1.12** and still works fine. 
Also, I'm using **go modules** here.

## Installation

```bash
go get -u github.com/eps90/dynamotest
```

With Go 1.11 or later you can use the following, optionally appending the version number:

```bash
GO111MODULE=on go get github.com/eps90/dynamotest@v0.0.1
```

## Usage

Let's start with quick example of how this packages can be used in the easiest way. 
First, create a **migration config** for the table you want to test.
```json
{
  "TableName": "pets",
  "AttributeDefinitions": [
    {
      "AttributeName": "PK",
      "AttributeType": "S"
    },
    {
      "AttributeName": "SK",
      "AttributeType": "S"
    }
  ],
  "KeySchema": [
    {
      "AttributeName": "PK",
      "KeyType": "HASH"
    },
    {
      "AttributeName": "SK",
      "KeyType": "RANGE"
    }
  ],
  "ProvisionedThroughput": {
    "ReadCapacityUnits": 10,
    "WriteCapacityUnits": 10
  }
}
```
The syntax is identical to the one known from CloudFormation's DynamoDB table schema definition.

Let's also create some basic fixtures for our new table:

```json
{
  "table": "pets",
  "items": [
    {
      "PK": "pets",
      "SK": "pet_323132311",
      "name": "First pet"
    },
    {
      "PK": "pets",
      "SK": "pet_7878998732",
      "name": "Second pet"
    }
  ]
}
```

Having such file the lib will insert all rows defined under `items` property into table called as in `table` key. 
No types definitions are needed.

Finally, given our migrations are in `migrations` directory and fixtures in `fixtures` directory, we can set up our test table:

```go
package main
import (
    "fmt"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/eps90/dynamotest"
)

func main() {
    dynamoSvc := dynamoDbInstance()
    dynamoTester := dynamotest.NewDefaultDynamoTester(
        dynamoSvc,
        "./migrations",
        "./fixtures",
    )

    dynamoTester.MustLoadFixtures()
    fmt.Println("Table created and fixtures loaded!")
}

func dynamoDbInstance() *dynamodb.DynamoDB {
    return dynamodb.New(session.New())
}
``` 

### What happens under the hood?

* The package loads all files in directories provided as 2nd and 3rd arguments to constructor function
* For each table defined in a **fixture** dynamotest will clear and create a new table with a given name using a `TableNameResolver`. 
Depending on what implementation is used the name is different. By default the `TimestampTableNameResolver` is used together with
`MemoizedTableNameResolver` which means it will create a table with timestamp appended and the name is memoized, that means 
that if you ask for a table name anytime you will get it same within the instance. For table name resolvers, please see below.
* Each fixture is loaded in default order (on linux is alphabetical order) if you don't provide a list of fixtures

## Configuring and extending 

### Loading migration and fixtures files

The default instance of `DynamoTester` looks JSON files in provided directories recursively. 
Type responsible for that is `FilesystemDirectoryLoader` which is constructed with two parameters: directory and extension.
To create custom Loader, e.g. for YAML migrations use the following:
```go
loader := NewFilesystemDirectoryLoader("config/migrations", "yml")
```

> **Unsupported yet possible:** Since it uses glob behind the curtains, you can pass several extensions if you want to
 ```go
loader := NewFilesystemDirectoryLoader("config/migrations", "(json|yml)")
``` 

Having that, you can replace the loader:
```go
dynamoTester := dynamotest.NewDefaultDynamoTester(/* ... */)
dynamoTester.FixturesLoader = NewFilesystemDirectoryLoader("config/fixtures", "yml")
dynamoTester.Migrator.MigrationsLoader = NewFilesystemDirectoryLoader("config/migrations", "yml")
```

### Loading particular fixtures

You can pass a list of fixtures to be executed. Remember to not drop extension from the name.
```go
dynamoTester.MustLoadFixtures("mypets", "subdir/also_pets")
```

### Resolving table names

As already written before, the `DynamoTester` generates a table name appending the timestamp.
You can fetch resolved table name in the following way:

```go
t := dynamotester.TableNameFor("pets")
fmt.Println(t)
// Outputs: pets_1568231521
```

Moreover for given input name result will be memoized:
```go
t1 := dynamotester.TableNameFor("pets")
t2 := dynamotester.TableNameFor("pets")

if t1 != t2 {
    // will never happen
}
```
By default this job is being made by two structs:
* `TimestampTableNameResolver`
* `MemoizedTableNameResolver`

Both implement `TableNameResolver` interface so feel free to implement your own and replace is in `DynamoTester` by
```go
dynamoTester.TableNameResolver = &MyCustomTableNameResolver{}
```

**Package comes with two more implementations** which are:
* `DefaultTableNameResolver` which doesn't change the input table name (it remains the same)
* `RandomTableNameResolver` which appends random string (other than timestamp) to input table name

Feel free to take a look at API docs for more.

## TODOs and other plans
There are couple of things to be done and I'm completely aware of it. I exported this lib to make it usable in couple of projects already.

From things I know **for sure** they have to be done:

- [ ] Complete comment-based documentation to generate nice online godocs. That includes: 
    - [ ] Comments for exported types and functions
    - [ ] Examples, especially for reusable and replaceable parts  
    - [ ] Testing instructions
- [ ] Integrate with a CI and code quality tools
- [ ] Add integration tests against local DynamoDB
- [ ] Add license
- [ ] Create a Makefile for the project with most repeating actions
- [ ] Drop `pkg/errors` package in favor of `xerrors` or built-in `errors` package
- [ ] Run operations concurrently (performance; may affect API)

Also I'm considering few things:
- [ ] Separate things to separate packages (BC)
- [ ] Extract _migration_ and _fixtures_ related things into separate repositories (BC)
