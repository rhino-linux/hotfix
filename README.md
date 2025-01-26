### Issues Tracker

To report issues or propose new features for this repository, visit [our tracker](https://github.com/rhino-linux/tracker).

# rhino-hotfix
A hotfix utility for Rhino Linux

## Usage

```
Usage: rhino-hotfix <input> [-h]                                                                                

  Input format:
    <hotfix>[@<metalink>]

    <hotfix> (optional):  
      - `hotfix`: Fetch a specific hotfix.  
      - `hotfix@<metalink>`: Fetch from a specific repo, branch, or PR.  

    <metalink> (optional):  
      - `[user/repo]`: Use a specific repo.  
      - `[user/repo:branch]` or `[:branch]`: Use a specific branch.  
      - `[user/repo#PR]` or `[#PR]`: Use a specific PR. 
      Note: only branch or PR can be used, not both. 

  Examples:
    rhino-hotfix                    # List hotfixes from rhino-linux/hotfix. 
    rhino-hotfix :branch            # List hotfixes from an upstream branch.
    rhino-hotfix hotfix             # Fetch a hotfix from rhino-linux/hotfix.
    rhino-hotfix hotfix@#42         # Fetch a hotfix from PR #42 upstream.
    rhino-hotfix @user/repo#99      # List hotfixes from PR #99 downstream.
```

## Creating hotfixes

Hotfix scripts should conform to the following standards:
1. The top of the script should start with the shebang `#!/usr/bin/env bash`
2. All script actions should be contained within functions, with the main function being called `hotfix`
3. All variables should be localized within the functions

Here are two examples:
```bash
#!/usr/bin/env bash

hotfix() {
  local names messages
  names=("test1" "test2")
  messages=("This is the first message" "This is the second")

  for i in "${!names[@]}"; do
    echo "${names[i]}: ${messages[i]}"
  done
}
```

```bash
#!/usr/bin/env bash

hotfix() {
  local names messages
  names=("test1" "test2")
  messages=("This is the first message" "This is the second")

  for i in "${!names[@]}"; do
    subprocess "${names[i]}" "${messages[i]}"
  done
  unset -f subprocess
}

subprocess() {
  local name="${1}" message="${2}"
  echo "${name}: ${message}"
}
```

If a script ever needs to `cd` to a location, the variable `${SRCDIR}` is available to them to return to their starting points.

Once written, the hotfix can be placed in `scripts/`, and must end with the prefix `.sh`. Then, you can run:
```bash
go run manager.go add -t <name> -d "<description>" -s scripts/<scriptname>.sh
git add hotfixes.json scripts/<scriptname>.sh
```
Name and script name do not necessarily need to match. Name will be what is called to run the hotfix (`rhino-hotfix <name>`).
