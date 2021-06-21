#!/usr/bin/env fish

function version_prompt
	echo "Semantic version: "
end

read --prompt version_prompt VERSION

if test (echo $VERSION | grep '^v')
	echo Use the raw semantic version, without a v prefix
	exit
end

set REV (git rev-parse --short HEAD)
echo Tagging $REV as v$VERSION
git tag --annotate v$VERSION -m "Release v$VERSION"
echo Be sure to: git push --tags
echo

set DISTDIR dist/v$VERSION
mkdir -p $DISTDIR

for pair in linux/amd64 darwin/amd64
	set GOOS   (echo $pair | cut -d'/' -f1)
	set GOARCH (echo $pair | cut -d'/' -f2)
	set FNAME  fastly-exporter-$VERSION-$GOOS-$GOARCH
	set BIN    $DISTDIR/$FNAME
	set TGZ    $DISTDIR/$FNAME.tar.gz
	echo $BIN
	env CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build -o $BIN -ldflags="-X main.programVersion=$VERSION" github.com/peterbourgon/fastly-exporter/cmd/fastly-exporter
	tar -C $DISTDIR --create --gzip --verbose --file $TGZ $FNAME
	rm $BIN
end
