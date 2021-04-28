# Tail: A wrapper for `$ tail` CLI

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

This library wraps and executes command line tools, providing a line-based scanning, `$ tail`-like experience.

## support environment

linux & windows should be fine.
But tests are still limited.

## API

- tail --line x => []string, error
- tail --follow => io.Reader, <-chan error

It might be convenient to have io.Reader for streaming data. But passing errors will lose benefits for debugging and error's functionalities, such as `Unwrap` since Go 1.13.

Functions can be interrupted by canceling the context.

## Testing

Little tests. Help wanted.
