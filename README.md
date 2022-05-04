# sbom
Sbom is a CLI tool to generate the Software Bill of Material (SBOM)  for container images.

It takes image string or image tarball as input.

Examples:
  Input image string
	sbom scan alpine:latest
	sbom scan alpine@sha256:51103b3f2993cbc1b45ff9d941b5d461484002792e02aa29580ec5282de719d4

  Input tarfile
	sbom scan --tar /path/to/tarfile/alpine.tar
