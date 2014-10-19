#!/bin/sh

cwd=$(pwd)
codes=/codeds
prog=abstee

testoutdir=_testoutputs

[ ! -d $testoutdir ] && mkdir -p $testoutdir

testfile=$(awk 'BEGIN{
    t = strftime("%F_%T")
    gsub(/:/, "_", t)
    print "testoutput_" t
}')

echo $testfile

for c in $cwd$codes/*; do
    $prog $c 1>>$testoutdir/$testfile
done

