# link-check

This linter is a golang markdown link verifier.
The goal of this linter is to verify links mentioned in markdown documentation are working properly, alive and not stale. 

This linter is easily used and envoked with MegaLinter as it contains a MegaLinter descriptor.

> **_NOTE:_**  Originally created and currently running on the [CSE
playbook](https://github.com/microsoft/code-with-engineering-playbook).

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

## Configuration

Currently, there are several parameters which can be customized, more to come in the future. :) For more details, see [linkcheck.json](configuration/linkcheck.json):
1. exclude_links - links which we want to skip and the link check to ignore.
2. only_errors - true by default - print only errors or all the logs including successful links.
3. project_path - ability to scan an entire folder instead of files one by one as a givne list.
4. serial - false by default - work parallel (use goroutines) or serializable.