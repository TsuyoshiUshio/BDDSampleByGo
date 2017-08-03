BDD with Go Using Ginkgo with Azure
===

I had a mob programming session with my colleague. I'd like to have some experiment for me before I forget what I learnt today. 

I'd like to enable BDD testing, JUnit style reporting, and Mocking. 


# 1. Setting up Ginkgo

[Ginkgo](http://onsi.github.io/ginkgo/#shared-example-patterns)is a golang BDD testing framework. 

## 1.1. GOPATH and directory structure

I create these directory structure. I create a `/Users/ushio/Codes/gosample` directory. I'd love to isolate directory for GOPATH. 

```
/Users/ushio/Codes/gosample/
        src/
          github.com/
             TsuyoshiUshio/
                  BDDSampleByGo/

```

I usually use setenv.sh like this. `source setenv.sh`.

```
export GOPATH=/Users/ushio/Codes/gosample
export PATH=$GOPATH/bin:$PATH
```
## 1.2. Library installation

Now you are ready to install a lot of libraries. 

```
$ go get github.com/onsi/ginkgo/ginkgo
$ go get github.com/onsi/gomega
```
Now you can use `ginkgo` command

I use [glide](https://github.com/Masterminds/glide) for package management.  You can install it via one liner on the github or by homebrew. 

I'm not sure it is a good practice, but I install ginkgo both on `go get` and `glide get`. I want to use ginkgo command. It is on the bin directory of the library. I don't want to set path on vendor directory. If you know the best practice, please let me know.

```
$ glide get github.com/onsi/ginkgo/ginkgo
$ glide get github.com/onsi/gomega
```

glide will install the libraries on the vendor directory. It is suitable for automated testing. 

## 1.3. bootstrap ginkgo

If you want to use ginkgo, you need to bootstrap it.  

```
ginkgo bootstrap
```

Then you can see `DIR_NAME_suite_test.go`  file. I create a cmd directory. then I execute `ginko bootstrap`. It helps to use the BDD testing framework enable to use from `go test`

_cmd_suite_test.go_

```
package cmd_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cmd Suite")
}
```
Now you are ready to go write a test!

# 2. Writing a test 

I'm an extream programming enthusiast. I'd like write test code before writing a production code. Let's start with a test code. 

## 2.1. Behavior first.

I want to create a simple command enable us to access Azure Key Vault. Start with very simple test case.

_key_vault_test.go_

```
package cmd_test

import (
	. "github.com/onsi/ginkgo"
)

var _ = Describe("KeyVaultClient", func() {
	Context("When it has a secret", func() {
		It("returns the secret value", func() {

		})
	})
})
```

It doesn't include any testing. Just an behavior. However you are ready to execute test. Just do it.

```
$ go test
Running Suite: Cmd Suite
========================
Random Seed: 1501646102
Will run 1 of 1 specs

•
Ran 1 of 1 Specs in 0.000 seconds
SUCCESS! -- 1 Passed | 0 Failed | 0 Pending | 0 Skipped PASS
ok  	github.com/TsuyoshiUshio/BDDSampleByGo/cmd	0.011s
```

It works. Let's add actual testing!

## 2.2. Adding test code

I add a test for this. 

```
package cmd_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("KeyVaultClient", func() {
	Context("When I request a secret", func() {
		client := KeyVaultClient{}
		It("returns the secret value", func() {
			Expect(client.GetSecretValue()).To(Equal("password"))
		})
	})
```

However, I can't compile it because I dont' have KeyVaultClient. Let's create it. 

_key_vault.go_

```
package cmd

type KeyVaultClient struct {
}

func (c *KeyVaultClient) GetSecretValue() string {
	return "password"
}

```

The directory structure

```
$ tree -I vendor
.
├── README.md
├── cmd
│   ├── cmd_suite_test.go
│   ├── key_vault.go
│   └── key_vault_test.go
├── glide.lock
└── glide.yaml

```

On the root directory of my repo, I can test it.

```
$ go test -v cmd/*
=== RUN   TestCmd
Running Suite: Cmd Suite
========================
Random Seed: 1501648188
Will run 1 of 1 specs

•
Ran 1 of 1 Specs in 0.000 seconds
SUCCESS! -- 1 Passed | 0 Failed | 0 Pending | 0 Skipped --- PASS: TestCmd (0.00s)
PASS
ok  	command-line-arguments	0.011s
```
It passes. If you change test little bit like this

```
Expect(client.GetSecretValue()).To(Equal("not password"))
```

It fails as expected.

```
$ go test -v cmd/*
=== RUN   TestCmd
Running Suite: Cmd Suite
========================
Random Seed: 1501648285
Will run 1 of 1 specs

• Failure [0.001 seconds]
KeyVaultClient
/Users/ushio/Codes/gosample/src/github.com/TsuyoshiUshio/BDDSampleByGo/cmd/key_vault_test.go:16
  When I request a secret
  /Users/ushio/Codes/gosample/src/github.com/TsuyoshiUshio/BDDSampleByGo/cmd/key_vault_test.go:15
    returns the secret value [It]
    /Users/ushio/Codes/gosample/src/github.com/TsuyoshiUshio/BDDSampleByGo/cmd/key_vault_test.go:14

    Expected
        <string>: password
    to equal
        <string>: not password

    /Users/ushio/Codes/gosample/src/github.com/TsuyoshiUshio/BDDSampleByGo/cmd/key_vault_test.go:13
------------------------------


Summarizing 1 Failure:

[Fail] KeyVaultClient When I request a secret [It] returns the secret value 
/Users/ushio/Codes/gosample/src/github.com/TsuyoshiUshio/BDDSampleByGo/cmd/key_vault_test.go:13

Ran 1 of 1 Specs in 0.001 seconds
FAIL! -- 0 Passed | 1 Failed | 0 Pending | 0 Skipped --- FAIL: TestCmd (0.00s)
FAIL
exit status 1
FAIL	command-line-arguments	0.012s
```
## 2.3. Adding a testing report

I'm a DevOps guy. I need a unit test report. 

You can add it on the bootstrap file. Ginkgo proivides a customer reporter for generating JUnit compatible XML output.

_cmd_suite_test.go_

```
package cmd_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
)

func TestCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "Cmd Suite", []Reporter{junitReporter})
}
```
if you execute the `go test -v cmd/*` the behavior is the same, however, you can see the, "junit.xml" file on the cmd directory. 

```
$ ls
cmd_suite_test.go	key_vault.go
junit.xml		key_vault_test.go

$ cat junit.xml 
<?xml version="1.0" encoding="UTF-8"?>
  <testsuite tests="1" failures="0" time="0.000146342">
      <testcase name="KeyVaultClient When I request a secret returns the secret value" classname="Cmd Suite" time="8.1732e-05"></testcase>
```

If you want to go parallele testing per node, you can write like this. Please refer [Generating JUnit XML Output](http://onsi.github.io/ginkgo/#the-spec-runner).

```
junitReporter := reporters.NewJUnitReporter(fmt.Sprintf("junit_%d.xml", config.GinkgoConfig.ParallelNode))
```
# 3. Writing a Mock 

I just implemented a fake test. I'd like to implement actual implementation. If you want to implement it, you need to add Azure SDK for your code. 

```
$ go get github.com/Azure/go-autorest/autorest
$ go get -u github.com/Azure/azure-sdk-for-go/arm/keyvault 
```

**NOTE: I want to use glide for this. However, I've got the issue. I report it. [glide can't solve the dependency](https://github.com/Azure/azure-sdk-for-go/issues/713)**
Also, you can't do like this. `go get github.com/Azure/azure-sdk-for-go` you'll get 

```
$ go get -u github.com/Azure/azure-sdk-for-go 
package github.com/Azure/azure-sdk-for-go: no buildable Go source files in /Users/ushio/Codes/gosample/src/github.com/Azure/azure-sdk-for-go
```

I wrote a Mock Sample blog. [Go Mock](https://github.com/TsuyoshiUshio/GoMock)