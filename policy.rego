package policy

ssh {
	m := input.management
    not m.ssh.shutdown
}

telnet {
	m := input.management
    m.telnet.shutdown
}


# in the following example we will say 'allow' is true if there are no violations
default allow = false

allow = true {                                      # allow is true if...
    count(violation) == 0                           # there are zero violations.
}

violation["ssh"] {
	m := input.management
    m.ssh.shutdown # revese logic from before.  SSH is in violation (true) if it is shutdown
}

violation["telnet"] {
	m := input.management
    not m.telnet.shutdown # Telnet is in violation (true) if it is not shutdown
}

