# grget
[![Build Status](https://travis-ci.org/codeskyblue/grget.svg)](https://travis-ci.org/codeskyblue/grget)

very simple go online build.

## Limits
* The repo must be in github.com
* Must be golang repo
* Repo must have folder `Godeps` in top dir.

## Usage

```
curl http://grget.shengxiang.me/codeskyblue/minicdn/master/linux/amd64 -o minicdn
```

There is also another way, use a script. This will auto check OS and ARCH. (not working in windows)

```
sh grins.sh codeskyblue/minicdn
```

## grcli

Like `apt-get install`, So I create a tools named `grcli`

	$ go get -v github.com/codeskyblue/grget/grcli
	$ grcli install codeskyblue/fswatch
	# now you can run fswatch

## Notification
Please donot hack my machine.

## LICENSE
[MIT](LICENSE)
