#!/usr/bin/env sh

IMAGE_TO_INSPECT=${1:?'first argument is needed (expected to be url to container image, e.g. "docker.io/library/alpine:latest")'}

# examples
#
# manifestList image:
#   * docker.io/library/alpine
# NON-manifestList image:
#   * docker.io/bitnami/redis-exporter:1.43.0-debian-11-r4
# arm only:
#   * datadog/agent-arm64:7.29.1

docker run -i --rm --entrypoint=bash quay.io/skopeo/stable -s <<EOF 
    dnf install -y jq >/dev/null 2>&1;

    echo "";
    echo "  inspecting image: ${IMAGE_TO_INSPECT}";

    _rawInspect=\$( skopeo inspect --raw docker://${IMAGE_TO_INSPECT} );
    _defaultInspect=\$( skopeo inspect docker://${IMAGE_TO_INSPECT} );
    _defaultInspectDigest=\$( echo "\${_defaultInspect}" | jq --raw-output '.Digest' );
    _defaultInspectArch=\$( echo "\${_defaultInspect}" | jq --raw-output '.Architecture' );

    _isManifest=\$(echo "\${_rawInspect}" | jq '. | has("manifests")');

    if test "\${_isManifest}" = "true"; then
        echo "    * detected image to be manifestList image"
    else
        echo "    * detected image to be a NON-manifestList image"
    fi

    echo "";

    _manifestListDigest="";
    _amd64Digest=""

    if test "\${_isManifest}" = "true"; then
        _manifestListDigest=\$( skopeo manifest-digest <(echo "\${_rawInspect}") );
        _amd64Digest=\$( echo "\${_rawInspect}" | jq --raw-output ' .manifests[] | select(.platform.architecture == "amd64") | .digest ' );
    else
        _manifestListDigest="<null>                  (note: this image is a NON-manifestList image!)"
    fi

    if test "\${_isManifest}" = "false" && test "\${_defaultInspectArch}" = "amd64"; then
        _amd64Digest=\${_defaultInspectDigest};
    fi

    if test "\${_isManifest}" = "false" && test "\${_defaultInspectArch}" != "amd64"; then
        _amd64Digest="<null>                  (note: no amd64 digest found)                  ";
    fi

    echo "  ┌── Digest Summary ─────────────────────────────────────────────────────────────────────┐"
    echo "  │ manifestList: \${_manifestListDigest} │"
    echo "  ├───────────────────────────────────────────────────────────────────────────────────────┤"
    echo "  │        amd64: \${_amd64Digest} │"
    echo "  └───────────────────────────────────────────────────────────────────────────────────────┘"

    echo "";
EOF
