# bareclone
Simple cli to setup a bare clone repository with worktrees

# Usage
Copy binary and configure the global `.gitconfig` for an alias
For example if binary is in `$HOME/bareclone`, add this:
```.gitconfig
;.gitconfig

[alias]
    bare-clone= "!$HOME/bareclone"
```

Then you can clone a repository like this:
```bash
got bare-clone <REPO URL>
```


