# aidev

AI developer: asks GPT-4 to modify an entire folder full of files.

**Warning**: This README and even parts of the code were written by aidev itself, via GPT-4. This is a meta-experiment. Always review the generated code before using it in your projects.

## Usage

You can invoke the program with no parameters and provide the prompt on stdin:

```
aidev
```

When prompted, enter your desired change request and end with EOF (Ctrl+D on Unix-like systems or Ctrl+Z on Windows).

The general workflow of using aidev involves the following steps:

1. (Optional) Specify the directories containing the code you want to modify using the `-d` option. By default, aidev will use the current directory.
2. Define slices in your `.aidev` configuration file to include or exclude specific parts of your codebase.
3. Provide a prompt describing the changes you want to make using the `-p` option.
4. Select a slice to use with the `-s` option.
5. The AI will process the request and generate modified versions of the relevant files.
6. Modified files will have a ".draft" extension added to their original file names. Review these files and, if the changes are satisfactory, replace the original files with the modified versions.

Remember to always review the generated ".draft" files before replacing the original files to ensure the changes made by the AI are correct and meet your requirements.

By default, aidev will ignore any files with a ".draft" extension or containing ".draft." in their name.

## Slicing

Slicing allows you to include or exclude specific parts of your codebase when using aidev. By defining slices in your `.aidev` configuration file, you can control which files or directories are considered by the AI when processing your request.

To use slicing, add a `slice` directive followed by a unique name in your `.aidev` file. Then, use the `only`, `ignore`, and `unignore` directives within the slice to specify the include and exclude patterns for that slice. To apply a specific slice when running aidev, use the `-s` option followed by the slice name.

### Examples

#### Basic slicing

In your `.aidev` file:

```
slice frontend
    only frontend/*.js
    only frontend/*.css
```

To use the `frontend` slice:

```
aidev -s frontend -p "Add a new function called 'helloWorld' that prints 'Hello, World!'"
```

#### Multiple slices

In your `.aidev` file:

```
slice frontend
    only frontend/*.js
    only frontend/*.css

slice backend
    only backend/*.go
```

To use the `backend` slice:

```
aidev -s backend -p "Add a new function called 'helloWorld' that prints 'Hello, World!'"
```

#### Combining slices with global patterns

In your `.aidev` file:

```
ignore *.log

slice frontend
    only frontend/*.js
    only frontend/*.css

slice backend
    only backend/*.go
```

The `*.log` pattern will be ignored globally, while the `frontend` and `backend` slices will still apply their specific include patterns.

## Environment Variables

The following environment variables can be set to configure aidev:

### Required

- `OPENAI_API_KEY`: Your OpenAI API key. This is required to authenticate with the OpenAI API.

### Optional

- `OPENAI_ORG`: Your OpenAI organization ID. This is optional and used for billing purposes.
- `AIDEV_SAVE_CODE`: File name to save combined code to. Defaults to an empty string.
- `AIDEV_SAVE_RESP`: File name to save the AI response to. Defaults to an empty string.
- `AIDEV_SAVE_PROMPT`: File name to save the prompt to. Defaults to an empty string.

## Usage

```
aidev [options]
```

## Configuration

The `.aidev` configuration file allows you to specify include and exclude patterns for the files and directories that should be considered by the AI. You can place this file in any directory, and the configuration will apply to that directory and its subdirectories.

Example of a `.aidev` file:

```
ignore LICENSE
ignore internal
ignore _data
ignore *.txt
ignore go.mod

slice foo
    only foo*.go
```

## Examples

### Basic usage

```
aidev -p "Add a new function called 'helloWorld' that prints 'Hello, World!'"
```

### Specify directories to include (optional)

```
aidev -d src -d lib -p "Add a new function called 'helloWorld' that prints 'Hello, World!'"
```

### Include and exclude patterns

```
aidev -i "*.go" -x "test_*" -p "Add a new function called 'helloWorld' that prints 'Hello, World!'"
```

### Save combined code and response to files

```
aidev -C code.txt -R response.txt -p "Add a new function called 'helloWorld' that prints 'Hello, World!'"
```

### Use GPT-3.5 instead of GPT-4

```
aidev -gpt35 -p "Add a new function called 'helloWorld' that prints 'Hello, World!'"
```

## Options

-d: Add code directory (defaults to ., can specify multiple times)

-i: Include only this glob pattern (can specify multiple times)

-x: Exclude this glob pattern (can specify multiple times, in case of conflict with -i longest pattern wins)

-gpt4: Use GPT-4 (default)

-gpt4-32k: Use GPT-4 32k

-gpt35: Use GPT 3.5

For more options and details, run `aidev -h`.

