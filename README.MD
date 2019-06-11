# Diversification

Writer takes 16-byte 2K3DES master key as command line argument, derives from it 
a new key based on UID of the tag, and then writes the new key to the tag. 
If second argument is provided, then it's used as the current key. 
If no second argument is provided, then the default key (BREAKMEIFYOUCAN!) is used.

Reader takes master key and reads the whole card with derived key.

This works only for Mifare Ultralight C tags.

## Dependencies
I use my [fork](https://github.com/graynk/freefare) of [fuzxxl/freefare](https://github.com/fuzxxl/freefare) 
that supports key diversification and has some ad-hoc fixes to work with latest
 version of [libfreefare](https://github.com/nfc-tools/libfreefare) (as well as some minor bug fixes).
 
Obviously this means that building and installing [libnfc](https://github.com/nfc-tools/libnfc/) and 
[libfreefare](https://github.com/nfc-tools/libfreefare) from 
master is needed too. I'd have provided a specific version that's required, but 
for some reason they have decided 
to not bump the version for the past 4 years despite continuous changes to the API, 
so expect things to be broken.

## Usage
```
sh build.sh
output/writer 00112233445566778899AABBCCDDEEFF FFEEDDCCBBAA99887766554433221100
output/reader FFEEDDCCBBAA99887766554433221100