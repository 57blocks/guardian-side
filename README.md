# Guardian Side

`guardian-side` is a CLI tool designed for downloading bloom filter files periodically, particularly useful for services
in need of managing bloom filters for blockchain nodes like Geth. It supports retries with a delay if the download fails
and allows flexible output file paths depending on user needs.

## Features

* Automatically downloads bloom filter files every 24 hours.
* Configurable output directory for storing the downloaded bloom filter.
* Supports retrying download with delays in case of failures.
* Customizable output path based on system type (Linux, MacOS).

## Installation

To install guardian-side from source, you must have Go (1.16 or above) installed.

1. Clone the repository:

```shell
git clone https://github.com/piplabs/guardian-side.git
```

2. Navigate into the project directory:

```shell
cd guardian-side
```

3. Install the dependencies:

```shell
go mod tidy
```

4. Build the application:

```shell
go build -o guardian-side cmd/*.go
```

5. Move the binary into your PATH or run it locally:

```shell
mv guardian-side /usr/local/bin/
```

Now you can execute `guardian-side` as a CLI tool on your terminal.

## Usage

Once the `guardian-side` has been installed, you can invoke it by running:

```shell
guardian-side [flags]
```

### Flags

The tool allows customization of where the bloom filter is saved by specifying flags. By default, the output directory
is automatically configured depending on the user's OS (`$HOME/geth/guardian` on Linux
or `$HOME/Library/Story/geth/guardian` on MacOS).

You can override this by providing your own output directory using the `-o` or `--output-dir` flag.

### Available flags:

* `-o`, `--output-dir`: The directory to store the bloom filter files. (default: OS-specific,
  e.g., `$HOME/geth/guardian` for Linux)

### Examples

1. *Basic usage (use default path)*: To run the program using the default output path for your system (
   i.e., `$HOME/geth/guardian` on Linux, `$HOME/Library/Story/geth/guardian` on Mac):

```shell
guardian-side
```

2. *Specifying custom output directory*: To specify a custom output directory for bloom filters, use the `-o`
   or `--output-dir` flag:

```shell
guardian-side -o /path/to/custom/directory
```

3. *Running the downloader in the background*: Since this tool is designed to run periodically, you can run it in the
   background using the following:

```shell
nohup guardian-side > downloader.log 2>&1 &
```

This runs the downloader in the background, redirecting the output to a log file.