# Binary Installation
+ **Windows**: 
    + **Step 1**: Download the pre-built binaries from [here](https://github.com/meeranh/recycleMe/releases/download/v0.1/recycleMe-x86_64_windows.exe).
    + **Step 2**: Add binary to your PATH environment variable.
    + **Step 3**: Run `recycleMe-x86_64_windows.exe yourPlainTextFile.md` from the command line.

+ **Linux**: 
    + **Step 1**: Download the pre-built binaries from [here](https://github.com/meeranh/recycleMe/releases/download/v0.1/recycleMe-x86_64_linux)
    + **Step 2**: Run `chmod 777 recycleMe-x86_64_linux` to make the binary executable.
    + **Step 2**: Run `./recycleMe-x86_64_linux yourPlainTextFile.md` from the command line.

# Building From Source
+ **Step 1**: Clone the repository and change directory to the cloned repository.
```unix
git clone https://github.com/meeranh/recycleMe.git
cd recycleMe
```

+ **Step 2**: Install dependencies and build the binary
```unix
go mod tidy
go build main.go
```

# Purpose
+ Content writers save a lot of time by using NLP tools such as ChatGPT, but if they are expected to produce non-AI flagged content, they will have to spend a considerable amount of time rephrasing their content. 
+ The purpose of this tool is to simplify this exact process by taking in a plain text file and returning a rephrased version of the content all in the command line.
+ Also, this was my first Golang project, so I wanted to build something simple and useful.
