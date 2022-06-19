# How to run
`go run .\<<language>>\<<version>>\dict_extractor.go` reads from a bi-directional dictionary and outputs an Espanso .yml config file.

This POC helps with:
- English-French bi-directional translation
- English-Chinese(simplified) bi-directional translation (todo)

The "config" in the .go scripts can be edited to generate translation configs for other language pairs.

The package structure roughly matches Espanso packages' structure. This is why we have semver package folders (see the [medicald-docs package](https://github.com/espanso/hub/tree/main/packages/medical-docs) for an example).

Feel free to raise a PR here for new language combinations, or more optimised Espanso .yml scripts.

Feel free to raise a PR like [this one](https://github.com/espanso/hub/pull/32/files) if you want to add packages to Espanso.