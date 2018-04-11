#!/bin/bash

#!/bin/bash
die() {
    echo >&2 "$@"
    exit 1
}

[ "$#" -eq 2 ] || die "Usage: rpmizer.sh <version> <build_number>"

specname=smshubbeat.spec
tmpdir=_tmp
libname=smshubbeat-$1

homesrc=`pwd`
mkdir -p $tmpdir/$libname/src/github.com/tata/smshubbeat
cp -Rf beater/ $tmpdir/$libname/src/github.com/tata/smshubbeat/
cp -Rf config/ $tmpdir/$libname/src/github.com/tata/smshubbeat/
cp -Rf main.go $tmpdir/$libname/src/github.com/tata/smshubbeat/
cp -Rf main_test.go $tmpdir/$libname/src/github.com/tata/smshubbeat/
cp -Rf Makefile $tmpdir/$libname/src/github.com/tata/smshubbeat/
cp -Rf vendor/ $tmpdir/$libname/src/github.com/tata/smshubbeat/

mkdir -p rpmbuild/{BUILD,BUILDROOT,RPMS,SOURCES,SPECS,SRPMS}
\cp -f $specname rpmbuild/SPECS
tar zcvf rpmbuild/SOURCES/$libname.tgz -C $tmpdir $libname --exclude=.git
rm -rf $tmpdir

cd rpmbuild
githome=`pwd`/BUILD/$libname
beathome=src/github.com/tata/smshubbeat
rpmbuild --define "homesrc $homesrc" --define "_topdir `pwd`" --define "githome $githome" --define "buildpath $githome/$beathome" --define "version $1" --define "buildnumber $2" --define "debug_package %{nil}" -ba SPECS/$specname
