# EOS to OPA

This tool will convert an EOS config to a input.json file. You can then run OPA against that file with a policy to query if policies are being met. For example, if you have a config that has enabled telnet. Create a rego file that will evaluate the telnet field of the input files, and can then be queried to see if telnet has been enabled. This becomes more powerful when you have multiple devices input and can run a single query against an array of configs.

**Currently working on parsing the management section of the configuration**

## Example files
input-example.json - example output file parsed through app ready for input into opa
policy.rego - example rego file for policy
output.json - example output from running opa on input and policy
config files


Example runs:

```
# What is the output of the SSH policy run against the input
$ opa eval -i input-example.json -d policy.rego "data.policy.ssh"

# Is SSH enabled?
$ opa eval --fail-defined -i input-example.json -d policy.rego "data.policy.ssh" | jq '.result[0].expressions[0].value'
true

```


REPL:

```
$ eosopa opa run policy.rego repl.input:input-example.json
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
- [ ] Export to proper JSON format
- [ ] Add more features to SSH, API, and Telnet
- [ ] Add security


