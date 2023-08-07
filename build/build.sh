#!/usr/bin/env bash

function build() {
	ROOT=$(dirname "$0")
	JS_ROOT=$ROOT/../web/public/js
	NAME="edge-admin"
	DIST=$ROOT/"../dist/${NAME}"
	OS=${1}
	ARCH=${2}
	TAG=${3}

	if [ -z "$OS" ]; then
		echo "usage: build.sh OS ARCH"
		exit
	fi
	if [ -z "$ARCH" ]; then
		echo "usage: build.sh OS ARCH"
		exit
	fi
	if [ -z "$TAG" ]; then
		TAG="community"
	fi

	# checking environment
	echo "checking required commands ..."
	commands=("zip" "unzip" "go" "find" "sed")
	for cmd in "${commands[@]}"; do
		if [ "$(which "${cmd}")" ]; then
			echo "checking ${cmd}: ok"
		else
			echo "checking ${cmd}: not found"
			return
		fi
	done

	VERSION=$(lookup-version "$ROOT"/../internal/const/const.go)
	ZIP="${NAME}-${OS}-${ARCH}-${TAG}-v${VERSION}.zip"

	# build edge-api
	APINodeVersion=$(lookup-version "$ROOT""/../../EdgeAPI/internal/const/const.go")
	echo "building edge-api v${APINodeVersion} ..."
	EDGE_API_BUILD_SCRIPT=$ROOT"/../../EdgeAPI/build/build.sh"
	if [ ! -f "$EDGE_API_BUILD_SCRIPT" ]; then
		echo "unable to find edge-api build script 'EdgeAPI/build/build.sh'"
		exit
	fi

	cd "$ROOT""/../../EdgeAPI/build" || exit
	echo "=============================="
	./build.sh "$OS" "$ARCH" $TAG
	echo "=============================="
	cd - || exit

    # generate files
    echo "generating files ..."
	go run -tags $TAG "$ROOT"/../cmd/edge-admin/main.go generate
	if [ "$(which uglifyjs)" ]; then
    	echo "compress to component.js ..."
    	uglifyjs --compress --mangle -- "${JS_ROOT}"/components.src.js > "${JS_ROOT}"/components.js
    	uglifyjs --compress --mangle -- "${JS_ROOT}"/utils.js > "${JS_ROOT}"/utils.min.js
    else
    	echo "copy to component.js ..."
    	cp "${JS_ROOT}"/components.src.js "${JS_ROOT}"/components.js
    	cp "${JS_ROOT}"/utils.js "${JS_ROOT}"/utils.min.js
    fi

	# create dir & copy files
	echo "copying ..."
	if [ ! -d "$DIST" ]; then
		mkdir "$DIST"
		mkdir "$DIST"/bin
		mkdir "$DIST"/configs
		mkdir "$DIST"/logs
	fi

	cp -R "$ROOT"/../web "$DIST"/
	rm -f "$DIST"/web/tmp/*
	rm -rf "$DIST"/web/public/js/components
	rm -f "$DIST"/web/public/js/components.src.js
	cp "$ROOT"/configs/server.template.yaml "$DIST"/configs/

	# change _plus.[ext] to .[ext]
	if [ "${TAG}" = "plus" ]; then
		echo "converting filenames ..."
		exts=("html" "js" "css")
		for ext in "${exts[@]}"; do
			pattern="*_plus."${ext}
			find "$DIST"/web/views -type f -name "$pattern" | \
				while read filename; do
					mv "${filename}" "${filename/_plus."${ext}"/."${ext}"}"
				done
		done
	fi

	EDGE_API_ZIP_FILE=$ROOT"/../../EdgeAPI/dist/edge-api-${OS}-${ARCH}-${TAG}-v${APINodeVersion}.zip"
	cp "$EDGE_API_ZIP_FILE" "$DIST"/
	cd "$DIST"/ || exit
	unzip -q "$(basename "$EDGE_API_ZIP_FILE")"
	rm -f "$(basename "$EDGE_API_ZIP_FILE")"
	cd - || exit

	# build
	echo "building ${NAME} ..."
	env GOOS="$OS" GOARCH="$ARCH" go build -trimpath -tags $TAG -ldflags="-s -w" -o "$DIST"/bin/${NAME} "$ROOT"/../cmd/edge-admin/main.go

	# delete hidden files
	find "$DIST" -name ".DS_Store" -delete
	find "$DIST" -name ".gitignore" -delete
	find "$DIST" -name "*.less" -delete
	#find "$DIST" -name "*.css.map" -delete
	#find "$DIST" -name "*.js.map" -delete

	# zip
	echo "zip files ..."
	cd "${DIST}/../" || exit
	if [ -f "${ZIP}" ]; then
		rm -f "${ZIP}"
	fi
	zip -r -X -q "${ZIP}" ${NAME}/
	rm -rf ${NAME}
	cd - || exit

	echo "[done]"
}

function lookup-version() {
	FILE=$1
	VERSION_DATA=$(cat "$FILE")
	re="Version[ ]+=[ ]+\"([0-9.]+)\""
	if [[ $VERSION_DATA =~ $re ]]; then
		VERSION=${BASH_REMATCH[1]}
		echo "$VERSION"
	else
		echo "could not match version"
		exit
	fi
}

build "$1" "$2" "$3"
