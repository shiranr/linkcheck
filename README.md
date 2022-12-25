# link-check

This linter is a golang markdown link verifier.
Originally created and currently running on the [CSE
playbook](https://github.com/microsoft/code-with-engineering-playbook).

Contains Megalinter descriptor.

## Process
When running only linkcheck:
1. By default, the checker will use the configuration linkcheck.json located [here](configuration/linkcheck.json).
   - We can pass a new configuration file by using `--config PATH` 
2. We will be working on the directory the command was running from.
   - We can pass files instead by `linkcheck readme.md readme2.md`
3. Scan the directory and search for *.md files.
4. Open each file and extract links from it.
5. Divide the links into one of the categories: Email, Folder, URL.
6. Analyze link.

## Installation

By default, the checker will use the configuration linkcheck.json located [here](configuration/linkcheck.json).

```shell
go install github.com/shiranr/linkcheck@latest

linkcheck
linkcheck README.md
linkcheck --config linkcheck.json README.MD
```
