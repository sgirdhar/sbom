<!--
Add a note mentioning external libraties used. Mention copyright and License under which those are published.
-->

# sbom
Sbom is a CLI tool to generate the Software Bill of Material (SBOM) 
for container images.

It takes image string or image tarball as input.

## Usage:
  sbom [command]

## Examples:
  ### Possible Inputs
	sbom scan my_image:my_tag			Get the image from dockerHub. If no tag specified, default is latest.
	sbom scan my_image@my_digest			Get the image from dockerHub. No default value for digest.
	sbom scan --tar /path/to/tarfile/my_image.tar	Docker tar or OCI tar.

  ### Output Formats
	sbom scan my_image:my_tag				Default output is human readable summary table
	sbom scan my_image:my_tag --output cyclonedx		Generates CycloneDX xml output
	sbom scan my_image:my_tag --output cyclonedx-json	Generates CycloneDX json output
	
  ### Compare SBOM
	sbom scan my_image:my_tag --compare /path/to/cyclonedx-json/my_sbom.json	Generates SBOM and compares with CycloneDX json file provided

## Available Commands:
  help    :    Help about any command
  
  scan    :    Generate SBOM using image name or image tarfile

## Flags:
  -h, --help      :	help for sbom
  
  -v, --version   :	version for sbom

#### Use "sbom [command] --help" for more information about a command.
