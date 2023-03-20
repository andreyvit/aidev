# aidev

AI developer: asks GPT-4 to modify an entire folder full of files.

## Workflow

The general workflow of using aidev involves the following steps:

1. Specify the directories containing the code you want to modify using the `-d` option.
2. Provide a prompt describing the changes you want to make using the `-p` option.
3. The AI will process the request and generate modified versions of the relevant files.
4. Modified files will have a ".draft" extension added to their original file names. Review these files and, if the changes are satisfactory, replace the original files with the modified versions.

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

### Specify directories to include

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

