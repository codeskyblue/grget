#! /bin/sh
#
# add-repo.sh
# Copyright (C) 2015 hzsunshx <hzsunshx@onlinegame-14-51>
#
# Distributed under terms of the MIT license.
#

cd $(dirname $0)
grep "/$1\$" repos.txt | uniq -c | sort -k1,1r | head -n1 | awk '{print $2}'
