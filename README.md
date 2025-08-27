# go-musthave-diploma-tpl

Template repository for the individual diploma project of the **"Go Developer"** course.

## Getting Started

1. Clone this repository into any suitable directory on your computer.  
2. In the root of the repository, run the command:

   ```bash
   go mod init <name>
   ```

   where `<name>` is the address of your GitHub repository without the `https://` prefix.  
   This will initialize your project as a Go module.

## Updating the Template

To be able to receive updates for autotests and other parts of the template, add the template repository as a remote:

```bash
git remote add -m master template https://github.com/yandex-praktikum/go-musthave-diploma-tpl.git
```

To update the autotest code, run:

```bash
git fetch template && git checkout template/master .github
```

Then commit and push the retrieved changes into your repository.
