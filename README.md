# Template

To create a new repository from this template repository for Go projects,
using the [GitHub CLI](https://github.com/cli/cli) run:

```bash
gh repo create <name> --template heaths/template-golang --public --clone
```

This will create a new repo with the given `<name>` in GitHub, copy the
`heaths/template-golang` files into that repo, and clone it into a
subdirectory of the current directory named `<name>`.
