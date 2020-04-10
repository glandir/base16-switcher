# Base16 Color Scheme Switcher

A tool to download and apply base16 color schemes and templates.

Inspired by [Base16 Universal Manager](https://github.com/pinpox/base16-universal-manager),
which seems to be incomplete.

## Installation

```
go get github.com/glandir/base16-switcher
go install github.com/glandir/base16-switcher
mkdir -p "${XDG_CONFIG_HOME:-~/.config}/base16-switcher"
touch "${XDG_CONFIG_HOME:-~/.config}/base16-switcher/config.yaml"
# Fill out config.yaml as described below
```

## Configuration

The configuration file
`${XDG_CONFIG_HOME-~/.config}/base16-switcher/config.yaml`
must exist for `base16-switcher` to run.

Template:
```yaml
default-colorscheme: "default"
scheme-sources:
    - "https://raw.githubusercontent.com/chriskempson/base16-schemes-source/master/list.yaml"

applications:
    i3:
        url: "https://github.com/khamer/base16-i3"
        files:
            default: "$XDG_CONFIG_HOME/i3/colors.config"
        hooks:
            - "cd \"$XDG_CONFIG_HOME\"/i3; cat other.config colors.config > config"
    xresources:
        url: "https://github.com/pinpox/base16-xresources"
        files:
            default: "$XDG_CONFIG_HOME/xresources/colors.xresources"
        hooks:
            - "xrdb -merge $XDG_CONFIG_HOME/colors.xresources"
            - "i3-msg reload"
```

`default-colorscheme` names the default color scheme.

`scheme-sources` lists urls to lists of color scheme repositories.
The centrally maintained list is at
<https://raw.githubusercontent.com/chriskempson/base16-schemes-source/master/list.yaml>.
Additional lists can be added if necessary.

`applications` lists templates and where to apply them.

`url` must be a template git repository.

`files` lists which of the template files in that repository to use and where to store the result.
Custom commands to finish or apply the configuration can be added under `hooks`.
They are listed under specific applications to allow logical grouping,
but all hooks are run in a separate pass after all templates have been applied.

## Usage

```
Usage: base16-switcher <command>

Commands:
  update
    Update templates and schemes.

  list
    List available color schemes.

  apply [<scheme-name>]
    Apply the named color scheme or use the default if none is specified.
```

### Update

`update` clones or pulls the git repositories
- containing color schemes,
	listed in the lists in `scheme-sources`
- containing templates,
	listed under `url` for each `application`

### List

`list` lists all locally available color schemes.

### Apply

`apply [<scheme-name>]` applies the named color scheme.

`<scheme-name>` must name one of the results of `base16-switcher list`.
If no `scheme-name` is supplied, `default-colorscheme` is used.

`apply` does the following (in this order):

- For each `application`:
	- For each entry `<name>: "<destination-path>"` in `files`:
		- Apply the color scheme to the file `templates/<name>.mustache`
			in the associated git repository.
		- Store the result in `<destination-path>`.
			Environment variables of the forms
			`$varname` and `${varname}` are substituted.
			`~` is not expanded; use `$HOME` instead.
- For each `application` (in random order):
	- For each entry `<shell-string>` in `hooks` (in order),
		execute `sh -c "<shell-string>"`.


## Example wrapper

To conveniently choose a color scheme to apply,
pass the output of `list` to a selection menu such as `rofi` or `dmenu`:

```sh
#!/bin/sh

material=$(base16-switcher list | rofi -dmenu -matching fuzzy)
# or with dmenu:
# material=$(base16-switcher list | dmenu)

if [ -z "$material" ]; then
	exit -1
fi

notify-send "Switching to color scheme $material"
base16-switcher apply "$material"
```
