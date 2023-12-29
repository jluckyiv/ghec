# ghuc

Gloomhaven upgrade calculator (cli)

- character or summons
- multiple int (2-6)
- level int (1-9)
- prior enhancements int (0-3)

## For Summons

- move: 100
- attack: 100
- range: 50
- HP: 50

## Base +1 Enhancement

- move: 30
- attack: 50
- range: 30
- shield: 100
- push: 30
- pull: 30
- pierce: 30
- retaliate: 50
- heal: 30
- target: 50

## Base Other Effects Cost

- poison: 75
- wound: 75
- muddle: 50
- immobilize: 100
- disarm: 150
- curse: 75
- strengthen: 50
- bless: 50
- jump: 50
- element: 100
- any element: 150

Double the cost for any ability with multiple targets.

Add attack hex: 200 / number of hexes currently targeted by attack

## Level of Ability Card

1. 0
2. 25
3. 50
4. 75
5. 100
6. 125
7. 150
8. 175
9. 200

## Number of Previous Enhancements

1. 75
2. 150
3. 225

Example: The Brute wishes to enhance the top action of his ability card with a
+1 on the attack. The base cost of this is 50 gold, but it is doubled because
the attack targets multiple enemies. In addition, this is a level 3 ability
card, so another 50 gold is added, for a total of 150 gold. Afterwards, the
Brute now wants to add an attack hex to the action This would normally cost 66
gold (200 gold divided by three existing hexes, rounded down), but again, 50
gold is added because the card is level 3, and now an additional 75 gold is
added because the action was previously enhanced once, for a total 191 gold.
