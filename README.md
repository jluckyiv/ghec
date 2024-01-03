# ghec (Gloomhaven Enhancement CLI)

The CLI calculates the cost to enhance
[Gloomhaven](https://cephalofair.com/pages/gloomhaven)
ability cards.
Enhancements are subcommands. Parameters cover the card level,
the number of previous enhancements, and the number of hexes.
Summons enhancements are under the `summons` subcommand.

This CLI is a work in progress and not comprehensively tested.
The test suite covers the [example from the rulebook](#example-from-the-rulebook).

There are non-CLI solutions to calculate enhancement costs and some are
listed in the
[GloomTactics Gloomhaven Enhancement guide](https://gloomtactics.blogspot.com/2020/08/gloomhaven-enhancement-guide.html).
There are non-CLI solutions like
the [The Arcane Library](https://ninjawithkillmoon.github.io/utilities/enhancementCalculator),
the [Gloomhaven Companion app](https://play.google.com/store/apps/details?id=tomkatcreative.gloomhavenenhancementcalc&pli=1),
and the [Gloomhaven Enhancement Calculator Excel spreadsheet](https://boardgamegeek.com/filepage/145329/enhancement-calculator).

I like the convenience of a CLI, which I can run while using
[Gloomhaven Secretariat](https://ghs.champonthis.de/).

## Installation

Install using `go install`:

```sh
go install github.com/jluckyiv/ghec/ghec@latest
```

I don't anticipate other people using this CLI, so I'm not publishing
release binaries. If you want to use this CLI and don't have Go installed,
file an issue and I'll publish binaries.

## Usage

Usage is `ghec <enhancement> [flags]`. The enhancement subcommands are
available via `ghec --help`. Each command takes the same flags, which are
`--level`, `--previous`, and `--targets`. The default flag values are
`--level 1 --previous 0 --targets 0`,
which presume a level 1 card with no previous enhancements and a single target.

Summons enhancements are under the `summons` subcommand.

The `--config` flag is currently unused, but the [Cobra](https://cobra.dev/) CLI
package includes it by default, so I left it as a placeholder.

```sh
ghec attack # increase attack on a level 1 card with no previous enhancements
ghec bless # add bless to a level 1 card with no previous enhancements
ghec summons move # increase move on a level 1 summons card with no previous enhancements
ghec tui # run the TUI
```

The `ghec tui` command runs a terminal user interface (TUI) made with
[Bubble Tea](https://github.com/charmbracelet/bubbletea). Use `/` to filter
the list of base enhancements. Use `j` and `k` to move the cursor. Use `q` to
quit. The number keys select the corresponding card level. The `p` and `P`
keys increment and decrement the number of previous enhancements. To change
the number of targets, use `+` and `-`. The title bar shows the current
status and cost. Use `esc` to clear the search bar, clear the modifiers, and
to quit.

## Example from the rulebook

The rulebook has an example set of enhancements, which is reproduced below.
The test suite also covers these two examples. In the first example,

> The Brute wishes to enhance the top action of his ability card with a +1 on
> the attack. The base cost of this is 50 gold, but it is doubled because
> the attack targets multiple enemies. In addition, this is a level 3 ability
> card, so another 50 gold is added, for a total of 150 gold.

The `ghec attack` command calculates the example above:

```sh
ghec attack --level 3 --targets 3 # 150
```

In the TUI, use the `+` key to set the number of targets to 3, then use the `3
` key to set the card level to 3. Filter or scroll to the "Attack" enhancement.

> Afterwards, the Brute now wants to add an attack hex to the action. This
> would normally cost 66 gold (200 gold divided by three existing hexes,
> rounded down), but again, 50 gold is added because the card is level 3, and
> now an additional 75 gold is added because the action was previously
> enhanced once, for a total 191 gold.

For the second example, use the `ghec hex` command with the `--previous` flag.

```sh
ghec hex --level 3 --targets 3 --previous 1 # 191
```

In the TUI, the level and targets are already set. Press `p` to increment the
number of previous enhancements to 1, then filter or scroll to the "Add hex"
enhancement. The cost is shown in the title bar.

### Multiple targets and adding a hex

In the examples above, the attack is an AoE attack covering three hexes, so
the `--targets` flag is set to 3. The `--targets` numerical value does not
directly affect the cost of the first enhancement ("Attack"), but it tells the
CLI to apply the multiple targets cost multiplier (2x). In the second example,
the CLI uses the `--targets` numerical value to calculate the cost of the additional
hex (200/number of current hexes).

Using a single `-- targets` flag for both the multiple targets doubler
and to calculate the premium for adding a hex seems simpler than having a separate
`--multiple` boolean flag and `--targets` integer flag, which would have to be
combined for the second example above.
