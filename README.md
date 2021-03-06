# EOS to OPA

This tool will convert an Arista EOS config to a JSON file for use with [Open Policy Agent](https://www.openpolicyagent.org/docs/latest/). You can then run OPA against that file with a policy to query if policies are being met. For example, if you have a config that has enabled telnet. Create a rego file that will evaluate the telnet field of the input files, and can then be queried to see if telnet has been enabled. This becomes more powerful when you have multiple devices input and can run a single query against an array of configs.

**Currently working on parsing the management section of the configuration**

## Example files
input-example.json - example output file parsed through app ready for input into opa
policy.rego - example rego file for policy
output.json - example output from running opa on input and policy
config files

## Example runs:
```
go run main.go -input="dmz-lf18.config" 1> test.json      
opa eval --fail-defined -i test.json -d policy.rego "data.policy.violation[x]"
echo $?

```

```
$ go run main.go 1> test.json

# What is the output of the SSH policy run against the input
$ opa eval -i input-example.json -d policy.rego "data.policy.ssh"

# Is SSH enabled?
$ opa eval --fail-defined -i input-example.json -d policy.rego "data.policy.ssh" | jq '.result[0].expressions[0].value'
true

# Is telnet enabled?
opa eval -i test.json -d policy.rego "data.policy.telnet"


```

The policy defined in policy.rego has the following in the ssh section:

```
ssh {
	m := input.management
  not m.ssh.shutdown
}
```

This tests if SSH is not shutdown. If it is it will be undefined. Given our input.json file has ssh in a no shutdown state (default) it will return true when we evaluate it:

```
$ opa eval -i test.json -d policy.rego "data.policy.ssh"
{
  "result": [
    {
      "expressions": [
        {
          "value": true,
          "text": "data.policy.ssh",
          "location": {
            "row": 1,
            "col": 1
          }
        }
      ]
    }
  ]
}
```

In contrast our policy wants to validate that telnet is shutdown (i.e. no telnet) so our policy looks like this:

```
telnet {
	m := input.management
  m.telnet.shutdown
}
```

Since this evaluates to `false` the output is undefined:

```
opa eval -i test.json -d policy.rego "data.policy.telnet"
{}
```

We can use the `--fail` to send a non-zero exit code if the policy fails

```
$ opa eval --fail -i test.json -d policy.rego "data.policy.ssh"
$ echo $?
0

$ opa eval --fail -i test.json -d policy.rego "data.policy.telnet"
$ echo $?
1
```

REPL:

```
$ opa run policy.rego repl.input:input-example.json
OPA 0.25.2 (commit 4c6e524, built at 2020-12-08T16:56:55Z)

Run 'help' to see a list of commands and check for updates.

> data.policy
{
  "hello": false,
  "ssh": true,
  "telnet": true
}
```

## TODO
- [x] Export to proper JSON format
- [ ] Add more features to SSH, API, and Telnet
- [ ] Refactor parseManagement to take a device and return a device
- [x] How to handle configlets? Configlet object? - will pass it in as text file
- [ ] Add security
- [x] Expand JSON to have hostname
- [ ] handle multiple switch configs 
- [x] parse hostname
- [ ] test code
- [ ] reorg file structure
  - parser stuff into different files under pkg
  - test cases under ./test
  - optional move main under cmd
  - remove parsing code from main
- [ ] add snmp community parsing for checking community strings
- [ ] static user accounts vs radius
- [ ] inbound filter for management access
- [ ] logging turned on and correct
- [ ] dns configured correctly
- [ ] ntp configured correctly
- [ ] test passwords against common hashes
- [ ] config validation as a service
- [ ] Create constructor for EOS Device with defaults set (i.e. logging on)
- [ ] allow output filename
- [x] EOS Software / HW model number in opa 
    - This could enable CVE check
- [ ] CloudVision to build Rego files to enforce policy dynamically
- [ ] Handle 'no' commands

