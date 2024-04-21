# Project Title

COMP4911 P2P Client

## Requirements

- Golang installation
- Support for bash scripts (Optional)

## Installation

1. Clone the git repository
2. Build project with go or with build script which will create 3 separate folders containing executables and download folders (macOS or Linux)
3. run the compiled executable

   ```sh
   go build . -o p2p_client 
   ```
   ```sh
   ./build.sh
   ```
   
## Usage

The program runs in a simple repl. follow the configuration prompts. Afterwards, enter the name of a file you want
if there are users with the file, then you will be shown a list of them. simply enter the number of the peer you
wish to get the file from and the download will proceed.
