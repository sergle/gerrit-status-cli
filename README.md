# gerrit-status-cli

cli tool to view list of active review in Gerrit

[[https://raw.githubusercontent.com/sergle/gerrit-status-cli/master/img/screenshot.png]]

1. Review status
  * ⨂  cannot be merged
  * ✔  approved (+2 mark)
  * +  looks good, but someone needs apporove it (+1 mark)
  * -  someones prefer you didn't submit this (-1 mark)
  * ✘  do not submit (-2 mark)

2. Project name (can be aliases to short via config file)
3. Branch name
4. Review owner
5. Subject (shortened)
6. Last updated time
7. CI username, color used to mark verified/failed status
8. List of reviewers, color used to show their marks

First section shows reviews owned by user.
Second section shows reviews where user is in reviewers list


Config file
```ini
[gerrit]
user = my-login
# to generate password: Gerrit -> Settings -> HTTP Password 
password = my-password
host = example.com:8443
# CI username (global for all projects)
ci = jenkins
# limit number of concurrent requests
connections = 15

# 256 color terminal
[color]
verified = 77
bad = 160
plus = 215
ok = 120
absent = 238
title = 110

# aliases to show instead of project name (when they are long)
[project]
alias=short:some-very-long-project-name
alias=p2:another-very-long-name
# override CI username for specific project
ci=p2:robot
```

