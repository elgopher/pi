# How to install dependencies on macOS

## Install C compiler

On the latest macOS, just type clang on your terminal and a dialog would appear if you don't have clang compiler.
Follow the instruction to install it.

You might find the following error when executing clang.

```
xcrun: error: invalid active developer path (/Library/Developer/CommandLineTools), missing xcrun at: /Library/Developer/CommandLineTools/usr/bin/xcrun
```

In this case, run `xcode-select --install` and install commandline tools.

