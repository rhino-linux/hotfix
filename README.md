### Issues Tracker

To report issues or propose new features for this repository, visit [our tracker](https://github.com/rhino-linux/tracker).

# rhino-hotfix
A hotfix utility for Rhino Linux

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
