# Go-magic

`magic` is the library behind `file`, the tool that guess file format.

**go-magic** is an implementation in pure go.

Spoiler: the files descriptions in `magic` has no grammar, just a man page, and some ambiguity. It's ugly.
Most of files described were used before you born, it's archeology (and cute).
But, `file` is the *de facto* standard, and I respect the elder.

See https://www.darwinsys.com/file/

## TODO

* [x] parse all the magic files
* [ ] handle tests hierarchy
* [ ] read values (little endian, big endian…)
* [ ] handle simple comparator
* [ ] serialize parsed files
* [ ] handle complex comparator
* [ ] handle strange commands
* [ ] functional tests with files available on my laptop
* [ ] is it fast enough ?
