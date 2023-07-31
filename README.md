# Golang API Testing Project

This Project for Golang project designed for API testing. The project is set up to help you create and run automated tests for Rancher Login APIs.

## Prerequisites

Before you can build and run the API testing project, ensure you have the following prerequisites installed on your system:

1. **Golang**: Make sure you have Golang installed. You can download and install it from the official Golang website: [https://golang.org/](https://golang.org/)

2. **Setup Rancher API Step BY Step**
```bash
1. Login into rancher UI with login password.
2. Select User Avatar -> API & Keys from the User Settings menu in the upper-right.
3. Click create API Key.
4. Enter a description for the API key and select an expiration period or a scope.
5. Click on create button.
7. Step Result: Your API Key is created. Your API Endpoint, Access Key, Secret Key, and Bearer Token are displayed.
   Use the Bearer Token to authenticate with Rancher CLI.
8. Copy the information displayed to a secure location.
9. Once API is created go back to Account and API Keys page and search for your key.
10.At your API key on extreme right you will find three dots (Option). Click on it and select "View in API"
11.Another page will open check of operation section at right top corner. Click on edit button from the operations section.
12.API request will be open with all the details along with cURL command.
```

## Getting Started

Follow these steps to set up the project and run the API tests:

1. **Clone project from git**:

```bash
# git clone https://github.com/Pk8753/Rancher_API_Automation.git
```

2. **Create mod file**:

```bash
#  go mod init <mod name>
```


3. **Get json schema to validate**:

```bash
#  go get github.com/xeipuuv/gojsonschema
```


4. **Get assert module for assertions**:

```bash
#  go get -u github.com/stretchr/testify/assert
```


5. **Run test to see result on cli**:

```bash
goto test folder.
# cd src/test/


Run test.
#  go test -v
```


6. **Run test to generate json report**:

```bash
By running below file a json file will gerenrate at the path /src/test

#  go test -json > test_output.json
```



