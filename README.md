# Project Template Extension

A [GitHub CLI] extension to format a project template.

## Example

To create a new repository from a template and format it:

```bash
gh repo clone <name> --template <template> --clone
cd <template>
gh template apply
```

## License

Licensed under the [MIT](LICENSE.txt) license.

[GitHub CLI]: https://github.com/cli/cli
