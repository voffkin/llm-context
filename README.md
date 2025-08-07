# LLM Context Generator

This command-line utility reads specified directories or files, processes their content, and consolidates it into a single file. This is particularly useful for providing context to Large Language Models (LLMs).

## Features

*   Processes files and directories recursively.
*   Filters out comments and empty lines.
*   Wraps file content in XML-style tags for better readability and parsing.
*   Configuration via a simple `.llm-context.json` file.

## Configuration

The utility is configured via a `.llm-context.json` file in the root of your project.

**Configuration options:**

*   `dir`: An array of directory paths to process recursively.
*   `file`: An array of individual file paths to process.
*   `output`: The path to the output file.
*   `cut_comments`: A boolean value (`true` or `false`) that determines whether to remove comments and empty lines from the files.
*   `exceptions`: A map of filenames to language identifiers. This is useful for files without extensions or with non-standard names.

**Example `.llm-context.json`:**

```json
{
  "dir": ["app"],
  "file": [
    "config/one.rb",
    "config/two.rb"
  ],
  "output": "llm-context.xml",
  "cut_comments": true,
  "exceptions": {
    "Procfile.dev": "",
    "Gemfile": "ruby",
    "Rakefile": "ruby",
    "Capfile": "ruby"
  }
}
```

## Build

To build the application, you can use the provided `Makefile`:

```bash
make build
```

This will compile the application and place the binary in the `bin` directory.

## Usage

After building the application, you can run it from the root of your project:

```bash
./bin/llm-context
```

The utility will then generate the output file as specified in your configuration.

## Global Access

To run the `llm-context` utility from any directory, you can add its `bin` directory to your system's `PATH` environment variable.

1.  From the root of the `llm-context` project directory, add the `bin` directory to your shell's configuration file.


    **For Zsh (`.zshrc`):**

    ```bash
    export PATH="$PATH:$HOME/path_to_folder/bin"
    source ~/.zshrc
    ```

2.  Now you can run the utility from any directory:

    ```bash
    llm-context
    ```

    If a `.llm-context.json` file exists in the current directory, the utility will generate the context file.

## Example Output

Given the following files:

**`config/one.rb`**
```ruby
puts "hello world"
```

**`config/two.js`**
```javascript
function example() {
	// log 'example'
	console.log('example')
}
```

The generated `llm-context.md` would look like this:

```xml
<file name="config/one.rb" lang="rb">
puts "hello world"
</file>

<file name="config/two.js" lang="js">
function example() {
console.log('example')
}
</file>
```